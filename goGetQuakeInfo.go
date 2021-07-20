package main

import (
	"bufio"
	"bytes"
	"crypto/tls"
	"flag"
	"fmt"
	jsoniter "github.com/json-iterator/go"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)
var(
	searchString=flag.String("s","","quake search string")
	size=flag.String("n","10","result size")
	json=jsoniter.ConfigCompatibleWithStandardLibrary
	apiKey=""
	userAgent ="Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/84.0.4147.125 Safari/537.36"
	deepSearchUrl="https://quake.360.cn/api/v3/scroll/quake_service"
)
type ServiceInfo struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    []struct {
		Time      time.Time `json:"time"`
		Transport string    `json:"transport"`
		Service   struct {
			HTTP struct {
				HTMLHash string `json:"html_hash"`
				Favicon  struct {
					Hash     string `json:"hash"`
					Location string `json:"location"`
					Data     string `json:"data"`
				} `json:"favicon"`
				Robots          string `json:"robots"`
				SitemapHash     string `json:"sitemap_hash"`
				Server          string `json:"server"`
				Body            string `json:"body"`
				XPoweredBy      string `json:"x_powered_by"`
				MetaKeywords    string `json:"meta_keywords"`
				RobotsHash      string `json:"robots_hash"`
				Sitemap         string `json:"sitemap"`
				Path            string `json:"path"`
				Title           string `json:"title"`
				Host            string `json:"host"`
				SecurityText    string `json:"security_text"`
				StatusCode      int    `json:"status_code"`
				ResponseHeaders string `json:"response_headers"`
			} `json:"http"`
			Version  string `json:"version"`
			Name     string `json:"name"`
			Product  string `json:"product"`
			Banner   string `json:"banner"`
			Response string `json:"response"`
		} `json:"service"`
		Images     []interface{} `json:"images"`
		OsName     string        `json:"os_name"`
		Components []interface{} `json:"components"`
		Location   struct {
			DistrictCn  string    `json:"district_cn"`
			ProvinceCn  string    `json:"province_cn"`
			Gps         []float64 `json:"gps"`
			ProvinceEn  string    `json:"province_en"`
			CityEn      string    `json:"city_en"`
			CountryCode string    `json:"country_code"`
			CountryEn   string    `json:"country_en"`
			Radius      float64   `json:"radius"`
			DistrictEn  string    `json:"district_en"`
			Isp         string    `json:"isp"`
			StreetEn    string    `json:"street_en"`
			Owner       string    `json:"owner"`
			CityCn      string    `json:"city_cn"`
			CountryCn   string    `json:"country_cn"`
			StreetCn    string    `json:"street_cn"`
		} `json:"location"`
		Asn       int    `json:"asn"`
		Hostname  string `json:"hostname"`
		Org       string `json:"org"`
		OsVersion string `json:"os_version"`
		IsIpv6    bool   `json:"is_ipv6"`
		IP        string `json:"ip"`
		Port      int    `json:"port"`
	} `json:"data"`
	Meta struct {
		Total        int    `json:"total"`
		PaginationID string `json:"pagination_id"`
	} `json:"meta"`
}
//X-QuakeToken: API Key
//Content-Type: application/json
func sendRequest(client *http.Client,data map[string]string)([]byte,error){
	dataBytes,err:=json.Marshal(data)
	if err!=nil{
		return nil,err
	}
	request,err:=http.NewRequest("POST",deepSearchUrl,bytes.NewBuffer(dataBytes))
	defer request.Body.Close()
	if err!=nil{
		return nil,err
	}
	defer request.Body.Close()
	request.Header.Set("Content-Type","application/json")
	request.Header.Set("X-QuakeToken",apiKey)
	request.Header.Add("User-Agent",userAgent)
	response,err:=client.Do(request)
	if err!=nil{
		return nil,err
	}
	if response.Body!=nil {
		defer response.Body.Close()
		body, err := ioutil.ReadAll(response.Body)
		if err != nil {
			return nil, err
		}
		return body, nil
	}
	return nil,nil
}
func main() {
	flag.Parse()
	if *searchString==""{
		fmt.Println("bad input search string")
		os.Exit(1)
	}
	client:=&http.Client{
		Timeout:time.Duration(10)*time.Second,
		Transport: &http.Transport{
			//参数未知影响，目前不使用
			//TLSHandshakeTimeout: time.Duration(timeout) * time.Second,
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},}
	outputFile,err:=os.OpenFile("output.txt",os.O_CREATE|os.O_TRUNC|os.O_RDWR,0666)
	if err!=nil{
		fmt.Println("fail to open output file")
		os.Exit(1)
	}
	defer outputFile.Close()
	buff:=bufio.NewWriter(outputFile)
	data:=make(map[string]string)
	honeyPots:=` AND NOT type:"蜜罐"`
	data["query"]=*searchString+honeyPots
	data["size"]=*size
	result,err:=sendRequest(client,data)
	if err!=nil{
		fmt.Println(err.Error())
		os.Exit(1)
	}
	var serviceInfo ServiceInfo
	err=json.Unmarshal(result,&serviceInfo)
	if err!=nil{
		fmt.Println(err.Error())
		os.Exit(1)
	}
	if serviceInfo.Code!=0{
		fmt.Println("bad api key")
		os.Exit(1)
	}
	allInfo:=""
	for _,single:=range serviceInfo.Data{
		output:=strings.Split(single.Service.HTTP.Favicon.Location,single.IP)[0]+single.IP+":"+strconv.Itoa(single.Port)
		fmt.Fprint(buff,output+"\r\n")
		buff.Flush()
		server:=""
		if len(single.Service.HTTP.Server)==0{
			server="noServer"
		}else {
			server=single.Service.HTTP.Server
		}
		tempStr:=output+"\t"+server+"\t"+single.Service.HTTP.Title+"\r\n"
		allInfo+=tempStr
	}
	outputFile.WriteString("\r\n"+allInfo)
}