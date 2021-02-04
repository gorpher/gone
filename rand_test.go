package gone

import (
	"math"
	"testing"
)

func TestRandInt(t *testing.T) {
	var i = 1
	for i < math.MaxInt16 {
		randInt := RandInt(0, i)
		if randInt > i {
			t.Errorf("rand int error randInt(0,%d)=%d", i, randInt)
		}
		randInt32 := RandInt32(0, int32(i))
		if randInt32 > int32(i) {
			t.Errorf("rand int error randInt(0,%d)=%d", i, randInt32)
		}
		randInt64 := RandInt64(0, int64(i))
		if randInt64 > int64(i) {
			t.Errorf("rand int error randInt(0,%d)=%d", i, randInt64)
		}
		size := 1 >> 10
		ints := RandInts(0, i, size)
		if len(ints) != size {
			t.Errorf("rand ints error RandInts(0,%d,%d)=%d", i, size, ints)
		}
		i++
	}
}

func TestRandLowerAndUpper(t *testing.T) {
	lower := RandLower(math.MaxInt16)
	if len(lower) != math.MaxInt16 {
		t.Error("RandLower() has error")
	}
	upper := RandUpper(math.MaxInt16)
	if len(upper) != math.MaxInt16 {
		t.Error("RandUpper() has error")
	}
}

func BenchmarkRandString(b *testing.B) {
	//BenchmarkRandString-4   	    5454	    189408 ns/op
	var s string
	for i := 0; i < b.N; i++ {
		s = RandString(math.MaxInt16)
		if len(s) != math.MaxInt16 {
			b.Error("RandString() has error")
		}
	}
}

func BenchmarkRandAlphaString(b *testing.B) {
	//BenchmarkRandAlphaString-4   	    1153	    952326 ns/op
	var s string
	for i := 0; i < b.N; i++ {
		s = RandAlphaString(math.MaxInt16)
		if len(s) != math.MaxInt16 {
			b.Error("RandAlphaString() has error")
		}
	}
}

func BenchmarkRandBytes(b *testing.B) {
	var s []byte
	for i := 0; i < b.N; i++ {
		s = RandBytes(math.MaxInt16)
		if len(s) != math.MaxInt16 {
			b.Error("RandBytes() has error")
		}
	}
}
