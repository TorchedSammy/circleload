# Circleload
> 📩 Command line osu! beatmap downloader

Circleload is a versatile CLI to easily download osu! beatmaps.
It supports filtering by gamemode and ranked status, different mirrors and
multiple downloads.

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
Basic usage is very simple. Circleload can take 3 types of arguments:
- A search query: `./circleload "1,000,000 times"`
- Beatmap (set) URL: `./circleload https://osu.ppy.sh/beatmapsets/1588063`
- Beatmap set ID: `./circleload 1588063`  
It can also take a variable amount of these, for mass download use.

## Filters
For a search query, Circleload supports providing filters.
Supported filters:
- `mode`: Beatmap gamemode. Accepts `osu`, `standard`, `taiko`, `mania`, and `catch`
- `status`: Ranked status. Accepts `ranked`, `loved`, `graveyard`, `pending`, `approved`, `qualified` and `wip`

# License
MIT

