# osiris

[![Build Status](https://ci.gryffyn.io/api/badges/gryffyn/osiris/status.svg?ref=refs/heads/main)](https://ci.gryffyn.io/gryffyn/osiris)  
A tool for renaming films / tv series based on named regex capture groups.

## installation

`go install git.gryffyn.io/gryffyn/osiris@latest`

*or*

```shell
git clone https://git.gryffyn.io/gryffyn/osiris
cd osiris
go build
```

## usage
```
Usage:
  osiris [OPTIONS] [regex] [filename...]

Application Options:
  -d, --dryrun    don't modify files
  -s, --silent    don't print file names
      --no-color  disables colored output
  -f, --film      uses film output format
  -y, --year=     release year override
  -t, --title=    release title override

Help Options:
  -h, --help      Show this help message
```

## regex parameters

the regex argument takes in a regex with specifically-named named capture groups.

| name    | example               | required | series | film | description                 |
|---------|-----------------------|:--------:|--------|------|-----------------------------|
| title   | `(?P<title>\w+)`      |   ✅    |    ✅    |   ✅   | title of the series/film    |
| year    | `(?P<year>\d{4})`     |    ❌    |    ✅    |   ✅   | release year                |
| ep      | `(?P<ep>S\d{2}E\d{2})` |   ✅    |   ✅     |    ❌  | episode number (ex. S01E01) |
| eptitle | `(?P<eptitle>\w+)`    |    ❌    |    ✅    |    ❌  | episode title               |
| scene   | `(?P<scene>[\w-\.]+)`  |    ❌    |    ✅    |   ✅   | scene / release info        |

## output naming format

odin follows a standard-ish format, namely

`Series Title - S01E01 - Episode Title (SCENE INFO).ext`

or in the case of a film (`-f, --film`)

`Film Title (YEAR) (SCENE INFO).ext`

## example

```shell
$ osiris -d '(?P<title>\w+)\.(?P<ep>S\d{2}E\d{2})\.(?P<eptitle>[\w\.]+)\.(?P<scene>WEB-.+)' Title.S01E01.Episode.Title.WEB-DL.H264.SCENE.mkv
Title.S01E01.Episode.Title.WEB-DL.H264.SCENE.mkv -> ./Title - S01E01 - Episode Title (WEB-DL H264 SCENE).mkv
```
