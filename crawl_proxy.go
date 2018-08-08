package main

import (
	"log"
	"net/http"
	"regexp"
	"strconv"
	"sync"
	"time"

	"github.com/PuerkitoBio/goquery"
)

const (
	//请求代理页面的深度
	//建议<=3
	MAXPAGE int = 1
)

var wg sync.WaitGroup

var reg_ip *regexp.Regexp = regexp.MustCompile(`^[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}$`)

type Proxy struct {
	IP        string `json:"ip"`
	Port      string `json:"port"`
	Anonymous string `json:"anonymous"`
	SSL       string `json:"ssl"`
}
type result_list chan Proxy
type task_list chan string

func catch_panic() {
	if r := recover(); r != nil {
		log.Println("[*]Error caught: %v", r)
	}
}
func request(url_ch chan string, website string) *http.Response {
	url := <-url_ch
	http_client := &http.Client{
		Timeout: time.Duration(5 * time.Second),
	}
	r, _ := http.NewRequest("GET", url, nil)
	func() {
		(*r).Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:59.0) Gecko/20100101 Firefox/59.0")
		(*r).Header.Set("DNT", "1")
		(*r).Header.Set("Host", website)
		(*r).Header.Set("Referer", "http://"+website+"/")
	}()
	req, err := http_client.Do(r)
	if err != nil {
		log.Println("[*]Request url error")
		return nil
	}
	defer catch_panic()
	return req
}

func process_xici() chan Proxy {
	result := make(result_list, 1500)
	index_url := "http://www.xicidaili.com"
	queue := make(task_list, 15)
	for i := 1; i <= MAXPAGE; i++ {
		queue <- (index_url + "/nn/" + strconv.Itoa(i))
		queue <- (index_url + "/nt/" + strconv.Itoa(i))
		queue <- (index_url + "/wn/" + strconv.Itoa(i))
		queue <- (index_url + "/wt/" + strconv.Itoa(i))
	}
	qsize := len(queue)
	wg.Add(qsize)
	for i := 0; i < qsize; i++ {
		go func() {
			defer catch_panic()
			defer wg.Done()
			resp := request(queue, "www.xicidaili.com")
			defer resp.Body.Close()
			doc, _ := goquery.NewDocumentFromResponse(resp)
			doc.Find("tr").Each(func(i int, contentSelection *goquery.Selection) {
				_IP := contentSelection.Find("td").Eq(1).Text()
				if reg_ip.Find([]byte(_IP)) != nil {
					_Port := contentSelection.Find("td").Eq(2).Text()
					_Anonymous := contentSelection.Find("td").Eq(4).Text()
					if _Anonymous == "高匿" {
						_Anonymous = "1"
					} else {
						_Anonymous = "0"
					}
					_SSL := contentSelection.Find("td").Eq(5).Text()
					if _SSL == "HTTPS" {
						_SSL = "1"
					} else {
						_SSL = "0"
					}
					append_proxy := Proxy{
						IP:        _IP,
						Port:      _Port,
						Anonymous: _Anonymous,
						SSL:       _SSL,
					}
					result <- append_proxy
				}

			})
		}()
	}
	return result
}
func process_66ip() chan Proxy {
	result := make(result_list, 100)
	index_url := "http://www.66ip.cn/"
	queue := make(task_list, 10)
	for i := 1; i <= MAXPAGE; i++ {
		queue <- (index_url + strconv.Itoa(i) + ".html")
	}
	qsize := len(queue)
	wg.Add(qsize)
	for i := 0; i < qsize; i++ {
		go func() {
			defer catch_panic()
			defer wg.Done()
			resp := request(queue, "www.66ip.cn")
			defer resp.Body.Close()
			doc, _ := goquery.NewDocumentFromResponse(resp)
			doc.Find("tr").Each(func(i int, contentSelection *goquery.Selection) {
				_IP := contentSelection.Find("td").Eq(0).Text()
				if reg_ip.Find([]byte(_IP)) != nil {
					_Port := contentSelection.Find("td").Eq(1).Text()
					_Anonymous := contentSelection.Find("td").Eq(3).Text()
					if _Anonymous == "高匿代理" {
						_Anonymous = "1"
					} else {
						_Anonymous = "0"
					}
					append_proxy := Proxy{
						IP:        _IP,
						Port:      _Port,
						Anonymous: _Anonymous,
						SSL:       "0",
					}
					result <- append_proxy
				}
			})
		}()
	}
	return result
}
func process_proxylist() chan Proxy {
	result := make(result_list, 2000)
	index_url := "https://list.proxylistplus.com/Fresh-HTTP-Proxy-List-"
	queue := make(task_list, 6)
	for i := 1; i <= MAXPAGE; i++ {
		queue <- index_url + strconv.Itoa(i)
	}
	qsize := len(queue)
	wg.Add(qsize)
	for i := 0; i < qsize; i++ {
		go func() {
			defer catch_panic()
			defer wg.Done()
			resp := request(queue, "list.proxylistplus.com")
			defer resp.Body.Close()
			doc, _ := goquery.NewDocumentFromResponse(resp)
			doc.Find("tr").Each(func(i int, contentSelection *goquery.Selection) {
				_IP := contentSelection.Find("td").Eq(1).Text()
				if reg_ip.Find([]byte(_IP)) != nil {
					_Port := contentSelection.Find("td").Eq(2).Text()
					_Anonymous := contentSelection.Find("td").Eq(3).Text()
					if _Anonymous == "anonymous" {
						_Anonymous = "1"
					} else {
						_Anonymous = "0"
					}
					_SSL := contentSelection.Find("td").Eq(6).Text()
					if _SSL == "yes" {
						_SSL = "1"
					} else {
						_SSL = "0"
					}
					append_proxy := Proxy{
						IP:        _IP,
						Port:      _Port,
						Anonymous: _Anonymous,
						SSL:       _SSL,
					}
					result <- append_proxy
				}
			})
		}()
	}
	return result
}
func get_proxy() []Proxy {
	a1 := process_xici()
	a2 := process_66ip()
	a3 := process_proxylist()
	wg.Wait()
	close(a1)
	close(a2)
	close(a3)
	size1 := len(a1)
	size2 := len(a2)
	size3 := len(a3)
	data_to_return := make([]Proxy, size1+size2+size3)
	for i := 0; i < size1; i++ {
		data_to_return[i] = <-a1
	}
	for i := size1; i < size1+size2; i++ {
		data_to_return[i] = <-a2
	}
	for i := size1 + size2; i < size1+size2+size3; i++ {
		data_to_return[i] = <-a3
	}
	return data_to_return
}
