// Copyright 2023 Changkun Ou <changkun.de>. All rights reserved.
// Use of this source code is governed by a MIT license that
// can be found in the LICENSE file.

package enhance

import (
	"math"
	"sync"
)

// fromLinear2sRGB converts a given value from linear space to
// sRGB space.
func fromLinear2sRGB(v float64) float64 {
	if !useLut {
		return linear2sRGB(v)
	}
	if v <= 0 {
		return 0
	}
	if v == 1 {
		return 1
	}
	i := v * lutSize
	ifloor := int(i) & (lutSize - 1)
	v0 := float64(lin2sRGBLUfloat64[ifloor])
	v1 := float64(lin2sRGBLUfloat64[ifloor+1])
	i -= float64(ifloor)
	return v0*(1.0-i) + v1*i
}

// fromsRGB2Linear converts a given value from linear space to
// sRGB space.
func fromsRGB2Linear(v float64) float64 {
	if !useLut {
		return sRGB2linear(v)
	}
	if v <= 0 {
		return 0
	}
	if v >= 1 {
		return 1
	}

	i := v * lutSize
	ifloor := int(i) & (lutSize - 1)
	v0 := float64(sRGB2linLUfloat64[ifloor])
	v1 := float64(sRGB2linLUfloat64[ifloor+1])
	i -= float64(ifloor)
	return v0*(1.0-i) + v1*i
}

const (
	lutSize = 1024 // keep a power of 2
)

var (
	once              sync.Once
	useLut            = true
	lin2sRGBLUfloat64 [lutSize + 1]float64
	sRGB2linLUfloat64 [lutSize + 1]float64
)

func init() {
	// Initialize look up table.
	once.Do(func() {
		for i := 0; i < lutSize; i++ {
			lin2sRGBLUfloat64[i] = linear2sRGB(float64(i) / lutSize)
			sRGB2linLUfloat64[i] = sRGB2linear(float64(i) / lutSize)
		}
		lin2sRGBLUfloat64[lutSize] = lin2sRGBLUfloat64[lutSize-1]
		sRGB2linLUfloat64[lutSize] = sRGB2linLUfloat64[lutSize-1]
	})
}

func sRGB2linear(v float64) float64 {
	if v <= 0.04045 {
		v /= 12.92
	} else {
		v = math.Pow((v+0.055)/1.055, 2.4)
	}
	return v
}

func linear2sRGB(v float64) float64 {
	if v <= 0.0031308 {
		v *= 12.92
	} else {
		v = 1.055*math.Pow(v, 1.0/2.4) - 0.055
	}
	return v
}
