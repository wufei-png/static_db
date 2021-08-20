package test_db

import "testing"

func BenchmarkDB(b *testing.B) {
	for i := 0; i < b.N; i++ {
		test()
	}
}
