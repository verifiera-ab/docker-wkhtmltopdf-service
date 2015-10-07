package main

import (
	"net/http"
        "fmt"
        "os/exec"
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
        fmt.Println(request.Method, request.URL, "(accepted)")
        cmd := exec.Command("/usr/local/bin/wkhtmltopdf", "http://www.google.com/", "-")
        response.Header().Set("Content-Type", "application/pdf")
        cmd.Stdout = response
        cmd.Run()
}
