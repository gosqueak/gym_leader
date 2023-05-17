package main

import (
	"fmt"
	"os"

	"github.com/gosqueak/leader/team"
)

func main() {
	mode := os.Args[1]

	switch mode {

	case "export":
		teamfileFP := os.Args[2]
		jsonFileFP := os.Args[3]

		fmt.Printf("exporting %v to %v", teamfileFP, jsonFileFP)

		tm, err := team.Load(teamfileFP)
		if err != nil {
			panic(err)
		}

		tm.SaveJSON(jsonFileFP)

	default:
		fmt.Printf("Error: Unknown mode '%v'\n", mode)
	}
}
