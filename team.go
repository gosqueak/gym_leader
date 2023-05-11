package leader

type jwtAudience struct{}

type serviceEndpoint struct {
	Url     string
	Methods map[string]bool
}

type serviceNodeAttributes struct {
	Url         string
	JWTAudience jwtAudience
	Endpoints   []serviceEndpoint
}

type ServiceNode struct {
	Dependents   map[*ServiceNode]bool
	Dependencies map[*ServiceNode]bool
	Attrs        serviceNodeAttributes
}

func (s *ServiceNode) Uses(other *ServiceNode) {
	s.Dependencies[other] = true
	other.usedBy(s)
}

func (s *ServiceNode) usedBy(other *ServiceNode) {
	s.Dependents[other] = true
}