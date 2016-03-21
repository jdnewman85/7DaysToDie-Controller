package main

import (
	"fmt"
	"strconv"
)

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

