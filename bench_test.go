package columnar

import (
	"encoding/binary"
	"testing"
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
