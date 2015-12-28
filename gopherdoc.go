package main

import (
	"bufio"
	"bytes"
	"flag"
	"log"
	"net"
	"os/exec"
	"strings"
)

const serverString string = "gopherdoc"

func main() {
	var port = flag.String("port", "7000", "Default port to bind to")
	flag.Parse()

	server, err := net.Listen("tcp", "localhost:"+*port)
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
	// We expect paths
	// FIXME: lynx chomps the first 2 bytes of the request, other agents MAY
	// NOT, so we may have to strip "/1" from requests
	if path == "/" || path == "" {
		return ("iSorry, index listing isn't supported yet :[")
	}
	split := strings.Split(path, "/")
	if len(split) == 0 || len(split) > 3 {
		return ("iInvalid URL!")
	}
	paths := split[1:]
	res := formatGoDoc(paths...)
	return res.String()
}

func formatGoDoc(paths ...string) bytes.Buffer {
	var res bytes.Buffer
	cmd := exec.Command("godoc", paths...)
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
