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

package sdk

import (
    "testing"
    "log"
    "fmt"
    "net/http"
    "io/ioutil"
    "strings"
    "sync"
    "io"
    "bytes"
    "net/url"
)

const (
    API_TEST = "http://localhost:3000"
    CLIENT_ID = 123456
    CLIENT_SECRET = "client secret"
    USER_CODE = "valid code with refresh token"
)

func Test_URL_for_authentication_is_properly_returned(t *testing.T) {

    expectedUrl := "https://auth.mercadolibre.com.ar/authorization?response_type=code&client_id=123456&redirect_uri=http%3A%2F%2Fsomeurl.com"

    url := GetAuthURL(CLIENT_ID, AUTH_URL_MLA, "http://someurl.com")

    if url != expectedUrl {
        log.Printf("Error: The URL is different from the one that was expected.")
        log.Printf("expected %s", expectedUrl)
        log.Printf("obtained %s", url)
        t.FailNow()
    }

}


func Test_Generic_Client_Is_Returned_When_No_UserCODE_is_given(t *testing.T) {

    client, _ := Meli(CLIENT_ID, "", CLIENT_SECRET, "htt://www.example.com")

    if client.auth != ANONYMOUS {
        log.Printf("Error: Client is not ANONYMOUS")
        t.FailNow()
    }

}

func Test_FullAuthenticated_Client_Is_Returned_When_UserCODE_And_ClientId_is_given(t *testing.T) {

    config := MeliConfig{

        ClientId: CLIENT_ID,
        UserCode: USER_CODE,
        Secret: CLIENT_SECRET,
        CallBackUrl: "http://www.example.com",
        HttpClient: MockHttpClient{},
        TokenRefresher: MockTockenRefresher{},
    }

    client, _ := MeliClient(config)

    if client == nil || client.auth == ANONYMOUS {
        log.Printf("Error: Client is not a full one")
        t.FailNow()
    }

}

func Test_GET_public_API_sites_works_properly ( t *testing.T){

    client, err := newTestAnonymousClient(API_TEST)

    if err != nil {
        log.Printf("Error:%s\n", err)
        t.FailNow()
    }
    //Public APIs do not need Authorization
    resp, err := client.Get("/sites")

    if err != nil {
        log.Printf("Error:%s\n", err)
        t.FailNow()
    }

    if resp.StatusCode != http.StatusOK {
        log.Printf("Error:Status was different from the expected one %s\n", err)
        t.FailNow()
    }

    body, err := ioutil.ReadAll(resp.Body)

    if err != nil || string(body) == ""{
        t.FailNow()
    }
}

func Test_GET_private_API_users_works_properly (t *testing.T){

    client, err := newTestClient(CLIENT_ID, USER_CODE, CLIENT_SECRET, "https://www.example.com", API_TEST)

    _, err = client.Get("/users/me")

    if err != nil {
        fmt.Printf("Error: %s\n", err)
        t.FailNow()
    }
}

func Test_POST_a_new_item_works_properly_when_token_IS_EXPIRED(t *testing.T){

    client, err := newTestClient(CLIENT_ID, USER_CODE, CLIENT_SECRET, "https://www.example.com", API_TEST)

    body := "{\"foo\":\"bar\"}"
    resp, err := client.Post("/items", body)

    if err != nil {
        log.Printf("Error while posting a new item %s\n", err)
        t.FailNow()
    }

    if resp.StatusCode != http.StatusCreated {
        log.Printf("Error while posting a new item status code: %d\n", resp.StatusCode)
        t.FailNow()
    }
}

func Test_POST_a_new_item_works_properly_when_token_IS_NOT_EXPIRED (t *testing.T){

    client, err := newTestClient(CLIENT_ID, USER_CODE, CLIENT_SECRET, "https://www.example.com", API_TEST)

    body := "{\"foo\":\"bar\"}"
    resp, err := client.Post("/items", body)

    if err != nil {
        log.Printf("Error while posting a new item %s\n", err)
        t.FailNow()
    }

    if resp.StatusCode != http.StatusCreated {
        log.Printf("Error while posting a new item status code: %d\n", resp.StatusCode)
        t.FailNow()
    }
}

func Test_PUT_a_new_item_works_properly_when_token_IS_NOT_EXPIRED (t *testing.T){

    client, err := newTestClient(CLIENT_ID, USER_CODE, CLIENT_SECRET, "https://www.example.com", API_TEST)

    body := "{\"foo\":\"bar\"}"
    resp, err := client.Put("/items/123", body)

    if err != nil {
        log.Printf("Error while posting a new item %s\n", err)
        t.FailNow()
    }

    if resp.StatusCode != http.StatusOK {
        log.Printf("Error while putting a new item. Status code: %d\n", resp.StatusCode)
        t.FailNow()
    }
}

