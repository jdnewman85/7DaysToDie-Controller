package main

import (
	"bufio"
	"fmt"
	"net"
	"strings"
)

type stateFunc func()

var (
	serverConn net.Conn
	serverBuffer *bufio.Reader


	state string = "state_password"
	stateMap = map[string]stateFunc {
		"state_password": state_password,
		"state_login": state_login,
		"state_run": state_run,
		"state_debug": state_debug,
	}

	line string
	err error
)

func main() {
	fmt.Println("BEGIN")

	//Server Connect
	serverConn, err = net.Dial("tcp", "7days.prototekokc.com:8081")
	if err != nil {
		panic(err)
	}
	defer serverConn.Close()

	serverBuffer := bufio.NewReader(serverConn)


	//Read those lines
	fmt.Println("Waiting for password")
	for line, err = serverBuffer.ReadString('\n'); err == nil; line, err = serverBuffer.ReadString('\n') {

		//Remove extra whitespace and skip empty lines
		line := strings.TrimSpace(line)
		if line == "" {
			continue
		}

		//Main State Machine
		stateMap[state]()


		fmt.Printf("SERVER: '%s'\n", line)
	}


	fmt.Println("END")
}

func state_password() {
	if strings.Contains(line, "Please enter password:") {
		fmt.Fprintf(serverConn, "%s\n", "TrumpTikAdmin%")
		fmt.Println("Waiting for logon confirm.")
		state = "state_login"
	}
}

func state_login() {
	if strings.Contains(line, "Logon successful.") {
		fmt.Println("Waiting for commands.")
		state = "state_run"
	}
}

func state_run() {
	switch {
	case strings.Contains(line, "blargy"):
		fmt.Println("blargy command!")
		fmt.Fprintf(serverConn, "say \"%s\"\n", "Yay blargs!")
	case strings.Contains(line, "SendMeHome"):
		fmt.Println("Home command!")
		fmt.Fprintf(serverConn, "%s\n", "tele scienthsine 2956 82 1220")
	}
}

func state_debug() {

}

