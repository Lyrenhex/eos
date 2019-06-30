# Eos Official App Server

## Installing

### GNU / Linux

#### Debian, and Debian-based distros

1. Install golang and git: `$ sudo apt install golang-go git` and set `PATH` to include `$GOHOME/bin`
2. Download and install the Eos server software: `$ go install github.com/lyrenhex/eos`
3. Create `data` folder and set up [configuration file](documentation/tech.md), `~/eos/data/` and `~/eos/config.json`
4. Run: `$ eos`

#### Windows

1. Install golang from [golang.org](https://golang.org/dl/) and git from [git-scm.com](https://git-scm.com/downloads) and set `PATH` to include `%GOHOME%/bin`
2. Download and install the Eos server software: `>> go install github.com/lyrenhex/eos`
3. Create `data` folder and set up [configuration file](documentation/tech.md), `%USERPROFILE%/eos/data/` and `%USERPROFILE%/eos/config.json`
4. Run: `>> eos`

## Updating

To update the software, re-run `go install github.com/lyrenhex/eos`.
