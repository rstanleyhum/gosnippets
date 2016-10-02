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
    "fmt"
	"log"
	"net/http"
    "io/ioutil"
)

func handler(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "text/plain")
	w.Write([]byte("This is an example server.\n"))
}

func redirect(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "https://casdev.cc.columbia.edu/cas/login?service=https://pedscareerdev.cumc.columbia.edu/accounts/caslogin/", 301)
}

func validateCAS(w http.ResponseWriter, r *http.Request) {
    v := r.URL.Query()
    log.Printf("%v\n", v)
    log.Printf(r.URL.Host)
    log.Printf(r.URL.RawQuery)
    newURL := fmt.Sprintf("https://casdev.cc.columbia.edu/cas/serviceValidate?service=https://pedscareerdev.cumc.columbia.edu/accounts/caslogin/&ticket=%s", v["ticket"][0])
    log.Printf("%v\n", newURL)
    res, err := http.Get(newURL)
    if err != nil {
        log.Fatal(err)
    }
    newurlbody, err := ioutil.ReadAll(res.Body)
    res.Body.Close()
    if err != nil {
        log.Fatal(err)
    }
    log.Printf("%s", newurlbody)
}

func main() {
	http.HandleFunc("/", handler)
	http.HandleFunc("/login/", redirect)
    http.HandleFunc("/accounts/caslogin/", validateCAS)
	log.Printf("About to listen on 10443. Go to https://localhost:10443/")
	err := http.ListenAndServeTLS(":10443", "ssl/stanleyhum.crt", "ssl/stanleyhum.key", nil)
	if err != nil {
		log.Fatal(err)
	}
}

