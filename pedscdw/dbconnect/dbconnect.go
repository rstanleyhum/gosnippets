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
package main

import (
	_ "github.com/go-sql-driver/mysql"

	"database/sql"
	"fmt"
	"log"
)

func getRows(db *sql.DB) {
	var id int
	var username string
	var key string
	var year int
	var data string

	q := "Select id, username, year, 'key', data from cdwkeystore_reviewitems where username = ? and year = ?"

	rows, err := db.Query(q, "rsh2117", 2011)
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
		fmt.Printf("%v:%v:%v:%v:%v\n", id, username, year, key, data)
	}

	if err = rows.Err(); err != nil {
		log.Fatal(err)
	}
}

func main() {
	db, err := sql.Open("mysql",
		"stanley:hum@tcp(192.168.99.100:3306)/pediatrics")
	if err != nil {
		log.Fatalf("Error in opening sql")
	}
	defer db.Close()

	fmt.Printf("Successfully opened pediatrics\n")

	getRows(db)

	fmt.Printf("Finished getting rows\n")
}
