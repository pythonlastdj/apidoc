// Copyright 2015 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package main

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/caixw/apidoc/app"
	"github.com/caixw/apidoc/input"
	"github.com/caixw/apidoc/output"
)

type config struct {
	Version string          `json:"version"` // 产生该配置文件的程序版本号，主版本号不同，表示不兼容
	Input   *input.Options  `json:"input"`
	Output  *output.Options `json:"output"`
}

// 从配置文件中加载配置项。
func loadConfig() (*config, error) {
	wd, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	data, err := ioutil.ReadFile(filepath.Join(wd, app.ConfigFilename))
	if err != nil {
		return nil, err
	}

	cfg := &config{}
	if err = json.Unmarshal(data, cfg); err != nil {
		return nil, err
	}

	if err := cfg.Input.Init(); err != nil {
		return nil, err
	}

	if err := cfg.Output.Init(); err != nil {
		return nil, err
	}

	return cfg, nil
}

// 在当前目录下产生个默认的配置文件。
func genConfigFile() error {
	wd, err := os.Getwd()
	if err != nil {
		return err
	}

	fi, err := os.Create(filepath.Join(wd, app.ConfigFilename))
	if err != nil {
		return err
	}
	defer fi.Close()

	lang, err := input.DetectDirLang(wd)
	if err != nil { // 不中断，仅作提示用。
		app.Warn(err)
	}

	cfg := &config{
		Version: app.Version,
		Input: &input.Options{
			Dir:       "./",
			Recursive: true,
			Lang:      lang,
		},
		Output: &output.Options{
			Type: "html",
			Dir:  "./apidoc",
		},
	}
	data, err := json.MarshalIndent(cfg, "", "    ")
	_, err = fi.Write(data)
	return err
}