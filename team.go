package main

import (
	"encoding/json"
	"fmt"
)

type Team map[Name]*ServiceNode

func NewTeam() Team {
	return make(Team)
}

func LoadTeamJSON(j string) Team {
	var team Team

	err := json.Unmarshal([]byte(j), &team)
	if err != nil {
		panic(err)
	}

	return team
}

func (t Team) NewServiceNode(n Name) *ServiceNode {
	node := newServiceNode(n)
	t[node.Name] = node
	return node
}

func (t Team) String() string {
	b, err := json.Marshal(t)
	fmt.Println(err)
	return string(b)
}

func (t Team) JSON() string {
	return t.String()
}

type Service interface {
	GetAttr(key string)
}

type ServiceNode struct {
	Name         Name                    `json:"name"`
	Dependents   Dependents              `json:"dependents"`
	Dependencies Dependencies            `json:"dependencies"`
	Url          string                  `json:"url"`
	JWTInfo      JWTInfo                 `json:"jwtInfo"`
	Endpoints    map[string]EndpointInfo `json:"endpoints"`
}

func newServiceNode(n Name) *ServiceNode {
	return &ServiceNode{
		Name:         n,
		Dependents:   make(Dependents, 0),
		Dependencies: make(Dependencies, 0),
		Endpoints:    make(map[string]EndpointInfo),
	}
}

func (s *ServiceNode) Uses(other *ServiceNode) {
	s.Dependencies = append(s.Dependencies, other.Name)
	other.usedBy(s)
}

func (s *ServiceNode) usedBy(other *ServiceNode) {
	s.Dependents = append(s.Dependents, other.Name)
}

type Name string
type Dependents []Name
type Dependencies []Name

type EndpointInfo struct {
	Methods map[string]bool `json:"methods"`
}

type JWTInfo struct {
	AudienceName string `json:"audienceName"`
	IssuerName   string `json:"issuerName"`
}

func main() {
	boolMap := func(methods ...string) map[string]bool {
		m := make(map[string]bool)
		for _, method := range methods {
			m[method] = true
		}
		return m
	}

	team := NewTeam()

	///
	sx := team.NewServiceNode("steelix")
	sx.Url = "0.0.0.0:8081"
	sx.Endpoints = map[string]EndpointInfo{
		"/jwtkeypub":    {Methods: boolMap("GET")},
		"/register":     {Methods: boolMap("POST")},
		"/logout":       {Methods: boolMap("POST")},
		"/login":        {Methods: boolMap("POST")},
		"/apitokens":    {Methods: boolMap("GET")},
		"/accesstokens": {Methods: boolMap("GET")},
	}
	sx.JWTInfo.IssuerName =  "steelix"
	sx.JWTInfo.AudienceName = "steelix"

	///
	kf := team.NewServiceNode("klefki")
	kf.Url = "0.0.0.0:8083"
	kf.Endpoints = map[string]EndpointInfo{
		"/": {Methods: boolMap("POST", "PATCH", "DELETE")},
	}
	kf.JWTInfo.AudienceName = "klefki"

	///
	ak := team.NewServiceNode("alakazam")
	ak.Url = "0.0.0.0:8082"
	ak.Endpoints = map[string]EndpointInfo{
		"/ws": {Methods: boolMap("GET")},
	}
	ak.JWTInfo.AudienceName = "alakazam"

	//dependency config
	kf.Uses(sx)
	ak.Uses(sx)

	fmt.Println(team)
}
