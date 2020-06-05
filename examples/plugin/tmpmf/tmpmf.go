package tmpmf

import (
	//"fmt"

	"github.com/bksworm/gst"
)

type TempMedianFilter struct {
	gst.VideoIPTransformPlugin
	line int
}

func NewTempMedianFilter(e *gst.Element) *TempMedianFilter {
	lp := &TempMedianFilter{}
	lp.VideoFilterPlugin.Element = *e
	return lp
}

//draws horithontal black line at the midle of frame
func (lp *TempMedianFilter) TransformIP(vf *gst.VideoFrame) error {

	//y := vf.Plane(0)

	// grayImg := &image.Gray{
	// 	Pix:    y.Pixels,
	// 	Stride: y.Stride,
	// 	Rect:   image.Rect(y.Width, y.Height, 0, 0),
	// }
	// //FilterGrayIP(grayImg)
	return nil
}

type bufRing struct {
	Buffers [][]byte
	pos     int
	nBufs   int
}

func newBufRing(nBufs int) (res *bufRing) {
	res = new(bufRing)
	res.nBufs = nBufs
	res.Buffers = make([][]byte, nBufs)
	return res
}

func (br *bufRing) Put(b []byte) {
	dst := &br.Buffers[br.pos]
	if dst == nil || len(*dst) < len(b) {
		*dst = make([]byte, len(b))
	}
	copy(*dst, b)
	br.pos += 1
	if br.pos >= br.nBufs {
		br.pos = 0
	}
}

func (br *bufRing) Line(n int) (res []byte) {
	res = make([]byte, br.nBufs)
	for i := 0; i < br.nBufs; i++ {
		if len(br.Buffers[i]) > n {
			res[i] = (br.Buffers[i][n])
		} else {
			res[i] = 0
		}
	}
	return res
}

func (br *bufRing) LineMean(n int) byte {
	line := br.Line(n)
	InsertionSortU8(line)
	n = len(line)
	if (n & 1) == 1 {
		n = n / 2
	} else {
		n = n/2 - 1
	}
	return line[n]
}

const buf5RingSize = 5

type buf5Ring struct {
	Buffers [][]byte
	pos     int
}

func newBuf5Ring() (res *buf5Ring) {
	res = new(buf5Ring)
	res.Buffers = make([][]byte, buf5RingSize)
	return res
}

func (br *buf5Ring) Put(b []byte) {
	dst := &br.Buffers[br.pos]
	if dst == nil || len(*dst) < len(b) {
		*dst = make([]byte, len(b))
	}
	copy(*dst, b)
	br.pos += 1
	if br.pos >= buf5RingSize {
		br.pos = 0
	}
}

func (br *buf5Ring) Line(n int) (res [buf5RingSize]byte) {
	for i := 0; i < buf5RingSize; i++ {
		src := &br.Buffers[i]
		//fmt.Printf("%d", len(*src))
		if src == nil || len(*src) <= n {
			res[i] = 0
		} else {
			res[i] = (*src)[n]
		}
	}
	return res
}

func (br *buf5Ring) LineMean(n int) byte {
	line := PartialSort5Byte(br.Line(n))
	return line[buf5RingSize/2]
}

func InsertionSortU8(v []byte) {
	for j := 1; j < len(v); j++ {
		// Invariant: v[:j] contains the same elements as
		// the original slice v[:j], but in sorted order.
		key := v[j]
		i := j - 1
		for i >= 0 && v[i] > key {
			v[i+1] = v[i]
			i--
		}
		v[i+1] = key
	}
}

func PartialSort5Int(a [5]int) [5]int {
	a[2], a[3] = SortU8(a[2], a[3])
	a[1], a[2] = SortU8(a[1], a[2])
	a[2], a[3] = SortU8(a[2], a[3])
	a[4] = MaxU8(a[1], a[4])
	a[0] = MinU8(a[0], a[3])
	a[2], a[0] = SortU8(a[2], a[0])
	a[2] = MaxU8(a[4], a[2])
	a[2] = MinU8(a[2], a[0])
	return a
}

func SortU8(a int, b int) (int, int) {
	d := a - b
	m := ^(d >> 8)
	b += d & m
	a -= d & m
	return a, b
}

func MaxU8(a int, b int) int {
	d := a - b
	m := ^(d >> 8)
	return b + (d & m)
}

func MinU8(a int, b int) int {
	d := a - b
	m := ^(d >> 8)
	return a - (d & m)
}

func PartialSort5Byte(a [5]byte) [5]byte {
	a[2], a[3] = SortByte(a[2], a[3])
	a[1], a[2] = SortByte(a[1], a[2])
	a[2], a[3] = SortByte(a[2], a[3])
	a[4] = MaxByte(a[1], a[4])
	a[0] = MinByte(a[0], a[3])
	a[2], a[0] = SortByte(a[2], a[0])
	a[2] = MaxByte(a[4], a[2])
	a[2] = MinByte(a[2], a[0])
	return a
}

func SortByte(a byte, b byte) (byte, byte) {
	if b < a {
		return b, a
	}
	return a, b
}

func MaxByte(a byte, b byte) byte {
	if a > b {
		return a
	}
	return b
}

func MinByte(a byte, b byte) byte {
	if a < b {
		return a
	}
	return b
}
