# Eos (version 2) Official App Server

## Installing

### GNU / Linux

#### Debian, and Debian-based distros

1. Install golang and git: `$ sudo apt install golang-go git`
2. Downlaod the Eos server software: `$ git clone https://gitlab.com/lyrenhex/eos-v2`
3. Set path to the Eos folder: `$ cd eos-v2`
4. Create `data` folder and set up [configuration file](documentation/tech.md), `data/config.json`
5. Run: `$ go run server.go structs.go func.go`

#### Windows

1. Install golang from [https://golang.org] and git from [https://git-scm.com]
2. Download the Eos server software: `>> git clone https://gitlab.com/lyrenhex/eos-v2`
3. Set path to the Eos folder: `>> cd eos-v2`
4. Create `data` folder and set up [configuration file](documentation/tech.md), `data/config.json`
5. Run: `>> go run server.go structs.go func.go`

## Updating

To update the software, simply perform a `git fetch` in the Eos folder.