package main

import (
	"flag"
	"fmt"
	upm_local_proxy "github.com/Toxic2k/upm-local-proxy"
	"github.com/Toxic2k/upm-local-proxy/settings"
	"github.com/Toxic2k/upm-local-proxy/tray/icon"
	"github.com/gen2brain/dlgs"
	"github.com/getlantern/systray"
	"github.com/rs/zerolog"
	"io"
	"net/http"
	"os"
)

type trayItems struct {
	mQuit      chan struct{}
	mSetup     chan struct{}
	mResetAuth chan struct{}
	mAutoRun   *systray.MenuItem
}

var logger zerolog.Logger

func main() {

	logPathPtr := flag.String("log", "", "logfile")
	flag.Parse()

	consoleOutput := zerolog.ConsoleWriter{Out: os.Stderr, NoColor: true, TimeFormat: "15:04"}
	var multiOutput io.Writer

	if logPathPtr != nil && *logPathPtr != "" {
		logFile, err := os.OpenFile(*logPathPtr, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			_, err = dlgs.Error("Error", fmt.Sprintf("Logfile open error: %s", err.Error()))
			if err != nil {
				panic(err)
			}
			return
		}

		fileOutput := zerolog.ConsoleWriter{Out: logFile, NoColor: true, TimeFormat: "15:04"}
		multiOutput = zerolog.MultiLevelWriter(consoleOutput, fileOutput)
	} else {
		multiOutput = consoleOutput
	}

	logger = zerolog.New(multiOutput).With().Timestamp().Logger()

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
		logger.Info().Msgf("Starting proxy at %s", proxyHost)
		err := http.ListenAndServe(proxyHost, upm_local_proxy.ReverseProxy(cfg, logger))
		if err != nil {
			logger.Fatal().Err(err).Msg(err.Error())
		}
	}()

	systray.Run(func() {
		onReady(cfg, fmt.Sprintf("http://%s/", proxyHost))
	}, onExit)

}

func onExit() {
	logger.Info().Msg("onExit")
}

func onReady(cfg *settings.Config, serverHost string) {
	logger.Info().Msg("onReady")

	systray.SetIcon(icon.Data)
	systray.SetTitle("UPM local proxy")

	items := trayItems{}

	items.mAutoRun = systray.AddMenuItem("AutoRun", "Run application on computer start")

	check, err := autoRunCheck()
	if err != nil {
		logger.Error().Err(err).Msg("autorun check error")
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
				logger.Error().Err(err).Msg("autorun error")
			}
		case <-items.mQuit:
			systray.Quit()
		}
	}
}
