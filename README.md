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