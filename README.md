# MercadoLibre's Golang SDK

This is the official GO SDK for MercadoLibre's Platform.

## How do I install it?

You can download the latest build at:
    https://github.com/mercadolibre/go-sdk/archive/master.zip

How do I install it using go:

Just run the following command within your $GOPATH

```bash
go get github.com/mercadolibre/sdk
```

And that's it!

## How do I start using it?

You can start integrating with api.mercadolibre by doing the following:

```go
client, err := sdk.Meli(CLIENT_ID, USER_CODE, CLIENT_SECRET, REDIRECT_URL)
```
**CLIENT_ID:** is the id that was given to you when you registered your application by using *Application Manager*
Go to "http://developers.mercadolibre.com/register-your-application/" for help

**CLIENT_SECRET:** is a secret which is created during the complition of the step above.

**USER_CODE:** is a code which is asigned to a specific user when this one tryies to access a private mercadolibre API. It means, that
to get this code, you will need to authenticate and authorize the user.

**REDIRECTA_URL:** This is the url where the user will be redirected once it is authorized by mercadolibre.

To obtain the **USER_CODE** named above, you will have to redirect the user to a specific url. To build this url, you can use
the following code:

```go
url := sdk.GetAuthURL(CLIENT_ID, sdk.MLA, "https://www.example.com")
```

As a result, you will need to somehow make the user to enter his/her credentials in that URL. Once mercadolibre api authenticates
the user, a redirection url will be returned and the **USER_CODE** will come attached to it.
 "(i.e https://www.example.com?code=TG-57f2b6c7e4b08aea0070353e-214509008)"

**Warning**: This **USER_CODE** needs to be parsed and kept by your application in order to be used for later instantiate the Meli client.

Now you can instantiate another ```Meli``` object, but this time ** this object will allow you to access the private API and also will manage the
 token refreshing, so you do not need to worrie about this handshake**


There are some design considerations worth to mention.
This SDK is just a thin layer on top of an http client to handle all the OAuth WebServer flow for you.


## Making GET calls to public API

```go
//USER_CODE can be empty since neither authorization nor authentication is needed.
client, err := sdk.Meli(CLIENT_ID, USER_CODE, CLIENT_SECRET, "www.example.com")
resp, err := client.Get("/users/me")

if err != nil {
	log.Printf("Error %s\n", err.Error())
}
userInfo, _:= ioutil.ReadAll(resp.Body)
fmt.Printf("response:%s\n", userInfo)

```


## Making GET calls to a private API
```go
var client *sdk.Meli

if  client, err := sdk.Meli(CLIENT_ID, "", CLIENT_SECRET, redirectURL); err != nil {
    log.Printf("Error: ", err.Error())
    return
}

var response *http.Response
if response, err = client.Get("/users/me"); err != nil {
    log.Printf("Error: ", err.Error())
    return
}

/*
 IF the API requires authorization you need to redirect the user.
 Once the user enters his/her credentials, you need to use the **USER_CODE** to instantiate a new client, but this
 time, it will be able to query private APIs.
*/
if response.StatusCode == http.StatusForbidden {

    url := sdk.GetAuthURL(CLIENT_ID, sdk.MLA, "www.example.com")
    log.Printf("Returning Authentication URL:%s\n", url)
    http.Redirect(w, r, url, 301)
}

/***** ONCE THE USER WAS REDIRECTED AND A USER_CODE WAS PROVIDED, THEN AGAIN..**+***/
if  client, err := sdk.Meli(CLIENT_ID, CODE_JUST_OBTEINED, CLIENT_SECRET, redirectURL); err != nil {
    log.Printf("Error: ", err.Error())
    return
}
```


## Making POST calls

```go
client, err := sdk.Meli(CLIENT_ID, CLIENT_CODE, CLIENT_SECRET, "https://www.example.com")

body :=	"{\"title\":\"Item de test - No Ofertar\",\"category_id\":\"MLA1912\",\"price\":10,\"currency_id\":\"ARS\",\"available_quantity\":1,\"buying_mode\":\"buy_it_now\",\"listing_type_id\":\"bronze\",\"condition\":\"new\",\"description\": \"Item:,  Ray-Ban WAYFARER Gloss Black RB2140 901  Model: RB2140. Size: 50mm. Name: WAYFARER. Color: Gloss Black. Includes Ray-Ban Carrying Case and Cleaning Cloth. New in Box\",\"video_id\": \"YOUTUBE_ID_HERE\",\"warranty\": \"12 months by Ray Ban\",\"pictures\":[{\"source\":\"http://upload.wikimedia.org/wikipedia/commons/f/fd/Ray_Ban_Original_Wayfarer.jpg\"},{\"source\":\"http://en.wikipedia.org/wiki/File:Teashades.gif\"}]}"

resp, err = client.Post("/items", body)

if err != nil {
    log.Printf("Error %s\n", err.Error())
}
userInfo, _= ioutil.ReadAll(resp.Body)
fmt.Printf("response:%s\n", userInfo)

```
## Making PUT calls

```go
client, err := sdk.Meli(CLIENT_ID, CLIENT_CODE, CLIENT_SECRET, "https://www.example.com")
change := "{\"available_quantity\": 6}"

resp, err = client.Put("/items/" + item.Id, &change)

if err != nil {
    log.Printf("Error %s\n", err.Error())
}
userInfo, _= ioutil.ReadAll(resp.Body)
fmt.Printf("response:%s\n", userInfo)
```
## Making DELETE calls

```go
client, err := sdk.Meli(CLIENT_ID, CLIENT_CODE, CLIENT_SECRET, "https://www.example.com")
client.Delete("/items/123")
```

## Community

You can contact us if you have questions using the standard communication channels described in the [developer's site](http://developers-forum.mercadolibre.com/)

## I want to contribute!

That is great! Just fork the project in github. Create a topic branch, write some code, and add some tests for your new code.
You can find some examples by taking a look at the main.go file.

To run the tests run ```make test```.

Thanks for helping!