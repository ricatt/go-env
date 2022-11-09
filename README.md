# go-env

## Introduction
This is just a small hobby-package, inspired after many discussions of "how do we want to load configuration into the project".

## Instructions


## Example
```env
BASE_URL=http://example.com
```
```go
package main

import "github.com/ricatt/go-env/env"

type Config struct {
    BaseURL string `env:"BASE_URL"`
}

func main() {
    var config Config
    env.Load(&config, env.Config{})
    // todo: profit
}
```