# Osiris

[![Build Status](https://ci.gryffyn.io/api/badges/gryffyn/osiris/status.svg?ref=refs/heads/main)](https://ci.gryffyn.io/gryffyn/osiris)  
A tool for renaming films / tv series based on named regex capture groups and Go templates.

## Installation

`go install git.gryffyn.io/gryffyn/osiris@latest`

*or*

```shell
git clone https://git.gryffyn.io/gryffyn/osiris
cd osiris
go build
```

## Usage
```
Usage:
  osiris [OPTIONS] [filename...]

Application Options:
  -d, --dryrun      don't modify files
  -s, --silent      don't print file names
      --no-color    disables colored output
  -f, --film        uses film output format
  -Y, --seriesyear  whether series year is output
  -S, --scene       whether scene info is output
  -y, --year=       release year override
  -t, --title=      release title override
  -c, --config=     config file (default ~/.config/osiris/osiris.yml)
  -r, --regex=      input regex pattern
  -p, --preset=     preset input regex

Help Options:
  -h, --help        Show this help message
```

## Input Regex

the regex argument takes in a regex with specifically-named named capture groups.

| name    | Output Template Name | example                | required | series | film | description                 |
|---------|----------------------|------------------------|:--------:|--------|------|-----------------------------|
| title   | Title                | `(?P<title>\w+)`       |    ✅     | ✅      | ✅    | title of the series/film    |
| year    | Year                 | `(?P<year>\d{4})`      |    ❌     | ✅      | ✅    | release year                |
| ep      | Episode              | `(?P<ep>S\d{2}E\d{2})` |    ✅     | ✅      | ❌    | episode number (ex. S01E01) |
| eptitle | EpisodeTitle         | `(?P<eptitle>\w+)`     |    ❌     | ✅      | ❌    | episode title               |
| scene   | Scene                | `(?P<scene>[\w-\.]+)`  |    ❌     | ✅      | ✅    | scene / release info        |

## Output Templates

### Default

| Type   | Template                                                                                                                                                               | Example                                   |
|--------|------------------------------------------------------------------------------------------------------------------------------------------------------------------------|-------------------------------------------|
| Series | {{ .Title }}{{if .Options.SeriesYear}} ({{ .Year }}){{end}} - {{ .Episode}}{{if EpisodeTitle}} - {{ .EpisodeTitle}}{{end}}{{if .Options.Scene }} ({{ .Scene }}){{end}} | Series Title - S01E01 - Episode Title.ext |
| Film   | {{ .Title }} ({{ .Year }}){{if .Options.Scene }} ({{ .Scene }}){{end}}                                                                                                 | Film Title (YEAR).ext                     |


## Config File Example

The configuration file is by default loaded from a file named `osiris.yml` the current user's configuration directory (`~/.config/osiris` on Linux). The config file location can be overridden with the `-c` flag.

```yaml
---
seriesYear: true
scene: true
templates:
  series: "{{ .Title }}{{if .Options.SeriesYear}} ({{ .Year }}){{end}} - {{ .Episode}}{{if EpisodeTitle}} - {{ .EpisodeTitle}}{{end}}{{if .Options.Scene }} ({{ .Scene }}){{end}}"
  film: "{{ .Title }} ({{ .Year }}){{if .Options.Scene }} ({{ .Scene }}){{end}}"
regex:
  series: '(?P<title>[\w\.]+)\.(?P<ep>S\d{2}E\d{2})\.(?P<eptitle>[\w\.-]+)(?P<scene>1080p.+)'
  film:
  custom:
    yify: '(?P<title>[\w\.]+)\.(?P<year>\d{4})\.(?P<scene>(?:2160|1080|720)p[\w\.]+YIFY)'
```

## Usage Example

```shell
$ osiris -d '(?P<title>\w+)\.(?P<ep>S\d{2}E\d{2})\.(?P<eptitle>[\w\.]+)\.(?P<scene>WEB-.+)' Title.S01E01.Episode.Title.WEB-DL.H264.SCENE.mkv
Title.S01E01.Episode.Title.WEB-DL.H264.SCENE.mkv -> Title - S01E01 - Episode Title.mkv
```
