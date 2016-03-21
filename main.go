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



type Point struct {
	x float64
	y float64
	z float64
}

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

	x, y, z float64 //TODO Point types?
	u, v, w float64

	//TEMP OPT Stop using 64 bits, and do the conversions where needed
	health uint64
	deaths uint64
	zKills uint64
	pKills uint64
	score uint64
	level uint64

	//Added stats
	blinkLocations []*Point
}




const (
	//SIMPLE REGULAR EXPRESSIONS
	//datetime_regex
	year_regex = `(?P<year>\d\d\d\d)`
	month_regex = `(?P<month>\d\d)`
	day_regex = `(?P<day>\d\d)`

	hour_regex = `(?P<hour>\d\d)`
	minute_regex = `(?P<minute>\d\d)`
	second_regex = `(?P<second>\d\d)`

	millisecond_regex = `(?P<millisecond>\d+.\d*)`

	//*_command_regex
	ip_regex = `(?:(?P<ip>\d+[.]\d+[.]\d+[.]\d+)(?::(?P<port>\d+))?)`
	steamid_regex = `(?P<steamid>\d+)`

	by_regex = `(?P<by>\w+)`
	//hotword_regex = `(?P<hotword>/w+)` //Gets any hotword
	admin_hotword_regex = `(?P<hotword>re)`
	player_hotword_regex = `(?P<hotword>pm)`

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

	keystonesnum_regex = `(?P<keystonesnum>\d+)`
	keystonesprotected_regex = `(?P<keystonesprotected>(False)|(True))`
	keystoneshardness_regex = `(?P<keystoneshardness>\d+)`
)

