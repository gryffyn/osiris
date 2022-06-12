package main

import (
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
	} `yaml:"regex,omitempty"`
}

func (c *config) Parse(data []byte) error {
	return yaml.Unmarshal(data, c)
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
