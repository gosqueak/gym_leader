package main

import "github.com/gosqueak/leader/team"

func main() {
	team := team.Load("gosqueak.Teamfile")
	team.SaveJSON("gosqueak.Teamfile.json")
}