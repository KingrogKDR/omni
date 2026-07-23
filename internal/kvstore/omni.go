package kvstore

import (
	"fmt"
	"os"
	"path/filepath"
)

type Omni struct {
	cfg     Config
	backend Backend
}

func Open(cfgs ...Config) (*Omni, error) {

	cfg := DefaultConfig()
	if len(cfgs) > 0 {
		cfg = cfgs[0] // override with user config
	}

	var o Omni

	if cfg.InMemOnly {
		o.backend = NewMemoryStore()
	} else {
		home, err := os.UserHomeDir()
		if err != nil {
			return nil, err
		}

		if cfg.DataDir == "" {
			cfg.DataDir = defaultDataDir
		}
		cfg.DataDir = filepath.Join(home, cfg.DataDir)

		err = os.MkdirAll(cfg.DataDir, 0o755)
		if err != nil {
			return nil, err
		}
		o.backend, err = NewPersistentStore(cfg.DataDir)
		if err != nil {
			return nil, err
		}
	}

	o.cfg = cfg

	fmt.Printf("Omni opened with cfg: %+v\n", cfg)
	return &o, nil
}

func (o *Omni) Put(key string, value []byte) error {
	return o.backend.Put(key, value)
}

func (o *Omni) Get(key string) ([]byte, error) {
	return o.backend.Get(key)
}

func (o *Omni) Delete(key string) error {
	return o.backend.Delete(key)
}

func (o *Omni) View() error {
	return o.backend.View()
}

func (o *Omni) Close() error {
	return o.backend.Close()
}
