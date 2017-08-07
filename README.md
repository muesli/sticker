sticker
=======

A Golang lib to generate placeholder images with text.

## Installation

Make sure you have a working Go environment. See the [install instructions](http://golang.org/doc/install.html).

To install sticker, simply run:

    go get github.com/muesli/sticker

To compile it from source:

    cd $GOPATH/src/github.com/muesli/sticker
    go get -u -v
    go build && go test -v

## Example
```go
gen, err := sticker.NewImageGenerator(sticker.Options{
    TTFPath:         "/usr/share/fonts/TTF/Roboto-Bold.ttf",
    MarginRatio:     0.2,
    Foreground:      color.RGBA{0x96, 0x96, 0x96, 0xff},
    Background:      color.RGBA{0xcc, 0xcc, 0xcc, 0xff},
    BackgroundImage: img,
})
img, err := gen.NewPlaceholder("Lorem ipsum!", 400, 200)
...
```

![example placeholder](example/lorem.png) ![example placeholder with background image](example/motivation.png)

If you supply a background image, it will automatically be cropped and scaled (while maintaining its original aspect ratio) to the desired output size.

## Development

API docs can be found [here](http://godoc.org/github.com/muesli/sticker).

[![Build Status](https://secure.travis-ci.org/muesli/sticker.png)](http://travis-ci.org/muesli/sticker)
[![Coverage Status](https://coveralls.io/repos/github/muesli/sticker/badge.svg?branch=master)](https://coveralls.io/github/muesli/sticker?branch=master)
[![Go ReportCard](http://goreportcard.com/badge/muesli/sticker)](http://goreportcard.com/report/muesli/sticker)
