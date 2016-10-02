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
	"github.com/dgrijalva/jwt-go"
	"time"
)

func main() {
	mytoken := "eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9.eyJzdWIiOiJzaWQ6Nzc1ZGJhM2IzZDdlZjljZjBjM2UzNzk5ZWZlYzk5NzQiLCJpZHAiOiJtaWNyb3NvZnRhY2NvdW50IiwidmVyIjpbIjMiLCIzIl0sImlzcyI6Imh0dHBzOi8vb3VjaGFwcHdlYi5henVyZXdlYnNpdGVzLm5ldCIsImF1ZCI6Imh0dHBzOi8vb3VjaGFwcHdlYi5henVyZXdlYnNpdGVzLm5ldCIsImV4cCI6MTYwNjMyMzE5NSwibmJmIjoxNDUyNTMxMTk1fQ.EEw_HOw7BTKSS2grwhE317yJBddNB4OQs87GWJjRbtY"
	signingKey := []byte("E12CF7EE4AFC7CC3057FBA589A03DC564F4275D730C100467456F57EB6D5F86E")
	//signingKey := []byte("PlVTyVibvISUnNWGsiRqoaGLTUCzBJ47")
	token := jwt.New(jwt.SigningMethodHS256)
	token.Claims["foo"] = "bar"
	token.Claims["exp"] = time.Now().Add(time.Hour * 1).Unix()
	tokenString, err := token.SignedString(signingKey)
	if(err != nil) {
		fmt.Println("Error in creating tokenString")
		fmt.Print(err)
	} else {
		fmt.Printf("Token: %v\n", tokenString)
	}

	ctoken, err := jwt.Parse(mytoken, func(token *jwt.Token) (interface{}, error) {
		return signingKey, nil
	})

	if err == nil && ctoken.Valid {
		fmt.Println("Your token is valid")
	} else {
		fmt.Println("Your token is not valid.")
	}
}

