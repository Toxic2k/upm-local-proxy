package main

import (
	"github.com/Toxic2k/upm-local-proxy"
	"github.com/Toxic2k/upm-local-proxy/settings"
	"github.com/gen2brain/dlgs"
	"os"
	"path/filepath"
)

func setupProject(cfg *settings.Config, serverHost string) {

	fn, res, err := dlgs.File("Select project", "Select unity project to add scoped registry", true)
	if err != nil {
		logger.Error().Err(err).Msg("select project error")
		return
	}
	if !res {
		return
	}
	manifestPath := filepath.Join(fn, "Packages", "manifest.json")
	if _, err = os.Stat(manifestPath); os.IsNotExist(err) {
		_, err = dlgs.Error("Error", "selected folder is not a Unity 3d Project")
		if err != nil {
			logger.Error().Err(err).Msg("error popup error")
		}
		return
	}

	err = upm_local_proxy.SetupManifest(manifestPath, cfg, serverHost)
	if err != nil {
		logger.Error().Err(err).Msg("setupManifest error")
	}

}
