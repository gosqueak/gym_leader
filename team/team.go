package team

import (
	"encoding/json"
	"io"
	"os"
	"net/http"
)

type Team map[string]*Service

func Load(fp string) Team {
	return fromTeamFile(fp)
}

func Download(url string) Team {
	r, err := http.Get(url)
	if err != nil {
		panic(err)
	}

	b, err := io.ReadAll(r.Body)
	if err != nil {
		panic(err)
	}

	return fromJSON(b)
}

func fromJSON(j []byte) (t Team) {
	err := json.Unmarshal(j, &t)
	if err != nil {
		panic(err)
	}

	return t
}

func (t Team) String() string {
	return string(t.JSON())
}

func (t Team) JSON() []byte {
	b, err := json.MarshalIndent(t, "", "   ")
	if err != nil {
		panic(err)
	}

	return b
}

func (t Team) SaveJSON(fp string) {
	b := t.JSON()
	err := os.WriteFile(fp, b, 0644)
	if err != nil {
		panic(err)
	}
}

type Service struct {
	Name         string                  `json:"name"`
	Url          string                  `json:"url"`
	Dependents   []string                `json:"dependents"`
	Dependencies []string                `json:"dependencies"`
	JWTInfo      JWTInfo                 `json:"jwtInfo"`
	Endpoints    map[string]EndpointInfo `json:"endpoints"`
}

func (s *Service) usedBy(other *Service) {
	s.Dependents = append(s.Dependents, other.Name)
}

type EndpointInfo struct {
	Methods []string `json:"methods"`
}

type JWTInfo struct {
	AudienceName string `json:"audienceName"`
	IssuerName   string `json:"issuerName"`
}
