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
	"encoding/csv"
	"io"
	"os"
	"log"
	"flag"
	"io/ioutil"
	"strings"
	"encoding/json"
	"fmt"
	"bufio"
)

type OldRecordIndex struct {
	MRN string
	Original24Hr string
	Author string
	ObsDateTime string
}


func main() {
	oldRecords := make(map[string]string)

	oldDataName := flag.String("old", "old.csv", "Old Data Filename")
	newDataName := flag.String("cur", "current.csv", "Current Data Filename")
	outputDataName := flag.String("out", "output.csv", "Output Data Filename")

	fo, err := os.Create(*outputDataName)
	if err != nil {
		log.Fatal(err)
	}

	defer func() {
		if err:= fo.Close(); err != nil {
			log.Fatal(err)
		}
	}()

	w := bufio.NewWriter(fo)
	writer := csv.NewWriter(w)

	oldData, err := ioutil.ReadFile(*oldDataName)
	if err != nil {
		log.Fatal(err)
	}
	oldDataString := string(oldData)
	oldDataReader := csv.NewReader(strings.NewReader(oldDataString))
	for {
		record, err := oldDataReader.Read()
		if err == io.EOF {
			break
		}

		if err != nil {
			log.Fatal(err)
		}

		oldr := OldRecordIndex{record[0], record[1], record[2], record[3]}
		oldr_json, err := json.Marshal(oldr)
		if err != nil {
			log.Fatal(err)
		}
		oldr_json_string := string(oldr_json)

		outputRecord := make([]string, 0, 5)
		outputRecord = append(outputRecord, record[0], record[1], record[2], record[3], record[6])

		_, ok := oldRecords[oldr_json_string]
		if !ok {
			oldRecords[oldr_json_string] = record[6]
		} else {
			log.Println(oldr_json, oldRecords[oldr_json_string])
		}

		writer.Write(outputRecord)
		writer.Flush()

	}

	fmt.Println(len(oldRecords))


	newData, err := ioutil.ReadFile(*newDataName)
	if err != nil {
		log.Fatal(err)
	}
	newDataString := string(newData)
	newDataReader := csv.NewReader(strings.NewReader(newDataString))

	count := 0
	firstRecord := true
	for {
		record, err := newDataReader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}

		if firstRecord {
			firstRecord = false
			continue
		}

		newr := OldRecordIndex{record[0], record[12], record[17], record[18]}
		newr_json, err := json.Marshal(newr)
		if err != nil {
			log.Fatal(err)
		}

		newr_json_string := string(newr_json)

		outputRecord := make([]string, 0, 5)
		outputRecord = append(outputRecord, record[0], record[12], record[17], record[18])

		_, ok := oldRecords[newr_json_string]
		if !ok {
			count += 1
			writer.Write(outputRecord)
			writer.Flush()
		}

	}
	writer.Flush()
	fmt.Println(count)
}
