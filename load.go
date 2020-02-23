package upm_local_proxy

import "github.com/Toxic2k/upm-local-proxy/settings"

func LoadConfig() (*settings.Config, bool, error) {
	cfgFn, err := settings.GetConfigFinename()
	if err != nil {
		return nil, false, err
	}

	cfg, def, err := settings.LoadConfig(cfgFn)
	if err != nil {
		return nil, def, err
	}

	for i := 0; i < len(cfg.Registries); i++ {
		if cfg.Registries[i].SavedToken != "" {
			cfg.Registries[i].Token = cfg.Registries[i].SavedToken
		}
	}

	return cfg, def, nil
}
