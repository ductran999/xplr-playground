package identity

import "time"

type Agent struct {
	ID                string
	RegistrationToken string
	Version           string
	Metadata          AgentMetadata
	RegisteredAt      time.Time
}

type AgentMetadata struct {
	NodeName   string
	PodName    string
	Namespace  string
	Hostname   string
	K8SVersion string
}
