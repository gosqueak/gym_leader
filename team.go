package main

import (
	"encoding/json"
	"fmt"
)

type Team map[string]*ServiceNode

func (t Team) String() string {
	b, err := json.Marshal(t)
	fmt.Println(err)
	return string(b)
}

func (t Team) JSON() string {
	return t.String()
}

type ServiceNode struct {
	Name         string                  `json:"name"`
	Dependents   Dependents              `json:"dependents"`
	Dependencies Dependencies            `json:"dependencies"`
	Url          string                  `json:"url"`
	JWTInfo      JWTInfo                 `json:"jwtInfo"`
	Endpoints    map[string]EndpointInfo `json:"endpoints"`
}

func NewServiceNode(name string) *ServiceNode {
	return &ServiceNode{
		Name:         name,
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

type Dependents []string
type Dependencies []string

type EndpointInfo struct {
	Methods map[string]bool `json:"methods"`
}

type JWTInfo struct {
	AudienceName string `json:"audienceName"`
	IssuerName   string `json:"issuerName"`
}

func main() {
	team := ParseTeamfileString(teamString)
	fmt.Println(team)
}
