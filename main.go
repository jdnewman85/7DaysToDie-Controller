package main

import (
	"bufio"
	"fmt"
	"math"
	"math/rand"
	"net"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type stateFunc func()

type triggerFunc func(map[string]string)

type trigger struct {
	reString string
	re *regexp.Regexp
	callback triggerFunc
}

type Player struct {
	id string
	name string
	steamId string

	lastUpdate int64
	online bool
	remote bool
	ip string
	ping uint64

	x, y, z float64
	u, v, w float64

	//TEMP OPT Stop using 64 bits, and do the conversions where needed
	health uint64
	deaths uint64
	zKills uint64
	pKills uint64
	score uint64
	level uint64
}

type Point struct {
	x float64
	y float64
	z float64
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

var (
	//NOTE: Storing a reference from name, steamid, and id in same map
	playerMap map[string]*Player = make(map[string]*Player)
	milliseconds int64
)

var (
	spawnPoints []Point
	mainBaseHorde bool
	zombieNum int64
	zombies = []int{7, 11, 12,
			1, 3, 4, 5, 13, 14, 15, 16, 17,
			19, 8, 9, 10,
			20, 21, 18, 18, 2, 2,
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
					matchMap[matchNames[j]] = match
				}
				//Run the trigger
				triggers[i].callback(matchMap)
			}
		}


		//fmt.Printf("SERVER: '%s'\n", line)
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
	case "pm storeSpawnPoint", "pm ssp":
		if reMatchMap["steamid"] != "" {
			tempPlayer := playerMap[reMatchMap["steamid"]]
			fmt.Fprintf(serverConn, "pm %s \"Point Added\"\n", tempPlayer.name)
			spawnPoints = append(spawnPoints, Point{tempPlayer.x, tempPlayer.y, tempPlayer.z})
			fmt.Printf("SPAWN POINTS: \n%v\n\n", spawnPoints)
		}
	case "pm resetSpawnPoints", "pm rsp":
		if reMatchMap["steamid"] != "" {
			tempPlayer := playerMap[reMatchMap["steamid"]]
			fmt.Fprintf(serverConn, "pm %s \"Spawn Points Reset\"\n", tempPlayer.name)
			spawnPoints = nil
			fmt.Printf("SPAWN POINTS: \n%v\n\n", spawnPoints)
		}
	case "pm spawnTest", "pm st":
		fmt.Printf("SPAWN TEST\n")
		if reMatchMap["steamid"] != "" {
			tempPlayer := playerMap[reMatchMap["steamid"]]
			fmt.Fprintf(serverConn, "pm %s \"Spawning zombies!\"\n", tempPlayer.name)

			//Store current position
			tempX := tempPlayer.x
			tempY := tempPlayer.y
			tempZ := tempPlayer.z

			for _, tempPoint := range spawnPoints {
				fmt.Fprintf(serverConn, "tele %s %d %d %d\n", tempPlayer.id, int(tempPoint.x), int(tempPoint.y)+5, int(tempPoint.z))
				//time.Sleep(time.Second/2)
				fmt.Fprintf(serverConn, "se %s 17\n", tempPlayer.id)
			}

			//Restore position
			fmt.Fprintf(serverConn, "tele %s %d %d %d\n", tempPlayer.id, int(tempX), int(tempY), int(tempZ))
		}
	case "pm mainBaseHorde", "pm mbh":
		fmt.Printf("Main Base Horde On\n")
		mainBaseHorde = true
	case "pm stopMainBaseHorde", "pm smbh":
		fmt.Printf("Main Base Horde Off\n")
		mainBaseHorde = false
	}
}

func tick_trigger(reMatchMap map[string]string) {
	fmt.Printf("TICK!\n")
	milliseconds, _ = strconv.ParseInt(reMatchMap["milliseconds"], 10, 64)
	zombieNum, _ = strconv.ParseInt(reMatchMap["zombienum"], 10, 64)
}

