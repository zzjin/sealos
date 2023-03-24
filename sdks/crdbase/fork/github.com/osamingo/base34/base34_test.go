package base34

import (
	"math"
	"math/rand"
	"testing"
	"time"
)

func TestMustNewEncoder(t *testing.T) {
	t.Parallel()

	enc := MustNewEncoder("rpshnaf39w472bcdeg65jkm8oqi1tuvxyz")
	if enc == nil {
		t.Error("should not be nil")
	}

	func() {
		defer func() {
			if err := recover(); err == nil {
				t.Error("should not be nil")
			}
		}()
		MustNewEncoder("")
		t.Error("should be panic")
	}()

	func() {
		defer func() {
			if err := recover(); err == nil {
				t.Error("should not be nil")
			}
		}()
		MustNewEncoder("test")
		t.Error("should be panic")
	}()
}

func TestNewEncoder(t *testing.T) {
	t.Parallel()

	enc, err := NewEncoder("rpshnaf39w472bcdeg65jkm8oqi1tuvxyz")
	if err != nil {
		t.Error("should be nil")
	}
	if enc == nil {
		t.Error("should not be nil")
	}

	_, err = NewEncoder("")
	if err == nil {
		t.Error("should not be nil")
	}

	_, err = NewEncoder("test")
	if err == nil {
		t.Error("should not be nil")
	}
}

func TestEncoder_Encode(t *testing.T) {
	t.Parallel()

	bc := map[uint64]string{
		0:              "1",
		57:             "2p",
		math.MaxUint8:  "8i",
		math.MaxUint16: "2opi",
		math.MaxUint32: "3sizkni",
		math.MaxUint64: "8qtr74ui5erii",
	}

	enc := MustNewEncoder(StandardSource)
	if id := enc.Encode(0); id != "1" {
		t.Error("should be", "1")
	}

	for k, v := range bc {
		if o := enc.Encode(k); o != v {

			t.Error("should be", k, o, v)
		}
	}
}

func TestEncoder_Decode(t *testing.T) {
	t.Parallel()

	bc := map[uint64]string{
		0:              "1",
		57:             "2p",
		math.MaxUint8:  "8i",
		math.MaxUint16: "2opi",
		math.MaxUint32: "3sizkni",
		math.MaxUint64: "8qtr74ui5erii",
	}

	enc := MustNewEncoder(StandardSource)
	if _, err := enc.Decode(""); err == nil {
		t.Error("should not be nil")
	}

	if _, err := enc.Decode("0"); err == nil {
		t.Error("should not be nil")
	}

	for k, v := range bc {
		r, err := enc.Decode(v)
		if err != nil {
			t.Error("should be nil")
		}
		if r != k {
			t.Error("should be", k)
		}
	}
}

func BenchmarkEncoder_Encode(b *testing.B) {
	s := rand.New(rand.NewSource(time.Now().UnixNano()))
	enc := MustNewEncoder(StandardSource)

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		enc.Encode(uint64(s.Int63()))
	}
}

func BenchmarkEncoder_Decode(b *testing.B) {
	bc := map[uint64]string{
		0:              "1",
		57:             "2p",
		math.MaxUint8:  "8i",
		math.MaxUint16: "2opi",
		math.MaxUint32: "3sizkni",
		math.MaxUint64: "8qtr74ui5erii",
	}

	l := len(bc)
	s := rand.New(rand.NewSource(time.Now().UnixNano()))
	enc := MustNewEncoder(StandardSource)

	vs := make([]string, 0, l)
	for k := range bc {
		vs = append(vs, bc[k])
	}

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := enc.Decode(vs[s.Intn(l)])
		if err != nil {
			b.Fatal(err)
		}
	}
}
