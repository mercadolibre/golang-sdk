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


**

This package allows you to interact with the Mercadolibre open platform API.
There are two main structures:
1) Client
2) Authorization

1) - This structure keeps within the secret to be used for generating the token to be sent when calling to the private APIs.
     This also provides several methods to call either public and private APIs

2) - This structure keeps the tokens and their expiration time and has to be passed by param each time a call has to be performed to any private API.
*/

package sdk

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"sync"
	"time"
)

const (
	AuthURLMLA = "https://auth.mercadolibre.com.ar" // Argentina
	AuthURLMLB = "https://auth.mercadolivre.com.br" // Brasil
	AuthURLMco = "https://auth.mercadolibre.com.co" // Colombia
	AuthURLMcr = "https://auth.mercadolibre.com.cr" // Costa Rica
	AuthURLMec = "https://auth.mercadolibre.com.ec" // Ecuador
	AuthURLMlc = "https://auth.mercadolibre.cl"     // Chile
	AuthURLMLM = "https://auth.mercadolibre.com.mx" // Mexico
	AuthURLMlu = "https://auth.mercadolibre.com.uy" // Uruguay
	AuthURLMlv = "https://auth.mercadolibre.com.ve" // Venezuela
	AuthURLMpa = "https://auth.mercadolibre.com.pa" // Panama
	AuthURLMpe = "https://auth.mercadolibre.com.pe" // Peru
	AuthURLMpt = "https://auth.mercadolivre.pt"     // Portugal
	AuthURLMrd = "https://auth.mercadolibre.com.do" // Dominicana
	AuthURlCBT = "https://global-selling.mercadolibre.com" // CBT

	AuthoricationCode = "authorization_code"
	APIURL            = "https://api.mercadolibre.com"
	RefreshToken      = "refresh_token"
)

var publicClient = &Client{apiURL: APIURL, auth: anonymous, httpClient: MeliHTTPClient{}, tokenRefresher: MeliTokenRefresher{}}
var clientByUser map[string]*Client
var clientByUserMutex sync.Mutex
var anonymous = Authorization{}
var authMutex = &sync.Mutex{}

var debugEnable = false //Set this true if you want to see debug messages

func init() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	clientByUser = make(map[string]*Client)
}

/*GetAuthURL function returns the URL for the user to authenticate and authorize*/
func GetAuthURL(clientID int64, baseSite, callback string) string {

	authURL := newAuthorizationURL(baseSite + "/authorization")
	authURL.addResponseType("code")
	authURL.addClientId(clientID)
	authURL.addRedirectURI(callback)

	return authURL.string()
}

type MeliConfig struct {
	ClientID       int64
	UserCode       string
	Secret         string
	CallBackURL    string
	HTTPClient     HTTPClient
	TokenRefresher TokenRefresher
}

/*Meli function returns a Client which can be used to call mercadolibre API.

client id, code and secret are generated when registering your application by using Application Manager

Please, visit the following link for further information: http://developers.mercadolibre.com/application-manager/

If userCode is empty, a client will be returned, but this one is only able to query
the public mercadolibre API.

If userCode has a value, then a full authenticated client will be returned. This one is able to query either public and private
mercadolibre API.
*/
func Meli(clientID int64, userCode string, secret string, callBackURL string) (*Client, error) {

	config := MeliConfig{
		ClientID:       clientID,
		UserCode:       userCode,
		Secret:         secret,
		CallBackURL:    callBackURL,
		HTTPClient:     MeliHTTPClient{},
		TokenRefresher: MeliTokenRefresher{},
	}

	return MeliClient(config)
}

/**
This function allows you to be more specific on the config you prefer giving to the sdk Client.
In case you want to use your own HttpClient or your TokenRefresher policy, you can use the following.
*/
func MeliClient(config MeliConfig) (*Client, error) {

	//If userCode is not provided, then a generic client is returned.
	//This client can be used only to access public API
	if strings.Compare(config.UserCode, "") == 0 {
		return publicClient, nil
	}

	//If we are here, userCode was provided, so a full client is going to be set up, to allow full access to either private
	//and public API
	clientByUserMutex.Lock()
	defer clientByUserMutex.Unlock()

	//The same client is going to be returned if the same applicationId and userCode is provided.
	key := strconv.FormatInt(config.ClientID, 10) + config.UserCode

	var client *Client
	client = clientByUser[key]

	if client == nil {

		client = &Client{
			id:             config.ClientID,
			code:           config.UserCode,
			secret:         config.Secret,
			redirectURL:    config.CallBackURL,
			apiURL:         APIURL,
			httpClient:     config.HTTPClient,
			tokenRefresher: config.TokenRefresher,
		}

		if debugEnable {
			log.Printf("Building a client: %p for clientid:%d code:%s\n", client, config.ClientID, config.UserCode)
		}

		auth, err := client.authorize()

		if err != nil {
			if debugEnable {
				log.Printf("error: %s", err.Error())
			}
			return nil, err
		}

		clientByUser[key] = client
		client.auth = *auth
	}

	return client, nil
}

