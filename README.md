# Circleload
> 📩 Command line osu! beatmap downloader

Circleload is a command line tool to easily download osu! beatmaps.  
It is a work in progress! There are some things that aren't handled properly.

Circleload downloads maps from unofficial mirrors, since downloading from osu.ppy.sh
requires logging in.

## Install
For users with Go installed (minimum version 1.17):
`go install github.com/TorchedSammy/Circleload@latest`

## Usage
Run the `Circleload` binary to see usage. A simple example would be
```
Circleload 1077483 -d .
```  
The `-d` option changes the directory to download to (in this case making it our current directory.)

# License
MIT

