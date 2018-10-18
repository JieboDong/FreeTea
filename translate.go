package main

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
)

var host = "https://fanyi.sogou.com/reventondc/api/sogouTranslate"
var pid = "" //搜狗翻译申请的pid
var key = "" //key

func main() {

	err, content := translate("测试数据", "mg", "zh-CHS")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(content)

}

//发送翻译请求
func translate(text string, to string, from string) (error, string) {

	param := url.Values{"q": {text}, "from": {from}, "to": {to}, "pid": {pid}, "salt": {randStr()}, "sign": {sign(text)}}
	res, err := http.PostForm(host, param)
	if err != nil {
		return err, ""
	}
	//读取完数据关闭回复主体
	defer res.Body.Close()
	//读取响应数据主体
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return err, ""
	}
	return nil, string(body)

}

//签名
func sign(text string) string {

	sign := md5.New()
	waitStr := pid + text + randStr() + key
	sign.Write([]byte(waitStr))
	signStr := sign.Sum(nil)
	return hex.EncodeToString(signStr)
}

//生成随机字符  测试暂定固定的，可以自己重写一下
func randStr() string {

	return "suijizifucuan"

}
