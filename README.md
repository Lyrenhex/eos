# Eos (version 2) Official App Server

## Installing

### GNU / Linux

#### Debian, and Debian-based distros

1. Install golang and git: `$ sudo apt install golang-go git`
2. Downlaod the Eos server software: `$ git clone https://gitlab.com/lyrenhex/eos-v2`
3. Set path to the Eos folder: `$ cd eos-v2`
4. Create `data` folder and set up [configuration file](documentation/tech.md), `data/config.json`
5. Install the backend modules: `$ go install gitlab.com/lyrenhex/eos-v2`
6. Build the binary: `$ go build gitlab.com/lyrenhex/eos-v2`
  - The binary should be built in the root Eos folder (that is, the folder which contains the `data` and `webclient` folders).
7. Run: `$ ./eos-v2`

#### Windows

1. Install golang from [https://golang.org] and git from [https://git-scm.com]
2. Download the Eos server software: `>> git clone https://gitlab.com/lyrenhex/eos-v2`
3. Set path to the Eos folder: `>> cd eos-v2`
4. Create `data` folder and set up [configuration file](documentation/tech.md), `data/config.json`
5. Install the backend modules: `>> go install gitlab.com/lyrenhex/eos-v2`
6. Build the binary: `>> go build gitlab.com/lyrenhex/eos-v2`
  - The binary should be built in the root Eos folder (that is, the folder which contains the `data` and `webclient` folders).
7. Run: `>> go run server.go structs.go func.go`

## Updating

To update the software, simply perform a `git pull` in the Eos folder, re-run the backend module installation to update them, and rebuild the binary.