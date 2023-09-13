# go-env

## About the project
This is just a small hobby-package, inspired after many discussions of "how do we want to load configuration into the project".

## Table of Content
1. [How to use](#how-to-use)
   1. [Supported types](#supported-types)
   2. [Attributes](#attributes)
2. [Examples](#examples)
   1. [Single value](#single-value)
   2. [Multi-level structs](#multi-level-structs)
   3. [Slices](#slices)
---
## How to use
### Supported types
 - String
 - Boolean
 - Int
 - Uint
 - Int8
 - Uint8
 - Int32
 - Uint32
 - Int64
 - Uint64
 - Float
 - Slices of the previously mentioned types.

### Attributes
`Force`: `bool` Forces a value to exist, will throw error if it comes up empty. Default: `false`

`EnvironmentFiles`: `[]string{}` A list of files where you wish to fetch your environment from.

`ErrorOnMissingFile`: `bool` Will throw error if any of the provided files are missing. Default: `false`

This attributes are functions which can be applied to the `Load`-function.
```go
package main

import "github.com/ricatt/go-env"

type Config struct {
    BaseURL string `env:"BASE_URL"`
}

func main() {
    var config Config
    env.Load(&config, env.EnvironmentFiles(".env"), env.Force(true))
}
```


## Examples

### Single value
.env
```env
BASE_URL=http://example.com
```
main.go
```go
package main

import "github.com/ricatt/go-env"

type Config struct {
    BaseURL string `env:"BASE_URL"`
}

func main() {
    var config Config
    env.Load(&config)
    // todo: profit
}
```

### Multi-level structs
There is support for multi-level structs, for easier structure and sorting of your data. The following example is based
on the idea that you connect to a third-party service. It is then easy to add this into a struct of their own. You could
take this further and make a complete struct just for the service you are using and use that one as an argument in your
function. But I digress..!

.env
```env
SERVICE_NAME=your-service
EXTERNAL_SERVICE_URL=http://external-service.com
EXTERNAL_SERVICE_USERNAME=username
EXTERNAL_SERVICE_PASSWORD=password
```
main.go
```go
package main

import "github.com/ricatt/go-env"

type ExternalService struct {
	URL string `env:"EXTERNAL_SERVICE_URL"`
	Username string `env:"EXTERNAL_SERVICE_USERNAME"`
	Password string `env:"EXTERNAL_SERVICE_PASSWORD"`
}

type Config struct {
    ServiceName string `env:"SERVICE_NAME"`
    ExternalService ExternalService
}

func main() {
    var config Config
    env.Load(&config)
	useExternalService(config.ExternalService)
}

func useExternalService(creds ExternalService) {
	// do stuff
}
```

### Slices
We have support for slices, for the times you need to specify a list in environment. Useful for when you are using the
parameter-store at AWS (for example) and just wish to upload a new list and restart your application.

.env
```env
MATCH_INT=1,2,3,4
```
main.go
```go
package main

import (
	"fmt"
	"github.com/ricatt/go-env"
	"slices"
)

type Config struct {
    MatchInt []int `env:"MATCH_INT"`
}

func main() {
    var config Config
    env.Load(&config)
    if slices.Contains(config.MatchInt, 42) {
		fmt.Println("Found the meaning of life, the universe and everything.")
    } else {
		fmt.Println("Oh no, not again...")
    }
}
```
