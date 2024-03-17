package internal

import (
	"crypto/ed25519"
	"encoding/base64"
	"errors"
	"fmt"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
	"log"
	"sync/atomic"
)

type DataSnapshot struct {
	Users map[string]ed25519.PublicKey
	Teams map[string][]string
}

type DataMonitor struct {
	reload   chan struct{}
	stop     chan struct{}
	snapshot atomic.Value
}

func NewDataMonitor(path string) (*DataMonitor, error) {
	v := viper.New()
	v.SetConfigFile(path)
	if err := v.ReadInConfig(); err != nil {
		return nil, err
	}

	m := &DataMonitor{
		reload: make(chan struct{}),
		stop:   make(chan struct{}),
	}

	snapshot, err := createSnapshot(v)
	if err != nil {
		return nil, err
	}

	m.snapshot.Store(snapshot)
	go m.monitor(v)

	v.OnConfigChange(func(event fsnotify.Event) {
		m.reload <- struct{}{}
	})
	v.WatchConfig()

	return m, nil
}

func (m *DataMonitor) monitor(v *viper.Viper) {
	for {
		select {
		case <-m.stop:
			fmt.Println("stop")
			return
		case <-m.reload:
			updatedSnapshot, err := createSnapshot(v)
			if err != nil {
				log.Println(err)
			} else {
				m.snapshot.Store(updatedSnapshot)
			}
		}
	}
}

func (m *DataMonitor) GetCurrent() DataSnapshot {
	return m.snapshot.Load().(DataSnapshot)
}

func (m *DataMonitor) Close() error {
	close(m.stop)
	return nil
}

func createSnapshot(v *viper.Viper) (DataSnapshot, error) {
	type rawDataFileContent struct {
		Teams map[string][]string `yaml:"teams"`
		Users []struct {
			Name string `yaml:"name"`
			Key  string `yaml:"key"`
		} `yaml:"users"`
	}

	content := rawDataFileContent{}
	if err := v.Unmarshal(&content); err != nil {
		return DataSnapshot{}, err
	}

	users := make(map[string]ed25519.PublicKey, len(content.Users))
	for _, u := range content.Users {
		key, err := base64.StdEncoding.DecodeString(u.Key)
		if err != nil {
			return DataSnapshot{}, err
		}
		users[u.Name] = key
	}

	for _, members := range content.Teams {
		for _, m := range members {
			if _, ok := users[m]; !ok {
				return DataSnapshot{}, errors.New("user in team does not exist")
			}
		}
	}

	return DataSnapshot{Users: users, Teams: content.Teams}, nil
}
