# engi

`import github.com/guregu/engi`

engi (演技) is a multi-platform 2D game library for Go, forked from [ENGi v0.6.0](http://ajhager.com/engi).

## Documentation

[godoc.org](http://godoc.org/github.com/guregu/engi)

## Status

*SUPER ALPHA*. Especially the audio bits.

## Differences from original engi

* Mostly working audio support
* Uses newer GLFW, fixes VSync issues.
* JS version is broken

## Audio

`Sounds` are loaded entirely in to memory, designed for sound effects. Files ending with `.flac`, and `.wav` will be loaded as Sounds.

`Streams` are streamed, designed for background music. Files ending with `.flac-s`. This system is pretty dumb/hacky so I will fix it eventually. 


## Desktop

The desktop backend depends on glfw3, but includes the source code and links it statically. If you are having linker errors on Windows, I suggest using [TDM-GCC](http://tdm-gcc.tdragon.net/download) instead of MinGW as your cgo compiler. Linux may need `xorg-dev`. 

## Web

This fork has broken GopherJS support for the time being. 

## Install

```bash
go get -u github.com/guregu/engi
```

## Success stories

* [HOT PLUG](http://hotplug.kawaii.solutions)

## Other libraries

* [Ebiten](http://hajimehoshi.github.io/ebiten/)
* [paked/engi](https://github.com/paked/engi): another engi fork, focus on Entity Component Systems.
* [ajhager/engi](https://github.com/ajhager/engi): the original
