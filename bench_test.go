package columnar

import (
	"crypto/rand"
	"encoding/binary"
	"strconv"
	"testing"
	"unsafe"
)

func BenchmarkUInt64(b *testing.B) {
	b.Run("Encode", func(b *testing.B) {
		buf := make([]byte, 8*1024)
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			for j := 0; j < 1024; j++ {
				binary.LittleEndian.PutUint64(buf[j*8:j*8+8], uint64(j))
			}
		}
		b.SetBytes(8 * 1024)
	})

	b.Run("Decode", func(b *testing.B) {
		buf := make([]byte, 8*1024)
		for j := 0; j < 1024; j++ {
			binary.LittleEndian.PutUint64(buf[j*8:j*8+8], uint64(j))
		}
		b.ResetTimer()

		var x uint64
		for i := 0; i < b.N; i++ {
			for j := 0; j < 1024; j++ {
				x += binary.LittleEndian.Uint64(buf[j*8 : j*8+8])
			}
		}
		b.SetBytes(8 * 1024)
	})

	b.Run("Both", func(b *testing.B) {
		buf := make([]byte, 8*1024)
		var x uint64
		b.ResetTimer()

		for i := 0; i < b.N; i++ {
			for j := 0; j < 1024; j++ {
				binary.LittleEndian.PutUint64(buf[j*8:j*8+8], uint64(j))
			}
			for j := 0; j < 1024; j++ {
				x += binary.LittleEndian.Uint64(buf[j*8 : j*8+8])
			}
		}
		b.SetBytes(8 * 1024)
	})
}

func BenchmarkLocality(b *testing.B) {
	const max = 20
	buf := make([]byte, 8*2<<max)
	_, err := rand.Read(buf)
	if err != nil {
		b.Fatal(err)
	}
	for words := 1; words < 2<<max; words = words << 2 {
		b.Run(strconv.Itoa(words), func(b *testing.B) {
			b.Run("local", func(b *testing.B) {
				var x uint64
				for i := 0; i < b.N; i++ {
					for j := 0; j < words; j++ {
						x += binary.LittleEndian.Uint64(buf[0:8])
					}
				}
				b.SetBytes(int64(8 * words))
			})
			b.Run("non-local", func(b *testing.B) {
				var x uint64
				for i := 0; i < b.N; i++ {
					for j := 0; j < words; j++ {
						x += binary.LittleEndian.Uint64(buf[j*8 : (j+1)*8])
					}
				}
				b.SetBytes(int64(8 * words))
			})
		})
	}
}

func BenchmarkAlignment(b *testing.B) {
	var buf [15]byte
	_, err := rand.Read(buf[0:])
	if err != nil {
		b.Fatal(err)
	}
	b.Run("unsafe", func(b *testing.B) {
		for offset := 0; offset < 8; offset++ {
			b.Run(strconv.Itoa(offset), func(b *testing.B) {
				var x uint64
				for i := 0; i < b.N; i++ {
					a := unsafe.Pointer(&buf[offset])
					b := *(*uint64)(a)
					x += b
				}
				b.SetBytes(8)
			})
		}
	})
	b.Run("encoding", func(b *testing.B) {
		for offset := 0; offset < 8; offset++ {
			b.Run(strconv.Itoa(offset), func(b *testing.B) {
				var x uint64
				for i := 0; i < b.N; i++ {
					a := binary.LittleEndian.Uint64(buf[offset:])
					x += a
				}
				b.SetBytes(8)
			})
		}
	})
}
