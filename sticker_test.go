/*
 * sticker - a Golang lib to generate placeholder images with text
 *     Copyright (c) 2017, Christian Muehlhaeuser <muesli@gmail.com>
 *
 *   For license see LICENSE
 */

package sticker

import (
	"image"
	"image/color"
	"os"
	"testing"

	_ "image/png"
)

func TestPlaceholder(t *testing.T) {
	gen, err := NewImageGenerator(Options{
		TTFPath:     "/usr/share/fonts/TTF/Roboto-Bold.ttf",
		MarginRatio: 0.2,
		Foreground:  color.RGBA{0x96, 0x96, 0x96, 0xff},
		Background:  color.RGBA{0xcc, 0xcc, 0xcc, 0xff},
	})
	if err != nil {
		t.Error(err)
	}

	img, err := gen.NewPlaceholder("Lorem ipsum!", 400, 200)
	if err != nil {
		t.Error(err)
	}

	f, err := os.Open("example/lorem.png")
	if err != nil {
		t.Error(err)
	}
	defer f.Close()

	loremimg, _, err := image.Decode(f)
	if err != nil {
		t.Error(err)
	}

	rimg := img.(*image.RGBA)
	rloremimg := loremimg.(*image.RGBA)
	for i := 0; i < len(rimg.Pix); i++ {
		if i >= len(rloremimg.Pix) || rimg.Pix[i] != rloremimg.Pix[i] {
			t.Error("Expected generated image to match example/lorem.png, but it doesn't")
			break
		}
	}
}

func BenchmarkPlaceholder(b *testing.B) {
}
