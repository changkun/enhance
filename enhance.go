// Copyright 2023 Changkun Ou <changkun.de>. All rights reserved.
// Use of this source code is governed by a MIT license that
// can be found in the LICENSE file.

package enhance

import (
	"image"
	"image/color"
	"image/draw"
	"math"
)

// Color RGB from 0 to 1.
type Color struct {
	R, G, B float64
}

func NewColor(c color.RGBA) Color {
	return Color{
		R: float64(c.R) / 255.0,
		G: float64(c.G) / 255.0,
		B: float64(c.B) / 255.0,
	}
}

func (c Color) ToRGBA() color.RGBA {
	return color.RGBA{
		R: uint8(c.R * 255),
		G: uint8(c.G * 255),
		B: uint8(c.B * 255),
		A: 255,
	}
}

// Pixel changes the color of a pixel based on the given parameters.
func Pixel(p Color, params Params) Color {
	brightness := clamp(params.Brightness) - 0.5
	contrast := clamp(params.Contrast) - 0.5
	saturation := clamp(params.Saturation) - 0.5
	temperature := clamp(params.Temperature) - 0.5
	tint := clamp(params.Tint) - 0.5

	c := Color{
		R: fromsRGB2Linear(p.R),
		G: fromsRGB2Linear(p.G),
		B: fromsRGB2Linear(p.B),
	}
	c = applyTemperatureTintEffect(c, temperature, tint)
	c = applyBrightnessEffect(c, brightness)
	c = applyContrastEffect(c, contrast)
	c = applySaturationEffect(c, saturation)
	return Color{
		R: clamp(fromLinear2sRGB(c.R)),
		G: clamp(fromLinear2sRGB(c.G)),
		B: clamp(fromLinear2sRGB(c.B)),
	}
}

func imageToRGBA(src image.Image) *image.RGBA {
	// No conversion needed if image is an *image.RGBA.
	if dst, ok := src.(*image.RGBA); ok {
		return dst
	}

	// Use the image/draw package to convert to *image.RGBA.
	b := src.Bounds()
	dst := image.NewRGBA(image.Rect(0, 0, b.Dx(), b.Dy()))
	draw.Draw(dst, dst.Bounds(), src, b.Min, draw.Src)
	return dst
}

func clamp(v float64) float64 { return math.Max(0, math.Min(v, 1.0)) }

// Y'UV (BT.709) to linear RGB
// Values are from https://en.wikipedia.org/wiki/YUV
func yuv2rgb(c Color) Color {
	return Color{
		R: +1.00000*c.R + 0.00000*c.G + 1.28033*c.B,
		G: +1.00000*c.R - 0.21482*c.G - 0.38059*c.B,
		B: +1.00000*c.R + 2.12798*c.G + 0.00000*c.B,
	}
}

// Linear RGB to Y'UV (BT.709)
// Values are from https://en.wikipedia.org/wiki/YUV
func rgb2yuv(c Color) Color {
	return Color{
		R: +0.21260*c.R + 0.71520*c.G + 0.07220*c.B,
		G: -0.09991*c.R - 0.33609*c.G + 0.43600*c.B,
		B: +0.61500*c.R - 0.55861*c.G - 0.05639*c.B,
	}
}

func applyTemperatureTintEffect(c Color, temperature, tint float64) Color {
	const scale = 0.10
	cc := rgb2yuv(c)
	cc = Color{
		R: (cc.R),
		G: (cc.G - temperature*scale + tint*scale),
		B: (cc.B + temperature*scale + tint*scale),
	}
	cc = yuv2rgb(cc)
	return Color{
		R: clamp(cc.R),
		G: clamp(cc.G),
		B: clamp(cc.B),
	}
}

func applyBrightnessEffect(c Color, brightness float64) Color {
	const scale = 1.5
	return Color{
		R: math.Pow(c.R, 1.0/(1.0+scale*brightness)),
		G: math.Pow(c.G, 1.0/(1.0+scale*brightness)),
		B: math.Pow(c.B, 1.0/(1.0+scale*brightness)),
	}
}

func applyContrastEffect(c Color, contrast float64) Color {
	const pi4 = 3.14159265358979 * 0.25
	contrastCoef := math.Tan(contrast+1) * pi4
	return Color{
		R: fromsRGB2Linear(math.Max(0, (fromLinear2sRGB(c.R)-0.5)*contrastCoef+0.5)),
		G: fromsRGB2Linear(math.Max(0, (fromLinear2sRGB(c.G)-0.5)*contrastCoef+0.5)),
		B: fromsRGB2Linear(math.Max(0, (fromLinear2sRGB(c.B)-0.5)*contrastCoef+0.5)),
	}
}

func rgb2h(c Color) float64 {
	r := c.R
	g := c.G
	b := c.B
	M := math.Max(math.Max(r, g), b)
	m := math.Min(math.Min(r, g), b)

	h := 0.0
	if M == m {
		h = 0.0
	} else if m == b {
		h = 60.0*(g-r)/(M-m) + 60.0
	} else if m == r {
		h = 60.0*(b-g)/(M-m) + 180.0
	} else if m == g {
		h = 60.0*(r-b)/(M-m) + 300.0
	}

	h /= 360.0
	if h < 0.0 {
		h += 1
	} else if h > 1.0 {
		h -= 1
	}
	return h
}

func rgb2s4hsv(c Color) float64 {
	r := c.R
	g := c.G
	b := c.B
	M := math.Max(math.Max(r, g), b)
	m := math.Min(math.Min(r, g), b)

	if M < 1e-14 {
		return 0.0
	}
	return (M - m) / M
}
func rgb2hsv(c Color) Color {
	r := c.R
	g := c.G
	b := c.B

	v := math.Max(math.Max(r, g), b)
	h := rgb2h(c)
	s := rgb2s4hsv(c)
	return Color{R: h, G: s, B: v}
}
func hsv2rgb(c Color) Color {
	h := c.R
	s := c.G
	v := c.B

	if s < 1e-14 {
		return Color{R: v, G: v, B: v}
	}

	h6 := h * 6.0
	i := int(math.Floor(h6)) % 6
	f := h6 - float64(i)
	p := v * (1 - s)
	q := v * (1 - (s * f))
	t := v * (1 - (s * (1 - f)))

	r, g, b := 0., 0., 0.
	switch i {
	case 0:
		r = v
		g = t
		b = p
	case 1:
		r = q
		g = v
		b = p
	case 2:
		r = p
		g = v
		b = t
	case 3:
		r = p
		g = q
		b = v
	case 4:
		r = t
		g = p
		b = v
	case 5:
		r = v
		g = p
		b = q
	}
	return Color{R: r, G: g, B: b}
}

func applySaturationEffect(c Color, saturation float64) Color {
	hsv := rgb2hsv(Color{
		R: clamp(c.R),
		G: clamp(c.G),
		B: clamp(c.B),
	})
	hsv.G = hsv.G * (saturation + 1)
	return hsv2rgb(hsv)
}