/**
HTTP Methods
Given that error handling for all the HTTP Methods is pretty the same, then an interface Callback is define, which is
going to be called by the handler to execute the different HTTP Methods, then check the response and handle the error
*/
type Callback interface {
	Call(apiURL string) (*http.Response, error)
}

func httpErrorHandler(client *Client, resource string, httpMethod Callback) (*http.Response, error) {

	var apiURL *AuthorizationURL
	var err error

	if apiURL, err = getAuthorizedURL(client, resource); err != nil {
		if debugEnable {
			log.Printf("Error %s", err)
		}
		return nil, err
	}

	var resp *http.Response
	if resp, err = httpMethod.Call(apiURL.string()); err != nil {
		if debugEnable {
			log.Printf("Error while calling url: %s \n Error: %s", apiURL.string(), err)
		}
		return nil, err
	}

	return resp, nil
}

/*
HTTP Methods to be called by httpErrorHandler
*/
type HTTPGet struct {
	httpClient HTTPClient
}

func (callback HTTPGet) Call(url string) (*http.Response, error) {
	return callback.httpClient.Get(url)
}

type HTTPPost struct {
	httpClient HTTPClient
	body       string
}

func (callback HTTPPost) Call(url string) (*http.Response, error) {
	return callback.httpClient.Post(url, "application/json", bytes.NewReader([]byte(callback.body)))
}

type HTTPPut struct {
	httpClient HTTPClient
	body       string
}

func (callback HTTPPut) Call(url string) (*http.Response, error) {
	return callback.httpClient.Put(url, strings.NewReader(callback.body))
}

type HTTPDelete struct {
	httpClient HTTPClient
}

func (callback HTTPDelete) Call(url string) (*http.Response, error) {
	return callback.httpClient.Delete(url, nil)
}

type Client struct {
	apiURL         string
	id             int64
	secret         string
	code           string
	redirectURL    string
	auth           Authorization
	httpClient     HTTPClient
	tokenRefresher TokenRefresher
}

/*
This method returns an Authorization object which contains the needed tokens
to interact with ML API
*/
func (client *Client) authorize() (*Authorization, error) {

	authURL := newAuthorizationURL(client.apiURL + "/oauth/token")
	authURL.addGrantType(AuthoricationCode)
	authURL.addClientId(client.id)
	authURL.addClientSecret(client.secret)
	authURL.addCode(client.code)
	authURL.addRedirectURI(client.redirectURL)

	var resp *http.Response
	var err error
	if resp, err = client.httpClient.Post(authURL.string(), "application/json", *(new(io.Reader))); err != nil {
		if debugEnable {
			log.Printf("Error when posting: %s", err)
		}
		return nil, err
	}

	body, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New(fmt.Sprintf("Ups, there was an error: %s", body))
	}

	authorization := new(Authorization)
	if err := json.Unmarshal(body, authorization); err != nil {
		if debugEnable {
			log.Printf("Error while receiving the authorization %s %s", err.Error(), body)
		}
		return nil, err
	}

	authorization.ReceivedAt = time.Now().Unix()
	return authorization, nil
}

func (client *Client) refreshToken() error {
	return client.tokenRefresher.RefreshToken(client)
}

func (client *Client) Get(resourcePath string) (*http.Response, error) {

	return httpErrorHandler(client, resourcePath, HTTPGet{httpClient: client.httpClient})
}

func (client *Client) Post(resourcePath string, body string) (*http.Response, error) {

	return httpErrorHandler(client, resourcePath, HTTPPost{httpClient: client.httpClient, body: body})
}

func (client *Client) Put(resourcePath string, body string) (*http.Response, error) {

	return httpErrorHandler(client, resourcePath, HTTPPut{httpClient: client.httpClient, body: body})
}

func (client *Client) Delete(resourcePath string) (*http.Response, error) {

	return httpErrorHandler(client, resourcePath, HTTPDelete{httpClient: client.httpClient})
}

func (client Client) IsAuthorized() bool {

	return (client.auth != anonymous)
}

/*
This method returns the URL + Token to be used by each HTTP request.
If Token needs to be refreshed, then this method will send a POST to ML API to refresh it.
*/
func getAuthorizedURL(client *Client, resourcePath string) (*AuthorizationURL, error) {

	finalURL := newAuthorizationURL(client.apiURL + resourcePath)
	var err error

	if client.auth != anonymous {

		authMutex.Lock()

		if client.auth.isExpired() {

			if debugEnable {
				log.Printf("Token has expired....Refreshing it...\n")
			}

			err := client.refreshToken()

			if err != nil {
				if debugEnable {
					log.Printf("Error while refreshing token %s\n", err.Error())
				}
				return nil, err
			}
		}

		authMutex.Unlock()
		finalURL.addAccessToken(client.auth.AccessToken)
	}

	return finalURL, err
}

type Authorization struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int16  `json:"expires_in"`
	ReceivedAt   int64
	RefreshToken string `json:"refresh_token"`
	Scope        string `json:"scope"`
}

func (auth Authorization) isExpired() bool {
	if debugEnable {
		log.Printf("received at:%d expires in: %d\n", auth.ReceivedAt, auth.ExpiresIn)
	}
	return ((auth.ReceivedAt + int64(auth.ExpiresIn)) <= (time.Now().Unix() + 60))
}

