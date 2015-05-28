package main

import (
	"bufio"
	"fmt"
	"net"
	"regexp"
	"strings"
)

type stateFunc func()

type triggerFunc func(map[string]string)

type trigger struct {
	reString string
	re *regexp.Regexp
	callback triggerFunc
}

var (
	serverConn net.Conn
	serverBuffer *bufio.Reader

	triggers = []trigger{
		trigger{"^Please enter password:$", nil, password_trigger},
		trigger{"^Logon successful.$", nil, login_trigger},
		trigger{".*blargy.*", nil, blargy_trigger},
		trigger{".*SendMeHome.*", nil, home_trigger},
	}
)

func main() {
	fmt.Println("BEGIN")

	//Compile trigger regexps
	for i := range triggers {
		triggers[i].re = regexp.MustCompile(triggers[i].reString) //OPT Method?
	}

	//Server Connect
	var err error
	serverConn, err = net.Dial("tcp", "7days.prototekokc.com:8081")
	if err != nil {
		panic(err)
	}
	defer serverConn.Close()

	serverBuffer := bufio.NewReader(serverConn)

	//Read those lines
	fmt.Println("Waiting for password")
	for line, err := serverBuffer.ReadString('\n'); err == nil; line, err = serverBuffer.ReadString('\n') {

		//Remove extra whitespace and skip empty lines
		line := strings.TrimSpace(line)
		if line == "" {
			continue
		}

		//Check for triggers
		//OPT Function, or something?
		for i := range triggers {
			matches := triggers[i].re.FindStringSubmatch(line)
			if matches != nil {
				//Assemble a matches map
				matchMap := make(map[string]string)
				matchNames := triggers[i].re.SubexpNames()
				for j, match := range matches {
					matchMap[matchNames[j]] = match
				}
				//Run the trigger
				triggers[i].callback(matchMap)
			}
		}


		fmt.Printf("SERVER: '%s'\n", line)
	}


	fmt.Println("END")
}

func password_trigger(reMatchMap map[string]string) {
	fmt.Fprintf(serverConn, "%s\n", "TrumpTikAdmin%")
	fmt.Println("Waiting for logon confirm.")
}

func login_trigger(reMatchMap map[string]string) {
	fmt.Println("Waiting for commands.")
}

func blargy_trigger(reMatchMap map[string]string) {
	fmt.Println("blargy command!")
	fmt.Fprintf(serverConn, "say \"%s\"\n", "Yay blargs!")
}

func home_trigger(reMatchMap map[string]string) {
	fmt.Println("Home command!")
	fmt.Fprintf(serverConn, "%s\n", "tele scienthsine 2956 82 1220")
}

