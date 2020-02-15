package main

import (
	"fmt"
	upm_local_proxy "github.com/Toxic2k/upm-local-proxy"
	"github.com/Toxic2k/upm-local-proxy/settings"
	"github.com/Toxic2k/upm-local-proxy/tray/icon"
	"github.com/gen2brain/dlgs"
	"github.com/getlantern/systray"
	"log"
	"net/http"
)

type trayItems struct {
	mQuit      chan struct{}
	mSetup     chan struct{}
	mResetAuth chan struct{}
	mAutoRun   *systray.MenuItem
}

func main() {

	cfg, def, err := upm_local_proxy.LoadConfig()
	if err != nil {
		_, err = dlgs.Error("Error", fmt.Sprintf("Config load error: %s", err.Error()))
		if err != nil {
			panic(err)
		}
		return
	}

	if def {
		_, err = dlgs.Info("Hi!", fmt.Sprintf("Default config was written to %s", cfg.Filename))
		if err != nil {
			panic(err)
		}
	}

	for {
		if repoLogin(cfg) {
			break
		}
	}

	proxyHost := fmt.Sprintf("localhost:%d", cfg.Port)

	go func() {
		log.Printf("Starting proxy at %s", proxyHost)
		log.Fatal(http.ListenAndServe(proxyHost, upm_local_proxy.ReverseProxy(cfg)))
	}()

	systray.Run(func() {
		onReady(cfg, fmt.Sprintf("http://%s/", proxyHost))
	}, onExit)

}

func onExit() {
	log.Printf("onExit\n")
}

func onReady(cfg *settings.Config, serverHost string) {
	log.Printf("onReady\n")

	systray.SetIcon(icon.Data)
	systray.SetTitle("UPM local proxy")

	items := trayItems{}

	items.mAutoRun = systray.AddMenuItem("AutoRun", "Run application on computer start")

	check, err := autoRunCheck()
	if err != nil {
		log.Printf("autorun check error: %v", err)
	}
	if check {
		items.mAutoRun.Check()
	}

	systray.AddSeparator()
	items.mSetup = systray.AddMenuItem("Setup Project", "Setup Unity Project").ClickedCh
	items.mResetAuth = systray.AddMenuItem("Reset Authorization", "Reset Authorization").ClickedCh
	systray.AddSeparator()
	items.mQuit = systray.AddMenuItem("Quit", "Exit").ClickedCh

	go onClicks(items, cfg, serverHost)

}

func onClicks(items trayItems, cfg *settings.Config, serverHost string) {
	for {
		select {
		case <-items.mSetup:
			setupProject(cfg, serverHost)
		case <-items.mResetAuth:
			cfg.ResetAuth()
			for {
				if repoLogin(cfg) {
					break
				}
			}
		case <-items.mAutoRun.ClickedCh:
			var err error
			if items.mAutoRun.Checked() {
				err = autoRunToggle(false)
				items.mAutoRun.Uncheck()
			} else {
				err = autoRunToggle(true)
				items.mAutoRun.Check()
			}
			if err != nil {
				log.Printf("autorun error: %v", err)
			}
		case <-items.mQuit:
			systray.Quit()
		}
	}
}
