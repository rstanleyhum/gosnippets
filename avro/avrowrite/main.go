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
	"github.com/linkedin/goavro"
	"log"
	"os"
	"io/ioutil"
)

func main() {
	var err error

	avroFilename := "/Users/stanley/Desktop/analysis/data2.avro"
	avroSchemaFilename := "/Users/stanley/Dev/pythonwork/github.com/humrs/avrotry/user.avsc"

	avscFile, err := ioutil.ReadFile(avroSchemaFilename)
	if err != nil {
		log.Fatal(err)
	}
	avsc := string(avscFile)

	f, err := os.Create(avroFilename)
	if err != nil {
		log.Fatal(err)
	}

	fw, err := goavro.NewWriter(
        goavro.BlockSize(10), // example; default is 10
		goavro.Compression(goavro.CompressionNull),
        goavro.WriterSchema(avsc),
        goavro.ToWriter(f))
    if err != nil {
        log.Fatal("cannot create Writer: ", err)
    }
	
	record, err := goavro.NewRecord(goavro.RecordSchema(avsc))
	if err != nil {
		log.Fatal(err)
	}
	record.Set("name", "stanley")
	record.Set("favorite_number", int32(32))
	record.Set("favorite_color", "red")
	fw.Write(record)
	
	record, err = goavro.NewRecord(goavro.RecordSchema(avsc))
	record.Set("name", "david")
	record.Set("favorite_number", int32(64))
	record.Set("favorite_color", "blue")
	
	fw.Write(record)
	
	fw.Close()
	f.Close()
}