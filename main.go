package main

import (
	"encoding/json"
	"errors"
	"flag"
	"github.com/coreos/go-iptables/iptables"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
)

// 规则
type Rule struct {
	Id         int    `json:"id,omitempty"`
	Source     string `json:"source"`      // 源地址
	Subnet     string `json:"subnet"`      // 子网掩码
	Protocol   string `json:"protocol"`    // 协议
	TargetPort string `json:"target_port"` // 目标端口
	Action     string `json:"action"`      // 动作:accept、drop
}

// 请求格式
type ReqBody struct {
	Rules []Rule `json:"rules"` //
}

// 消息格式
type Message struct {
	Total int    `json:"total"` // 长度
	Data  []Rule `json:"data"`  // 内容
}

var (
	filepath, port *string
	ipt            *iptables.IPTables
	err            error
)

func init() {
	// 参数设置
	port = flag.String("p", "8089", "端口")   // web端口
	filepath = flag.String("f", "", "配置文件") // 配置文件
	flag.Parse()
	// 获取iptables实例
	ipt, err = iptables.New()
	if err != nil {
		log.Fatal(err)
		return
	}
}

// 校验规则
func checkRule(r Rule) bool {
	return isIPv4(r.Source) && subnetRange(r.Subnet) && portRange(r.TargetPort)
}

// 判断是否是ipv4
func isIPv4(ip string) bool {
	pattern := `^((25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\.){3}(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)$`
	regex := regexp.MustCompile(pattern)
	return regex.MatchString(ip)
}

// 判断掩码范围
func subnetRange(subnet string) bool {
	number, err := strconv.Atoi(subnet)
	if err != nil {
		return false
	}
	return 1 <= number && number <= 32
}

// 判断端口范围
func portRange(port string) bool {
	number, err := strconv.ParseInt(port, 10, 16)
	if err != nil {
		return false
	}
	return 1 <= number && number <= 65535
}

// 获取INPUT下所有规则
func getRules() (error, []Rule) {
	var r []Rule
	rules, err := ipt.List("filter", "INPUT")
	if err != nil {
		return err, nil
	}
	temp := 0
	for _, item := range rules {
		sourceRe := regexp.MustCompile(`-s (\S+)`) // 源地址提取
		sourceMatch := sourceRe.FindStringSubmatch(item)
		protocolRe := regexp.MustCompile(`-p (\S+)`) // 协议提取
		protocolMatch := protocolRe.FindStringSubmatch(item)
		portRe := regexp.MustCompile(`--dport (\S+)`) // 端口提取
		portMatch := portRe.FindStringSubmatch(item)
		actionRe := regexp.MustCompile(`-j (\S+)`) // 动作提取
		actionMatch := actionRe.FindStringSubmatch(item)
		if len(sourceMatch) > 1 {
			arr := strings.Split(sourceMatch[1], "/")
			temp = temp + 1
			r = append(r, Rule{
				Id:         temp,
				Source:     arr[0],
				Subnet:     arr[1],
				Protocol:   protocolMatch[1],
				TargetPort: portMatch[1],
				Action:     actionMatch[1],
			})
		} else {
			if len(protocolMatch) > 1 && len(portMatch) > 1 && len(actionMatch) > 1 {
				temp = temp + 1
				r = append(r, Rule{
					Id:         temp,
					Source:     "0.0.0.0",
					Subnet:     "0",
					Protocol:   protocolMatch[1],
					TargetPort: portMatch[1],
					Action:     actionMatch[1],
				})
			}
		}
	}
	return nil, r
}

// 新增INPUT规则
func addRules(r []Rule) error {
	for _, item := range r {
		if !checkRule(item) {
			return errors.New("校验失败")
		}
		if item.Action == "ACCEPT" {
			err = ipt.InsertUnique("filter", "INPUT", 1, "-s", item.Source+"/"+item.Subnet, "-p", item.Protocol, "--dport", item.TargetPort, "-j", "ACCEPT")
		}
		if item.Action == "DROP" {
			err = ipt.AppendUnique("filter", "INPUT", "-s", item.Source+"/"+item.Subnet, "-p", item.Protocol, "--dport", item.TargetPort, "-j", "DROP")
		}
		if err != nil {
			return err
		}
	}
	return nil
}

// 删除INPUT下规则
func delRules(r []Rule) error {
	for _, item := range r {
		if !checkRule(item) {
			return errors.New("校验失败")
		}
		err = ipt.Delete("filter", "INPUT", "-s", item.Source+"/"+item.Subnet, "-p", item.Protocol, "--dport", item.TargetPort, "-j", item.Action)
		if err != nil {
			return err
		}
	}
	return nil
}

func getRulesView(c *gin.Context) {
	ip_addr := c.ClientIP() // 客户端ip
	err, data := getRules()
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    -1,
			"message": "查询失败",
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"code": 0,
			"ip":   ip_addr,
			"message": Message{
				Total: len(data),
				Data:  data,
			},
		})
	}
}

func addRulesView(c *gin.Context) {
	var body ReqBody
	err = c.BindJSON(&body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    -1,
			"message": "错误的请求",
		})
	} else {
		err = addRules(body.Rules)
		if err != nil {
			c.JSON(http.StatusOK, gin.H{
				"code":    -1,
				"message": "添加失败",
			})
		} else {
			c.JSON(http.StatusOK, gin.H{
				"code":    0,
				"message": "添加成功",
			})
		}
	}

}

func delRulesView(c *gin.Context) {
	var body ReqBody
	err = c.BindJSON(&body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    -1,
			"message": "错误的请求",
		})
	} else {
		err = delRules(body.Rules)
		if err != nil {
			c.JSON(http.StatusOK, gin.H{
				"code":    -1,
				"message": "删除失败",
			})
		} else {
			c.JSON(http.StatusOK, gin.H{
				"code":    0,
				"message": "删除成功",
			})
		}
	}
}

// 首页
func index(c *gin.Context) {
	c.HTML(http.StatusOK, "index.html", nil)
}

func main() {
	if *filepath != "" {
		_, err = os.Stat(*filepath)
		if err != nil && os.IsNotExist(err) {
			log.Fatal(err)
			os.Exit(0)
		}
		data, err := ioutil.ReadFile(*filepath)
		if err != nil {
			log.Fatal(err)
			return
		}
		var r ReqBody
		if err = json.Unmarshal(data, &r); err != nil {
			log.Fatal(err)
			return
		}
		if err = addRules(r.Rules); err != nil {
			log.Fatal(err)
			return
		}
	} else {
		r := gin.Default()
		r.LoadHTMLGlob("/usr/local/share/rule/html/*.html")    // html模板资源/usr/local/share/rule/html
		r.Static("/layui", "/usr/local/share/rule/html/layui") // 静态资源目录/usr/local/share/rule/html/layui
		r.GET("/", index)
		r.GET("/rule/api/v1/get", getRulesView)
		r.POST("/rule/api/v1/add", addRulesView)
		r.POST("/rule/api/v1/del", delRulesView)
		r.Run(":" + *port)
	}
}
