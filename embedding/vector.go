package embedding

import (
	"errors"
	"math"
)

var (
	errZeroDimension     = errors.New("vector has zero dimension")
	errZeroVector        = errors.New("zero vector")
	errDimensionMismatch = errors.New("dimension mismatch")
)

func VectorMagnitude(v []float32) (float32, error) {
	if len(v) == 0 {
		return 0, errZeroDimension
	}
	var sum float32
	for _, x := range v {
		sum += x * x
	}
	return sum, nil
}

func VectorNormalize(v []float32, result []float32) error {
	if len(v) == 0 {
		return errZeroDimension
	}
	if len(result) != len(v) {
		return errDimensionMismatch
	}
	var err error
	var magnitude float32

	magnitude, err = VectorMagnitude(v)
	if err != nil {
		return err
	}
	if magnitude == 0 {
		return errZeroVector
	}

	var inverse float32 = 1.0 / magnitude

	// BCE hint, see https://go101.org/optimizations/5-bce.html
	result = result[:len(v)]
	for i := 0; i < len(v); i++ {
		result[i] = v[i] * inverse
	}
	return nil
}

func VectorAdd(v1, v2 []float32, result []float32) error {
	if len(v1) != len(v2) {
		return errDimensionMismatch
	}
	if len(result) != len(v1) {
		return errDimensionMismatch
	}

	// BCE hint, see https://go101.org/optimizations/5-bce.html
	result = result[:len(v1)]
	v2 = v2[:len(v1)]
	for i := 0; i < len(v1); i++ {
		result[i] = v1[i] + v2[i]
	}
	return nil
}

func VectorSubtract(v1, v2 []float32, result []float32) error {
	if len(v1) != len(v2) {
		return errDimensionMismatch
	}
	if len(result) != len(v1) {
		return errDimensionMismatch
	}
	// BCE hint, see https://go101.org/optimizations/5-bce.html
	result = result[:len(v1)]
	v2 = v2[:len(v1)]

	for i := range v1 {
		result[i] = v1[i] - v2[i]
	}
	return nil
}

func VectorDotProduct(v1, v2 []float32) (float32, error) {
	if len(v1) != len(v2) {
		return 0, errDimensionMismatch
	}
	var sum float32

	// BCE hint, see https://go101.org/optimizations/5-bce.html
	v2 = v2[:len(v1)]
	for i := 0; i < len(v1); i++ {
		sum += v1[i] * v2[i]
	}
	return sum, nil
}

func VectorEuclideanDistance(v1, v2 []float32) (float32, error) {
	if len(v1) != len(v2) {
		return 0, errDimensionMismatch
	}
	var sum float32

	// BCE hint, see https://go101.org/optimizations/5-bce.html
	v2 = v2[:len(v1)]
	for i := 0; i < len(v1); i++ {
		d := v1[i] - v2[i]
		sum += d * d
	}
	return float32(math.Sqrt(float64(sum))), nil
}

func VectorCosineSimilarity(v1, v2 []float32) (float32, error) {
	var x, y float32
	var dot float32
	var norm1, norm2 float32

	if len(v1) != len(v2) {
		return 0, errDimensionMismatch
	}

	// BCE hint, see https://go101.org/optimizations/5-bce.html
	v2 = v2[:len(v1)]
	for i := 0; i < len(v1); i++ {
		x = v1[i]
		y = v2[i]
		dot += x * y
		norm1 += x * x
		norm2 += y * y
	}

	if norm1 == 0 || norm2 == 0 {
		return 0, errZeroVector
	}

	var cosine = dot / (float32(math.Sqrt(float64(norm1))) * float32(math.Sqrt(float64(norm2))))
	return cosine, nil
}
