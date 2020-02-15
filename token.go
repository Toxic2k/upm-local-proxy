package upm_local_proxy

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/Toxic2k/upm-local-proxy/settings"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

type loginJson struct {
	Id       string    `json:"_id"`
	Name     string    `json:"name"`
	Password string    `json:"password"`
	Type     string    `json:"type"`
	Roles    []string  `json:"roles"`
	Date     time.Time `json:"date"`
}

type loginResponseJson struct {
	Token string `json:"token"`
}

func GetToken(cfg *settings.ConfigRegistry) error {
	uri := fmt.Sprintf("%s://%s/-/user/org.couchdb.user:%s", cfg.Url.Scheme, cfg.Url.Host, cfg.Login)
	log.Print(uri)

	js := loginJson{
		Id:       fmt.Sprintf("org.couchdb.user:%s", cfg.Login),
		Name:     cfg.Login,
		Password: cfg.Pass,
		Type:     "user",
		Roles:    []string{},
		Date:     time.Now(),
	}

	ba, err := json.Marshal(js)
	if err != nil {
		return err
	}

	r := bytes.NewReader(ba)

	req, err := http.NewRequest(http.MethodPut, uri, r)
	if err != nil {
		return err
	}

	req.Header.Set("content-type", "application/json")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}

	if res.StatusCode == http.StatusCreated {
		ba, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return err
		}

		var jsRes loginResponseJson
		err = json.Unmarshal(ba, &jsRes)
		if err != nil {
			return err
		}

		cfg.Token = jsRes.Token
		return nil
	}

	return fmt.Errorf("wrong response status: %s", res.Status)
}
