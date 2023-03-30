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
    "bytes"
    "errors"
    "fmt"
    "log"
    "os"
    "path"
    "regexp"
    "strings"
    "text/template"

    "github.com/fatih/color"
    "github.com/jessevdk/go-flags"
)

type releaseOptions struct {
    Scene      bool
    SeriesYear bool
}

type release struct {
    Title        string
    Year         string
    Episode      string
    EpisodeTitle string
    Scene        string
    Options      releaseOptions
}

type args struct {
    Dryrun     bool   `short:"d" long:"dryrun" description:"don't modify files"`
    Silent     bool   `short:"s" long:"silent" description:"don't print file names"`
    NoColor    bool   `long:"no-color" description:"disables colored output"`
    Film       bool   `short:"f" long:"film" description:"uses film output format"`
    SeriesYear bool   `short:"Y" long:"seriesyear" description:"whether series year is output"`
    Scene      bool   `short:"S" long:"scene" description:"whether scene info is output"`
    Year       string `short:"y" long:"year" description:"release year override"`
    Title      string `short:"t" long:"title" description:"release title override"`
    ConfigFile string `short:"c" long:"config" description:"config file (default ~/.config/osiris/osiris.yml)"`
    Regex      string `short:"r" long:"regex" description:"input regex pattern"`
    Preset     string `short:"p" long:"preset" description:"preset input regex"`
    Positional struct {
        Filename []string `positional-arg-name:"filename" required:"true"`
    } `positional-args:"true"`
}

var (
    seriesTemplate = "{{ .Title }}{{if .Options.SeriesYear}} ({{ .Year }}){{end}} - {{ .Episode}}{{if ." +
        "EpisodeTitle}} - {{ .EpisodeTitle}}{{end}}{{if .Options.Scene }} ({{ .Scene }}){{end}}"
    filmTemplate = "{{ .Title }} ({{ .Year }}){{if .Options.Scene }} ({{ .Scene }}){{end}}"
)

func main() {
    var args args
    _, err := flags.Parse(&args)

    var e *flags.Error
    if ok := errors.Is(err, e); ok {
        if e.Type == flags.ErrHelp {
            os.Exit(0)
        } else {
            os.Exit(1)
        }
    }

    var data []byte
    cfgFile, err := ConfigFile()
    if err != nil {
        log.Fatalln(err)
    }
    if args.ConfigFile != "" {
        cfgFile = args.ConfigFile
    }

    data, err = os.ReadFile(cfgFile)
    if err != nil {
        log.Fatalln(err)
    }

    cfg := config{}
    err = cfg.Parse(data)
    if err != nil {
        log.Fatalf("Error reading config file (%s): %v", cfgFile, err)
    }
    cfg.Argparse(&args)

    var re *regexp.Regexp

    if args.Preset != "" {
        pre, err := getCustomPreset(&cfg, &args.Preset, args.Film)
        if err != nil {
            log.Fatalln(err)
        }
        re = regexp.MustCompile(*pre)
    } else {
        if args.Film {
            if cfg.Regex.Film == nil || *cfg.Regex.Film == "" {
                log.Fatalln("Regex must be provided by `-r` flag or in osiris.yml")
            }
            re, err = regexp.Compile(*cfg.Regex.Film)
        } else {
            if cfg.Regex.Series == nil || *cfg.Regex.Series == "" {
                log.Fatalln("Regex must be provided by `-r` flag or in osiris.yml")
            }
            re, err = regexp.Compile(*cfg.Regex.Series)
        }
    }

    if err != nil {
        log.Fatalf("Failed to parse regex: %v\n", err)
    }

    for _, filename := range args.Positional.Filename {
        newfilename := getFilename(filename, re, &cfg, args.Year, args.Title, args.Film)
        if !args.Silent {
            printRename(&filename, &newfilename) //nolint:gosec
        }
        if !args.Dryrun {
            renameFile(&filename, &newfilename) //nolint:gosec
        }
    }
}

func getFilename(filepath string, re *regexp.Regexp, cfg *config, year, title string, film bool) string {
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

    scene := regexp.MustCompile(`\D(\.)\D`).ReplaceAllStringFunc(metadata["scene"], func(s string) string {
        return strings.ReplaceAll(s, ".", " ")
    })

    r := release{
        Title:        strings.TrimSpace(strings.ReplaceAll(metadata["title"], ".", " ")),
        Year:         strings.TrimSpace(metadata["year"]),
        Episode:      strings.TrimSpace(strings.ToUpper(metadata["ep"])),
        EpisodeTitle: strings.TrimSpace(strings.ReplaceAll(metadata["eptitle"], ".", " ")),
        Scene:        strings.TrimSpace(scene),
        Options: releaseOptions{
            Scene:      *cfg.Scene,
            SeriesYear: *cfg.SeriesYear,
        },
    }
    if r.Year == "" && year != "" {
        r.Year = year
    }
    if r.Title == "" && title != "" {
        r.Title = title
    }

    var tmpl string

    if film {
        if *cfg.Templates.Film != "" {
            tmpl = *cfg.Templates.Film
        } else {
            tmpl = filmTemplate
        }
    } else {
        if *cfg.Templates.Series != "" {
            tmpl = *cfg.Templates.Series
        } else {
            tmpl = seriesTemplate
        }
    }

    releaseTemplate, _ := template.New("release").Parse(tmpl)
    var newname bytes.Buffer
    err := releaseTemplate.Execute(&newname, &r)
    if err != nil {
        log.Fatalln(err)
    }

    return fmt.Sprintf("%s/%s%s", fpath, newname.String(), fext)
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

func getCustomPreset(cfg *config, preset *string, film bool) (*string, error) {
    switch film {
    case true:
        for k, v := range cfg.Regex.Custom.Film {
            if *k == *preset {
                return v, nil
            }
        }
    case false:
        for k, v := range cfg.Regex.Custom.Series {
            if *k == *preset {
                return v, nil
            }
        }
    }

    return nil, fmt.Errorf("preset '%s' does not exist", *preset)
}

func fileExists(path string) bool {
    _, err := os.Stat(path)
    return !errors.Is(err, os.ErrNotExist)
}
