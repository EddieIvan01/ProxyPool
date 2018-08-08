package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type data struct {
	Code    int     `json:"code"`
	Proxies []Proxy `json:"proxies"`
}

var db *sql.DB = initDB()

func initHTTP() {
	defer db.Close()
	server := http.NewServeMux()
	server.HandleFunc("/", index)
	server.HandleFunc("/proxy", proxy_api)
	log.Println("[*]Server start listen on 127.0.0.1:2333")
	err := http.ListenAndServe("0.0.0.0:2333", server)
	if err != nil {
		log.Fatal("[*]Listen on Port:2333 error")
	}
}
func index(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "OK, it works")
	fmt.Fprintln(w, "Now you can get data on /proxy")
	log.Println("[*]Client connect successfully")
}
func action_get(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	cc := queueDT(db)
	code := 200
	if len(cc) == 0 {
		code = 201
	}
	var proxy_list []Proxy
	for i := range cc {
		proxy_list = append(proxy_list, i.Proxy)
	}
	var dc data = data{
		Code:    code,
		Proxies: proxy_list,
	}
	ret, _ := json.Marshal(dc)
	log.Printf("[*]Push response: %s", string(ret))
	fmt.Fprintln(w, string(ret))
}
func action_reflush(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	checkDB(db)
	cc := checkAll(get_proxy())
	for _, i := range cc {
		replaceDT(db, i)
	}
	log.Println("[*]Reflush success")
	fmt.Fprintln(w, "[*]Reflush success")
}
func proxy_api(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	if _, ok := r.Form["act"]; !ok || len(r.Form) > 1 {
		dd, _ := json.Marshal(data{
			Code:    202,
			Proxies: nil,
		})
		fmt.Fprintln(w, string(dd))
		return
	}
	action, err1 := r.Form["act"]
	if !err1 {
		fmt.Fprintln(w, "Params error")
	}
	db := initDB()
	if action[0] == "get" {
		action_get(db, w, r)
	} else {
		if action[0] == "reflush" {
			action_reflush(db, w, r)
		} else {
			dd, _ := json.Marshal(data{
				Code:    202,
				Proxies: nil,
			})
			fmt.Fprintln(w, string(dd))
		}
	}
}
