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

const (
	//datetime_regex
	year_regex = `(?P<year>\d\d\d\d)`
	month_regex = `(?P<month>\d\d)`
	day_regex = `(?P<day>\d\d)`

	hour_regex = `(?P<hour>\d\d)`
	minute_regex = `(?P<minute>\d\d)`
	second_regex = `(?P<second>\d\d)`

	millisecond_regex = `(?P<millisecond>\d+.\d*)`

	//command_regex
	by_regex = `(?P<by>\w+)`
	ip_regex = `(?:(?P<ip>\d+[.]\d+[.]\d+[.]\d+)(?::(?P<port>\d+))?)`
	steamid_regex = `(?P<steamid>\d+)`

	//tick_regex
	gametime_regex = `(?P<gametime>\d+.\d\d)m`
	fps_regex = `(?P<fps>\d+.\d\d)`
	heap_regex = `(?P<heap>\d+.\d)MB`
	maxheap_regex = `(?P<maxheap>\d+.\d)MB`
	chunks_regex = `(?P<chunks>\d+)`
	cgo_regex = `(?P<cgo>\d+)`
	playernum_regex = `(?P<playernum>\d+)`
	zombienum_regex = `(?P<zombienum>\d+)`
	entitynum_regex = `(?P<entitynum>\d+)`
	entitynumpar_regex = `\((?P<entitynumpar>\d+)\)`
	itemnum_regex = `(?P<itemnum>\d+)`

	//player_regex
	listnum_regex = `(?P<listnum>\d+)`
	playerid_regex = `(?P<playerid>\d+)`
	playername_regex = `(?P<playername>\w+)`

	playerx_regex = `(?P<playerx>-?\d+.\d)`
	playery_regex = `(?P<playery>-?\d+.\d)`
	playerz_regex = `(?P<playerz>-?\d+.\d)`
	playeru_regex = `(?P<playeru>-?\d+.\d)`
	playerv_regex = `(?P<playerv>-?\d+.\d)`
	playerw_regex = `(?P<playerw>-?\d+.\d)`

	playerremote_regex = `(?P<playerremote>(?:True)|(?:False))`
	playerhealth_regex = `(?P<playerhealth>\d+)`
	playerdeaths_regex = `(?P<playerdeaths>\d+)`
	playerzkills_regex = `(?P<playerzkills>\d+)`
	playerpkills_regex = `(?P<playerpkills>\d+)`
	playerscore_regex = `(?P<playerscore>\d+)`
	playerlevel_regex = `(?P<playerlevel>\d+)`
	playersteamid_regex = `(?P<playersteamid>\d+)`
	playerip_regex = `(?P<playerip>\d+[.]\d+[.]\d+[.]\d+)`
	playerping_regex = `(?P<playerping>\d+)`

)

var (
	date_regex = fmt.Sprintf(`(?P<date>%s-%s-%s)`, year_regex, month_regex, day_regex)
	time_regex = fmt.Sprintf(`(?P<time>%s:%s:%s)`, hour_regex, minute_regex, second_regex)
	datetime_regex = fmt.Sprintf(`%sT%s %s`, date_regex, time_regex, millisecond_regex)

	client_regex = fmt.Sprintf(`(?:client %s)`, steamid_regex)
	ipclient_regex = fmt.Sprintf(`(?:%s|%s)`, ip_regex, client_regex)

	//player_regex
	playerposition_regex = fmt.Sprintf(`(?P<playerposition>\(%s, %s, %s\))`, playerx_regex, playery_regex, playerz_regex)
	playerrotation_regex = fmt.Sprintf(`(?P<playerrotation>\(%s, %s, %s\))`, playeru_regex, playerv_regex, playerw_regex)

	command_regex = fmt.Sprintf(`^%s INF Executing command '(?P<command>.+)' (?:by %s )?from %s$`, datetime_regex, by_regex, ipclient_regex)

	tick_regex = fmt.Sprintf(`^%s INF Time: %s FPS: %s Heap: %s Max: %s Chunks: %s CGO: %s Ply: %s Zom: %s Ent: %s %s Items: %s$`, datetime_regex, gametime_regex, fps_regex, heap_regex, maxheap_regex, chunks_regex, cgo_regex, playernum_regex, zombienum_regex, entitynum_regex, entitynumpar_regex, itemnum_regex)

	//Each player in a player list (lp) command
	player_regex = fmt.Sprintf(`^%s. id=%s, %s, pos=%s, rot=%s, remote=%s, health=%s, deaths=%s, zombies=%s, players=%s, score=%s, level=%s, steamid=%s, ip=%s, ping=%s$`, listnum_regex, playerid_regex, playername_regex, playerposition_regex, playerrotation_regex, playerremote_regex, playerhealth_regex, playerdeaths_regex, playerzkills_regex, playerpkills_regex, playerscore_regex, playerlevel_regex, playersteamid_regex, playerip_regex, playerping_regex)
)

var (
	serverConn net.Conn
	serverBuffer *bufio.Reader

	triggers = []trigger{
		trigger{`^Please enter password:$`, nil, password_trigger},
		trigger{`^Logon successful.$`, nil, login_trigger},
		trigger{command_regex, nil, command_trigger},
		trigger{tick_regex, nil, tick_trigger},
		trigger{player_regex, nil, player_trigger},
	}
)

func main() {
	fmt.Println("BEGIN")
	//fmt.Printf("\n\n\n%s\n\n\n", player_regex)

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

func command_trigger(reMatchMap map[string]string) {
	fmt.Printf("Recieved a command: %s\n", reMatchMap["command"])
	switch reMatchMap["command"] {
	case "pm lp":
		fmt.Fprintf(serverConn, "%s\n", "lp")
	case "pm storeWavePoint", "pm swp":
	}
}

func tick_trigger(reMatchMap map[string]string) {
	fmt.Printf("TICK!\n")
}

func player_trigger(reMatchMap map[string]string) {
	//fmt.Printf("Player Here!\n\n%v\n\n", reMatchMap)
}
