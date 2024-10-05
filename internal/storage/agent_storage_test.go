package storage

import "testing"

func TestMemStorage_setStorage(t *testing.T) {
	type fields struct {
		Metrics Metrics
	}
	type args struct {
		r float64
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &MemStorage{
				Metrics: tt.fields.Metrics,
			}
			m.setStorage(tt.args.r)
		})
	}
}
