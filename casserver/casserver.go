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
    "encoding/xml"
    "fmt"
	"html/template"
	"log"
	"net/http"
)

type casUser struct {
    XMLName xml.Name `xml:"cas:user"`
    Value string `xml:",chardata"`
}

type casAuthenticationFailure struct {
    XMLName xml.Name `xml:"cas:authenticationFailure"`
    Code string `xml:"code,attr"`
    Value string `xml:",chardata"`
}

type casAuthenticationSuccess struct {
    XMLName xml.Name `xml:"cas:authentcationSuccess"`
    User casUser `xml:"cas:user"`
}

type casServiceResponse struct {
    XMLName xml.Name `xml:"cas:serviceResponse"`
    AuthFailure []casAuthenticationFailure `xml:"cas:authenticationFailure"`
    AuthSuccess []casAuthenticationSuccess `xml:"cas:authenticationSuccess"`
}

func homePage(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain")
	w.Write([]byte("The cas server is running."))
}

func getServiceURL(r *http.Request) (string, bool) {
    v := r.URL.Query()
    service, ok := v["service"]
    if !ok || service[0] == "" {
        return "", false
    }
    return service[0], true
}

func getTicket(r *http.Request) (string, bool) {
    v:= r.URL.Query()
    ticket, ok := v["ticket"]
    if !ok || ticket[0] == "" {
        return "", false
    }
    return ticket[0], true
}

func doSendTicketGetForm(w http.ResponseWriter, callbackURL string) {
    t, err := template.ParseFiles("sendticket.gtpl")
    if err != nil {
        log.Fatal(err)
    }
    
    err = t.Execute(w, callbackURL)
    if err != nil {
        log.Fatal(err)
    }
}

func getTicketFromForm(r *http.Request) (string, bool) {
    err := r.ParseForm()
    if err != nil {
        log.Fatal(err)
    }
    
    ticket, ok := r.Form["ticket"]
    if !ok || ticket[0] == "" {
        return "", false
    }
    return ticket[0], true
}

func loginForm(w http.ResponseWriter, r *http.Request) {
	log.Printf("method: %v\n", r.Method)
    service, ok := getServiceURL(r)
    if !ok {
        log.Printf("URL: %v\n", r.URL)
        log.Printf("No service value")
    }
    callbackURL := fmt.Sprintf("/login?service=%s", service)
    
	if r.Method == "GET" {
		doSendTicketGetForm(w, callbackURL)
        return
	}
    
    ticket, ok := getTicketFromForm(r)
    if !ok {
        log.Printf("No ticket found. Redirecting back to GET\n")
        doSendTicketGetForm(w, callbackURL)
        return
    }
    
    url := fmt.Sprintf("%s?ticket=%s", service, ticket)
    http.Redirect(w, r, url, 301)
}

func doValidateTicketGetForm(w http.ResponseWriter, callbackURL, ticket string) {
    type validateData struct {
        CallbackURL string
        Ticket string
    }
    
    var data validateData
    data.CallbackURL = callbackURL
    data.Ticket = ticket
    
    t, err := template.ParseFiles("validateticket.gtpl")
    if err != nil {
        log.Fatal(err)
    }
    
    err = t.Execute(w, data)
    if err != nil {
        log.Fatal(err)
    }
}

func getAuthFailure(r *http.Request) (casAuthenticationFailure, bool) {
    err := r.ParseForm()
    if err != nil {
        log.Fatal(err)
    }
    
    var f casAuthenticationFailure
    
    authFailure, ok := r.Form["authFailure"]
    if !ok || authFailure[0] == "" {
        return f, false
    }
    
    f.Code = authFailure[0]
    
    authReason, ok2 := r.Form["failureReason"]
    if ok2 {
        f.Value = authReason[0]
    } 
    
    return f, true
}

func getAuthSuccess(r *http.Request) (casAuthenticationSuccess, bool) {
    err := r.ParseForm()
    if err != nil {
        log.Fatal(err)
    }
    
    var s casAuthenticationSuccess
    
    casuser, ok := r.Form["casuser"]
    if !ok || casuser[0] == "" {
        return s, false
    }
    
    var u casUser
    u.Value = casuser[0]
    
    s.User = u
    
    return s, true
}
func serviceValidate(w http.ResponseWriter, r *http.Request) {
	log.Printf("method: %v\n", r.Method)
    service, ok := getServiceURL(r)
    if !ok {
        log.Printf("URL: %v\n", r.URL)
        log.Printf("No service value")
    }
    callbackURL := fmt.Sprintf("/serviceValidate?service=%s", service)
    
    ticket, ok := getTicket(r)
    if !ok {
        log.Printf("URL: %v\n", r.URL)
        log.Printf("No ticket value")
    }
    
    if r.Method == "GET" {
    	doValidateTicketGetForm(w, callbackURL, ticket)
	    return
	}
    
    var v casServiceResponse
    f, hasFailure := getAuthFailure(r)
    if hasFailure {
        v.AuthFailure = append(v.AuthFailure, f)
    }
    
    var s casAuthenticationSuccess
    s, hasSuccess := getAuthSuccess(r)
    if hasSuccess {
        v.AuthSuccess = append(v.AuthSuccess, s)
    }
    
    output, err := xml.MarshalIndent(v, "", "    ")
    if err != nil {
        fmt.Printf("error: %v\n", err)
    }

    w.Header().Set("Content-Type", "text/plain")
	w.Write(output)
}
    

func main() {
	http.HandleFunc("/", homePage)
	http.HandleFunc("/login", loginForm)
	http.HandleFunc("/serviceValidate", serviceValidate)
	log.Printf("Starting on localhost:5001")
	err := http.ListenAndServe(":5001", nil)
	if err != nil {
		log.Fatal(err)
	}
}
