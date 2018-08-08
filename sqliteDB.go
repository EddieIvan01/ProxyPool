package main

import (
	"database/sql"
	"log"
	"os"
	"strconv"
	"sync"

	_ "github.com/mattn/go-sqlite3"
)

type queueProxy struct {
	//To delete data from database
	id int
	Proxy
}

func returnDB() (*sql.DB, bool) {
	isExist := false
	f, err := os.Open("ProxyPool.db")
	defer f.Close()
	if !(err != nil && os.IsNotExist(err)) {
		isExist = true
	}
	db, err := sql.Open("sqlite3", "./ProxyPool.db")
	if err != nil {
		log.Fatal("[*]Create DB error")
	}
	return db, isExist
}

func createTB(db *sql.DB) {
	sql_query := `
        CREATE TABLE IF NOT EXISTS proxy(
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        ip TEXT UNIQUE NOT NULL,
        port TEXT NOT NULL,
        anonymous TEXT,
        ssl TEXT
    );
    `
	db.Exec(sql_query)
}
func insertDT(db *sql.DB, ip string, port string, anonymous string, ssl string) {
	sqlry, _ := db.Prepare("INSERT INTO proxy(ip, port, anonymous, ssl) values(?, ?, ?, ?);")
	defer sqlry.Close()
	_, err := sqlry.Exec(ip, port, anonymous, ssl)
	if err != nil {
		log.Println("[*]Insert into DB error")
	}
}
func queueDT(db *sql.DB) chan queueProxy {
	proxy_for_server := make(chan queueProxy, 500)
	defer close(proxy_for_server)
	rows, err := db.Query("select * from proxy;")
	defer rows.Close()
	if err != nil {
		log.Println("[*]Get data from DB error")
	}
	var id int
	var ip string
	var port string
	var anonymous string
	var ssl string
	for rows.Next() {
		err := rows.Scan(&id, &ip, &port, &anonymous, &ssl)
		if err != nil {
			log.Println("[*]Get data from DB error")
		}
		var cc = queueProxy{
			Proxy: Proxy{
				IP:        ip,
				Port:      port,
				Anonymous: anonymous,
				SSL:       ssl,
			},
			id: id,
		}
		proxy_for_server <- cc
	}
	return proxy_for_server
}
func checkDB(db *sql.DB) {
	var wg sync.WaitGroup
	ch := queueDT(db)
	wg.Add(8)
	sqlry := "DELETE FROM proxy WHERE id="
	for i := 0; i < 8; i++ {
		go func() {
			defer wg.Done()
			for {
				if len(ch) == 0 {
					break
				}
				sh := <-ch
				if !_isWorks(sh) {
					db.Exec(sqlry + strconv.Itoa(sh.id))
				}
			}
		}()
	}
	wg.Wait()
}
func checkAll(p []Proxy) []Proxy {
	var wg sync.WaitGroup
	wg.Add(8)
	result := make([]Proxy, 0)
	ss := make(chan Proxy, len(p))
	for _, i := range p {
		ss <- i
	}
	for i := 0; i < 8; i++ {
		go func() {
			defer wg.Done()
			for {
				if len(ss) == 0 {
					break
				}
				sdc := <-ss
				if isWorks(sdc) {
					result = append(result, sdc)
				}
			}
		}()
	}
	wg.Wait()
	return result
}
func replaceDT(db *sql.DB, p Proxy) {
	sqlry, _ := db.Prepare("REPLACE INTO proxy(ip, port, anonymous, ssl) values(?, ?, ?, ?)")
	defer sqlry.Close()
	_, err := sqlry.Exec(p.IP, p.Port, p.Anonymous, p.SSL)
	if err != nil {
		log.Println("[*]Replace data error")
	}
}
func initDB() *sql.DB {
	db, isExist := returnDB()
	if !isExist {
		createTB(db)
		cc := checkAll(get_proxy())
		for _, i := range cc {
			insertDT(db, i.IP, i.Port, i.Anonymous, i.SSL)
		}
	}
	return db
}