var (
	//COMPLEX REGULAR EXPRESSIONS
	date_regex = fmt.Sprintf(`(?P<date>%s-%s-%s)`, year_regex, month_regex, day_regex)
	time_regex = fmt.Sprintf(`(?P<time>%s:%s:%s)`, hour_regex, minute_regex, second_regex)
	datetime_regex = fmt.Sprintf(`%sT%s %s`, date_regex, time_regex, millisecond_regex)

	client_regex = fmt.Sprintf(`(?:client %s)`, steamid_regex)
	ipclient_regex = fmt.Sprintf(`(?:%s|%s)`, ip_regex, client_regex)

	//player_regex
	playerposition_regex = fmt.Sprintf(`(?P<playerposition>\(%s, %s, %s\))`,
		playerx_regex, playery_regex, playerz_regex)
	playerrotation_regex = fmt.Sprintf(`(?P<playerrotation>\(%s, %s, %s\))`,
		playeru_regex, playerv_regex, playerw_regex)

	//TODO Make command a simple regex above?
	//admin_command_regex
	admin_command_regex = fmt.Sprintf(`^%s INF Executing command '%s (?P<command>.+)' (?:by %s )?from %s$`,
		datetime_regex, admin_hotword_regex, by_regex, ipclient_regex)
	//player_command_regex
	player_command_regex = fmt.Sprintf(`^%s INF Denying command '%s (?P<command>.+)' from client %s$`,
		datetime_regex, player_hotword_regex, playername_regex)

	//tick_regex
	tick_regex = fmt.Sprintf(`^%s INF Time: %s FPS: %s Heap: %s Max: %s Chunks: %s CGO: %s Ply: %s Zom: %s Ent: %s %s Items: %s$`,
		datetime_regex, gametime_regex, fps_regex, heap_regex, maxheap_regex, chunks_regex, cgo_regex, playernum_regex, zombienum_regex, entitynum_regex, entitynumpar_regex, itemnum_regex)

	//Each player in a player list (lp) command
	player_regex = fmt.Sprintf(`^%s. id=%s, %s, pos=%s, rot=%s, remote=%s, health=%s, deaths=%s, zombies=%s, players=%s, score=%s, level=%s, steamid=%s, ip=%s, ping=%s$`,
		listnum_regex, playerid_regex, playername_regex, playerposition_regex, playerrotation_regex, playerremote_regex, playerhealth_regex, playerdeaths_regex, playerzkills_regex, playerpkills_regex, playerscore_regex, playerlevel_regex, playersteamid_regex, playerip_regex, playerping_regex)

	//Trigger to prepare for list of land protection blocks
	keystonetrigger_regex = fmt.Sprintf(`^Player "%s \(%s\)" owns %s keystones \(protected: %s, current hardness multiplier: %s\)$`,
		playername_regex, playersteamid_regex, keystonesnum_regex, keystonesprotected_regex, keystoneshardness_regex)
	keystoneendtrigger_regex = fmt.Sprintf(`^Total of %s keystones in the game$`,
		keystonesnum_regex)
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



//TRIGGERS CALLBACKS
func password_trigger(reMatchMap map[string]string) {
	fmt.Fprintf(serverConn, "%s\n", "7Blackrocks1")
	fmt.Println("Waiting for logon confirm.")
}

func login_trigger(reMatchMap map[string]string) {
	fmt.Println("Waiting for commands.")
}

//ADMIN COMMANDS
func admin_command_trigger(reMatchMap map[string]string) {
	//OPT should go ahead and get player object above case, or provide a simple access function for it
	fmt.Printf("Recieved an admin command: %s\n", reMatchMap["command"])
	switch reMatchMap["command"] {
	case "whoami":
		fmt.Printf("reMatchMap:'\n%v\n'\n", reMatchMap)
		if reMatchMap["steamid"] != "" {
			player := playerMap[reMatchMap["steamid"]]
			fmt.Fprintf(serverConn, "pm %s \"You are %s, an admin.\"\n", player.name, player.name)
		}
	case "lp":
		fmt.Fprintf(serverConn, "%s\n", "lp")
	case "storeSpawnPoint", "pm ssp":
		if reMatchMap["steamid"] != "" {
			player := playerMap[reMatchMap["steamid"]]
			fmt.Fprintf(serverConn, "pm %s \"Point Added\"\n", player.name)
			spawnPoints = append(spawnPoints, Point{player.x, player.y, player.z})
			fmt.Printf("SPAWN POINTS: \n%v\n\n", spawnPoints)
		}
	case "resetSpawnPoints", "pm rsp":
		if reMatchMap["steamid"] != "" {
			player := playerMap[reMatchMap["steamid"]]
			fmt.Fprintf(serverConn, "pm %s \"Spawn Points Reset\"\n", player.name)
			spawnPoints = nil
			fmt.Printf("SPAWN POINTS: \n%v\n\n", spawnPoints)
		}
	case "spawnTest", "pm st":
		fmt.Printf("SPAWN TEST\n")
		if reMatchMap["steamid"] != "" {
			player := playerMap[reMatchMap["steamid"]]
			fmt.Fprintf(serverConn, "pm %s \"Spawning zombies!\"\n", player.name)

			//Store current position
			tempX := player.x
			tempY := player.y
			tempZ := player.z

			for _, tempPoint := range spawnPoints {
				fmt.Fprintf(serverConn, "tele %s %d %d %d\n", player.id, int(tempPoint.x), int(tempPoint.y)+5, int(tempPoint.z))
				//time.Sleep(time.Second/2)
				fmt.Fprintf(serverConn, "se %s 17\n", player.id)
			}

			//Restore position
			fmt.Fprintf(serverConn, "tele %s %d %d %d\n", player.id, int(tempX), int(tempY), int(tempZ))
		}
	case "mainBaseHorde", "pm mbh":
		fmt.Printf("Main Base Horde On\n")
		mainBaseHorde = true
	case "stopMainBaseHorde", "pm smbh":
		fmt.Printf("Main Base Horde Off\n")
		mainBaseHorde = false
	case "blink":
		//TEMP TODO make different system
		fmt.Printf("Blink Command!\n")

		if reMatchMap["steamid"] != "" {
			player := playerMap[reMatchMap["steamid"]]

			if player.blinkLocations == nil {
				//No points
				//TODO TEMP Diff
				player.blinkLocations = append(player.blinkLocations, &Point{player.x, player.y, player.z})
			} else {
				blinkLocation := player.blinkLocations[0]
				//Points
				//TODO TEMP for now, just go to first
				//TODO z position?
				fmt.Fprintf(serverConn, "tele %s %d %d %d\n", player.id, int(blinkLocation.x), int(blinkLocation.y), int(blinkLocation.z)+2) //TODO TEMP Magic number
			}
		}
	}
}

//PLAYER COMMANDS
func player_command_trigger(reMatchMap map[string]string) {
	fmt.Printf("Recieved a player command: %s\n", reMatchMap["command"])
	switch reMatchMap["command"] {
	case "whoami":
		if reMatchMap["steamid"] != "" {
			playerName := reMatchMap["playername"] //OPT Need whitespace
			player := playerMap[playerName]

			fmt.Fprintf(serverConn, "pm %s \"You are %s, a player.\"\n", player.name, player.name)
		}
	case "blink":
					//TEMP TODO make different system
					fmt.Printf("Blink Command!\n")
					playerName := reMatchMap["playername"] //OPT Need whitespace
					player := playerMap[playerName]
					if player.blinkLocations == nil {
						//No points
						//TODO TEMP Diff
						player.blinkLocations = append(player.blinkLocations, &Point{player.x, player.y, player.z})
					} else {
						blinkLocation := player.blinkLocations[0]
						//Points
						//TODO TEMP for now, just go to first
						//TODO z position?
						fmt.Fprintf(serverConn, "tele %s %d %d %d\n", player.id, int(blinkLocation.x), int(blinkLocation.y), int(blinkLocation.z)+2) //TODO TEMP Magic number
					}
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
	player := playerMap[playerId]
	if player == nil {
		player = new(Player)
		playerMap[playerId] = player
		playerMap[reMatchMap["playername"]] = player
		playerMap[reMatchMap["playersteamid"]] = player

	}

	//Update info
	player.id = reMatchMap["playerid"]
	player.name = reMatchMap["playername"]
	player.steamId = reMatchMap["playersteamid"]
	player.ip = reMatchMap["playerip"]
	player.ping, _ = strconv.ParseUint(reMatchMap["playerping"], 10, 32)
	player.x, _ = strconv.ParseFloat(reMatchMap["playerx"], 64)
	player.y, _ = strconv.ParseFloat(reMatchMap["playery"], 64)
	player.z, _ = strconv.ParseFloat(reMatchMap["playerz"], 64)
	player.u, _ = strconv.ParseFloat(reMatchMap["playeru"], 64)
	player.v, _ = strconv.ParseFloat(reMatchMap["playerv"], 64)
	player.w, _ = strconv.ParseFloat(reMatchMap["playerw"], 64)
	player.remote, _ = strconv.ParseBool(reMatchMap["playerremote"])
	player.health, _ = strconv.ParseUint(reMatchMap["playerhealth"], 10, 32)
	player.deaths, _ = strconv.ParseUint(reMatchMap["playerdeaths"], 10, 32)
	player.zKills, _ = strconv.ParseUint(reMatchMap["playerzkills"], 10, 32)
	player.pKills, _ = strconv.ParseUint(reMatchMap["playerpkills"], 10, 32)
	player.score, _ = strconv.ParseUint(reMatchMap["playerscore"], 10, 32)
	player.level, _ = strconv.ParseUint(reMatchMap["playerlevel"], 10, 32)
	player.lastUpdate = milliseconds
	player.online = true
}

func keystone_trigger(reMatchMap map[string] string) {
	fmt.Printf("Listing keystones for: %s(%s) numbering %s protected? %s hardness %s\n",
	reMatchMap["playername"],
	reMatchMap["playersteamid"],
	reMatchMap["keystonesnum"],
	reMatchMap["keystonesprotected"],
	reMatchMap["keystoneshardness"])

	keystoneOwner = reMatchMap["playername"]
	triggers = keystoneTriggers
}

func keystoneend_trigger(reMatchMap map[string] string) {
	fmt.Printf("Total keystones: %s\n", reMatchMap["keystonesnum"])
	triggers = mainTriggers
}



//THREAD FUNCTIONS
func updatePlayerInfo_thread() {
	//Periodically sends listplayers command, to update player information
	for {
		fmt.Fprintf(serverConn, "%s\n", "lp")
		fmt.Fprintf(serverConn, "%s\n", "listlandprotection") //OPT TODO Move this
		time.Sleep(time.Second*1)

		//TEMP DEBUG
		/*
		fmt.Printf("\n\n\n%v\n", playerMap)
		for _, player := range playerMap {
			fmt.Printf("%v\n", player)
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
