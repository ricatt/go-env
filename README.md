# go-env

## Introduction
This is just a small hobby-package, inspired after many discussions of "how do we want to load configuration into the project".

## Example
.env
```env
BASE_URL=http://example.com
```
main.go
```go
package main

import "github.com/ricatt/go-env/pkg/env"

type Config struct {
    BaseURL string `env:"BASE_URL"`
}

func main() {
    var config Config
    env.Load(&config, env.Config{})
    // todo: profit
}
```

There is also support for multi-level structs, for easier organisation of your data.

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

import "github.com/ricatt/go-env/pkg/env"

type Config struct {
    ServiceName string `env:"NAME_URL"`
    ExternalService struct {
        URL string `env:"EXTERNAL_SERVICE_URL"`
        Username string `env:"EXTERNAL_SERVICE_USERNAME"`
        Password string `env:"EXTERNAL_SERVICE_PASSWORD"`
    }
}

func main() {
    var config Config
    env.Load(&config, env.Config{})
    // todo: profit
}
```