# Discordtranscriber
![img](https://i.imgur.com/lDkpxvE.png)
## Installing
`go get -u gitlab.com/koishi/discordtranscriber/...`

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
|------|-----------------------------|
| u    | Username                    |
| p    | Password                    |
| d    | Specify custom asset dir    |
| port | Server port (default: 8100) |
| t    | Account token               |

## Building

### Dependencies:
| Dependency                                                          | Reason                                              |
|---------------------------------------------------------------------|-----------------------------------------------------|
| [go-bindata](https://github.com/jteeuwen/go-bindata)                | Embedding webui inside executeable                  |
| [go-bindata-assetfs](https://github.com/elazarl/go-bindata-assetfs) | implementing http.FileSystem interface with bindata |

run `go generate` then `go build`
