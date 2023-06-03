package enhance_test

import (
	"fmt"
	"image/jpeg"
	"os"
	"testing"

	"changkun.de/x/enhance"
)

func TestImage(t *testing.T) {
	for j := 1; j <= 3; j++ {
		t.Run(fmt.Sprintf("%d", j), func(t *testing.T) {
			f, err := os.Open(fmt.Sprintf("testdata/%d.jpg", j))
			if err != nil {
				t.Fatal(err)
			}
			defer f.Close()
			img, err := jpeg.Decode(f)
			if err != nil {
				t.Fatal(err)
			}

			for i := 0.0; i < 1; i += 0.1 {
				p := enhance.Params{
					Brightness:  i,
					Contrast:    0.5,
					Saturation:  0.5,
					Temperature: 0.5,
					Tint:        0.5,
				}
				x := enhance.Image(img, p)
				ff, err := os.Create(fmt.Sprintf("testdata/%d/brightness-%.1f.jpg", j, p.Brightness))
				if err != nil {
					t.Fatal(err)
				}
				err = jpeg.Encode(ff, x, nil)
				if err != nil {
					t.Fatal(err)
				}
				ff.Close()

				p = enhance.Params{
					Brightness:  0.5,
					Contrast:    i,
					Saturation:  0.5,
					Temperature: 0.5,
					Tint:        0.5,
				}
				x = enhance.Image(img, p)
				ff, err = os.Create(fmt.Sprintf("testdata/%d/contrast-%.1f.jpg", j, p.Contrast))
				if err != nil {
					t.Fatal(err)
				}
				err = jpeg.Encode(ff, x, nil)
				if err != nil {
					t.Fatal(err)
				}
				ff.Close()

				p = enhance.Params{
					Brightness:  0.5,
					Contrast:    0.5,
					Saturation:  i,
					Temperature: 0.5,
					Tint:        0.5,
				}
				x = enhance.Image(img, p)
				ff, err = os.Create(fmt.Sprintf("testdata/%d/saturation-%.1f.jpg", j, p.Saturation))
				if err != nil {
					t.Fatal(err)
				}
				err = jpeg.Encode(ff, x, nil)
				if err != nil {
					t.Fatal(err)
				}
				ff.Close()

				p = enhance.Params{
					Brightness:  0.5,
					Contrast:    0.5,
					Saturation:  0.5,
					Temperature: i,
					Tint:        0.5,
				}
				x = enhance.Image(img, p)
				ff, err = os.Create(fmt.Sprintf("testdata/%d/temperature-%.1f.jpg", j, p.Temperature))
				if err != nil {
					t.Fatal(err)
				}
				err = jpeg.Encode(ff, x, nil)
				if err != nil {
					t.Fatal(err)
				}
				ff.Close()

				p = enhance.Params{
					Brightness:  0.5,
					Contrast:    0.5,
					Saturation:  0.5,
					Temperature: 0.5,
					Tint:        i,
				}
				x = enhance.Image(img, p)
				ff, err = os.Create(fmt.Sprintf("testdata/%d/tint-%.1f.jpg", j, p.Tint))
				if err != nil {
					t.Fatal(err)
				}
				err = jpeg.Encode(ff, x, nil)
				if err != nil {
					t.Fatal(err)
				}
				ff.Close()
			}
		})
	}
}
