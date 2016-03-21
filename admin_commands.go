package main

import (
	"fmt"
	//"time"
)

//ADMIN COMMANDS
func admin_command_trigger(reMatchMap map[string]string) {
	//OPT should go ahead and get player object above case, or provide a simple access function for it
	fmt.Printf("Recieved an admin command: %s\n", reMatchMap["command"])

	fmt.Printf("reMatchMap:'\n%v\n'\n", reMatchMap)

	switch reMatchMap["command"] {
	case "whoami":
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
