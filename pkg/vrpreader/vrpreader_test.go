package vrpreader

import (
	"bytes"
	"io"
	"os"
	"testing"
)

func TestReadInstance(t *testing.T) {
	f, err := os.Open("27x12.txt")
	if err != nil {
		t.Fatalf("failed to open instance file: %v", err)
	}
	defer f.Close()

	buf := bytes.NewBuffer(nil)
	_, err = io.Copy(buf, f)
	if err != nil {
		t.Fatalf("f.Read(buf) failed: %v", err)
	}

	input := string(buf.Bytes())
	_, err = ReadInstance(input)
	if err != nil {
		t.Fatal("ReadInstance failed!")
	}
}

func benchmarkReadInstance(instance string, b *testing.B) {
	for n := 0; n < b.N; n++ {
		f, err := os.Open(instance)
		if err != nil {
			b.Fatalf("failed to open instance file: %v", err)
		}
		buf := bytes.NewBuffer(nil)
		_, err = io.Copy(buf, f)
		if err != nil {
			f.Close()
			b.Fatalf("io.Copy(buf, f) failed: %v", err)
		}
		f.Close()
		input := string(buf.Bytes())
		_, err = ReadInstance(input)
	}
}

func BenchmarkReadInstance27x12(b *testing.B) {
	benchmarkReadInstance("27x12.txt", b)
}
