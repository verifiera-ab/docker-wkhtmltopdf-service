package main

import (
	"net/http"
        "fmt"
        "os/exec"
        "encoding/json"
        "io"
)

func main() {
        const bindAddress = ":3000"
        http.HandleFunc("/", requestHandler)
        fmt.Println("Http server listening on", bindAddress)
        http.ListenAndServe(bindAddress, nil)
}

type documentRequest struct {
        Content string
        Options map[string]interface{}
}

func requestHandler(response http.ResponseWriter, request *http.Request) {
        if request.URL.Path != "/" {
                response.WriteHeader(http.StatusNotFound)
                fmt.Println(request.Method, request.URL, "404 not found")
                return
        }
        if request.Method != "POST" {
                response.Header().Set("Allow", "POST")
                response.WriteHeader(http.StatusMethodNotAllowed)
                fmt.Println(request.Method, request.URL, "405 not allowed")
                return
        }
        decoder := json.NewDecoder(request.Body)
        var req documentRequest
        if err := decoder.Decode(&req); err != nil {
                response.WriteHeader(http.StatusBadRequest)
                fmt.Println(request.Method, request.URL, "400 bad request (invalid JSON)")
                return
        }
        for key, element := range req.Options {
                fmt.Println("Option", key, "=", element)
        }
        cmd := exec.Command("/usr/local/bin/wkhtmltopdf", "-", "-")
        response.Header().Set("Content-Type", "application/pdf")
        cmd.Stdout = response
        stdin, err := cmd.StdinPipe()
        if err != nil {
                response.WriteHeader(http.StatusInternalServerError)
                fmt.Println(request.Method, request.URL, "500 internal server error")
                fmt.Println("Error dump =", err)
                return
        }
        cmd.Start()
        defer cmd.Wait()
        io.WriteString(stdin, req.Content)
        stdin.Close()
        fmt.Println(request.Method, request.URL, "200 OK")
}
