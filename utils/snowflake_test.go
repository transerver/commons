package utils

import "testing"

func TestSnowId(t *testing.T) {
	t.Log(GenerateIdWithDefaultNode().Int64())
	t.Log(GenerateIdWithDefaultNode().Int64())
	t.Log(GenerateIdWithDefaultNode().Int64())
	t.Log(GenerateIdWithDefaultNode().Int64())
	t.Log(GenerateIdWithDefaultNode().Int64())
	t.Log(GenerateIdWithDefaultNode().Int64())
	t.Log(GenerateIdWithDefaultNode().Int64())
	t.Log(GenerateIdWithDefaultNode().Int64())
	t.Log(GenerateIdWithDefaultNode().Int64())
	t.Log(GenerateIdWithDefaultNode().Int64())
	t.Log(GenerateIdWithDefaultNode().Int64())
	t.Log(GenerateIdWithDefaultNode().Step())
}

func BenchmarkGenerateIdWithDefaultNode(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		GenerateIdWithDefaultNode().Int64()
	}
}
