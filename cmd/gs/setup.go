package main

import (
	"errors"
	"os"
	"os/exec"
	"strings"
)

var (
	ErrAlreadyExists = errors.New("already exists")
)

var coreLibraries = []string{
	"github.com/kjbreil/goscript/pkg/core",
	"github.com/kjbreil/goscript/pkg/module",
}

func setupProject(directory string, modules map[string]string) error {
	// make sure directory exists
	err := os.MkdirAll(directory, os.ModePerm)
	if err != nil {
		return err
	}
	initModCmd := exec.Command("go", "mod", "init", "gs")
	initModCmd.Dir = directory

	var outbuf, errbuf strings.Builder // or bytes.Buffer
	initModCmd.Stdout = &outbuf
	initModCmd.Stderr = &errbuf

	err = initModCmd.Run()
	if err != nil {
		errMsg := errbuf.String()
		if !strings.Contains(errMsg, "already exists") {
			return errors.New(errMsg)
		}
	}

	for _, lib := range coreLibraries {
		err = runGoGet(directory, lib)
	}

	for _, mod := range modules {
		err = runGoGet(directory, mod+"@modules")
	}

	return nil
}

func runGoGet(directory, addr string) error {
	cmd := exec.Command("go", "get", addr)
	cmd.Dir = directory

	var outbuf, errbuf strings.Builder // or bytes.Buffer
	cmd.Stdout = &outbuf
	cmd.Stderr = &errbuf
	err := cmd.Run()
	if err != nil {
		errMsg := errbuf.String()
		if !strings.Contains(errMsg, "already exists") {
			return errors.New(errMsg)
		}
	}

	return nil
}
