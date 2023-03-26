package main

import (
	"fmt"
	"io"
	"net/http"
	"regexp"
	"time"
)

var (
	reMail = `\w+@\w+\.\w+\.?\w*`

	//s? 有或者没有s
	//+ 代表一次或多次
	//\s\S 各种字符
	//+? 贪婪模式
	reLink  = `href="https?://[\s\S]+?"`
	hotSpot = `title="(.*)".*</a><em>`
)

//go标准库 https://studygolang.com/pkgdoc

func HandleError(err error, where string) {
	if err != nil {
		fmt.Println(where, err)
	}
}

// 获取页面内容
func GetPageStr(url string) (pageStr string) {
	resp, err := http.Get(url)
	HandleError(err, "GetUrl")
	defer resp.Body.Close()
	bodyall, err := io.ReadAll(resp.Body)
	HandleError(err, "ioReadBody")
	pageStr = string(bodyall)
	return pageStr
}

func CanonicalFind(rule, pageStr string) (result [][]string) {
	//拿到正则对象过滤数据
	compile := regexp.MustCompile(rule)
	//根据正则对象取数据 -1 取全部
	result = compile.FindAllStringSubmatch(pageStr, -1)
	return result
}

func GetPub(url, rule, title string) {
	pageStr := GetPageStr(url)
	result := CanonicalFind(rule, pageStr)

	for _, v := range result {
		fmt.Printf("%s: %s\n", title, v[0])
	}
}

func GetNews(url, rule, title string) {
	pageStr := GetPageStr(url)
	result := CanonicalFind(rule, pageStr)
	NowD := time.Now()
	fmt.Println(title, NowD.Format("2006-01-02"), NowD.Weekday())
	for i, v := range result {
		fmt.Printf("T%d %s\n", i, v[1])
	}
}

func main() {
	//爬取邮箱
	GetPub("http://bbs.tianya.cn/post-funinfo-3328457-1.shtml", reMail, "Email")
	//爬超链接
	GetPub("http://bbs.tianya.cn/post-funinfo-3328457-1.shtml", reLink, "Link")
	//爬新闻
	GetNews("https://tech.163.com/internet/", hotSpot, "互联网快讯")

}

/*
Email: ty_qsn@163.com
Email: cynie@tom.com
Email: redsohu@qq.com
Email: 23253937@qq.com
.....
Email: 154944460@qq.com
Email: 623412487@qq.com
Email: shou198809@163.com
Email: 112241575@qq.com
Link: href="http://bbs.tianya.cn/m/post-funinfo-3328457-1.shtml"
Link: href="http://static.tianyaui.com/global/ty/TY.css"
Link: href="http://static.tianyaui.com/global/bbs/web/static/css/bbs_article_da57f9c.css"
Link: href="http://static.tianyaui.com/favicon.ico"
.....
Link: href="http://service.tianya.cn/jbts.html"
互联网快讯 2023-03-26 Sunday
T0 周鸿祎：每个行业都将拥有私有化GPT
T1 华为手机归来：能与苹果硬刚的还得是华为？
T2 小米2022年实现营收2800亿元  研发投入160亿元
T3 “华为擎云”品牌发布 华为终端商用市场驶入快车道
T4 华为P60系列发布 双星北斗卫星消息+鸿蒙3.1售价4488元起
T5 赚钱也要勒紧裤腰带过日子 微软谷歌等巨头都在大幅裁员
T6 美团联合创始王慧文退任美团执行董事，将放弃超级投票权
T7 美团2022年营收2200亿元同比增长23% 净利润28亿元
T8 OpenAI发布插件帮助ChatGPT连网 内容质量、安全性引担忧
T9 AIGC疯狂一夜！英伟达投下“核弹”、Google版ChatGPT开放，盖茨都震惊了
T10 马斯克想要在得州建乌托邦，遭不少居民反对：担忧废水排入河流
*/
