package main

import (
	"encoding/base64"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/gobuffalo/packr"
)

// ServerCommand holds the options for running the HTTP server.
type ServerCommand struct {
	Host string `short:"h" long:"host" description:"Host or address to bind the HTTP server to. (Default: 0.0.0.0)"`
	Port uint16 `short:"p" long:"port" description:"The port to bind the HTTP server to. (Default: 80)"`
}

// Execute runs the HTTP server of this tool.
func (cmd *ServerCommand) Execute(args []string) error {

	if cmd.Port == 0 {
		cmd.Port = 80
	}
	if cmd.Host == "" {
		cmd.Host = "0.0.0.0"
	}

	http.HandleFunc("/", handler)

	bindStr := fmt.Sprintf("%v:%v", cmd.Host, cmd.Port)

	log.Printf("Starting web server on %v\n", bindStr)

	return http.ListenAndServe(bindStr, nil)
}

func handler(resp http.ResponseWriter, req *http.Request) {

	if req.URL.Path == "/" {

		asset(resp, "index.html")

	} else if strings.HasPrefix(req.URL.Path, "/icon/") {

		icon(resp, req)

	} else {

		asset(resp, strings.Trim(req.URL.Path, "/"))

	}
}

func icon(resp http.ResponseWriter, req *http.Request) {

	link, ok := extractURL(req.URL.Path)

	if !ok {
		resp.WriteHeader(400)
		return
	}

	log.Print(link)

	img, err := getImage(link)

	if err != nil {
		resp.WriteHeader(404)
		log.Print(err)
		return
	}

	img = Unmask(img)

	imgBytes, err := imageBytes(img)

	if err != nil {
		resp.WriteHeader(500)
		log.Print(err)
		return
	}

	resp.Write(imgBytes)
}

func asset(resp http.ResponseWriter, assetName string) {

	box := packr.NewBox("./public")

	asset, err := box.MustBytes(assetName)

	if err != nil {
		resp.WriteHeader(404)
		return
	}

	resp.Write(asset)
}

func extractURL(path string) (link string, ok bool) {

	path = strings.Replace(path, "/icon/", "", -1)

	if link, ok = base64Decode(path); !ok {
		link = strings.Replace(path, "https:/", "", -1)
		link = strings.Replace(link, "http:/", "", -1)
		link = "https://" + strings.Trim(link, "/")
		ok = true
	}

	return
}

func base64Decode(str string) (string, bool) {
	data, err := base64.StdEncoding.DecodeString(str)
	if err != nil {
		return "", false
	}
	return string(data), true

}