/*
This struct allows adding all the params needed to the URL to be sent
to the ML API
*/
type AuthorizationURL struct {
	url bytes.Buffer
}

func (u *AuthorizationURL) addGrantType(value string) {
	u.add("grant_type=" + value)
}

func (u *AuthorizationURL) addClientId(value int64) {
	u.add("client_id=" + strconv.FormatInt(value, 10))
}

func (u *AuthorizationURL) addClientSecret(value string) {
	u.add("client_secret=" + url.QueryEscape(value))
}

func (u *AuthorizationURL) addCode(value string) {
	u.add("code=" + url.QueryEscape(value))
}

func (u *AuthorizationURL) addRedirectURI(uri string) {
	u.add("redirect_uri=" + url.QueryEscape(uri))
}

func (u *AuthorizationURL) addRefreshToken(t string) {
	u.add("refresh_token=" + url.QueryEscape(t))
}

func (u *AuthorizationURL) addResponseType(value string) {
	u.add("response_type=" + url.QueryEscape(value))
}

func (u *AuthorizationURL) addAccessToken(t string) {
	u.add("access_token=" + url.QueryEscape(t))
}

func (u *AuthorizationURL) string() string {
	return u.url.String()
}

func (u *AuthorizationURL) add(value string) {

	if !strings.Contains(u.url.String(), "?") {
		u.url.WriteString("?" + value)
	} else if strings.LastIndex("&", u.url.String()) >= u.url.Len() {
		u.url.WriteString(value)
	} else {
		u.url.WriteString("&" + value)
	}
}

func newAuthorizationURL(baseURL string) *AuthorizationURL {
	authURL := new(AuthorizationURL)
	authURL.url.WriteString(baseURL)
	return authURL
}

/**
This interface allows you to change or mock the way Meli client make HTTP Requests.
*/
type HTTPClient interface {
	Get(url string) (*http.Response, error)
	Post(url string, bodyType string, body io.Reader) (*http.Response, error)
	Put(url string, body io.Reader) (*http.Response, error)
	Delete(url string, body io.Reader) (*http.Response, error)
}

type MeliHTTPClient struct {
}

func (httpClient MeliHTTPClient) Get(url string) (*http.Response, error) {
	return http.Get(url)
}

func (httpClient MeliHTTPClient) Post(url string, bodyType string, body io.Reader) (*http.Response, error) {

	return http.Post(url, bodyType, body)
}

func (httpClient MeliHTTPClient) Put(url string, body io.Reader) (*http.Response, error) {

	return httpClient.executeHTTPRequest(http.MethodPut, url, body)
}

func (httpClient MeliHTTPClient) Delete(url string, body io.Reader) (*http.Response, error) {

	return httpClient.executeHTTPRequest(http.MethodDelete, url, body)

}

func (httpClient MeliHTTPClient) executeHTTPRequest(method string, url string, body io.Reader) (*http.Response, error) {

	req, err := http.NewRequest(method, url, body)

	if err != nil {
		if debugEnable {
			log.Printf("Error when creating %s request %s.", method, err.Error())
		}
		return nil, err
	}

	resp, err := http.DefaultClient.Do(req)

	if err != nil {
		if debugEnable {
			log.Printf("Error while calling url: %s\n Error: %s", url, err.Error())
		}
		return nil, err
	}

	return resp, nil
}

/**TokenRefresher is an interface which allows you to implement your own authentication/authorization mechanism.*/
type TokenRefresher interface {
	RefreshToken(*Client) error
}

/**MeliTokenRefresher implements ToeknRefresher interface.
This type is the default implementation provided by the SDK to deal with
Oauth token handling.
*/
type MeliTokenRefresher struct {
}

/**RefreshToken is a method which has side effects. This one, alters the token that is within the client.
Every time this method is called some locking mechanism has to be used to avoid concurrency problems when client param is modified.
*/
func (refresher MeliTokenRefresher) RefreshToken(client *Client) error {

	authorizationURL := newAuthorizationURL(client.apiURL + "/oauth/token")
	authorizationURL.addGrantType(RefreshToken)
	authorizationURL.addClientId(client.id)
	authorizationURL.addClientSecret(client.secret)
	authorizationURL.addRefreshToken(client.auth.RefreshToken)

	var resp *http.Response
	var err error

	if resp, err = client.httpClient.Post(authorizationURL.string(), "application/json", *(new(io.Reader))); err != nil {
		if debugEnable {
			log.Printf("Error: %s\n", err.Error())
		}
		return err
	}

	if resp.StatusCode != http.StatusOK {
		return errors.New("Refreshing token returned status code " + resp.Status)
	}

	body, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()

	if err := json.Unmarshal(body, &(client.auth)); err != nil {
		if debugEnable {
			log.Printf("Error while receiving the authorization %s %s", err.Error(), body)
		}
		return err
	}

	client.auth.ReceivedAt = time.Now().Unix()

	if debugEnable {
		log.Printf("auth received at: %d expires in:%d\n", client.auth.ReceivedAt, client.auth.ExpiresIn)
	}
	return nil
}
