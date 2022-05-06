package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"regexp"
	"strings"
)

const (
	HOST      = "localhost"
	CONN_PORT = "3333"
	CONN_TYPE = "tcp"
)

const HTTP_OK = `
HTTP/1.1 200 OK
Date: Fri, 06 Jul 2022 09:02:53 GMT
Server: Server
Last-Modified: Wed, 22 Jul 2009 19:15:56 GMT
Content-Length: 1024
Content-Type: text/html
Connection: Closed
`
const METHOD_NOT_ALLOWED = `
HTTP/1.1 405 NOT ALLOWED
Content-Type: text/html
Content-Length: 1024
Date: Thu, 5 May 2022 21:21:00 GMT 
Server: MyServer
`
const NOT_FOUND = `
HTTP/1.1 404 Not Found
Date: Thu, 5 May 2022 21:21:00 GMT 
Server: MyServer
Content-Type: text/html
Content-Length: 1024
`

const SERVER_ERROR = `
HTTP/1.1 500 Server Error
Date: Thu, 5 May 2022 21:21:00 GMT 
Server: MyServer
Content-Type: text/html
Content-Length: 1024
`

func main() {
	listen, err := net.Listen(CONN_TYPE, HOST+":"+CONN_PORT)

	if err != nil {
		fmt.Println("Error listening", err.Error())
		os.Exit(1)
	}

	defer listen.Close()

	fmt.Println("Listening on " + HOST + ":" + CONN_PORT)

	for {
		conn, err := listen.Accept()

		if err != nil {
			fmt.Println("Error accepting conn ", err.Error())
			os.Exit(1)
		}

		go handleRequest(conn)
	}
}

func handleRequest(conn net.Conn) {
	defer conn.Close()

	buf := make([]byte, 1024)

	_, err := conn.Read(buf)

	if err != nil {
		fmt.Println(err.Error())
	}
	handle, err := treatRequest(string(buf))
	if err == nil {
		conn.Write(append([]byte(HTTP_OK), handle...))
		return
	}
	
	if os.IsNotExist(err) {
		file, _ := findFile("not_found.html")
		conn.Write(append([]byte(NOT_FOUND), file...))
		return
	}

	fmt.Println(err.Error())
	file, _ := findFile("server_error.html")
	conn.Write(append([]byte(SERVER_ERROR), file...))
}

func treatRequest(content string) ([]byte, error) {
	scanner := bufio.NewScanner(strings.NewReader(content))
	for scanner.Scan() {
		text := scanner.Text()
		var re = regexp.MustCompile(`(?m)GET`) // NO SUPPORT FOT POST|PUT|PATCH|DELETE

		match := re.FindString(text)
		if match != "" {
			path := strings.Split(text, " ")[1]
			if path == "/" {
				return findFile("index.html")
			}
			return findFile(path[1:])
		}
		break //scanning only the first line
	}

	// no method found.
	file, err := findFile("method_not_allowed.html")
	if err != nil {
		return nil, err
	}
	return append([]byte(METHOD_NOT_ALLOWED), file...), nil
}

func findFile(path string) ([]byte, error) {
	f, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	return f, nil
}
