package cokacoin

import (
	"crypto/sha256"
	"fmt"
)

func GenerateMask(zeros int) []byte {
	full, half := zeros/2, zeros%2
	var mask []byte
	for i := 0; i < full; i++ {
		mask = append(mask, 0)
	}

	if half > 0 {
		mask = append(mask, 0xf)
	}
	return mask
}

func EasyHash(data ...interface{}) []byte {
	/*
		var f []byte
		for i := range data {
			f = append(f, data[i]...)
		}
		h := sha256.Sum256(f)
	*/

	h := sha256.New()

	// h has a Write method, so is a io.Writer
	fmt.Fprint(h, data...)
	//h.Write(data[i])

	return h.Sum(nil)

}

func GoodEnough(mask []byte, hash []byte) bool {
	for i := range mask {
		if hash[i] > mask[i] {
			return false
		}

	}

	//if hash[0] == 0 {
	//	return true
	//}
	return true
}

func DifficultHash(mask []byte, data ...interface{}) ([]byte, int32) {
	ln := len(data)
	data = append(data, nil)
	var i int32

	for {
		data[ln] = i
		hash := EasyHash(data...)
		if GoodEnough(mask, hash) {
			return hash, i
		}
		i++
	}

}
