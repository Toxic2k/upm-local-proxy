package main

import (
	"fmt"
	"github.com/Toxic2k/upm-local-proxy"
	"github.com/Toxic2k/upm-local-proxy/settings"
	"github.com/gen2brain/dlgs"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
)

func unityAuthSupported(fn string) (bool, error) {
	versionPath := filepath.Join(fn, "ProjectSettings", "ProjectVersion.txt")

	if _, err := os.Stat(versionPath); os.IsNotExist(err) { // too old unity
		return false, nil
	}

	ba, err := ioutil.ReadFile(versionPath)
	if err != nil {
		return false, err
	}

	var re = regexp.MustCompile(`(?m)m_EditorVersion: (\d+)\.(\d+)\.(\d+)`)
	tmp := re.FindAllStringSubmatch(string(ba), -1)
	if len(tmp) != 1 || len(tmp[0]) != 4 {
		return false, fmt.Errorf("wrong ProjectVersion.txt format")
	}

	major, err := strconv.ParseInt(tmp[0][1], 10, 64)
	if err != nil {
		return false, err
	}
	minor, err := strconv.ParseInt(tmp[0][2], 10, 64)
	if err != nil {
		return false, err
	}
	patch, err := strconv.ParseInt(tmp[0][3], 10, 64)
	if err != nil {
		return false, err
	}

	if major >= 2019 && minor >= 3 && patch >= 4 {
		return true, nil
	}

	return false, nil
}

func setupProject(cfg *settings.Config, serverHost string) {

	fn, res, err := dlgs.File("Select project", "Select unity project to add scoped registry", true)
	if err != nil {
		logger.Error().Err(err).Msg("select project error")
		return
	}
	if !res {
		return
	}

	authSupported, err := unityAuthSupported(fn)
	if err != nil {
		logger.Error().Err(err).Msg("setupManifest error")
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

	err = upm_local_proxy.SetupManifest(manifestPath, cfg, serverHost, authSupported)
	if err != nil {
		logger.Error().Err(err).Msg("setupManifest error")
	}

}