func Test_PUT_a_new_item_works_properly_when_token_IS_EXPIRED (t *testing.T){

    client, err := newTestClient(CLIENT_ID, USER_CODE, CLIENT_SECRET, "https://www.example.com", API_TEST)

    body := "{\"foo\":\"bar\"}"
    resp, err := client.Put("/items/123", body)

    if err != nil {
        log.Printf("Error while posting a new item %s\n", err)
        t.FailNow()
    }

    if resp.StatusCode != http.StatusOK {
        log.Printf("Error while putting a new item. Status code: %d\n", resp.StatusCode)
        t.FailNow()
    }
}

func Test_DELETE_an_item_returns_200_when_token_IS_NOT_EXPIRED (t *testing.T){

    client, err := newTestClient(CLIENT_ID, USER_CODE, CLIENT_SECRET, "https://www.example.com", API_TEST)

    resp, err := client.Delete("/items/123")

    if err != nil {
        log.Printf("Error while deleting an item %s\n", err)
        t.FailNow()
    }

    if resp.StatusCode != http.StatusOK {
        log.Printf("Error while putting a new item. Status code: %d\n", resp.StatusCode)
        t.FailNow()
    }
}

func Test_DELETE_an_item_returns_200_when_token_IS_EXPIRED (t *testing.T){

    client, err := newTestClient(CLIENT_ID, USER_CODE, CLIENT_SECRET, "https://www.example.com", API_TEST)

    resp, err := client.Delete("/items/123")

    if err != nil {
        log.Printf("Error while deleting an item %s\n", err)
        t.FailNow()
    }
    if resp.StatusCode != http.StatusOK {
        log.Printf("Error while putting a new item. Status code: %d\n", resp.StatusCode)
        t.FailNow()
    }
}

func Test_AuthorizationURL_adds_a_params_separator_when_needed(t *testing.T)  {

    auth := newAuthorizationURL(API_URL+ "/authorizationauth")
    auth.addGrantType(AUTHORIZATION_CODE)

    url := API_URL + "/authorizationauth?" + "grant_type=" + AUTHORIZATION_CODE

    if strings.Compare(url, auth.string()) != 0 {
        log.Printf("url was different from what was expected\n expected: %s \n obtained: %s \n", url, auth.string())
        t.FailNow()
    }
}

func Test_AuthorizationURL_adds_a_query_param_separator_when_needed(t *testing.T)  {

    auth := newAuthorizationURL(API_URL + "/authorizationauth")
    auth.addGrantType(AUTHORIZATION_CODE)
    auth.addClientId(1213213)

    url := API_URL + "/authorizationauth?" + "grant_type=" + AUTHORIZATION_CODE + "&client_id=1213213"

    if strings.Compare(url, auth.string()) != 0 {
        log.Printf("url was different from what was expected\n expected: %s \n obtained: %s \n", url, auth.string())
        t.FailNow()
    }
}

func Test_only_one_token_refresh_call_is_done_when_several_threads_are_executed(t *testing.T){

    client, err := newTestClient(CLIENT_ID, USER_CODE, CLIENT_SECRET, "https://www.example.com", API_TEST)

    if err != nil {
        log.Printf("Error during Client instantation %s\n", err)
        t.FailNow()
    }
    client.auth.ExpiresIn = 0

    wg.Add(100)
    for i := 0; i< 100 ; i++ {
       go callHttpMethod(client)
    }
    wg.Wait()

    if counter > 1 {
        t.FailNow()
    }
}

var counter = 0;
var m = sync.Mutex{}


type MockTockenRefresher struct {}

func (mock MockTockenRefresher) RefreshToken (client *Client) error{
    realRefresher := MeliTokenRefresher{}
    realRefresher.RefreshToken(client)
    m.Lock()
    counter++
    fmt.Printf("counter %d", counter)
    m.Unlock()
    return nil
}

var wg sync.WaitGroup

func callHttpMethod(client *Client){
    defer wg.Done()
    client.Get("/users/me")
}

/*
Clients for testing purposes
 */
func newTestAnonymousClient(apiUrl string) (*Client, error) {

    client := &Client{apiUrl:apiUrl, auth:ANONYMOUS, httpClient:MockHttpClient{}}

    return client, nil
}

func newTestClient(id int64, code string, secret string, redirectUrl string, apiUrl string) (*Client, error){

    client := &Client{id:id, code:code, secret:secret, redirectUrl:redirectUrl, apiUrl:apiUrl, httpClient:MockHttpClient{}, tokenRefresher:MockTockenRefresher{}}

    auth, err := client.authorize()

    if err != nil {
        return nil, err
    }

    client.auth = *auth

    return client, nil
}

type MockHttpClient struct{

}


