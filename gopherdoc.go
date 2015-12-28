package main

import (
	"bufio"
	"bytes"
	"log"
	"net"
	"os/exec"
	"strings"
)

const serverString string = "gopherdoc"

func main() {
	server, err := net.Listen("tcp", "localhost:7000")
	if err != nil {
		log.Panicln("Couldn't start listening: " + err.Error())
	}
	conns := clientConnection(server)
	log.Println("Now listening on port 7000...")
	for {
		go handleRequest(<-conns)
	}
}

func clientConnection(listener net.Listener) chan net.Conn {
	ch := make(chan net.Conn)
	go func() {
		for {
			client, err := listener.Accept()
			if err != nil {
				log.Println("Couldn't accept connection: " + err.Error())
				continue
			}
			ch <- client
		}
	}()
	return ch
}

func handleRequest(client net.Conn) {
	b := bufio.NewReader(client)
	for {
		line, err := b.ReadBytes('\n')
		if err != nil {
			break
		}
		time := time.Now().Format(time.RFC3339)
		path := strings.TrimSpace(string(line[:]))
		// TODO: Fix the logging here
		log.Println(path)

		// TODO: router goes here
		// FIXME: channel, not return
		res := goDocRouter(path)
		client.Write([]byte(res))
		client.Write([]byte("\n\n")) // Explicitly write this out
		client.Close()
	}
}

func goDocRouter(path string) string {
	// We assume the path will look either look like /package or /package/thing
	// e.g:
	// * /buf
	// * /buf/ScanLines
	// Anything else is an error, that means appended slashes! What do you think
	// this is, HTTP?
	// FIXME: So much room for failure here, RFC too vague like every other RFC
	// with a number less than 4000
	// FIXME: lynx chomps the first 2 bytes of the request, other agents MAY
	// NOT, so we may have to strip "/1" from requests
	if path == "/" || path == "" {
		return ("iSorry, index listing isn't supported yet :[")
	}
	lookup := strings.Split(path, "/")
	if len(lookup) == 2 {
		res := formatGoDoc(string(lookup[1]))
		return res.String()
	}
	if len(lookup) == 3 {
		return ("Request: " + lookup[1] + " - " + lookup[2])
	}
	// MAGIC GOES HERE
	return ("iInvalid URL!")
}

func formatGoDoc(path string) bytes.Buffer {
	var res bytes.Buffer

	cmd := exec.Command("godoc", path)
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		log.Fatal(err)
	}
	if err := cmd.Start(); err != nil {
		log.Fatal(err)
	}
	scanner := bufio.NewScanner(stdout)
	for scanner.Scan() {
		res.WriteString("i" + scanner.Text() + "\n")
	}
	return res
}
