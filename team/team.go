package team

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"regexp"
	"strings"
)

type Team struct {
	teamMap
	Teamfile string
}

func Load(fp string) Team {
	f, err := os.Open(fp)
	if err != nil {
		panic(fmt.Errorf("failed to open file: %w", err))
	}
	defer f.Close()

	b, err := io.ReadAll(f)
	if err != nil {
		panic(fmt.Errorf("failed to read file: %w", err))
	}

	return FromTeamfileStr(string(b))
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

	return FromTeamfileStr(string(b))
}

func FromTeamfileStr(s string) Team {
	// get string and remove all whitespace and semicolons
	re := regexp.MustCompile(`[\s\n\t;]`)
	stripped := re.ReplaceAllString(s, "")

	team := Team{teamMap: make(teamMap), Teamfile: s}

	parsedObject := parseObject(stripped)

	// build the Team struct
	for k, v := range parsedObject {
		// type assert service and add service name to service
		service, ok := v.(map[string]interface{})
		if !ok {
			panic(fmt.Errorf("invalid service format for service %q", k))
		}

		service["name"] = k
		team.teamMap[k] = &Member{}

		// marshall map to JSON then unmarshal into the new Service
		b, err := json.Marshal(service)
		if err != nil {
			panic(fmt.Errorf("failed to marshal service %q to JSON: %w", k, err))
		}
		if err := json.Unmarshal(b, team.teamMap[k]); err != nil {
			panic(fmt.Errorf("failed to unmarshal JSON into service %q: %w", k, err))
		}
	}

	// set up inter-service dependencies
	for _, service := range team.teamMap {
		service.configureDependencies()
	}

	return team
}


func (t *Team) String() string {
	return t.Teamfile
}

func (t *Team) Member(serviceName string) (s *Member) {
	return t.teamMap[serviceName]
}

type teamMap map[string]*Member

func (t teamMap) String() string {
	return string(t.AsJSON())
}

func (t teamMap) AsJSON() []byte {
	b, err := json.MarshalIndent(t, "", "   ")
	if err != nil {
		panic(fmt.Errorf("failed to marshal teamMap to JSON: %w", err))
	}
	return b
}

func (t teamMap) SaveJSON(fp string) {
	f, err := os.OpenFile(fp, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		panic(fmt.Errorf("failed to open file: %w", err))
	}
	defer f.Close()

	data := t.AsJSON()
	if _, err := f.Write(data); err != nil {
		panic(fmt.Errorf("failed to write JSON data to file: %w", err))
	}
}


type Member struct {
	Team          Team
	Name          string                  `json:"name"`
	Url           url                     `json:"url"`
	Endpoints     map[string]endpointInfo `json:"endpoints"`
	ListenAddress string                  `json:"listenAddress"`
	Dependents    []string                `json:"dependents"`
	Dependencies  []string                `json:"dependencies"`
	JWTInfo       jwtInfo                 `json:"jwtInfo"`
}

// Add the Service as a dependent to all depended on services
func (s *Member) configureDependencies() {
	for _, otherName := range s.Dependencies {
		s.Team.teamMap[otherName].addDependent(s)
	}
}

func (s *Member) addDependent(other *Member) {
	s.Dependents = append(s.Dependents, other.Name)
}

type endpointInfo struct {
	Methods []string `json:"methods"`
}

type jwtInfo struct {
	AudienceName string `json:"audienceName"`
	IssuerName   string `json:"issuerName"`
}

type url struct {
	Scheme string `json:"scheme"`
	Domain string `json:"domain"`
	Port   string `json:"port"`
	Path   string `json:"path"`
}

func (u url) String() string {
	var b strings.Builder
	put := func(s string) { b.WriteString(s) }

	if u.Scheme != "" {
		put(u.Scheme + "://")
	}

	put(u.Domain)

	if u.Port != "" {
		put(":" + u.Port)
	}

	put(u.Path)

	return b.String()
}
