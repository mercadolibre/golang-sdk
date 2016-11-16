/*
Copyright [2016] [mercadolibre.com]

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"sync"

	"github.com/gorilla/mux"
	"github.com/mercadolibre/golang-sdk/sdk"
)

const (
	clientID     = 2016679662291617
	clientSecret = "bA89yqE9lPeXwcZkOLBTdKGDXYFbApuZ"
	host         = "http://localhost:8080"
)

var userCode map[string]string
var userCodeMutex sync.Mutex

/*Main application method.*/
func main() {
	userCode = make(map[string]string)
	log.Fatal(http.ListenAndServe(":8080", getRouter()))
}

type item struct {
	ID string
}

type route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
}

type routes []route

/*getRouter returns a configured Router with all the paths and http supported methods*/
func getRouter() *mux.Router {

	routes := routes{

		route{
			"item",
			"GET",
			"/{userId}/items/{itemId}",
			getItem,
		},
		route{
			"item",
			"POST",
			"/{userId}/items/{itemId}",
			postItem,
		},
		route{
			"sites",
			"GET",
			"/{userId}/sites",
			getSites,
		},
		route{
			"me",
			"GET",
			"/{userId}/users/me",
			me,
		},
		route{
			"addresses",
			"GET",
			"/{userId}/users/addresses",
			addresses,
		},
		route{
			"index",
			"GET",
			"/",
			returnLinks,
		},
	}

	router := mux.NewRouter()

	for _, route := range routes {
		var handler http.Handler

		handler = route.HandlerFunc

		router.
			Methods(route.Method).
			Path(route.Pattern).
			Name(route.Name).
			Handler(handler)

	}

	return router
}

const userID = "userId"
const itemID = "itemId"

/*getItem example: performs a GET Method against items MELI API */
func getItem(w http.ResponseWriter, r *http.Request) {

	user := getParam(r, userID)
	productID := getParam(r, itemID)
	code := getUserCode(r)
	resource := "/items/" + productID
	redirectURL := host + "/" + user + "/items/" + productID

	//Getting a client to make the https://api.mercadolibre.com/items/MLU439286635
	client, err := sdk.Meli(clientID, code, clientSecret, redirectURL)

	var response *http.Response
	if response, err = client.Get(resource); err != nil {
		log.Printf("Error: ", err.Error())
		return
	}

	body, _ := ioutil.ReadAll(response.Body)
	fmt.Fprintf(w, "%s", body)
}

/*postItem example shows how to POST (publish) a new Items throught MELI Api*/
func postItem(w http.ResponseWriter, r *http.Request) {

	user := getParam(r, userID)
	productID := getParam(r, itemID)

	code := getUserCode(r)
	redirectURL := host + "/" + user + "/items/" + productID

	client, err := sdk.Meli(clientID, code, clientSecret, redirectURL)

	item := "{\"title\":\"Item de test - No Ofertar\",\"category_id\":\"MLA1912\",\"price\":10,\"currency_id\":\"ARS\",\"available_quantity\":1,\"buying_mode\":\"buy_it_now\",\"listing_type_id\":\"bronze\",\"condition\":\"new\",\"description\": \"Item:,  Ray-Ban WAYFARER Gloss Black RB2140 901  Model: RB2140. Size: 50mm. Name: WAYFARER. Color: Gloss Black. Includes Ray-Ban Carrying Case and Cleaning Cloth. New in Box\",\"video_id\": \"YOUTUBE_ID_HERE\",\"warranty\": \"12 months by Ray Ban\",\"pictures\":[{\"source\":\"http://upload.wikimedia.org/wikipedia/commons/f/fd/Ray_Ban_Original_Wayfarer.jpg\"},{\"source\":\"http://en.wikipedia.org/wiki/File:Teashades.gif\"}]}"

	response, err := client.Post("/items/", item)

	if err != nil {
		log.Printf("Error: ", err)
		return
	}
	printOutput(w, response)
}

/*getSites example shows how to GET a public MELI API*/
func getSites(w http.ResponseWriter, r *http.Request) {

	user := getParam(r, userID)
	code := getUserCode(r)
	resource := "/sites"

	redirectURL := host + "/" + user + resource
	client, err := sdk.Meli(clientID, code, clientSecret, redirectURL)

	var response *http.Response
	if response, err = client.Get(resource); err != nil {
		log.Printf("Error: ", err.Error())
		return
	}

	printOutput(w, response)
}

