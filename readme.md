[![GitHub commits since latest release (branch)](https://img.shields.io/github/commits-since/KnicKnic/go-powershell/latest.svg)](https://github.com/ten16thomasg/crypter/latest)

# Goal

The goal of this project is to showcase Crypter-API and Crypter-Clients, a containerized GO micro-service with portable GO binaries that work in tandem with Crypt to provide an end to end, client -> server, secrets management solution. I am running the server components of Crypt and Crypter-API on AWS EC2 with Docker, behind an HTTPS Nginx proxy. I also use AWS RDS to serve the Postgres database. Then As an added bonus I extend the capabilities of the Default CRUD Crypt Django Application written by Graham https://github.com/grahamgilbert/Crypt-Server by providing a front-end API to manage what data is sent to crypt.

## Features
**Crypter-API**
  - Dockerized GO API Endpoint serving as a trigger for Crypter Clients. Host in cloud of choice and simply docker-compose up to build.
    - Get Events from Crypter-API
    - Post Events to Crypter-API
    - Delete Events from Crypter-API

**Crypter-Client**
- LAPS 2 Crypt. Simply update config file with desired encryption strength. Deploy the binary and let crypter do the rest. 
  - Creates a local admin account. 
  - Rotates the credentials 
  - Sends key to Crypt Server

- Crypter Locker. Binary Subscribes to the crypter-api and locks or unlocks an endpoint based on status returned. 
  - Rotate local group membership
  - Rotate local admin account credentials
  - Clear cached logons
  - Restart Device 

# Status

It works

1. This project is not api stable, however I believe it will be simple if you do use the current api to migrate to any future changes.

**TODO:**

- [x] support for default loggers, like glog or log (in separate package)
- [x] add some code for easy back and forth json serialization of complex objects
- [X] more examples / tests
- [X] add MacOS and Linux Support
- [X] example / helper classes around exception
- [x] a doc overview
- [x] plumb through callback handler for each specific logging type (verbose, debug, warning, ...)

# Example Usage

<img src="crypter.png"
     alt="Markdown Monster icon"
     style="float: left; margin-right: 20px;" />

## How to run Crypter-API

Pull down the Crypter Docker Image
```shell script
docker pull github.com/ten16thomasg:crypter-api
```

Clone the Crypter github repo
```shell script
git clone github.com/ten16thomasg/crypter-api
```

Change Directory into ./crypter-api
```shell script
cd ./crypter-api
```

Install Crypter project dependencies

```shell script
go run github.com/ten16thomasg/crypter-api
```
or
```shell script
go build -o ./bin/main
```

Build the image using the command in our root directory i.e. ./crypter-api
```shell script
docker-compose up
```

### How to test
```shell script
go test -v server.go main.go handlers_test.go  -covermode=count  -coverprofile=./bin/coverage.out
```

### Initial Structure
```
    .
    ├── bin
    │   ├── coverage.out
    │   └── main
    ├── errors
    │   └── errors.go
    ├── handlers
    │   └── handlers.go
    ├── objects
    │   ├── event.go
    │   └── requests.go
    ├── store
    │   ├── postgres.go
    │   └── store.go
    ├── test
    │   ├── main.go
    │   └── test.go
    ├── .gitignore
    ├── docker-compose.yml
    ├── Dockerfile
    ├── go.mod
    ├── main.go
    ├── README.md
    └── server.go
```

#### Rest api
**Object: Event**
```go
package objects

import (
	"time"
)

// EventStatus defines the status of the event
type EventStatus string

const (
	// Some default event status
	Original    EventStatus = "original"
)

type TimeSlot struct {
	StartTime time.Time `json:"start_time,omitempty"`
	EndTime   time.Time `json:"end_time,omitempty"`
}

// Event object for the API
type Event struct {
	// Identifier
	ID string `gorm:"primary_key" json:"id,omitempty"`

	// General details
	Name        string `json:"name,omitempty"`
	Platform string `json:"Platform,omitempty"`
	Source     string `json:"source,omitempty"`
	State     string `json:"state,omitempty"`
	SerialNumber string `json:"serial_number,omitempty"`

	// Event slot duration
	Slot *TimeSlot `gorm:"embedded" json:"slot,omitempty"`

	// Change status
	Status EventStatus `json:"status,omitempty"`

	// Meta information
	CreatedOn     time.Time `json:"created_on,omitempty"`
	UpdatedOn     time.Time `json:"updated_on,omitempty"`
}
```

#### Endpoints

**Get all events**
```go
package main

import (
  "fmt"
  "net/http"
  "io/ioutil"
)

func main() {

  url := "http://localhost:8080/api/v1/events"
  method := "GET"

  client := &http.Client {
  }
  req, err := http.NewRequest(method, url, nil)

  if err != nil {
    fmt.Println(err)
    return
  }
  res, err := client.Do(req)
  if err != nil {
    fmt.Println(err)
    return
  }
  defer res.Body.Close()

  body, err := ioutil.ReadAll(res.Body)
  if err != nil {
    fmt.Println(err)
    return
  }
  fmt.Println(string(body))
}
###
```

**Get Single event**
```go
package main

import (
  "fmt"
  "net/http"
  "io/ioutil"
)

func main() {

  url := "http://localhost:8080/api/v1/event?id=1650030060-0643970700-9740512683"
  method := "GET"

  client := &http.Client {
  }
  req, err := http.NewRequest(method, url, nil)

  if err != nil {
    fmt.Println(err)
    return
  }
  res, err := client.Do(req)
  if err != nil {
    fmt.Println(err)
    return
  }
  defer res.Body.Close()

  body, err := ioutil.ReadAll(res.Body)
  if err != nil {
    fmt.Println(err)
    return
  }
  fmt.Println(string(body))
}
###
```

**Create an event**
```go
package main

import (
  "fmt"
  "strings"
  "net/http"
  "io/ioutil"
)

func main() {

  url := "http://localhost:8080/api/v1/event"
  method := "POST"

  payload := strings.NewReader(`{
    "name": "Desktop-12345",
    "Platform": "Windows",
    "slot": {
        "start_time": "2020-12-11T09:00:00+05:30",
        "end_time": "2020-12-11T15:00:00+05:30"
    },
    "source": "Crypter Client",
    "state": "Lock",
    "serial_number": "698HG356I"
}`)

  client := &http.Client {
  }
  req, err := http.NewRequest(method, url, payload)

  if err != nil {
    fmt.Println(err)
    return
  }
  req.Header.Add("Content-Type", "application/json")

  res, err := client.Do(req)
  if err != nil {
    fmt.Println(err)
    return
  }
  defer res.Body.Close()

  body, err := ioutil.ReadAll(res.Body)
  if err != nil {
    fmt.Println(err)
    return
  }
  fmt.Println(string(body))
}
###
```

**List at max 40 events**
```go
package main

import (
  "fmt"
  "net/http"
  "io/ioutil"
)

func main() {

  url := "http://localhost:8080/api/v1/events?limit=40"
  method := "GET"

  client := &http.Client {
  }
  req, err := http.NewRequest(method, url, nil)

  if err != nil {
    fmt.Println(err)
    return
  }
  res, err := client.Do(req)
  if err != nil {
    fmt.Println(err)
    return
  }
  defer res.Body.Close()

  body, err := ioutil.ReadAll(res.Body)
  if err != nil {
    fmt.Println(err)
    return
  }
  fmt.Println(string(body))
}
###
```

**Update event details**
```go
package main

import (
  "fmt"
  "strings"
  "net/http"
  "io/ioutil"
)

func main() {

  url := "http://localhost:8080/api/v1/event"
  method := "PUT"

  payload := strings.NewReader(`{
    "name": "Desktop-12345",
    "Platform": "Windows",
    "source": "Crypter Client",
    "state": "Unlock",
    "serial_number": "698HG356I"
}`)

  client := &http.Client {
  }
  req, err := http.NewRequest(method, url, payload)

  if err != nil {
    fmt.Println(err)
    return
  }
  req.Header.Add("Content-Type", "application/json")

  res, err := client.Do(req)
  if err != nil {
    fmt.Println(err)
    return
  }
  defer res.Body.Close()

  body, err := ioutil.ReadAll(res.Body)
  if err != nil {
    fmt.Println(err)
    return
  }
  fmt.Println(string(body))
}
###
```

**Delete the event**
```go
package main

import (
  "fmt"
  "net/http"
  "io/ioutil"
)

func main() {

  url := "http://localhost:8080/api/v1/events?limit=1"
  method := "DELETE"

  client := &http.Client {
  }
  req, err := http.NewRequest(method, url, nil)

  if err != nil {
    fmt.Println(err)
    return
  }
  res, err := client.Do(req)
  if err != nil {
    fmt.Println(err)
    return
  }
  defer res.Body.Close()

  body, err := ioutil.ReadAll(res.Body)
  if err != nil {
    fmt.Println(err)
    return
  }
  fmt.Println(string(body))
}package main

import (
  "fmt"
  "net/http"
  "io/ioutil"
)

func main() {

  url := "http://localhost:8080/api/v1/events?limit=1"
  method := "DELETE"

  client := &http.Client {
  }
  req, err := http.NewRequest(method, url, nil)

  if err != nil {
    fmt.Println(err)
    return
  }
  res, err := client.Do(req)
  if err != nil {
    fmt.Println(err)
    return
  }
  defer res.Body.Close()

  body, err := ioutil.ReadAll(res.Body)
  if err != nil {
    fmt.Println(err)
    return
  }
  fmt.Println(string(body))
}
###
```
## How to run Crypter-Clients
```bash
$ go env
...
GOARCH="amd64"
GOHOSTARCH="amd64"
GOHOSTOS="linux"
GOOS="linux"
...
```
#### Compile for Windows
Here’s the command you need to run to compile your Go project for a 64-bit Windows machine:

In this scenario, GOOS is windows, and GOARCH is amd64 indicating a 64-bit architecture. If you need to support a 32-bit architecture, all you need to do is change GOARCH to 386.
```bash
$ GOOS=windows GOARCH=386 go build -o bin/app-386.exe app.go
```
#### Compile for macOS
The GOARCH values for Windows are also valid for macOS, but in this case the required GOOS value is darwin:
**64-bit**
```bash
$ GOOS=darwin GOARCH=amd64 go build -o bin/app-amd64-darwin app.go
```
**32-bit**
```bash
$ GOOS=darwin GOARCH=386 go build -o bin/app-386-darwin app.go
```
#### Compile for Linux
To build your Go program for Linux, use linux as the value of GOOS and the appropriate GOARCH value for your target CPU architecture:
**64-bit**
```bash
$ GOOS=linux GOARCH=amd64 go build -o bin/app-amd64-linux app.go
```
**32-bit**
```
$ GOOS=linux GOARCH=386 go build -o bin/app-386-linux app.go
```


# Dependencies

github.com/iamacarpet/go-win64api
github.com/spf13/viper
github.com/gorilla/mux
github.com/ten16thomasg/crypter-api/handlers
github.com/ten16thomasg/crypter-api/store



# Docs


