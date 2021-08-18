package http_db

import "testing"

func BenchmarkHello(b *testing.B) {
	for i := 0; i < b.N; i++ {
		wf()
	}
}
