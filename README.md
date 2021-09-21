## Installation

Compile https://gitlab.com/gomidi/midicat and copy it to PATH (this is for windows compat) - also see the driver at https://gitlab.com/gomidi/midicatdrv.

Then go build this.

Or cross-compile for windows on macOS:

```
brew install mingw-w64
CC=x86_64-w64-mingw32-gcc GOOS=windows GOARCH=amd64 CGO_ENABLED=1 go build -o warryall.exe
```
