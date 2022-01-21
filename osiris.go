package main

import (
	"fmt"
	"log"
	"os"
	"path"
	"regexp"
	"strings"

	"github.com/fatih/color"
	"github.com/jessevdk/go-flags"
)

type Release struct {
	Title        string
	Year         string
	Episode      string
	EpisodeTitle string
	Scene        string
}

type args struct {
	Dryrun     bool   `short:"d" long:"dryrun" description:"don't modify files"`
	Silent     bool   `short:"s" long:"silent" description:"don't print file names"`
	NoColor    bool   `long:"no-color" description:"disables colored output"`
	Film       bool   `short:"f" long:"film" description:"uses film output format"`
	Year       string `long:"year" description:"release year override"`
	Title      string `long:"title" description:"release title override"`
	Positional struct {
		Regex    string `positional-arg-name:"regex" required:"true"`
		Filename string `positional-arg-name:"filename" required:"true"`
	} `positional-args:"true"`
}

func main() {
	var args args
	_, err := flags.Parse(&args)

	if e, ok := err.(*flags.Error); ok {
		if e.Type == flags.ErrHelp {
			os.Exit(0)
		} else {
			os.Exit(1)
		}
	}

	re, err := regexp.Compile(args.Positional.Regex)
	if err != nil {
		log.Fatalf("Failed to parse regex: %v\n", err)
	}

	fpath := path.Dir(args.Positional.Filename)
	fext := path.Ext(args.Positional.Filename)
	fnext := strings.TrimSuffix(args.Positional.Filename, fext)

	metadata := make(map[string]string)
	match := re.FindStringSubmatch(fnext)
	for k, name := range re.SubexpNames() {
		if k > 0 && k <= len(match) {
			metadata[name] = match[k]
		}
	}

	r := Release{
		Title:        strings.ReplaceAll(metadata["title"], ".", " "),
		Year:         metadata["year"],
		Episode:      metadata["ep"],
		EpisodeTitle: strings.ReplaceAll(metadata["eptitle"], ".", " "),
		Scene:        strings.ReplaceAll(metadata["scene"], ".", " "),
	}
	if r.Year == "" && args.Year != "" {
		r.Year = args.Year
	}
	if r.Title == "" && args.Title != "" {
		r.Title = args.Title
	}

	newname := fmt.Sprintf("%s (%s)", r.Title, r.Year)
	if !args.Film {
		newname = fmt.Sprintf("%s - %s", newname, r.Episode)
		if r.EpisodeTitle != "" {
			newname = fmt.Sprintf("%s - %s", newname, r.EpisodeTitle)
		}
	}
	if r.Scene != "" {
		newname = fmt.Sprintf("%s (%s)", newname, r.Scene)
	}
	newname = fmt.Sprintf("%s/%s%s", fpath, newname, fext)

	if !args.Silent {
		green := color.New(color.FgGreen).SprintFunc()
		fmt.Printf("%s %s %s\n", args.Positional.Filename, green("->"), newname)
	}

	if !args.Dryrun {
		if strings.TrimSpace(newname) == "" {
			log.Fatalln("Not renaming to empty string")
		}
		err = os.Rename(args.Positional.Filename, newname)
		if err != nil {
			log.Fatalf("Failed to rename: %v\n", err)
		}
	}
}
