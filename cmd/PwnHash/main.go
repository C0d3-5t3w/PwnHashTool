package main

import (
	"fmt"
	"os"
	"os/user"

	"github.com/C0d3-5t3w/PwnHashTool/internal/gui"
)

func checkRoot() bool {
	currentUser, err := user.Current()
	if err != nil {
		fmt.Println("Error getting current user:", err)
		os.Exit(1)
	}
	return currentUser.Uid == "0"
}

func main() {
	if !checkRoot() {
		fmt.Println("This program must be run as root")
		os.Exit(1)
	}
	gui.Launch(nil)
}
