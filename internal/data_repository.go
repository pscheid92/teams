package internal

import (
	"crypto/ed25519"
)

type DataRepository interface {
	UserExists(username string) bool
	GetUserPublicKey(username string) (ed25519.PublicKey, bool)
	GetTeamMembers(team string) ([]string, bool)
}

type YAMLFileDataRepository struct {
	monitor *DataMonitor
}

func NewYAMLFileDataRepository(monitor *DataMonitor) *YAMLFileDataRepository {
	return &YAMLFileDataRepository{monitor: monitor}
}

func (r *YAMLFileDataRepository) Snapshot() DataSnapshot {
	return r.monitor.GetCurrent()
}

func (r *YAMLFileDataRepository) UserExists(username string) bool {
	snapshot := r.Snapshot()
	_, ok := snapshot.Users[username]
	return ok
}

func (r *YAMLFileDataRepository) GetUserPublicKey(username string) (ed25519.PublicKey, bool) {
	snapshot := r.Snapshot()
	key, ok := snapshot.Users[username]
	return key, ok
}

func (r *YAMLFileDataRepository) GetTeamMembers(team string) ([]string, bool) {
	snapshot := r.Snapshot()
	members, ok := snapshot.Teams[team]
	return members, ok
}
