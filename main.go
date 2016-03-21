package main

import (
	"bufio"
	"fmt"
	"net"
	"regexp"
	"strings"
)

//GLOBALS
var (
	serverConn net.Conn
	serverBuffer *bufio.Reader

	//TRIGGER SETUP
	mainTriggers = []trigger{
		trigger{`^Please enter password:$`, nil, password_trigger},
		trigger{`^Logon successful.$`, nil, login_trigger},
		trigger{admin_command_regex, nil, admin_command_trigger},
		trigger{player_command_regex, nil, player_command_trigger},
		trigger{tick_regex, nil, tick_trigger},
		trigger{player_regex, nil, player_trigger},
		trigger{keystonetrigger_regex, nil, keystone_trigger},
	}

	keystoneTriggers = []trigger{
		//trigger{keystonetrigger_regex, nil, keystone_trigger},
		//trigger{keystoneendtrigger_regex, nil, keystoneend_trigger},
	}

	triggers []trigger

	keystoneOwner string
)

var (
	//NOTE: Storing a reference from name, steamid, and id in same map
	playerMap map[string]*Player = make(map[string]*Player)
	milliseconds int64
)

var (
	spawnPoints []Point
	mainBaseHorde bool
	zombieNum int64
	zombies = []int{ 1,  2,  3,  4,  5,  6,
			 7,  8,  9, 10, 11, 12,
			13, 14, 15, 16, 17, 18,
			19, 20,
			27,
			35, 36,
		}
)



//MAIN
func main() {
	triggers = mainTriggers
	fmt.Println("BEGIN")
	fmt.Println("keystonetrigger_regex:")
	fmt.Printf("\n\n\n%s\n\n\n", keystonetrigger_regex)

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

	//TEMP
	go updatePlayerInfo_thread()
	go mainBaseHorde_thread()

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
					matchName := matchNames[j]
					if matchName == "" {
						matchName = "$" //TODO TEMP Change to something else for the full trigger text?
					}
					matchMap[matchName] = match
				}
				//Run the trigger
				triggers[i].callback(matchMap)
			}
		}


		//fmt.Printf("SERVER: '%s'\n", line)
	}


	fmt.Println("END")
}

