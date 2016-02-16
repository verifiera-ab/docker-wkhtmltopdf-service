package main

import (
	"net/http"
	"net/url"
	"fmt"
	"os/exec"
	"encoding/json"
	"strings"
)

func main() {
	const bindAddress = ":3000"
	http.HandleFunc("/", requestHandler)
	fmt.Println("Http server listening on", bindAddress)
	http.ListenAndServe(bindAddress, nil)
}

type documentRequest struct {
	Url string
	Output string
	Options map[string]interface{}
	Cookies map[string]string
}

func logOutput(request *http.Request, message string) {
	ip := strings.Split(request.RemoteAddr, ":")[0]
	fmt.Println(ip, request.Method, request.URL, message)
}

func requestHandler(response http.ResponseWriter, request *http.Request) {
	if request.URL.Path != "/" {
		response.WriteHeader(http.StatusNotFound)
		logOutput(request, "404 not found")
		return
	}
	if request.Method != "POST" {
		response.Header().Set("Allow", "POST")
		response.WriteHeader(http.StatusMethodNotAllowed)
		logOutput(request, "405 not allowed")
		return
	}
	decoder := json.NewDecoder(request.Body)
	var req documentRequest
	if err := decoder.Decode(&req); err != nil {
		response.WriteHeader(http.StatusBadRequest)
		logOutput(request, "400 bad request (invalid JSON)")
		return
	}
	segments := make([]string, 0)
	for key, element := range req.Options {
		if element == true {
			// if it was parsed from the JSON as an actual boolean, 
			// convert to command-line single argument	(--foo)
			segments = append(segments, fmt.Sprintf("--%v", key))
		} else if element != false {
			// Otherwise, use command-line argument with value (--foo bar)
			segments = append(segments, fmt.Sprintf("--%v", key), fmt.Sprintf("%v", element))
		}
	}
	for key, value := range req.Cookies {
		segments = append(segments, "--cookie", key, url.QueryEscape(value))
	}
	var programFile string
	var contentType string
	switch req.Output {
		case "jpg":
			programFile = "/usr/local/bin/wkhtmltoimage"
			contentType = "image/jpeg"
			segments = append(segments, "--format", "jpg")
		case "png":
			programFile = "/usr/local/bin/wkhtmltoimage"
			contentType = "image/png"
			segments = append(segments, "--format", "png")
		default:
			// defaults to pdf
			programFile = "/usr/local/bin/wkhtmltopdf"
			contentType = "application/pdf"
	}
	segments = append(segments, req.Url, "-")
	fmt.Println("\tRunning:", programFile, strings.Join(segments, " "))
	cmd := exec.Command(programFile, segments...)
	response.Header().Set("Content-Type", contentType)
	cmd.Stdout = response
	cmd.Start()
	defer cmd.Wait()
	logOutput(request, "200 OK")
}
