package main

import "github.com/gosqueak/leader/team"

func main() {
	team := team.Load("Teamfile")
	team.SaveJSON("Teamfile.json")
}