package unity_upm_config

import (
	"bytes"
	"github.com/BurntSushi/toml"
	"io/ioutil"
)

type ConfigElement struct {
	Token      string `toml:"token"`
	Email      string `toml:"email"`
	AlwaysAuth bool   `toml:"alwaysAuth"`
}

type Config struct {
	NpmAuth map[string]ConfigElement `toml:"npmAuth"`
}

func NewConfig() *Config {
	return &Config{NpmAuth: map[string]ConfigElement{}}
}

func Load(fn string) (*Config, error) {
	ba, err := ioutil.ReadFile(fn)
	if err != nil {
		return nil, err
	}
	return loadBA(ba)
}

func loadBA(ba []byte) (*Config, error) {
	var cfg Config
	err := toml.Unmarshal(ba, &cfg)
	if err != nil {
		return nil, err
	}
	return &cfg, nil
}

func (cfg Config) saveBA() ([]byte, error) {
	buf := new(bytes.Buffer)
	if err := toml.NewEncoder(buf).Encode(cfg); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func (cfg Config) Save(fn string) error {
	ba, err := cfg.saveBA()
	if err != nil {
		return err
	}
	return ioutil.WriteFile(fn, ba, 0644)
}
