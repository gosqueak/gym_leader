package team_test

import (
	"os"
	"testing"

	"github.com/gosqueak/leader/team"
)

func TestMain(m *testing.M) {
	m.Run()
}

func TestLoad(t *testing.T) {
	passCases := []string{
		`s1(dependencies[]url (scheme "s" domain "d" port "p") listenAddress ":8080" endpoints(/(methods[GET,POST])) jwtInfo(issuerName "i" audienceName "a")) s2(dependencies[]url (scheme "s" domain "d" port "p") listenAddress ":8080" endpoints(/(methods[GET,POST])) jwtInfo(issuerName "i" audienceName "a"))`,
	}

	// write string to temp file
	fp := "TestTeamfile"

	f, err := os.OpenFile(fp, os.O_RDWR | os.O_CREATE | os.O_TRUNC, 0644)
	if err!= nil {
		t.Fatal(err)
	}

	defer f.Close()
	defer os.Remove(fp)

	for _, test := range passCases {
		f.Truncate(0)
		_, err := f.WriteString(test)
		if err != nil {
			t.Fatal(err)
		}

		// load team from temp team file
		_, err = team.Load(fp)

		if err != nil {
			t.Error(err)
		}
	}
}

func TestDownload(t *testing.T) {

}

func TestFromTeamfileStr(t *testing.T) {

}

func TestString(t *testing.T) {

}

func TestMember(t *testing.T) {

}

func TestSaveJSON(t *testing.T) {

}