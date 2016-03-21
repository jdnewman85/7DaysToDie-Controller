package main

import (
	"fmt"
)

//PLAYER COMMANDS
func player_command_trigger(reMatchMap map[string]string) {
	fmt.Printf("Recieved a player command: %s\n", reMatchMap["command"])

	fmt.Printf("reMatchMap:'\n%v\n'\n", reMatchMap)

	switch reMatchMap["command"] {
	//HOME
	case "blink":
		//TEMP TODO make different system
		fmt.Printf("Blink Command!\n")
		playerName := reMatchMap["playername"]
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

	//MISC
	case "reboot":
		playerName := reMatchMap["playername"]
		player := playerMap[playerName]

		fmt.Fprintf(serverConn, "say \"%s has started a reboot\"\n", player.name)
		go reboot_thread()
	case "whoami":
		playerName := reMatchMap["playername"]
		player := playerMap[playerName]

		fmt.Fprintf(serverConn, "pm %s \"You are %s, a player.\"\n", player.name, player.name)
	} //Switch end
}
