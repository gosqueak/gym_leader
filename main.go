package main

import (
	"fmt"
	"os"

	"github.com/gosqueak/leader/team"
)

func main() {
	mode := os.Args[1]

	switch mode {
	case "download":
		url := os.Args[2]
		fp := os.Args[3]

		fmt.Printf("downloading teamfile json from %v to %v\n", url, fp)
		tm := team.Download(url)
		tm.SaveJSON(fp)
	
	case "export":
		teamfileFP := os.Args[2]
		jsonFileFP := os.Args[3]

		fmt.Printf("exporting %v to %v", teamfileFP, jsonFileFP)
		
		tm := team.Load(teamfileFP)
		tm.SaveJSON(jsonFileFP)
	}
}