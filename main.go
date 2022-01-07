/*
 * DREGU - Docker Registry Utility
 * Copyright (c) 2020 DesertBit
 */

package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/desertbit/docker-registry-client/registry"
	"github.com/desertbit/grumble"
	"github.com/fatih/color"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

const (
	flagVerbose = "verbose"
	flagJSON    = "json"

	flagDockerConfig = "docker-config"
	flagRegistry     = "registry"
	flagUser         = "user"
	flagPassword     = "password"
)

// Create the grumble app.
var App = grumble.New(&grumble.Config{
	Name:        "dregu",
	Description: "Docker-REGistry-Utility offers a convenient API to a docker-registry",
	PromptColor: color.New(color.FgGreen, color.Bold),

	Flags: func(f *grumble.Flags) {
		f.Bool("v", flagVerbose, false, "verbose mode")
		f.Bool("j", flagJSON, false, "JSON log mode")

		dcDefault := ""
		uhd, err := os.UserHomeDir()
		if err == nil {
			dcDefault = filepath.Join(uhd, ".docker/config.json")
		}
		f.StringL(flagDockerConfig, dcDefault, "the config.json file of the local .docker dir")
		f.String("r", flagRegistry, "docker.wahtari.m", "the hostname of the docker registry")
		f.String("u", flagUser, "docker", "the user used to access the docker registry")
		f.String("p", flagPassword, "", "the password used to access the docker registry")
	},
})

var reg *registry.Registry

func main() {
	App.OnInit(func(a *grumble.App, f grumble.FlagMap) (err error) {
		// Check JSON logging.
		if !f.Bool(flagJSON) {
			log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
		}

		// Check verbose mode.
		if f.Bool(flagVerbose) {
			zerolog.SetGlobalLevel(zerolog.DebugLevel)
		} else {
			zerolog.SetGlobalLevel(zerolog.InfoLevel)
		}

		// Connect to the registry.
		if f.String(flagDockerConfig) != "" {
			reg, err = connectFromConfig(f.String(flagRegistry), f.String(flagDockerConfig))
		} else {
			reg, err = connectWithCreds(f.String(flagRegistry), f.String(flagUser), f.String(flagPassword))
		}
		if err != nil {
			return err
		}

		return nil
	})

	grumble.Main(App)
}

func connectFromConfig(registryHost, configFP string) (reg *registry.Registry, err error) {
	opts := registry.Options{
		Logf: func(format string, args ...interface{}) {},
	}

	// Try to parse the credentials from the docker config.
	opts.Username, opts.Password, err = parseCredentialsFromDockerConfig(configFP, registryHost)
	if err != nil {
		log.Debug().Err(err).Msg("parse credentials from docker config")
	}

	return connect(registryHost, opts)
}

func connectWithCreds(registryHost, user, password string) (reg *registry.Registry, err error) {
	opts := registry.Options{
		Username: user,
		Password: password,
		Logf:     func(format string, args ...interface{}) {},
	}

	return connect(registryHost, opts)
}

func connect(registryHost string, opts registry.Options) (reg *registry.Registry, err error) {
	// Prepare registry client using the options.
	reg, err = registry.NewCustom("https://"+registryHost, opts)
	if err != nil {
		err = fmt.Errorf("could not connect to docker registry at https://%s: %v", registryHost, err)
		return
	}

	return
}
