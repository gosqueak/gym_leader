package team

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
)

type Team map[string]*Service

func FromJSON(j []byte) (t Team) {
	err := json.Unmarshal(j, &t)
	if err != nil {
		panic(err)
	}

	return t
}

func Load(fp string) Team {
	return fromTeamFile(fp)
}

func (t Team) String() string {
	return string(t.JSON())
}

func (t Team) JSON() []byte {
	b, _ := json.Marshal(t)
	return b
}

func (t Team) SaveJSON(fp string) error {
	return ioutil.WriteFile(fp, t.JSON(), 0644)
}

type EndpointInfo struct {
	Methods []string `json:"methods"`
}

type JWTInfo struct {
	AudienceName string `json:"audienceName"`
	IssuerName   string `json:"issuerName"`
}

type Service struct {
	Name         string                  `json:"name"`
	Dependents   []string                `json:"dependents"`
	Dependencies []string                `json:"dependencies"`
	Url          string                  `json:"url"`
	JWTInfo      JWTInfo                 `json:"jwtInfo"`
	Endpoints    map[string]EndpointInfo `json:"endpoints"`
}

func (s *Service) uses(other *Service) {
	s.Dependencies = append(s.Dependencies, other.Name)
	other.usedBy(s)
}

func (s *Service) usedBy(other *Service) {
	s.Dependents = append(s.Dependents, other.Name)
}


func FetchJSON(url string) []byte {
	r, err := http.Get(url)
	if err != nil {
		panic(err)
	}

	b, err := io.ReadAll(r.Body)
	if err != nil {
		panic(err)
	}

	return b
}
