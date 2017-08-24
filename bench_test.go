package columnar

import (
	"testing"

	capnproto "zombiezen.com/go/capnproto2"
)

func BenchmarkUInt64(b *testing.B) {
	setup := func(b *testing.B, size int) capnproto.UInt64List {
		msg := &capnproto.Message{Arena: capnproto.SingleSegment(nil)}
		seg, err := msg.Segment(0)
		if err != nil {
			b.Fatal(err)
		}

		l, err := capnproto.NewUInt64List(seg, int32(size))
		if err != nil {
			b.Fatal(err)
		}
		return l
	}

	b.Run("Encode", func(b *testing.B) {
		l := setup(b, 1024)
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			for j := 0; j < 1024; j++ {
				l.Set(j, uint64(j))
			}
		}
		b.SetBytes(8 * 1024)
	})

	b.Run("Decode", func(b *testing.B) {
		l := setup(b, 1024)
			for j := 0; j < 1024; j++ {
				l.Set(j, uint64(j))
			}
		b.ResetTimer()

		var x uint64
		for i := 0; i < b.N; i++ {
			for j := 0; j < 1024; j++ {
				x += l.At(j)
			}
		}
		b.SetBytes(8 * 1024)
	})
}