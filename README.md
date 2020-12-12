# TinySRC GO SDK

This is a Go package for accessing [TinySRC API](http://api.tinysrc.me/).

[![License](https://img.shields.io/badge/License-MIT-blue.svg)](https://opensource.org/licenses/MIT)

## What is TinySRC?

Free Tiny URL Shortener Service with a detailed statistic and extended features.

## How to get Started

To get started working with TinySRC GO SDK, you need register and get your personal API Key first.  

## Installing TinySRC-GO-SDK
```go
go get github.com/dmitrypro77/tinysrc-api-sdk
```

## How to use TinySRC GO SDK

### Create new instance
```go
var err models.ErrorResponse
client, e := tinysrc.NewClient(context.Background(), "apiKey", http.DefaultClient)

if e != nil {
    panic(e.Error())
}
```

### Get Current User Info
```go
user, err := client.GetCurrentUser()

if len(err.Errors) > 0 {
    panic(err)
}

fmt.Println(user.Username)
fmt.Println(user.Email)
fmt.Println(user.ApiKey)
// etc ...
```

### Create shortened url
```go
linkRequest := models.LinkRequest{
    Url:            "http://test.com",
    AuthRequired:   0, // if auth required
    Password:       "", // optional if access to link by password needed
    ExpirationTime: "",  // optional time.Now().Format(tinysrc.DATE_FORMAT)
}

link, err := client.CreateShortLink(linkRequest)

if len(err.Errors) > 0 {
    // You Can access to validation messages like this
    fmt.Println(err.Validations)
    panic(err)
}

fmt.Print(link.Url)
fmt.Print(link.StatUrl)
fmt.Print(link.StatPassword)
// etc ...
```

### Get List URLs
```go
request := models.ListUrlsRequest{
    Limit: 10,
    Page:  1,
    Query: "", // optional if you need search by hash for example
}

urls, err := client.GetListUrls(request)

if len(err.Errors) > 0 {
    panic(err)
}

// Total Urls
fmt.Println(urls.Total)

for _, url := range urls.Data {
    fmt.Println(url.Url)
    //fmt.Println(url.Clicks)
    //fmt.Println(url.Bots)
    //fmt.Println(url.QRCode)
    // etc ...
}
```

### Get URL Details By Hash
```go
url, err := client.GetUrlByHash("test")

if len(err.Errors) > 0 {
    panic(err)
}

fmt.Println(url.Url)
fmt.Println(url.Bots)
fmt.Println(url.Clicks)
// etc...
```

### Activate/Deactivate Link
```go
status, err := client.SetActive("test", &models.LinkActivationRequest{Active: false})

if len(err.Errors) > 0 {
    panic(err)
}

fmt.Println(status)
```

### Get Statistic By Hash

```go
stat, err := client.GetStatByHash("test", models.StatRequest{
    Limit:     10,
    Page:      1,
    DateStart: time.Date(2015, 1,1,1,1,1,1, time.UTC),
    DateEnd:   time.Now(),
})

if len(err.Errors) > 0 {
    panic(err)
}

// Total Records
fmt.Println(stat.Total)

for _, s := range stat.Data {
    fmt.Println(s.Ip)
    fmt.Println(s.Bot)
    fmt.Println(s.Mobile)
    fmt.Println(s.Browser)
    fmt.Println(s.Os)
    fmt.Println(s.Platform)
    fmt.Println(s.Created)
    // etc...
}
```


## Tests
```go
go test --cover
```

## Contributing

Please note that this package still is in development.

**Pull requests are welcome**.

## Links

- https://tinysrc.me/

## License

[MIT]