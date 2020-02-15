package upm_local_proxy

import (
	"encoding/json"
	"github.com/Toxic2k/upm-local-proxy/settings"
	"io/ioutil"
)

type unityScopedRegistry struct {
	Name   string   `json:"name"`
	Url    string   `json:"url"`
	Scopes []string `json:"scopes"`
}

type unityScopedRegistries []unityScopedRegistry

func SetupManifest(fn string, cfg *settings.Config, serverHost string) error {

	ba, err := ioutil.ReadFile(fn)
	if err != nil {
		return err
	}

	var man map[string]interface{}
	err = json.Unmarshal(ba, &man)
	if err != nil {
		return err
	}

	var rgs unityScopedRegistries
	var scopes []string
	for _, r := range cfg.Registries {
		for _, s := range r.Scopes {
			scopes = append(scopes, s)
		}
	}
	rgs = append(rgs, unityScopedRegistry{
		Name:   "UPM local proxy",
		Url:    serverHost,
		Scopes: scopes,
	})

	man["scopedRegistries"] = rgs

	ba2, err := json.MarshalIndent(man, "", "  ")
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(fn, ba2, 0644)
	if err != nil {
		return err
	}

	return nil

}
