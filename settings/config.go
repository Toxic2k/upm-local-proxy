package settings

import (
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"net/url"
	"os"
	"path/filepath"
)

type ConfigRegistry struct {
	Name       string   `yaml:"name"`
	UrlString  string   `yaml:"url"`
	Scopes     []string `yaml:"scopes"`
	Login      string   `yaml:"login"`
	Email      string   `yaml:"email"`
	Pass       string   `yaml:"-"`
	Token      string   `yaml:"-"`
	SavedToken string   `yaml:"saved_token,omitempty"`
	Url        url.URL  `yaml:"-"`
}

type Config struct {
	Filename   string            `yaml:"-"`
	Port       int32             `yaml:"port"`
	Registries []*ConfigRegistry `yaml:"registries"`
}

func Default(fn string) (*Config, error) {
	reg := ConfigRegistry{
		Name:      RepoName,
		UrlString: RepoUrl,
		Scopes:    RepoScopes,
	}
	u, err := url.Parse(reg.UrlString)
	if err != nil {
		return nil, err
	}
	reg.Url = *u

	return &Config{
		Filename: fn,
		Port:     18080,
		Registries: []*ConfigRegistry{
			&reg,
		},
	}, nil
}

func Load(fn string) (*Config, error) {
	ba, err := ioutil.ReadFile(fn)
	if err != nil {
		return nil, err
	}
	var res Config
	err = yaml.Unmarshal(ba, &res)
	if err != nil {
		return nil, err
	}

	for _, r := range res.Registries {
		u, err := url.Parse(r.UrlString)
		if err != nil {
			return nil, err
		}
		r.Url = *u
	}

	res.Filename = fn

	return &res, nil
}

func (cfg Config) Save() error {
	ba, err := yaml.Marshal(cfg)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(cfg.Filename, ba, 0644)
}

const configFilename = "upm-local-proxy.yml"

func GetConfigFinename() (string, error) {
	homedir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(homedir, configFilename), nil
}

func LoadConfig(fn string) (*Config, bool, error) {
	if _, err := os.Stat(fn); os.IsNotExist(err) {
		cfg, err := Default(fn)
		if err != nil {
			return nil, false, err
		}
		err = cfg.Save()
		if err != nil {
			return nil, false, err
		}
		return cfg, true, nil
	}
	cfg, err := Load(fn)
	return cfg, false, err
}

func (cfg *Config) ResetAuth() {
	for _, registry := range cfg.Registries {
		registry.Login = ""
		registry.Email = ""
		registry.Pass = ""
		registry.Token = ""
		registry.SavedToken = ""
	}
}
