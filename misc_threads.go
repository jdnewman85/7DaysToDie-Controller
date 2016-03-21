package main

import (
	"fmt"
	"math"
	"math/rand"
	"time"
)

//THREAD FUNCTIONS
func updatePlayerInfo_thread() {
	//Periodically sends listplayers command, to update player information
	for {
		fmt.Fprintf(serverConn, "%s\n", "lp")
		fmt.Fprintf(serverConn, "%s\n", "listlandprotection") //OPT TODO Move this
		time.Sleep(time.Second * 1)

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

				for i := 0; i < 20; i++ {
					randAngle := rand.Float64() * math.Pi
					randDist := 50.0 + rand.Float64()*50.0
					randX := centerX + int(math.Cos(randAngle)*randDist)
					randZ := centerZ + int(math.Sin(randAngle)*randDist)
					randZombie := rand.Intn(len(zombies) - 1)
					fmt.Fprintf(serverConn, "tele %s %d %d %d\n", "Scienthsine", randX, centerY, randZ)
					fmt.Fprintf(serverConn, "se %s %d\n", "Scienthsine", zombies[randZombie])
				}

				//Restore
				fmt.Fprintf(serverConn, "tele %s %d %d %d\n", "Scienthsine", tempX, tempY, tempZ)
			}
			time.Sleep(time.Minute / 4)
		}
		time.Sleep(time.Second)
	}
}
