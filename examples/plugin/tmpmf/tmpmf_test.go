package tmpmf

import (
	"fmt"
	"testing"
)

const nBufs = 5

func Test_bufRing_new(t *testing.T) {
	br := newBufRing(nBufs)
	if len(br.Buffers) != nBufs {
		t.Error("wrong size")
		t.FailNow()
	}
}

func Test_bufRing_Put(t *testing.T) {
	br := newBufRing(nBufs)

	for i := 0; i < nBufs; i++ {
		bb := make([]byte, (i+1)*2)
		//fmt.Printf("%d, ", len(bb))
		br.Put(bb)
	}
	fmt.Println("\n---------------")
	for j, bb := range br.Buffers {
		//fmt.Printf("%d, ", len(bb))
		if len(bb) != (j+1)*2 {
			t.Error("wrong size of buffer ", j)
			t.FailNow()
		}
	}
	fmt.Println(" ")
}

const minSize = 2

func Test_bufRing_Line(t *testing.T) {
	br := newBufRing(nBufs)

	for i := 0; i < nBufs; i++ {
		n := byte((i + 1) * 2)
		bb := make([]byte, n)
		bb[minSize-1] = n
		if len(bb) > minSize {
			bb[minSize] = n + 1
		}
		fmt.Printf("%v, ", bb)
		br.Put(bb)
	}
	fmt.Println("\n---------------")
	line := br.Line(minSize - 1)
	fmt.Printf("%v, ", line)
	for j, b := range line {
		//fmt.Printf("%d, ", len(bb))
		n := byte((j + 1) * 2)
		if b != n {
			t.Error("wrong  buffer content ", j)
			t.FailNow()
		}
	}
	line = br.Line(minSize)
	fmt.Printf("%v, ", line)
	fmt.Println(" ")
}

func Test_bufRing_Content(t *testing.T) {
	br := newBufRing(nBufs)

	for i := 0; i < nBufs; i++ {
		n := byte((i + 1) * 2)
		bb := make([]byte, n)
		//fmt.Printf("%d, ", len(bb))
		bb[0] = n
		bb[1] = n + 1
		br.Put(bb)
	}
	//fmt.Println("\n---------------")
	for j, bb := range br.Buffers {
		//fmt.Printf("%d, ", len(bb))
		n := byte((j + 1) * 2)
		if bb[0] != n || bb[1] != n+1 {
			t.Error("wrong  buffer content ", j)
			t.FailNow()
		}
	}
	fmt.Println(" ")
}

func Test_bufRing_Sort(t *testing.T) {
	br := newBufRing(nBufs)

	for i := 0; i < nBufs; i++ {
		n := byte((i + 1) * 2)
		bb := make([]byte, n)
		bb[minSize-1] = n
		if len(bb) > minSize {
			bb[minSize] = n + 1
		}
		fmt.Printf("%v, ", bb)
		br.Put(bb)
	}
	fmt.Println("\n---------------")
	line := br.Line(minSize - 1)
	InsertionSortU8(line)
	fmt.Printf("%v, ", line)
}

func Test_bufRing_LineMean(t *testing.T) {
	br := newBufRing(nBufs)

	for i := 0; i < nBufs; i++ {
		n := byte((i + 1) * 2)
		bb := make([]byte, n)
		bb[minSize-1] = n
		if len(bb) > minSize {
			bb[minSize] = n + 1
		}
		fmt.Printf("%v, ", bb)
		br.Put(bb)
	}
	fmt.Println("\n---------------")
	mean := br.LineMean(minSize - 1)
	fmt.Printf("mean %d\n ", mean)
	//src := &br.Buffers[3]
	//fmt.Printf("len(br.Buffers[3]) %d - %d\n ", len(br.Buffers[3]), len(*src))
}

func Benchmark_bufRing_PartialSort5Int(b *testing.B) {
	for i := 0; i < b.N; i++ {
		b := [5]int{3, 2, 4, 1, 5}
		b = PartialSort5Int(b)
	}
}

func Benchmark_bufRing_PartialSort5Byte(b *testing.B) {
	for i := 0; i < b.N; i++ {
		b := [5]byte{3, 2, 4, 1, 5}
		b = PartialSort5Byte(b)
	}
}

func Benchmark_bufRing_LineMean(b *testing.B) {
	b.StopTimer()
	br := newBufRing(nBufs)
	for i := 0; i < nBufs; i++ {
		n := byte((i + 1) * 2)
		bb := make([]byte, n)
		bb[minSize-1] = n
		if len(bb) > minSize {
			bb[minSize] = n + 1
		}
		br.Put(bb)
	}
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		br.LineMean(minSize - 1)
	}
}

func Benchmark_buf5Ring_LineMean(b *testing.B) {
	b.StopTimer()
	br := newBuf5Ring()
	for i := 0; i < buf5RingSize; i++ {
		n := byte((i + 1) * 2)
		bb := make([]byte, n)
		bb[minSize-1] = n
		if len(bb) > minSize {
			bb[minSize] = n + 1
		}
		br.Put(bb)
	}
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		br.LineMean(minSize - 1)
	}
}
