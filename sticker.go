/*
 * sticker - a Golang lib to generate placeholder images with text
 *     Copyright (c) 2017, Christian Muehlhaeuser <muesli@gmail.com>
 *
 *   For license see LICENSE
 */

package sticker

import (
	"errors"
	"image"
	"image/color"
	"image/draw"
	"io/ioutil"

	"github.com/disintegration/imaging"
	"github.com/golang/freetype"
	"github.com/golang/freetype/truetype"
	"golang.org/x/image/font"
)

const (
	dpi         = 72.00
	maxFontSize = 512.00
)

// ImageGenerator lets you generate images & placeholders
type ImageGenerator struct {
	options Options
	font    *truetype.Font
}

// Options contains all the settings for an ImageGenerator
type Options struct {
	TTFPath         string
	TTF             []byte
	Foreground      color.RGBA
	Background      color.RGBA
	BackgroundImage image.Image
	MarginRatio     float64
}

var (
	// ErrInvalidDimensions gets returned when the requested image has an invalid width or height
	ErrInvalidDimensions = errors.New("values for width or height must be positive")

	// ErrMissingFontOption gets returned when there's no font specified in the options
	ErrMissingFontOption = errors.New("no font option specified")
)

// NewImageGenerator returns a new ImageGenerator
func NewImageGenerator(options Options) (*ImageGenerator, error) {
	if options.Foreground.A == 0 {
		options.Foreground = color.RGBA{0x96, 0x96, 0x96, 0xff}
	}
	if options.Background.A == 0 {
		options.Background = color.RGBA{0xcc, 0xcc, 0xcc, 0xff}
	}
	if options.MarginRatio < 0 {
		options.MarginRatio = 0.2
	}

	var ttf []byte

	if len(options.TTF) > 0 {
		ttf = options.TTF
	} else if options.TTFPath != "" {
		var err error

		ttf, err = ioutil.ReadFile(options.TTFPath)

		if err != nil {
			return nil, err
		}
	} else {
		return nil, ErrMissingFontOption
	}

	f, err := freetype.ParseFont(ttf)
	if err != nil {
		return nil, err
	}

	return &ImageGenerator{
		font:    f,
		options: options,
	}, nil
}

// NewPlaceholder returns a placeholder image with the given text, width & height
func (p *ImageGenerator) NewPlaceholder(text string, width, height int) (image.Image, error) {
	if width < 0 || height < 0 {
		return nil, ErrInvalidDimensions
	}
	if width == 0 && height == 0 {
		return nil, ErrInvalidDimensions
	}

	if width == 0 {
		width = height
	} else if height == 0 {
		height = width
	}

	img := image.NewRGBA(image.Rect(0, 0, width, height))

	c := freetype.NewContext()
	c.SetDPI(dpi)
	c.SetFont(p.font)
	c.SetSrc(image.NewUniform(color.RGBA{0, 0, 0, 0}))
	c.SetDst(img)
	c.SetClip(img.Bounds())
	c.SetHinting(font.HintingNone)

	// draw the background
	draw.Draw(img, img.Bounds(), image.NewUniform(p.options.Background), image.ZP, draw.Src)
	// draw background image
	if p.options.BackgroundImage != nil {
		bgimg := imaging.Fill(p.options.BackgroundImage, width, height, imaging.Center, imaging.Lanczos)
		draw.Draw(img, img.Bounds(), bgimg, image.ZP, draw.Src)
	}

	if text != "" {
		// draw with scaled fontsize to get the real text extent
		fontsize, actwidth := maxPointSize(text, c,
			int(float64(width)*(1.0-p.options.MarginRatio)),
			int(float64(height)*(1.0-p.options.MarginRatio)))

		actheight := c.PointToFixed(fontsize/2.0) / 64
		xcenter := (float64(width) / 2.0) - (float64(actwidth) / 2.0)
		ycenter := (float64(height) / 2.0) + (float64(actheight) / 2.0)

		// draw the text
		c.SetFontSize(fontsize)
		c.SetSrc(image.NewUniform(p.options.Foreground))
		_, err := c.DrawString(text, freetype.Pt(int(xcenter), int(ycenter)))
		if err != nil {
			return nil, err
		}
	}

	return img, nil
}

// maxPointSize returns the maximum point size we can use to fit text inside width and height
// as well as the resulting text-width in pixels
func maxPointSize(text string, c *freetype.Context, width, height int) (float64, int) {
	// never let the font size exceed the requested height
	fontsize := maxFontSize
	for int(c.PointToFixed(fontsize)/64) > height {
		fontsize -= 2
	}

	// find the biggest matching font size for the requested width
	var actwidth int
	for actwidth = width + 1; actwidth > width; fontsize -= 2 {
		c.SetFontSize(fontsize)

		textExtent, err := c.DrawString(text, freetype.Pt(0, 0))
		if err != nil {
			return 0, 0
		}

		actwidth = int(float64(textExtent.X) / 64)
	}

	return fontsize, actwidth
}
