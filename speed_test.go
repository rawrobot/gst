package gst

import (
	"testing"
	"unsafe"

	"github.com/bksworm/gst/testsuite"
)

const testArraySize = 256

type MemoryTS struct {
	tst       *testing.T
	bench     *testing.B
	testArray []byte
}

// SetUpSuite is called once before the very first test in suite runs
func (s *MemoryTS) SetUpSuite() {
	s.testArray = make([]byte, testArraySize)
	for i := 0; i < len(s.testArray); i++ {
		s.testArray[i] = byte(i)
	}
}

// TearDownSuite is called once after thevery last test in suite runs
func (s *MemoryTS) TearDownSuite() {
}

// SetUp is called before each test method
func (s *MemoryTS) SetUp() {
}

// TearDown is called after each test method
func (s *MemoryTS) TearDown() {
}

// Hook up  into the "go test" runner.
func TestIt(t *testing.T) {
	s := &MemoryTS{tst: t}
	testsuite.Run(t, s)
}

func BenchmarkIT(b *testing.B) {
	s := &MemoryTS{bench: b}
	testsuite.Bench(b, s)
}

//THERE ARE TESTS

func (s *MemoryTS) Test_BufCopy(t *testing.T) {
	dst := make([]byte, testArraySize)
	dst = buf_copy(s.testArray, dst, testArraySize)
	for i := 0; i < len(s.testArray); i++ {
		testsuite.AssertEqual(t, dst[i], s.testArray[i])
	}
	PrintMemUsage()
}

func (s *MemoryTS) Test_BufCopyRealloc(t *testing.T) {
	dst := make([]byte, testArraySize/2)
	dst = buf_copy(s.testArray, dst, testArraySize)
	for i := 0; i < len(s.testArray); i++ {
		testsuite.AssertEqual(t, dst[i], s.testArray[i])
	}
	PrintMemUsage()
}
func (s *MemoryTS) Test_BufCopyPool(t *testing.T) {
	bp := NewBytePool(2, testArraySize/2)
	for i := 0; i < 5; i++ {
		dst := bp.Get()
		//log.Printf("get len(dst)=%d", len(dst))
		dst = buf_copy(s.testArray, dst, testArraySize)
		dst[0] = 1
		//log.Printf("put len(dst)=%d", len(dst))
		bp.Put(dst)
		//log.Printf("bp.w=%d", bp.Width())
		//log.Printf("pooled %d", bp.NumPooled())
	}
	PrintMemUsage()
}

func (s *MemoryTS) Benchmark_StructOnStack(b *testing.B) {
	for i := 0; i < b.N; i++ {
		var buf [testArraySize]byte //to avoid malloc, but use stack
		bp := (unsafe.Pointer(&buf[0]))
		dst := (*[1 << 30]byte)(unsafe.Pointer(bp))
		copy(dst[:], s.testArray)
	}
}
func (s *MemoryTS) Benchmark_StructMalloc(b *testing.B) {
	for i := 0; i < b.N; i++ {
		go_malloc_copy_test(s.testArray, testArraySize)
	}
}

func (s *MemoryTS) Benchmark_BuffCopyMake(b *testing.B) {
	for i := 0; i < b.N; i++ {
		dst := make_buf_copy(s.testArray, testArraySize)
		dst[0] = 1
		dst = nil
	}
}

func (s *MemoryTS) Benchmark_BufCopy(b *testing.B) {
	for i := 0; i < b.N; i++ {
		s.BufCopy()
	}
}

func (s *MemoryTS) BufCopy() {
	dst := make([]byte, testArraySize)
	dst = buf_copy(s.testArray, dst, testArraySize)
	dst[0] = 1
	dst = nil
}

func (s *MemoryTS) Benchmark_BufCopyRealoc(b *testing.B) {
	for i := 0; i < b.N; i++ {
		dst := make([]byte, testArraySize/2)
		dst = buf_copy(s.testArray, dst, testArraySize)
		dst[0] = 1
		dst = nil
	}
}
func (s *MemoryTS) Benchmark_BufCopyPool(b *testing.B) {
	bp := NewBytePool(2, testArraySize/2)
	for i := 0; i < b.N; i++ {
		s.BufCopyPool(bp)
	}
}

func (s *MemoryTS) BufCopyPool(bp *BytePool) {
	dst := bp.Get()
	dst = buf_copy(s.testArray, dst, testArraySize)
	dst[0] = 1
	bp.Put(dst)
}
