package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/buger/jsonparser"
)

// Decimal 四舍五入函数
func Decimal(value float64) float64 {
	value, _ = strconv.ParseFloat(fmt.Sprintf("%.2f", value), 64)
	return value
}

// Usage 帮助
func Usage() {
	fmt.Println("Usage: bwgInfo <veid> <api_key>")
}

// getHostInfo 获取kiwivm api
func getHostInfo() string {
	Time1 := time.Now().UnixNano()
	args := os.Args
	if args == nil || len(args) != 3 {
		Usage()

		return "err -- No veid or api_key!"
	}
	var veid string = args[1]
	var apiKey string = args[2]

	url := "https://api.64clouds.com/v1/getServiceInfo?veid=" + veid + "&api_key=" + apiKey
	method := "GET"

	client := &http.Client{}
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		fmt.Println(err)
	}

	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
	}
	defer res.Body.Close()
	Time2 := time.Now().UnixNano()
	hostInfo, err := ioutil.ReadAll(res.Body)
	errmessage, err := jsonparser.GetString(hostInfo, "message")
	if err == nil {
		fmt.Println("err -- " + errmessage)
		return "err -- " + errmessage
	}
	hostname, err := jsonparser.GetString(hostInfo, "hostname")
	dataNextReset, err := jsonparser.GetInt(hostInfo, "data_next_reset")
	dataCounter, err := jsonparser.GetFloat(hostInfo, "data_counter")
	var dataResetTime string = time.Unix(dataNextReset, 0).Format("2006-01-02 15:04:05")
	var useData float64 = dataCounter / 1024 / 1024 / 1024
	currentTime := time.Now()
	var nowTime string = currentTime.Format("2006-01-02 15:04:05")
	var useTime = (Time2 - Time1) / 1000000
	var log string = fmt.Sprintf("Info获取时间：%s api耗时：%d ms 主机名：%s 流量使用情况：%.2f GB / 500 GB 下次重置流量日期：%s", nowTime, useTime, hostname, Decimal(useData), dataResetTime)
	fmt.Println(log)
	var s string = fmt.Sprintf("Info获取时间：%s \napi耗时：%d ms \n主机名：%s \n流量使用情况：%.2f GB / 500 GB \n下次重置流量日期：%s", nowTime, useTime, hostname, Decimal(useData), dataResetTime)
	return s
}

// main 启动http服务
func main() {
	mux := http.NewServeMux()
	mux.Handle("/", &myHandler{})
	mux.HandleFunc("/bye", sayBye)
	log.Println("Starting bwgHostInfo httpserver")
	log.Fatal(http.ListenAndServe(":1234", mux))
}

type myHandler struct{}

func (*myHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(getHostInfo()))
}
func sayBye(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("bye bye ,this is v2 httpServer"))
}
