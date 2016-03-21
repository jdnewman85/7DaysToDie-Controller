package main

import (
	"fmt"
)

func password_trigger(reMatchMap map[string]string) {
	fmt.Fprintf(serverConn, "%s\n", "7Blackrocks1")
	fmt.Println("Waiting for logon confirm.")
}

func login_trigger(reMatchMap map[string]string) {
	fmt.Println("Waiting for commands.")
}
