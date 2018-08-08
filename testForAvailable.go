package main

import (
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"regexp"
	"time"
)

var reg_proxy *regexp.Regexp = regexp.MustCompile(`[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}`)

func getTransportFieldURL(proxy_addr string) *http.Transport {
	url_proxy, _ := url.Parse(proxy_addr)
	transport := &http.Transport{
		Proxy: http.ProxyURL(url_proxy),
	}
	return transport
}
func requestUrl(proxy_addr string) (*http.Response, bool) {
	isOK := true
	transport := getTransportFieldURL(proxy_addr)
	client := &http.Client{
		Transport: transport,
		Timeout:   time.Duration(5 * time.Second),
	}
	r, _ := http.NewRequest("GET", "http://ip.chinaz.com/getip.aspx", nil)
	func() {
		(*r).Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:59.0) Gecko/20100101 Firefox/59.0")
		(*r).Header.Set("DNT", "1")
		(*r).Header.Set("Host", "ip.chinaz.com")
	}()
	resp, err := client.Do(r)
	if err != nil {
		isOK = false
	}
	return resp, isOK
}
func isWorks(testProxy Proxy) bool {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("[*]Discard invalid proxy")
		}
	}()
	var proxy_addr string
	if testProxy.SSL == "1" {
		proxy_addr = "https://" + testProxy.IP + ":" + testProxy.Port
	} else {
		proxy_addr = "http://" + testProxy.IP + ":" + testProxy.Port
	}
	response, isOK := requestUrl(proxy_addr)
	text, _ := ioutil.ReadAll(response.Body)
	current_ip := reg_proxy.Find(text)
	if testProxy.SSL == "1" {
		return isOK
	} else {
		return isOK && string(current_ip[:]) == testProxy.IP
	}
}
func _isWorks(testProxy queueProxy) bool {
	//function for delete data from database
	//can anyone tells me how to use a type to represent any struct
	defer func() {
		if r := recover(); r != nil {
			log.Printf("[*]Discard invalid proxy")
		}
	}()
	var proxy_addr string
	if testProxy.SSL == "1" {
		proxy_addr = "https://" + testProxy.IP + ":" + testProxy.Port
	} else {
		proxy_addr = "http://" + testProxy.IP + ":" + testProxy.Port
	}
	response, isOK := requestUrl(proxy_addr)
	text, _ := ioutil.ReadAll(response.Body)
	current_ip := reg_proxy.Find(text)
	if testProxy.SSL == "1" {
		return isOK
	} else {
		return isOK && string(current_ip[:]) == testProxy.IP
	}
}
