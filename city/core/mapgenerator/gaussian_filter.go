package mapgenerator

import (
	"math"
)

var (
	kernal    [][]float64
	kernalNum = 3
)

func GaussianFunc(x, y int, o float64) float64 {
	return (1.0 / (2.0 * math.Pi * o * o)) * math.Pow(math.E, ((-1.0)*(float64(x*x+y*y)/(2.0*o*o))))
}

func GaussianAvgArray(n int, o float64) [][]float64 {
	sum := 0.0
	arr := make([][]float64, 2*n+1)
	for i := 0; i < (2*n + 1); i++ {
		arr2 := make([]float64, 2*n+1)
		for j := 0; j < (2*n + 1); j++ {
			result := GaussianFunc(i-n, j-n, o)
			arr2[j] = result
			sum += result
		}
		arr[i] = arr2
	}

	for i := 0; i < (2*n + 1); i++ {
		thisArr := arr[i]
		for j := 0; j < (2*n + 1); j++ {
			thisArr[j] = thisArr[j] / sum
		}
		arr[i] = thisArr
	}

	return arr
}

func GaussianFilter(height_map *HeightMap) {
	filtered_data := make([]int, height_map.Width*height_map.Height)

	for x := 0; x < height_map.Width; x++ {
		for y := 0; y < height_map.Height; y++ {
			filtered_data[x+y*height_map.Width] = computeFilteredValue(height_map, x, y)
		}
	}

	height_map.data = filtered_data
}

func computeFilteredValue(height_map *HeightMap, x, y int) int {
	sum := 0.0
	for n := ((-1) * kernalNum); n <= kernalNum; n++ {
		for m := ((-1) * kernalNum); m <= kernalNum; m++ {
			px := checkBound(x+n, height_map.Width)
			py := checkBound(y+m, height_map.Height)
			sum += kernal[n+kernalNum][m+kernalNum] * float64(height_map.Value(px, py))
		}
	}

	return int(sum)
}

func checkBound(n, bound int) int {
	if n < 0 {
		return 0
	}
	if n >= bound {
		return bound - 1
	}
	return n
}

func init() {
	kernal = GaussianAvgArray(kernalNum, 20)
}
