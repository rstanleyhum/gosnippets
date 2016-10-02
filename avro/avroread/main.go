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
	"fmt"
	"io"
)

func main() {
	dumpReader(os.Stdin)
}

func dumpReader(r io.Reader) {
    fr, err := goavro.NewReader(goavro.BufferFromReader(r))
    if err != nil {
        log.Fatal("cannot create Reader: ", err)
    }
    defer func() {
        if err := fr.Close(); err != nil {
            log.Fatal(err)
        }
    }()

    for fr.Scan() {
        datum, err := fr.Read()
        if err != nil {
            log.Println("cannot read datum: ", err)
            continue
        }
        fmt.Println(datum)
    }
}