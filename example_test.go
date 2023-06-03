package enhance_test

import (
	"image/jpeg"
	"log"
	"os"

	"changkun.de/x/enhance"
)

func ExampleImage() {
	f, err := os.Open("testdata/1.jpg")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	img, err := jpeg.Decode(f)
	if err != nil {
		log.Fatal(err)
	}

	enhanced := enhance.Image(img, enhance.Params{
		Brightness:  .5,
		Contrast:    .5,
		Saturation:  .5,
		Temperature: .5,
		Tint:        .5,
	})
	f, err = os.Create("testdata/enhanced.jpg")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	if err := jpeg.Encode(f, enhanced, nil); err != nil {
		log.Fatal(err)
	}

	// Output:
	//
}
