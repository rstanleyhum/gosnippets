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
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	_ "github.com/go-sql-driver/mysql"
)

type cdwdb struct {
    db *sql.DB
    dbtype string
    dbname string
    username string
    passwd string
    hostip string
    port string
}

type node struct {
    Name string
    ParentName string
}


func (c *cdwdb) getTree(nodeID int) (results []node) {
	var name string
    var parentName string

	q := "SELECT node.name, parent.name FROM university_universityposition AS node, university_universityposition AS parent WHERE node.lft BETWEEN parent.lft AND parent.rght AND node.tree_id = 2 AND parent.tree_id = 2 AND parent.id = ? ORDER BY node.lft"

	var err error
	rows, err := c.db.Query(q, nodeID)
	if err != nil {
		log.Fatal(err)
		return
	}
	defer rows.Close()
    
    var mynode node
	for rows.Next() {
		err = rows.Scan(&name, &parentName)
		if err != nil {
			log.Fatal(err)
        }
        mynode = node{ Name: name, ParentName: parentName }
        fmt.Printf("%v\n", mynode)
		results = append(results, mynode)
	}

	if err = rows.Err(); err != nil {
		log.Fatal(err)
	}

	return results
}

func (c *cdwdb) UniPositions(w http.ResponseWriter, r *http.Request) {
	nodeID := r.URL.Query().Get("nodeid")

	id, err := strconv.Atoi(nodeID)

	results := c.getTree(id)

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
    
    var peds cdwdb
    
	peds.dbtype = "mysql"
    peds.dbname = "pediatrics"
    peds.username = "stanley"
    peds.passwd = "hum"
    peds.hostip = "192.168.99.100"
    peds.port = "3306"
    connectS := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", peds.username, peds.passwd,
        peds.hostip, peds.port, peds.dbname)
        
    peds.db, err = sql.Open(peds.dbtype, connectS)
	if err != nil {
		log.Fatalf("Error in opening sql")
	}
	defer peds.db.Close()

	http.HandleFunc("/cdwups/", peds.UniPositions)

	fmt.Printf("Server running on localhost:8081\n")
	fmt.Println(http.ListenAndServe(":8081", nil))
}
