// Copyright 2016 R. Stanley Hum
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

/*
	CDW Review Items Server Microservice
*/
package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	_ "github.com/go-sql-driver/mysql"
)

type ReviewItem struct {
	Id       int
	Username string
	Year     int
	Key      string
	Data     string
}

var (
	db *sql.DB
)

func getRows(db *sql.DB, s_username string, s_year string, s_key string) (results []ReviewItem) {
	var (
		id       int
		username string
		year     int
		key      string
		data     string
		item     ReviewItem
	)

	q := "Select id, username, year, cdwkeystore_reviewitems.key, data from cdwkeystore_reviewitems where username = ? and year = ?"

	var err error
	search_year := 2014
	if s_year != "" {
		search_year, err = strconv.Atoi(s_year)
	}
	rows, err := db.Query(q, s_username, search_year)
	if err != nil {
		log.Fatal(err)
		return
	}
	defer rows.Close()

	for rows.Next() {
		err = rows.Scan(&id, &username, &year, &key, &data)
		if err != nil {
			log.Fatal(err)
		}
		item = ReviewItem{id, username, year, key, data}
		results = append(results, item)
	}

	if err = rows.Err(); err != nil {
		log.Fatal(err)
	}

	return results
}

func ReviewItems(w http.ResponseWriter, r *http.Request) {
	username := r.URL.Query().Get("username")
	year := r.URL.Query().Get("year")
	key := r.URL.Query().Get("key")

	fmt.Printf("%v:%v:%v\n", username, year, key)

	results := getRows(db, username, year, key)

	w.Header().Set("Content-Type", "application.json; charset=utf-8")
	js, err := json.Marshal(results)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(js)
}

func main() {
	var err error
	db, err = sql.Open("mysql",
		"stanley:hum@tcp(192.168.99.100:3306)/pediatrics")
	if err != nil {
		log.Fatalf("Error in opening sql")
	}
	defer db.Close()

	http.HandleFunc("/reviewitems/", ReviewItems)

	fmt.Printf("Server running on localhost:8081\n")
	fmt.Println(http.ListenAndServe(":8081", nil))
}
