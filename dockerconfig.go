/*
 * DREGU - Docker Registry Utility
 * Copyright (c) 2020 DesertBit
 */

package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strings"

	cos "git.wahtari.m/nlab/common/pkg/os"
)

type DockerConfig struct {
	Auths       map[string]Registry `json:"auths"` // key: registry-url
	HTTPHeaders map[string]string   `json:"HttpHeaders"`
}

func ParesCredentialsFromDockerConfig(filePath, registry string) (user, password string, err error) {
	// Check, if a ~/.docker/config.json exists, containing the auth token.
	ex, err := cos.Exists(filePath)
	if err != nil || !ex {
		return
	}

	// Read the file data.
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		return
	}

	var dc DockerConfig
	err = json.Unmarshal(data, &dc)
	if err != nil {
		return
	}

	// Check, if for the registry credentials are defined.
	if authReg, ok := dc.Auths[registry]; ok {
		return authReg.Credentials()
	}

	return
}

type Registry struct {
	Auth string `json:"auth"`
}

func (r Registry) Credentials() (user, password string, err error) {
	data, err := base64.StdEncoding.DecodeString(r.Auth)
	if err != nil {
		return
	}

	parts := strings.Split(string(data), ":")
	if len(parts) != 2 {
		err = fmt.Errorf("unknown credentials format: %s", string(data))
		return
	}

	user = parts[0]
	password = parts[1]
	return
}