func player_trigger(reMatchMap map[string]string) {
	//fmt.Printf("Player Here!\n\n%v\n\n", reMatchMap)
	//Update player information
	playerId := reMatchMap["playerid"]
	tempPlayer := playerMap[playerId]
	if tempPlayer == nil {
		tempPlayer = new(Player)
		playerMap[playerId] = tempPlayer
		playerMap[reMatchMap["playername"]] = tempPlayer
		playerMap[reMatchMap["playersteamid"]] = tempPlayer

	}

	//Update info
	tempPlayer.id = reMatchMap["playerid"]
	tempPlayer.name = reMatchMap["playername"]
	tempPlayer.steamId = reMatchMap["playersteamid"]
	tempPlayer.ip = reMatchMap["playerip"]
	tempPlayer.ping, _ = strconv.ParseUint(reMatchMap["playerping"], 10, 32)
	tempPlayer.x, _ = strconv.ParseFloat(reMatchMap["playerx"], 64)
	tempPlayer.y, _ = strconv.ParseFloat(reMatchMap["playery"], 64)
	tempPlayer.z, _ = strconv.ParseFloat(reMatchMap["playerz"], 64)
	tempPlayer.u, _ = strconv.ParseFloat(reMatchMap["playeru"], 64)
	tempPlayer.v, _ = strconv.ParseFloat(reMatchMap["playerv"], 64)
	tempPlayer.w, _ = strconv.ParseFloat(reMatchMap["playerw"], 64)
	tempPlayer.remote, _ = strconv.ParseBool(reMatchMap["playerremote"])
	tempPlayer.health, _ = strconv.ParseUint(reMatchMap["playerhealth"], 10, 32)
	tempPlayer.deaths, _ = strconv.ParseUint(reMatchMap["playerdeaths"], 10, 32)
	tempPlayer.zKills, _ = strconv.ParseUint(reMatchMap["playerzkills"], 10, 32)
	tempPlayer.pKills, _ = strconv.ParseUint(reMatchMap["playerpkills"], 10, 32)
	tempPlayer.score, _ = strconv.ParseUint(reMatchMap["playerscore"], 10, 32)
	tempPlayer.level, _ = strconv.ParseUint(reMatchMap["playerlevel"], 10, 32)
	tempPlayer.lastUpdate = milliseconds
	tempPlayer.online = true
}

func updatePlayerInfo_thread() {
	//Periodically sends listplayers command, to update player information
	for {
		fmt.Fprintf(serverConn, "%s\n", "lp")
		time.Sleep(time.Second*1)

		//TEMP DEBUG
		/*
		fmt.Printf("\n\n\n%v\n", playerMap)
		for _, tempPlayer := range playerMap {
			fmt.Printf("%v\n", tempPlayer)
		}
		*/
	}
}

func mainBaseHorde_thread() {
	for {
		if mainBaseHorde {
			if zombieNum < 80 {
				//SPAWNZ
				//Record old location
				tempX := int(playerMap["Scienthsine"].x)
				tempY := int(playerMap["Scienthsine"].y)
				tempZ := int(playerMap["Scienthsine"].z)

				centerX := 1766
				centerY := 70
				centerZ := 2216

				for i := 0; i<20; i++ {
					randAngle := rand.Float64()*math.Pi
					randDist := 50.0+rand.Float64()*50.0
					randX := centerX+int(math.Cos(randAngle)*randDist)
					randZ := centerZ+int(math.Sin(randAngle)*randDist)
					randZombie := rand.Intn(len(zombies)-1)
					fmt.Fprintf(serverConn, "tele %s %d %d %d\n", "Scienthsine", randX, centerY, randZ)
					fmt.Fprintf(serverConn, "se %s %d\n", "Scienthsine", zombies[randZombie])
				}

				//Restore
				fmt.Fprintf(serverConn, "tele %s %d %d %d\n", "Scienthsine", tempX, tempY, tempZ)
			}
			time.Sleep(time.Minute/4)
		}
		time.Sleep(time.Second)
	}
}
