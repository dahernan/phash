// This program is free software: you can redistribute it and/or modify it
// under the terms of the GNU General Public License as published by the Free
// Software Foundation, either version 3 of the License, or (at your option)
// any later version.
//
// This program is distributed in the hope that it will be useful, but
// WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the GNU General
// Public License for more details.
//
// You should have received a copy of the GNU General Public License along
// with this program.  If not, see <http://www.gnu.org/licenses/>.

// Package phash is a simple wrapper around the pHash library.
package phash

/*
#cgo pkg-config: pHash

#include <stdlib.h>

void cimg_exception_mode_quiet();

typedef unsigned long long ulong64;

ulong64 ph_dct_imagehash_wrapper(const char* file);
int ph_hamming_distance(const ulong64 hash1,const ulong64 hash2);

ulong64* ph_dct_videohash_wrapper(const char *filename, int *l);

*/
import "C"
import (
	"reflect"
	"unsafe"
)

func init() {
	C.cimg_exception_mode_quiet()
}

// ImageHashDCT returns the perceptual hash of the image file fn.
func ImageHashDCT(fn string) (uint64, error) {
	cfn := C.CString(fn)
	defer C.free(unsafe.Pointer(cfn))
	hash, err := C.ph_dct_imagehash_wrapper(cfn)
	return uint64(hash), err
}

func VideoHash(fn string) ([]C.ulong64, error) {
	cfn := C.CString(fn)
	var l int

	defer C.free(unsafe.Pointer(cfn))

	array, err := C.ph_dct_videohash_wrapper(cfn, (*C.int)(unsafe.Pointer(&l)))
	if err != nil {
		return nil, err
	}

	length := l
	hdr := reflect.SliceHeader{
		Data: uintptr(unsafe.Pointer(array)),
		Len:  length,
		Cap:  length,
	}
	goSlice := *(*[]C.ulong64)(unsafe.Pointer(&hdr))

	return goSlice, err
}

// HammingDistance returns the Hamming distance between the two perceptual hashes.
func HammingDistance(h1, h2 uint64) int {
	const (
		m1  = 0x5555555555555555
		m2  = 0x3333333333333333
		h01 = 0x0101010101010101
		m4  = 0x0f0f0f0f0f0f0f0f
	)
	x := h1 ^ h2
	x -= (x >> 1) & m1
	x = (x & m2) + ((x >> 2) & m2)
	x = (x + (x >> 4)) & m4
	return int((x * h01) >> 56)
}

// for benchmarks
func hammingDistanceC(h1, h2 uint64) int {
	return int(C.ph_hamming_distance(C.ulong64(h1), C.ulong64(h2)))
}
