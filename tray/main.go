package main

import (
	"flag"
	"fmt"
	upm_local_proxy "github.com/Toxic2k/upm-local-proxy"
	"github.com/Toxic2k/upm-local-proxy/settings"
	"github.com/Toxic2k/upm-local-proxy/tray/icon"
	unity_upm_config "github.com/Toxic2k/upm-local-proxy/unity-upm-config"
	"github.com/gen2brain/dlgs"
	"github.com/getlantern/systray"
	"github.com/lxn/walk"
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

	tempDir := os.TempDir()
	err := os.Chdir(tempDir)
	if err != nil {
		panic(err)
	}
	err = walk.Resources.SetRootDirPath(tempDir)
	if err != nil {
		panic(err)
	}

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

		//multiOutput = zerolog.MultiLevelWriter(consoleOutput, logFile)
		multiOutput = logFile
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

	unityCfg, err := unity_upm_config.LoadConfig()
	if err != nil {
		_, err = dlgs.Error("Error", fmt.Sprintf("Unity Config load error: %s", err.Error()))
		if err != nil {
			panic(err)
		}
		return
	}

	for {
		if repoLogin(cfg, unityCfg) {
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
		onReady(cfg, unityCfg, fmt.Sprintf("http://%s/", proxyHost))
	}, onExit)

}

func onExit() {
	logger.Info().Msg("onExit")
}

func onReady(cfg *settings.Config, unityCfg *unity_upm_config.Config, serverHost string) {
	logger.Info().Msg("onReady")

	systray.SetIcon(icon.Data)
	systray.SetTitle(fmt.Sprintf("UPM local proxy %s", settings.VERSION))

	items := trayItems{}

	systray.AddMenuItem(fmt.Sprintf("UPM local proxy %s", settings.VERSION), "").Disable()
	systray.AddSeparator()

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

	go onClicks(items, cfg, unityCfg, serverHost)

}

func onClicks(items trayItems, cfg *settings.Config, unityCfg *unity_upm_config.Config, serverHost string) {
	for {
		select {
		case <-items.mSetup:
			setupProject(cfg, serverHost)
		case <-items.mResetAuth:
			cfg.ResetAuth()
			for {
				if repoLogin(cfg, unityCfg) {
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
