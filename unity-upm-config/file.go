package unity_upm_config

import (
	"os"
	"path/filepath"
)

const configFilename = ".upmconfig.toml"

func GetConfigFinename() (string, error) {
	homedir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(homedir, configFilename), nil
}

func LoadConfig() (*Config, error) {
	fn, err := GetConfigFinename()
	if err != nil {
		return nil, err
	}
	if _, err := os.Stat(fn); os.IsNotExist(err) {
		return &Config{NpmAuth: map[string]ConfigElement{}}, nil
	}
	return Load(fn)
}

func (cfg Config) SaveConfig() error {
	fn, err := GetConfigFinename()
	if err != nil {
		return err
	}
	return cfg.Save(fn)
}
