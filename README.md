# Discordtranscriber
![img](https://i.imgur.com/BygMSrr.png)
## Installing
`go get -u gitlab.com/Necroforger/discordtranscriber/cmd/...`

### Dependencies for voice synthesis
The following programs are required for voice synthesis to work

| Dependency                                     | Description                                              |
| ---------------------------------------------- | -------------------------------------------------------- |
| [espeak](http://espeak.sourceforge.net/)       | Export the espeak command line executeable to your path. |
| [ffmpeg](https://www.ffmpeg.org/download.html) | Needed to convert espeak output to opus                  |

## Usage

### Without a token
`discordtranscriber "username" "password"`

### With a token
`discordtranscriber -t "Bot token"`

### Interactive
If you launch the executeable by double clicking on it or by executing without arguments It will prompt you to enter your credentials.

## Flags
`args <username> <password>`


| Flag | Description                 |
| ---- | --------------------------- |
| u    | Username                    |
| p    | Password                    |
| d    | Specify custom asset dir    |
| port | Server port (default: 8100) |
| t    | Account token               |

## Building

### Dependencies:
| Dependency                                                          | Reason                                              |
| ------------------------------------------------------------------- | --------------------------------------------------- |
| [go-bindata](https://github.com/jteeuwen/go-bindata)                | Embedding webui inside executeable                  |
| [go-bindata-assetfs](https://github.com/elazarl/go-bindata-assetfs) | implementing http.FileSystem interface with bindata |

run `go generate` then `go build`
