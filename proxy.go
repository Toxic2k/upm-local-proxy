package upm_local_proxy

import (
	"bytes"
	"fmt"
	"github.com/Toxic2k/upm-local-proxy/settings"
	"io/ioutil"
	"net/http"
	"net/http/httputil"
	"strings"
)

func singleJoiningSlash(a, b string) string {
	aslash := strings.HasSuffix(a, "/")
	bslash := strings.HasPrefix(b, "/")
	switch {
	case aslash && bslash:
		return a + b[1:]
	case !aslash && !bslash:
		return a + "/" + b
	}
	return a + b
}

func findRepo(cfg *settings.Config, path string) *settings.ConfigRegistry {
	for i := len(cfg.Registries) - 1; i >= 1; i-- {
		reg := cfg.Registries[i]
		for _, scope := range reg.Scopes {
			if strings.Contains(path, scope) {
				return reg
			}
		}
	}
	return nil
}

func ReverseProxy(cfg *settings.Config) *httputil.ReverseProxy {
	director := func(req *http.Request) {
		r := findRepo(cfg, req.URL.Path)
		if r == nil {
			r = cfg.Registries[0]
		}

		req.URL.Scheme = r.Url.Scheme
		req.URL.Host = r.Url.Host
		req.URL.Path = singleJoiningSlash(r.Url.Path, req.URL.Path)
		req.Host = ""
		if r.Url.RawQuery == "" || req.URL.RawQuery == "" {
			req.URL.RawQuery = r.Url.RawQuery + req.URL.RawQuery
		} else {
			req.URL.RawQuery = r.Url.RawQuery + "&" + req.URL.RawQuery
		}
		if _, ok := req.Header["User-Agent"]; !ok {
			// explicitly disable User-Agent so it's not set to default value
			req.Header.Set("User-Agent", "")
		}
		//req.SetBasicAuth(user,pass)
		req.Header.Set("Authorization", "Bearer "+r.Token)
		req.Header.Del("Accept-Encoding")
	}

	modifyResponse := func(response *http.Response) error {
		contType := response.Header.Get("Content-Type")
		if !strings.HasPrefix(contType, "text/html") && !strings.HasPrefix(contType, "application/json") {
			return nil
		}

		r := findRepo(cfg, response.Request.URL.Path)
		if r == nil {
			r = cfg.Registries[0]
		}

		ba, err := ioutil.ReadAll(response.Body)
		if err != nil {
			return err
		}
		replaced := bytes.ReplaceAll(ba, []byte(r.UrlString), []byte(fmt.Sprintf("http://localhost:%d", cfg.Port)))

		response.Body = ioutil.NopCloser(bytes.NewReader(replaced))

		return nil
	}

	return &httputil.ReverseProxy{Director: director, ModifyResponse: modifyResponse}
}
