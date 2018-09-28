# Eos (version 2) Official App Server

## Installing

### GNU / Linux

#### Debian, and Debian-based distros

1. Install golang: `$ sudo apt install golang-go`
2. Downlaod the Eos server software: `$ go get gitlab.com/lyrenhex/eos-v2`
3. Install the Eos server software: `$ go install gitlab.com/lyrenhex/eos-v2`
4. Create `data` folder and set up [configuration file](documentation/tech.md), `data/config.json`
5. Run: `$ go run eos-v2`

#### Windows

1. Install golang from [https://golang.org]
2. Download the Eos server software: `>> go get gitlab.com/lyrenhex/eos-v2`
3. Install the Eos server software: `>> go install gitlab.com/lyrenhex/eos-v2`
4. Create `data` folder and set up [configuration file](documentation/tech.md), `data/config.json`
5. Run: `>> go run eos-v2`

## Setting up

To set up the application, a data folder should be created in the application's root. 