package imagehash

import (
	"sort"
)

// waveletCoefficients struct holds the coefficients used in wavelet transforms.
type waveletCoefficients struct {
	LowPass  []float64
	HighPass []float64
}

// Haar wavelet coefficients are specified here.
// https://wavelets.pybytes.com/wavelet/haar/
// ....┌──┐
// ....│  │
// . ──┘  │  ┌──
// .......│  │
// .......└──┘
var haar = waveletCoefficients{
	HighPass: []float64{0.5, -0.5},
	LowPass:  []float64{0.5, 0.5},
}

// DWT1d applies 1D Discrete Wavelet Transform on data using Haar wavelet coefficients.
func DWT1d(data []float64) {
	temp := make([]float64, len(data))
	half := len(data) / 2
	for i := 0; i < half; i++ {
		k := i * 2
		temp[i] = haar.LowPass[0]*data[k] + haar.LowPass[1]*data[k+1]
		temp[i+half] = haar.HighPass[0]*data[k] + haar.HighPass[1]*data[k+1]
	}
	copy(data, temp)
}

// DWT2d applies 2D Discrete Wavelet Transform on data at specified level.
func DWT2d(data [][]float64, level int) {
	dims := len(data)
	for k := 0; k < level; k++ {
		curlvl := 1 << k
		curdims := dims / curlvl
		row := make([]float64, curdims)
		for i := 0; i < curdims; i++ {
			copy(row, data[i])
			DWT1d(row)
			copy(data[i], row)
		}
		col := make([]float64, curdims)
		for j := 0; j < curdims; j++ {
			for i := 0; i < curdims; i++ {
				col[i] = data[i][j]
			}
			DWT1d(col)
			for i := 0; i < curdims; i++ {
				data[i][j] = col[i]
			}
		}
	}
}

// iDWT1d applies 1D Inverse Discrete Wavelet Transform on data using Haar wavelet coefficients.
func iDWT1d(data []float64) {
	temp := make([]float64, len(data))
	half := len(data) / 2
	for i := 0; i < half; i++ {
		k := i * 2
		temp[k] = (haar.LowPass[0]*data[i] + haar.HighPass[0]*data[i+half]) / haar.HighPass[0]
		temp[k+1] = (haar.LowPass[1]*data[i] + haar.HighPass[1]*data[i+half]) / haar.LowPass[0]
	}
	copy(data, temp)
}

// IDWT2d applies 2D Inverse Discrete Wavelet Transform on data at specified level.
func IDWT2d(data [][]float64, level int) {
	dims := len(data)
	for k := level - 1; k >= 0; k-- {
		curlvl := 1 << k
		curdims := dims / curlvl
		col := make([]float64, curdims)
		for j := 0; j < curdims; j++ {
			for i := 0; i < curdims; i++ {
				col[i] = data[i][j]
			}
			iDWT1d(col)
			for i := 0; i < curdims; i++ {
				data[i][j] = col[i]
			}
		}
		row := make([]float64, curdims)
		for i := 0; i < curdims; i++ {
			copy(row, data[i])
			iDWT1d(row)
			copy(data[i], row)
		}
	}
}

// floorp2 finds the largest power of 2 less than or equal to val.
func floorp2(val int) uint {
	val |= val >> 1
	val |= val >> 2
	val |= val >> 4
	val |= val >> 8
	val |= val >> 16
	return uint(val - (val >> 1))
}

// flatten transforms a 2D slice to a 1D slice.
func flatten(data [][]float64) []float64 {
	flat := make([]float64, len(data)*len(data))
	offset := 0
	for _, row := range data {
		copy(flat[offset:offset+len(row)], row)
		offset += len(row)
	}
	return flat
}

// median finds the median value in a 2D slice.
func median(data []float64) float64 {
	temp := make([]float64, len(data))
	copy(temp, data)
	sort.Float64s(temp)
	if len(temp)%2 == 1 {
		return data[len(temp)/2]
	} else {
		return 0.5 * (temp[len(temp)/2-1] + temp[len(temp)/2])
	}
}

// extractSquareRegion extracts a square region of the given width from the data.
func extractSquareRegion(data [][]float64, width uint) [][]float64 {
	excerpt := make([][]float64, width)
	for i := 0; i < int(width); i++ {
		excerpt[i] = make([]float64, width)
		copy(excerpt[i], data[i][:width])
	}
	return excerpt
}

// min returns the minimum of two comparable values.
func min[T int | uint | float64](x, y T) T {
	if x < y {
		return x
	}
	return y
}
