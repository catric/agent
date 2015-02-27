package daemon

import (
	"fmt"
	"os"
	"os/exec"
)

func init() {
	argc := len(os.Args)
	if argc == 1 {
		fmt.Println("Server running in console")
	} else if argc == 2 {
		if os.Args[1] == "--daemon=true" {
			cmd := exec.Command(os.Args[0], "--daemon")
			cmd.Start()
			fmt.Println("Server running in daemon . [PID]", cmd.Process.Pid)

			os.Exit(0)
		} else if os.Args[1] == "--daemon" {
			os.Stdin.Close()
			os.Stdout.Close()
			os.Stderr.Close()
			fmt.Println("Daemon Server Initializing...")
		} else {
			fmt.Println("Argument incorrect .")
			os.Exit(-1)
		}
	} else {
		fmt.Println("The number of arguments incorrect .")
		os.Exit(-1)
	}
}
