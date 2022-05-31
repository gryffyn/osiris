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

type release struct {
	Title        string
	Year         string
	Episode      string
	EpisodeTitle string
	Scene        string
}

type inputMetadata struct {
	Film  *bool
	Year  *string
	Title *string
}

type args struct {
	Dryrun     bool   `short:"d" long:"dryrun" description:"don't modify files"`
	Silent     bool   `short:"s" long:"silent" description:"don't print file names"`
	NoColor    bool   `long:"no-color" description:"disables colored output"`
	Film       bool   `short:"f" long:"film" description:"uses film output format"`
	Year       string `short:"y" long:"year" description:"release year override"`
	Title      string `short:"t" long:"title" description:"release title override"`
	Positional struct {
		Regex    string   `positional-arg-name:"regex" required:"true"`
		Filename []string `positional-arg-name:"filename" required:"true"`
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

	inmet := inputMetadata{
		Film:  &args.Film,
		Year:  &args.Year,
		Title: &args.Title,
	}

	re, err := regexp.Compile(args.Positional.Regex)
	if err != nil {
		log.Fatalf("Failed to parse regex: %v\n", err)
	}

	for _, filename := range args.Positional.Filename {
		newfilename := getFilename(filename, re, &inmet)
		if !args.Silent {
			printRename(&filename, &newfilename)
		}
		if !args.Dryrun {
			renameFile(&filename, &newfilename)
		}
	}
}

func getFilename(filepath string, re *regexp.Regexp, inputmetadata *inputMetadata) string {
	fpath := path.Dir(filepath)
	fext := path.Ext(filepath)
	fnext := strings.TrimSuffix(filepath, fext)

	metadata := make(map[string]string)
	match := re.FindStringSubmatch(fnext)
	for k, name := range re.SubexpNames() {
		if k > 0 && k <= len(match) {
			metadata[name] = match[k]
		}
	}

	r := release{
		Title:        strings.ReplaceAll(metadata["title"], ".", " "),
		Year:         metadata["year"],
		Episode:      metadata["ep"],
		EpisodeTitle: strings.TrimSpace(strings.ReplaceAll(metadata["eptitle"], ".", " ")),
		Scene:        strings.TrimSpace(strings.ReplaceAll(metadata["scene"], ".", " ")),
	}
	if r.Year == "" && *inputmetadata.Year != "" {
		r.Year = *inputmetadata.Year
	}
	if r.Title == "" && *inputmetadata.Title != "" {
		r.Title = *inputmetadata.Title
	}

	var newname string
	if !*inputmetadata.Film {
		newname = fmt.Sprintf("%s - %s", r.Title, r.Episode)
		if r.EpisodeTitle != "" {
			newname = fmt.Sprintf("%s - %s", newname, r.EpisodeTitle)
		}
	} else {
		newname = fmt.Sprintf("%s (%s)", r.Title, r.Year)
	}
	if r.Scene != "" {
		newname = fmt.Sprintf("%s (%s)", newname, r.Scene)
	}

	return fmt.Sprintf("%s/%s%s", fpath, newname, fext)
}

func printRename(filepath, newfilepath *string) {
	green := color.New(color.FgGreen).SprintFunc()
	fmt.Printf("%s %s %s\n", path.Base(*filepath), green("->"), path.Base(*newfilepath))
}

func renameFile(filepath, newfilepath *string) {
	if strings.TrimSpace(*newfilepath) == "" {
		log.Fatalln("Not renaming to empty string")
	}
	err := os.Rename(*filepath, *newfilepath)
	if err != nil {
		log.Fatalf("Failed to rename: %v\n", err)
	}
}
