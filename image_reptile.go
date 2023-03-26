package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"sync"
)

var (
	image = `img src="(https?://[\S]+(small[\w]+).(jpg|png))"`

	chimage  = make(chan *JobImage, 2000)
	chresult = make(chan *Result, 2000)
	wr       sync.WaitGroup

	wl sync.RWMutex
)

type JobImage struct {
	id  string
	url string
}

type Result struct {
	JobImage *JobImage
	stat     bool
}

//go标准库 https://studygolang.com/pkgdoc

func HandleError(err error, where string) {
	if err != nil {
		fmt.Println(where, err)
	}
}

// 获取页面内容
func GetPageByte(url string) (pageStr []byte) {
	resp, err := http.Get(url)
	HandleError(err, "GitUrl")
	defer resp.Body.Close()
	bodyall, err := io.ReadAll(resp.Body)
	HandleError(err, "ioReadBody")
	return bodyall
}

func CanonicalFind(rule, pageStr string) (result [][]string) {
	//拿到正则对象过滤数据
	compile := regexp.MustCompile(rule)
	//根据正则对象取数据 -1 取全部
	result = compile.FindAllStringSubmatch(pageStr, -1)
	return result
}

// 1.每一页作为一个协程爬取数据 90个
// 2.图片链接放入管道
// 3.协程下载

func GetMimageUrl(url, rule string) {
	pagebyte := GetPageByte(url)
	result := CanonicalFind(rule, string(pagebyte))
	for _, v := range result {
		//fmt.Println(v, "数据写入")
		chimage <- &JobImage{
			v[2],
			v[1],
		}
	}
	wr.Done()
}

func Downloadimage(chimage chan *JobImage, chresult chan *Result, assignment int) (stat bool) {
	for i := 0; i <= assignment; i++ {
		wr.Add(1)
		go func(chimage chan *JobImage) {
			defer wr.Done()
			for url := range chimage {
				Imageall := GetPageByte(url.url)
				err := os.WriteFile("./image/"+url.id+".jpg", Imageall, 0644)
				if err != nil {
					stat = false
				} else {
					stat = true
				}
				chresult <- &Result{
					url,
					stat,
				}
			}
		}(chimage)
	}
	return
}

func Dimgresult(chresult chan *Result, assignment int) {
	for i := 0; i <= assignment; i++ {
		wr.Add(1)
		go func(chresult chan *Result) {
			defer wr.Done()
			for resultv := range chresult {
				fmt.Printf("Job名称:%s,Job状态:%t\n", resultv.JobImage.id, resultv.stat)
			}
		}(chresult)
	}
}

func main() {
	for i := 90; i >= 2; i-- {
		wr.Add(1)
		go GetMimageUrl("http://www.netbian.com/shouji/meinv/index_"+strconv.Itoa(i)+".htm", image)
	}

	Downloadimage(chimage, chresult, 30)
	Dimgresult(chresult, 30)
	wr.Wait()
}


/*
felix@MacBook-Pro 0319 % go run main.go 
Job名称:small144547RAOe21642920347,Job状态:true
Job名称:small182730F7M9o1642933650,Job状态:true
Job名称:small011150VIkp81654794710,Job状态:true
Job名称:small1836361LnP91642934196,Job状态:true
Job名称:small1630510OXR11642926651,Job状态:true
Job名称:small005254P54et1654534374,Job状态:true
Job名称:small162840stF8U1642926520,Job状态:true
Job名称:small122209i76CR1642911729,Job状态:true
...
*/