func me(w http.ResponseWriter, r *http.Request) {

	user := getParam(r, userID)
	code := getUserCode(r)
	resource := "/users/me"

	redirectURL := host + "/" + user + resource
	client, err := sdk.Meli(clientID, code, clientSecret, redirectURL)

	if err != nil {
		log.Printf("Error: ", err.Error())
		return
	}

	/*Example
	  If the API to be called needs authorization and authentication (private api), the the authentication URL needs to be generated.
	  Once you generate the URL and call it, you will be redirected to a ML login page where your credentials will be asked. Then, after
	  entering your credentials you will obtained a CODE which will be used to get all the authorization tokens.
	*/

	var response *http.Response
	if response, err = client.Get("/users/me"); err != nil {
		log.Printf("Error: ", err.Error())
		return
	}

	if response.StatusCode == http.StatusForbidden {

		url := sdk.GetAuthURL(clientID, sdk.AuthURLMLA, host+"/"+user+"/users/me")
		log.Printf("Returning Authentication URL:%s\n", url)
		http.Redirect(w, r, url, 301)

	}

	printOutput(w, response)
}

/*
This method responses when clicking in addresses link. After that, this will call
https://api.mercadolibre.com/users/214509008/addresses?access_token=$ACCESS_TOKEN
to get the addresses of the user.
*/
func addresses(w http.ResponseWriter, r *http.Request) {

	user := getParam(r, userID)
	code := getUserCode(r)

	resource := "/users/" + user + "/addresses"
	redirectURL := host + "/" + user + "/users/addresses"

	client, err := sdk.Meli(clientID, code, clientSecret, redirectURL)

	if err != nil {
		log.Printf("Error: ", err.Error())
		return
	}

	var response *http.Response
	if response, err = client.Get(resource); err != nil {
		log.Printf("Error: ", err.Error())
		return
	}

	/*Example
	  If the API to be called needs authorization/authentication (private api), then the authentication URL needs to be generated.
	  Once you generate the URL and call it, you will be redirected to a ML login page where your credentials will be asked. Then, after
	  entering your credentials you will obtain a CODE which will be used to get all the authorization tokens.
	*/
	if response.StatusCode == http.StatusForbidden {
		url := sdk.GetAuthURL(clientID, sdk.AuthURLMLA, redirectURL)
		body, _ := ioutil.ReadAll(response.Body)
		log.Printf("Returning Authentication URL:%s\n", url)
		log.Printf("Error:%s", body)
		http.Redirect(w, r, url, 301)
	}

	printOutput(w, response)
}

/**
This method returns the code for a specific user if it was previously sent.
*/
func getUserCode(r *http.Request) string {

	user := getParam(r, userID)
	code := r.FormValue("code")

	userCodeMutex.Lock()
	defer userCodeMutex.Unlock()

	if strings.Compare(code, "") == 0 {
		code = userCode[user]
	} else {
		userCode[user] = code
	}

	return code
}

func getParam(r *http.Request, param string) string {

	pathParams := mux.Vars(r)
	value := pathParams[param]

	if strings.Compare(value, "") == 0 {
		log.Printf("%s is missing", param)
	}

	return value
}

func printOutput(w http.ResponseWriter, response *http.Response) {
	body, _ := ioutil.ReadAll(response.Body)
	fmt.Fprintf(w, "%s", body)
}

func returnLinks(w http.ResponseWriter, r *http.Request) {

	userID := "/214509008" //WARNING: REPLACE BY YOUR USER ID
	href := "href=" + host + userID

	var links bytes.Buffer
	links.WriteString("<a " + href + "/items/MLU439286635>" + host + "/items/MLU439286635</a><br>")
	links.WriteString("<a " + href + "/sites>" + host + "/sites</a><br>")
	links.WriteString("<a " + href + "/users/me>" + host + "/users/me</a><br>")
	links.WriteString("<a " + href + "/users/addresses>" + host + "/users/addresses</a><br>")

	fmt.Fprintf(w, "%s", links.String())
}
