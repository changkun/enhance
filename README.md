# enhance

Package `enhance` is a Go implementation for enhancing photographs (adjusting brightness, contrast, etc.)


<img width="100px" src="./testdata/1/brightness-0.7.jpg"><img width="100px" src="./testdata/1/contrast-0.7.jpg"><img width="100px" src="./testdata/1/saturation-0.7.jpg"><img width="100px" src="./testdata/1/temperature-0.7.jpg"><img width="100px" src="./testdata/1/tint-0.7.jpg">

<img width="100px" src="./testdata/2/brightness-0.7.jpg"><img width="100px" src="./testdata/2/contrast-0.7.jpg"><img width="100px" src="./testdata/2/saturation-0.7.jpg"><img width="100px" src="./testdata/2/temperature-0.7.jpg"><img width="100px" src="./testdata/2/tint-0.7.jpg">

<img width="100px" src="./testdata/3/brightness-0.7.jpg"><img width="100px" src="./testdata/3/contrast-0.7.jpg"><img width="100px" src="./testdata/3/saturation-0.7.jpg"><img width="100px" src="./testdata/3/temperature-0.7.jpg"><img width="100px" src="./testdata/3/tint-0.7.jpg">

This project aims to reproduce https://github.com/yuki-koyama/enhancer.

## Usage

```go
package main

import (
	"image"
	"image/jpeg"
	"log"
	"os"

	"changkun.de/x/enhance"
)

func main() {
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
}
```

## License

[MIT](./LICENSE)