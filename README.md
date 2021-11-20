# Circleload
> 📩 Command line osu! beatmap downloader

Circleload is a command line tool to easily download osu! beatmaps.  
It is a work in progress! There are some things that aren't handled properly.

https://user-images.githubusercontent.com/38820196/142741654-67cc917a-ef51-4032-882a-463c5a14df6f.mp4

## Install
### Prebuilt Binaries
See the [latest release](https://github.com/TorchedSammy/Circleload/releases/latest).
`amd64` is for 64bit and `386` is 32bit. `arm64` is for 64bit ARM like the M1.

### Manual build
For users with Go installed (minimum version 1.17):  
`go install github.com/TorchedSammy/Circleload@latest`

Or:  
```sh
git clone https://github.com/TorchedSammy/Circleload
cd Circleload
go get
go build
```

## Usage
Run the binary to see usage. A simple example would be
```
Circleload 1077483 -d .
```  
The `-d` option changes the directory to download to (in this case making it our current directory.)

# License
MIT

