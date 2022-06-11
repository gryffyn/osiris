package main

import "gopkg.in/yaml.v2"

type config struct {
	SeriesYear *bool `yaml:"seriesYear,omitempty"`
	Scene      *bool `yaml:"scene,omitempty"`
	Templates  struct {
		Series *string `yaml:"series,omitempty"`
		Film   *string `yaml:"film,omitempty"`
	} `yaml:"templates"`
}

func (c *config) Parse(data []byte) error {
	return yaml.Unmarshal(data, c)
}

func (c *config) Argparse(args *args) {
	if args.Scene {
		c.Scene = &args.Scene
	}
	if args.SeriesYear {
		c.SeriesYear = &args.SeriesYear
	}
}