func (httpClient MockHttpClient) Get(url string) (*http.Response, error){

    log.Printf("Getting url %s ", url)
    resp := new (http.Response)

    if strings.Contains(url,"/sites") {
        resp.Body = ioutil.NopCloser(bytes.NewReader([]byte("[{\"id\":\"MLA\",\"name\":\"Argentina\"},{\"id\":\"MLB\",\"name\":\"Brasil\"},{\"id\":\"MCO\",\"name\":\"Colombia\"},{\"id\":\"MCR\",\"name\":\"Costa Rica\"},{\"id\":\"MEC\",\"name\":\"Ecuador\"},{\"id\":\"MLC\",\"name\":\"Chile\"},{\"id\":\"MLM\",\"name\":\"Mexico\"},{\"id\":\"MLU\",\"name\":\"Uruguay\"},{\"id\":\"MLV\",\"name\":\"Venezuela\"},{\"id\":\"MPA\",\"name\":\"Panamá\"},{\"id\":\"MPE\",\"name\":\"Perú\"},{\"id\":\"MPT\",\"name\":\"Portugal\"},{\"id\":\"MRD\",\"name\":\"Dominicana\"}]\")))")))
        resp.StatusCode = http.StatusOK
    }

    if strings.Contains(url, "/users/me") {
        resp.Body = ioutil.NopCloser(bytes.NewReader([]byte("")))
        resp.StatusCode = http.StatusOK
    }

    return resp, nil
}


func (httpClient MockHttpClient) Post(uri string, bodyType string, body io.Reader) (*http.Response, error) {

    resp := new (http.Response)
    fullUri, _ := url.Parse(uri)

    if strings.Contains(uri,"/oauth/token") {

        grant_type := fullUri.Query().Get("grant_type")

        if strings.Compare(grant_type, "authorization_code") == 0 {
            log.Printf("auth")
            code := fullUri.Query().Get("code")

            if strings.Compare(code, "bad code") == 0  {

                log.Printf("reader")
                resp.Body = ioutil.NopCloser(bytes.NewReader([]byte("{\"message\":\"Error validando el parámetro code\",\"error\":\"invalid_grant\"}")))
                resp.StatusCode = http.StatusNotFound

            } else if strings.Compare(code, "valid code without refresh token") == 0 {

                log.Printf("valid code")
                resp.Body = ioutil.NopCloser(bytes.NewReader([]byte(
                "{\"access_token\" : \"valid token\"," +
                "\"token_type\" : \"bearer\"," +
                "\"expires_in\" : 10800," +
                "\"scope\" : \"write read\"}")))

                resp.StatusCode = http.StatusOK

            } else if strings.Compare(code, "valid code with refresh token") == 0 {

                log.Printf("valid code with refresh")
                resp.Body = ioutil.NopCloser(bytes.NewReader([]byte(
                "{\"access_token\":\"valid token\"," +
                "\"token_type\":\"bearer\"," +
                "\"expires_in\":10800," +
                "\"refresh_token\":\"valid refresh token\"," +
                "\"scope\":\"write read\"}")))

            }

        } else if strings.Compare(grant_type, "refresh_token") == 0 {

            refresh := fullUri.Query().Get("refresh_token")

            if strings.Compare(refresh, "valid refresh token") == 0 {
                log.Printf("valid code with refresh")
                resp.Body = ioutil.NopCloser(bytes.NewReader([]byte(
                "{\"access_token\":\"valid token\"," +
                "\"token_type\":\"bearer\"," +
                "\"expires_in\":10800," +
                "\"scope\":\"write read\"}")))
            }
        }

        resp.StatusCode = http.StatusOK

    } else if strings.Contains(uri,"/items") {

        access_token := fullUri.Query().Get("access_token")

        if strings.Compare(access_token, "valid token") == 0 {

            b, _ := ioutil.ReadAll(body)
            if b != nil && strings.Contains(string(b),"bar") {
                resp.StatusCode = http.StatusCreated
            } else {
                resp.StatusCode = http.StatusNotFound
            }
        }
    }

    return resp, nil
}

func (httpClient MockHttpClient) Put(uri string, body io.Reader) (*http.Response, error){

    resp := new (http.Response)
    fullUri, _ := url.Parse(uri)

    if strings.Contains(uri,"/items/123") {

        access_token := fullUri.Query().Get("access_token")

        if strings.Compare(access_token, "valid token") == 0 {

            b, _ := ioutil.ReadAll(body)
            if b != nil && strings.Contains(string(b),"bar") {
                resp.StatusCode = http.StatusOK
            } else {
                resp.StatusCode = http.StatusNotFound
            }

        } else if strings.Compare(access_token, "expired token") == 0 {
            resp.StatusCode = http.StatusNotFound
        } else {
            resp.StatusCode = http.StatusForbidden
        }
    }

    return resp, nil
}

func (httpClient MockHttpClient) Delete(uri string, body io.Reader) (*http.Response, error){

    resp := new (http.Response)
    fullUri, _ := url.Parse(uri)

    if strings.Contains(uri,"/items/123") {
        access_token := fullUri.Query().Get("access_token")

        if strings.Compare(access_token, "valid token") == 0 {
            resp.StatusCode = http.StatusOK
        } else if strings.Compare(access_token, "expired token") == 0 {
            resp.StatusCode = http.StatusNotFound
        } else {
            resp.StatusCode = http.StatusForbidden
        }
    }

    return resp, nil

}
