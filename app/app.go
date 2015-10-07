package main

import (
	"net/http"
        "fmt"
)

func main() {
        const bindAddress = ":3000"
        http.HandleFunc("/", requestHandler)
        fmt.Println("Http server listening on", bindAddress)
        http.ListenAndServe(bindAddress, nil)
}

func requestHandler(response http.ResponseWriter, request *http.Request) {
        if request.URL.Path != "/" {
                response.WriteHeader(http.StatusNotFound)
                fmt.Println(request.Method, request.URL, "(rejected)")
                return
        }
        fmt.Fprintf(response, "<h1>Hello, World!</h1>");
        fmt.Println(request.Method, request.URL, "(accepted)")
}
