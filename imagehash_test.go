package imagehash

import (
	"image"
	"image/jpeg"
	"image/png"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestImagehash_FromString(t *testing.T) {
	type fields struct {
		hash []byte
	}
	tests := []struct {
		name    string
		fields  fields
		hashstr string
		wantErr bool
	}{
		{
			name:    "error hash string",
			fields:  fields{},
			hashstr: "abcdefg",
			wantErr: true,
		},
		{
			name:    "11223344 hash",
			fields:  fields{[]byte{0x11, 0x22, 0x33, 0x44}},
			hashstr: "11223344",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var i Imagehash
			j := Imagehash{hash: tt.fields.hash}
			if err := i.FromString(tt.hashstr); (err != nil) != tt.wantErr {
				t.Fatalf("FromString() error = %v, wantErr %v", err, tt.wantErr)
			}
			assert.Equal(t, i.hash, j.hash)
		})
	}
}

func TestImagehash_String(t *testing.T) {
	type fields struct {
		hash []byte
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name:   "short string",
			fields: fields{[]byte{0x11, 0x22, 0x33, 0x44}},
			want:   "11223344",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			i := Imagehash{
				hash: tt.fields.hash,
			}
			assert.Equal(t, tt.want, i.String())
		})
	}
}

func TestImagehash_Whash(t *testing.T) {
	tests := []struct {
		name     string
		file     string
		level    uint
		distance int
		whash    string
		wantErr  assert.ErrorAssertionFunc
	}{
		{
			name:     "lenna.png_16×16 bits",
			file:     "testdata/lenna.png",
			level:    16,
			distance: 0,
			whash:    "cfbccfbc43f847e947fb5e7348e341e7414741c741cf40cf40ca40fe40f441f0",
			wantErr:  assert.NoError,
		},
		{
			name:     "lenna.png_8×8 bits",
			file:     "testdata/lenna.png",
			level:    8,
			distance: 0,
			whash:    "be98bd890b0b8f8c",
			wantErr:  assert.NoError,
		},
		{
			name:     "gopher.png_16×16 bits",
			file:     "testdata/gopher.png",
			level:    16,
			distance: 0,
			whash:    "01800fa01ff03cf03ff83ffc1ffc1ffc0ffc07fc07fc07fc07fe07f003800200",
			wantErr:  assert.NoError,
		},
		{
			name:     "gopher.png_8×8 bits",
			file:     "testdata/gopher.png",
			level:    8,
			distance: 0,
			whash:    "187c7e3e3e3e1c10",
			wantErr:  assert.NoError,
		},
		{
			name:     "gopher.png_512*512 bits",
			file:     "testdata/gopher.png",
			level:    512,
			distance: 0,
			whash:    "",
			wantErr:  assert.Error,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			i := Imagehash{}
			j := Imagehash{}
			err := j.FromString(tt.whash)
			assert.NoError(t, err)
			var img image.Image
			func() {
				file, err := os.Open(tt.file)
				assert.NoError(t, err)
				img, err = png.Decode(file)
				assert.NoError(t, err)
				defer file.Close()
			}()
			err = i.Whash(img, tt.level)
			if !tt.wantErr(t, err) {
				t.Fatalf("Whash() error = %v", err)
			}
			d, err := i.Distance(j)
			assert.NoError(t, err)
			assert.LessOrEqual(t, d, tt.distance)
		})
	}
}

func TestImagehash_Distance(t *testing.T) {
	type fields struct {
		hash []byte
	}
	type args struct {
		other Imagehash
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    int
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name:    "equal distance=0",
			fields:  fields{[]byte{0x11, 0x22, 0x33, 0x44}},
			args:    args{other: Imagehash{[]byte{0x11, 0x22, 0x33, 0x44}}},
			want:    0,
			wantErr: assert.NoError,
		},
		{
			name:    "not equal distance=6",
			fields:  fields{[]byte{0x11, 0x22, 0x33, 0x44}},
			args:    args{other: Imagehash{[]byte{0xFF, 0x22, 0x33, 0x44}}},
			want:    6,
			wantErr: assert.NoError,
		},
		{
			name:    "incompatible hashes",
			fields:  fields{[]byte{0x11, 0x22, 0x33, 0x44}},
			args:    args{other: Imagehash{[]byte{0x11, 0x22, 0x33, 0x44, 0x55}}},
			want:    0,
			wantErr: assert.Error,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			i := &Imagehash{
				hash: tt.fields.hash,
			}
			got, err := i.Distance(tt.args.other)
			if !tt.wantErr(t, err) {
				t.Fatalf("Distance() error = %v", err)
			}
			assert.Equal(t, tt.want, got)
		})
	}
}

func BenchmarkImagehash_Whash(b *testing.B) {
	i := Imagehash{}
	var img image.Image
	file, err := os.Open("testdata/lenna.png")
	assert.NoError(b, err)
	img, err = png.Decode(file)
	assert.NoError(b, err)
	file.Close()

	// Exclude the loading of the image from the test
	b.ResetTimer()
	// run the Fib function b.N times
	for n := 0; n < b.N; n++ {
		if hashErr := i.Whash(img, 8); hashErr != nil {
			b.Fatalf("Failed hashing %v", hashErr)
		}
	}
}

func BenchmarkImagehash_Whash6000x2000(b *testing.B) {
	i := Imagehash{}
	var img image.Image
	file, err := os.Open("testdata/6000x2000.jpg")
	assert.NoError(b, err)
	img, err = jpeg.Decode(file)
	assert.NoError(b, err)
	file.Close()

	// Exclude the loading of the image from the test
	b.ResetTimer()
	// run the Fib function b.N times
	for n := 0; n < b.N; n++ {
		if hashErr := i.Whash(img, 8); hashErr != nil {
			b.Fatalf("Failed hashing %v", hashErr)
		}
	}
}

func BenchmarkImagehash_Whash20MB(b *testing.B) {
	i := Imagehash{}
	var img image.Image
	file, err := os.Open("20MB.jpg")
	assert.NoError(b, err)
	img, err = jpeg.Decode(file)
	assert.NoError(b, err)
	file.Close()

	// Exclude the loading of the image from the test
	b.ResetTimer()
	// run the Fib function b.N times
	for n := 0; n < b.N; n++ {
		if hashErr := i.Whash(img, 8); hashErr != nil {
			b.Fatalf("Failed hashing %v", hashErr)
		}
	}
}
