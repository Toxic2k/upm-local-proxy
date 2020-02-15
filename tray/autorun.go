package main

import (
	"golang.org/x/sys/windows/registry"
	"log"
	"os"
)

const keyName = "upm-local-proxy"

func autoRunCheck() (bool, error) {
	key, _, err := registry.CreateKey(registry.CURRENT_USER, "SOFTWARE\\Microsoft\\Windows\\CurrentVersion\\Run", registry.ALL_ACCESS)
	if err != nil {
		return false, err
	}
	defer key.Close()

	_, _, err = key.GetStringValue(keyName)
	if err == registry.ErrNotExist {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return true, nil
}

func autoRunToggle(run bool) error {
	execPath := os.Args[0]
	log.Printf("execPath: %s", execPath)

	key, _, err := registry.CreateKey(registry.CURRENT_USER, "SOFTWARE\\Microsoft\\Windows\\CurrentVersion\\Run", registry.ALL_ACCESS)
	if err != nil {
		return err
	}
	defer key.Close()

	if run {
		err = key.SetStringValue(keyName, execPath)
		if err != nil {
			return err
		}
	} else {
		return key.DeleteValue(keyName)
	}

	return nil
}
