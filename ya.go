package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"runtime"
	"strings"
	"sync"
	"time"
)

var Headers map[string]string
var successCount int
var failCount int
var wait sync.WaitGroup
var url *string
var phone *string
var total *int
var param *string
var method *string
var stats runtime.MemStats

const (
	getServiceUrl  = "域名地址/technician/store/service"        //获取店铺服务
	getGoodsUrl    = "域名地址/technician/store/goods"          //获取店铺商品
	createOrderUrl = "域名地址/technician/store/dotOrder"       //点单创建订单
	payOrderUrl    = "域名地址/threeSides/dotStore/payDotOrder" //点单创建订单

)

//参数信息构造
type Request struct {
	Url    string
	Method string
	Param  string
}

func RequestUrl(info Request) map[string]interface{} {
	fmt.Println(time.StampNano)
	client := &http.Client{}
	req, err := http.NewRequest(info.Method, info.Url, strings.NewReader(info.Param))
	handerErr(err)
	for k, v := range Headers {
		req.Header.Add(k, v)
	}
	req.Header.Add("content-type", "application/x-www-form-urlencoded")
	response, err := client.Do(req)
	handerErr(err)
	defer response.Body.Close()
	msg, err := ioutil.ReadAll(response.Body)
	handerErr(err)
	parseBox := make(map[string]interface{})
	json.Unmarshal([]byte(msg), &parseBox)
	if response.StatusCode == 200 && info.Url != "域名地址/technician/userLogin" {
		successCount++
		wait.Done()
	} else if response.StatusCode != 200 && info.Url != "域名地址/technician/userLogin" {
		wait.Done()
	}
	// fmt.Println()

	return parseBox
}
func loginUrl(info Request) map[string]interface{} {
	client := &http.Client{}
	req, err := http.NewRequest(info.Method, info.Url, strings.NewReader(info.Param))
	handerErr(err)
	for k, v := range Headers {
		req.Header.Add(k, v)
	}
	req.Header.Add("content-type", "application/x-www-form-urlencoded")
	response, err := client.Do(req)
	handerErr(err)
	defer response.Body.Close()
	msg, err := ioutil.ReadAll(response.Body)
	handerErr(err)
	parseBox := make(map[string]interface{})
	json.Unmarshal([]byte(msg), &parseBox)
	if response.StatusCode == 200 && info.Url != "域名地址/technician/userLogin" {
		successCount++
	} else if response.StatusCode != 200 && info.Url != "域名地址/technician/userLogin" {
		failCount++
	}
	return parseBox
}

func init() {

	var login Request
	runtime.ReadMemStats(&stats)
	runtime.GOMAXPROCS(10)
	url = flag.String("url", "", "请求url地址")
	method = flag.String("method", "", "请求方式:post、get、put、delete")
	param = flag.String("param", "", "请求参数")
	total = flag.Int("total", 0, "并发数量")
	phone = flag.String("phone", "", "测试的用户手机号")
	flag.Parse()
	headerParam := make(map[string]string)
	login.Url = "域名地址/technician/userLogin"
	headerParam["uuid"] = "3232"
	Headers = headerParam
	login.Param = "mobile=" + *phone + "&code=32323"
	login.Method = "POST"
	loginResponse := loginUrl(login)
	if loginResponse["code"].(float64) != 0 {
		fmt.Println("登录失败:", login.Url, loginResponse["msg"])
		os.Exit(int(loginResponse["code"].(float64)))
	}
	info := loginResponse["data"].(map[string]interface{})
	headerParam["Authorization"] = info["token"].(string)
}

func main() {
	start := time.Now()
	// <-ch
	wait.Add(*total)
	for i := 1; i <= *total; i++ {
		// fmt.Println(runtime.Version(), "\r")
		// fmt.Println(runtime.NumCPU(), "\r")

		go RequestUrl(Request{Url: *url, Param: *param, Method: *method})
	}
	wait.Wait()
	// if len(ch) != 0 {
	// 	// go RequestUrl(Request{Url: *url, Param: *param, Method: *method}, ch)
	// }
	duration := time.Since(start)
	fmt.Printf("成功请求:%d\n", successCount)
	fmt.Printf("失败请求:%d\n", failCount)
	fmt.Printf("执行时长:%s\n", duration)
	fmt.Println("申请内存次数:", stats.Mallocs)
	fmt.Println("已申请且仍在使用的字节数:", stats.HeapAlloc)

	// }

}


func handerErr(msg error) {

	if msg != nil {
		panic(msg)
	}

}
