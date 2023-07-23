package imagehash

import (
	"encoding/hex"
	"errors"
	"image"
	"math/bits"
)

var (
	ErrBadHashSize   = errors.New("incompatible hashsize")
	ErrUnequalHashes = errors.New("hashes have unequal sizes")
	ErrNilImage      = errors.New("image is nil")
)

type Imagehash struct {
	hash []byte
}

func (i *Imagehash) String() string {
	return hex.EncodeToString(i.hash)
}

func (i *Imagehash) FromString(hashstr string) error {
	buf, err := hex.DecodeString(hashstr)
	if err != nil {
		return err
	}
	if i != nil {
		i.hash = buf
	} else {
		i = &Imagehash{hash: buf}
	}
	return nil
}

func (i *Imagehash) Distance(other Imagehash) (int, error) {
	hamming := 0
	if len(i.hash) != len(other.hash) {
		return 0, ErrUnequalHashes
	}
	for idx := 0; idx < len(i.hash); idx++ {
		hamming += bits.OnesCount8(i.hash[idx] ^ other.hash[idx])
	}
	return hamming, nil
}

func (i *Imagehash) Whash(image image.Image, hashsize uint) error {
	if image == nil {
		return ErrNilImage
	}
	data, err := grayscale(image)
	if err != nil {
		return err
	}
	if hashsize == 0 || hashsize > uint(len(data)) {
		return ErrBadHashSize
	}
	hashsize = floorp2(int(hashsize))
	i.hash = make([]byte, hashsize*hashsize/8)
	level := bits.Len(uint(len(data))) - 1
	hashlevel := bits.Len(hashsize) - 1
	DWT2d(data, level)
	data[0][0] = 0.0
	IDWT2d(data, level)
	DWT2d(data, level-hashlevel)
	excerpt := extractSquareRegion(data, hashsize)
	med := median(flatten(excerpt))
	i.pack(excerpt, med, hashsize)
	return nil
}

func (i *Imagehash) pack(excerpt [][]float64, med float64, hashsize uint) {
	var (
		acc byte
		k   uint
		j   uint
	)
	ctr := 0
	offset := 0
	for k = 0; k < hashsize; k++ {
		for j = 0; j < hashsize; j++ {
			if excerpt[k][j] > med {
				acc ^= 1
			}
			ctr++
			if ctr%8 == 0 {
				i.hash[offset] = acc
				offset++
				acc = 0
			}
			acc <<= 1
		}
	}
}
