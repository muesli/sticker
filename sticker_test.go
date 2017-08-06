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
		TTFPath:     "/usr/share/fonts/truetype/roboto/Roboto-Bold.ttf",
		MarginRatio: 0.2,
		Foreground:  color.RGBA{0x96, 0x96, 0x96, 0xff},
		Background:  color.RGBA{0xcc, 0xcc, 0xcc, 0xff},
	})
	if err != nil {
		t.Error(err)
		return
	}

	img, err := gen.NewPlaceholder("Lorem ipsum!", 400, 200)
	if err != nil {
		t.Error(err)
		return
	}

	f, err := os.Open("example/lorem.png")
	if err != nil {
		t.Error(err)
		return
	}
	defer f.Close()

	loremimg, _, err := image.Decode(f)
	if err != nil {
		t.Error(err)
		return
	}

	rimg := img.(*image.RGBA)
	rloremimg := loremimg.(*image.RGBA)
	if len(rimg.Pix) != len(rloremimg.Pix) {
		t.Error("Generated image has unexpected dimensions")
		return
	}

	errs := 0
	for i := 0; i < len(rimg.Pix); i++ {
		if rimg.Pix[i] != rloremimg.Pix[i] {
			errs++
			break
		}
	}

	if float64(errs)/float64(len(rimg.Pix)) > 0.01 {
		// The difference between the pre-generated image and the test case might vary slightly:
		// Roboto could have been updated or a newer freetype behaves different
		// We account for this by allowing 1% of pixels to differ between the two images
		t.Errorf("Expected generated image to match example/lorem.png, but it doesn't: %d pixel mismatches", errs)
	}
}

func BenchmarkPlaceholder(b *testing.B) {
}
