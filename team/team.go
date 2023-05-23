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

// cache a loaded teamfile for subsequent access
var loadedTeam Team
var isLoaded bool

func setLoaded(t Team) {
	loadedTeam = t
	isLoaded = true
}

type Team struct {
	Map
	Teamfile string
}

func Load(fp string) (Team, error) {
	if team, ok := Loaded(); ok {
		return team, nil
	}

	f, err := os.Open(fp)
	if err != nil {
		return Team{}, fmt.Errorf("failed to open file: %w", err)
	}
	defer f.Close()

	b, err := io.ReadAll(f)
	if err != nil {
		return Team{}, fmt.Errorf("failed to read file: %w", err)
	}

	team, err := FromTeamfileStr(string(b))

	if err == nil {
		setLoaded(team)
	}

	return team, err
}

func Loaded() (t Team, ok bool) {
	return loadedTeam, isLoaded
}

func Download(url string) (Team, error) {
	if team, ok := Loaded(); ok {
		return team, nil
	}

	r, err := http.Get(url)
	if err != nil {
		return Team{}, err
	}

	b, err := io.ReadAll(r.Body)
	if err != nil {
		return Team{}, err
	}

	team, err := FromTeamfileStr(string(b))

	if err == nil {
		setLoaded(team)
	}

	return team, err
}

func FromTeamfileStr(s string) (Team, error) {
	// get string and remove all whitespace and semicolons
	re := regexp.MustCompile(`[\s\n\t;]`)
	stripped := re.ReplaceAllString(s, "")

	team := Team{Map: make(Map), Teamfile: s}

	parsedObject := parseObject(stripped)

	// build the Team struct
	for k, v := range parsedObject {
		// type assert service and add service name to service
		service, ok := v.(map[string]interface{})
		if !ok {
			return Team{}, fmt.Errorf("invalid service format for service %q", k)
		}

		service["name"] = k
		team.Map[k] = &Member{}

		// marshall map to JSON then unmarshal into the new Service
		b, err := json.Marshal(service)
		if err != nil {
			return Team{}, fmt.Errorf("failed to marshal service %q to JSON: %w", k, err)
		}
		if err := json.Unmarshal(b, team.Map[k]); err != nil {
			return Team{}, fmt.Errorf("failed to unmarshal JSON into service %q: %w", k, err)
		}
	}

	// set up inter-service dependencies
	for _, service := range team.Map {
		service.configureDependencies()
	}

	return team, nil
}

func (t *Team) String() string {
	return t.Teamfile
}

func (t *Team) Member(serviceName string) (s *Member) {
	return t.Map[serviceName]
}

func (t *Team) SaveJSON(fp string) error {
	f, err := os.OpenFile(fp, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return fmt.Errorf("failed to open file: %w", err)
	}
	defer f.Close()

	b, err := json.MarshalIndent(t, "", "   ")
	if err != nil {
		return fmt.Errorf("failed to marshal Map to JSON: %w", err)
	}

	if _, err := f.Write(b); err != nil {
		return fmt.Errorf("failed to write JSON data to file: %w", err)
	}

	return nil
}

type Map map[string]*Member

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
	s.Dependencies = make([]string, 0)
	for _, otherName := range s.Dependencies {
		s.Team.Map[otherName].addDependent(s)
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
