/*
 * The MIT License (MIT)
 *
 * Copyright (c) 2022.
 *
 * Permission is hereby granted, free of charge, to any person obtaining a copy
 * of this software and associated documentation files (the "Software"), to deal
 * in the Software without restriction, including without limitation the rights
 * to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
 * copies of the Software, and to permit persons to whom the Software is
 * furnished to do so, subject to the following conditions:
 *
 * The above copyright notice and this permission notice shall be included in all
 * copies or substantial portions of the Software.
 *
 * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
 * IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
 * FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
 * AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
 * LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
 * OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
 * SOFTWARE.
 */

package main

import (
	"errors"
	"log"
	"os"
	"path"

	"github.com/BurntSushi/toml"
	"gopkg.in/yaml.v3"
)

type config struct {
	SeriesYear *bool `yaml:"seriesYear,omitempty"`
	Scene      *bool `yaml:"scene,omitempty"`
	Templates  struct {
		Series *string `yaml:"series,omitempty"`
		Film   *string `yaml:"film,omitempty"`
	} `yaml:"templates,omitempty"`
	Regex struct {
		Series *string `yaml:"series,omitempty"`
		Film   *string `yaml:"film,omitempty"`
		Custom struct {
			Series map[*string]*string `yaml:"series,omitempty"`
			Film   map[*string]*string `yaml:"film,omitempty"`
		} `yaml:"custom,omitempty"`
	} `yaml:"regex,omitempty"`
}

func (c *config) Parse(data []byte) error {
	err := yaml.Unmarshal(data, c)
	if err != nil {
		return toml.Unmarshal(data, c)
	}
	return err
}

// Argparse replaces config values with provided cli flag values
func (c *config) Argparse(args *args) {
	if args.Regex != "" {
		if args.Film {
			c.Regex.Film = &args.Regex
		} else {
			c.Regex.Series = &args.Regex
		}
	}
	if args.Scene {
		c.Scene = &args.Scene
	}
	if args.SeriesYear {
		c.SeriesYear = &args.SeriesYear
	}
}

func ConfigFile() (string, error) {
	cfgdir, err := os.UserConfigDir()
	if err != nil {
		log.Fatalln(err)
	}
	cfgFile := path.Join(cfgdir, "osiris", "osiris.yml")

	if !fileExists(cfgFile) {
		cfgFile = path.Join(cfgdir, "osiris", "osiris.yaml")
	}

	if !fileExists(cfgFile) {
		cfgFile = path.Join(cfgdir, "osiris", "osiris.toml")
	}

	if !fileExists(cfgFile) {
		return "", errors.New("config file does not exist")
	}

	return cfgFile, nil
}
