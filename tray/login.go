package main

import (
	"fmt"
	upm_local_proxy "github.com/Toxic2k/upm-local-proxy"
	"github.com/Toxic2k/upm-local-proxy/settings"
	"github.com/gen2brain/dlgs"
)

func repoLogin(cfg *settings.Config) bool {
	configChanged := false
	for _, r := range cfg.Registries {
		if r.Login == "" {
			eUser, res, err := dlgs.Entry("Auth", fmt.Sprintf("Enter your login for %s", r.Name), "")
			if err != nil {
				logger.Error().Err(err).Msg("login dialog error")
				break
			}
			if res {
				r.Login = eUser
				configChanged = true
			} else {
				break
			}
		}

		if r.Token == "" {
			ePass, res, err := dlgs.Password("Auth", fmt.Sprintf("Enter your password for %s", r.Name))
			if err != nil {
				logger.Error().Err(err).Msg("password dialog error")
				break
			}
			if res {
				r.Pass = ePass
			} else {
				break
			}

			err = upm_local_proxy.GetToken(r, logger)
			if err != nil {
				logger.Error().Err(err).Msgf("get token for %s", r.Name)
				r.Pass = ""
				_, err = dlgs.Error("Error", "unauthorized")
				if err != nil {
					panic(err)
				}
				return false
			}
			if settings.TokenAutoSave {
				configChanged = true
			}
		}
	}
	if configChanged {
		err := cfg.Save()
		if err != nil {
			_, err = dlgs.Error("Error", fmt.Sprintf("Config save error: %s", err.Error()))
			if err != nil {
				panic(err)
			}
		}
	}
	return true
}
