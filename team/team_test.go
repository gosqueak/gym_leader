package team_test

import (
	"net/http"
	"net/http/httptest"
	"os"
	"strconv"
	"testing"

	"github.com/gosqueak/leader/team"
)

var teamFileStrings = []string{
	`s1(dependencies[]url (scheme "s" domain "d" port "p") listenAddress ":8080" endpoints(/(methods[GET,POST])) jwtInfo(issuerName "i" audienceName "a")) s2(dependencies[]url (scheme "s" domain "d" port "p") listenAddress ":8080" endpoints(/(methods[GET,POST])) jwtInfo(issuerName "i" audienceName "a"))`,
}

func TestMain(m *testing.M) {
	m.Run()
}

func TestLoad(t *testing.T) {

	// write string to temp file
	fp := "TestTeamfile"

	f, err := os.OpenFile(fp, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		t.Fatal(err)
	}

	defer f.Close()
	defer os.Remove(fp)

	for _, test := range teamFileStrings {
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
	var i int

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(teamFileStrings[i]))
	}))
	defer server.Close()

	for i = range teamFileStrings {
		_, err := team.Download(server.URL + "?i=" + strconv.Itoa(i))

		if err != nil {
			t.Error(err)
		}
	}
}

func TestFromTeamfileStr(t *testing.T) {
	for _, s := range teamFileStrings {
		_, err := team.FromTeamfileStr(s)

		if err != nil {
			t.Error(err)
		}
	}
}

func TestString(t *testing.T) {
	tm, err := team.FromTeamfileStr(teamFileStrings[0])

	if err != nil {
		t.Fatal(err)
	}

	if tm.String() != teamFileStrings[0] {
		t.Fail()
	}
}

func TestMember(t *testing.T) {
	tm, err := team.FromTeamfileStr(teamFileStrings[0])

	if err != nil {
		t.Fatal(err)
	}

	if tm.Member("s1").Name != "s1" {
		t.Fail()
	}
}

func TestSaveJSON(t *testing.T) {
	t.Error("test not implemented")
}
