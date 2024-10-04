package handlers

import "testing"

func TestSendMetric(t *testing.T) {
	type args struct {
		metricType  string
		metricName  string
		metricValue string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "simple test #2",
			args: args{
				metricType:  "counter",
				metricName:  "test",
				metricValue: "1",
			},
			wantErr: false,
		},
		{name: "simple test #2",
			args: args{
				metricType:  "counter",
				metricName:  "test",
				metricValue: "1",
			},
			wantErr: false,
		},
		{
			name: "simple test #3",
			args: args{
				metricType:  "counter",
				metricName:  "test",
				metricValue: "1",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := SendMetric(tt.args.metricType, tt.args.metricName, tt.args.metricValue); (err != nil) != tt.wantErr {
				t.Errorf("SendMetric() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
