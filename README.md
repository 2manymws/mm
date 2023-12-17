# mm [![Go Reference](https://pkg.go.dev/badge/github.com/2manymws/mm.svg)](https://pkg.go.dev/github.com/2manymws/mm) [![build](https://github.com/2manymws/mm/actions/workflows/ci.yml/badge.svg)](https://github.com/2manymws/mm/actions/workflows/ci.yml) ![Coverage](https://raw.githubusercontent.com/k1LoW/octocovs/main/badges/2manymws/mm/coverage.svg) ![Code to Test Ratio](https://raw.githubusercontent.com/k1LoW/octocovs/main/badges/2manymws/mm/ratio.svg) ![Test Execution Time](https://raw.githubusercontent.com/k1LoW/octocovs/main/badges/2manymws/mm/time.svg)

mm is a **m**iddleware-**m**iddleware for multiple rules.

It provides middleware that changes the middlewares used based on the request.

## Usage

Prepare an instance that implements [`mm.Builder`](https://pkg.go.dev/github.com/2manymws/mm#Builder) interface.

Then, generate the middleware ( `func(next http.Handler) http.Handler` ) with [`mm.New`](https://pkg.go.dev/github.com/2manymws/mm#New)

```go
package main

import (
    "log"
    "net/http"

    "github.com/2manymws/mm"
)

func main() {
    r := http.NewServeMux()
    r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        w.Write([]byte("Hello World"))
    })

    var b mm.Builder = newMyBuilder()
    m := mm.New(b)

    log.Fatal(http.ListenAndServe(":8080", m(r)))
}
```
