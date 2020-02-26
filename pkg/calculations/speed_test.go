package calculations

import "testing"

func TestSpeedMilesPerHour(t *testing.T) {
	type args struct {
		timestamp1    int64
		timestamp2    int64
		distanceMiles float64
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{
			name: "boston to hartford",
			args: args{timestamp1: 1582466400, timestamp2: 1582471800, distanceMiles: 101.1},
			want: 67,
		},
		{
			name: "hartford to boston, verify absolute time",
			args: args{timestamp1: 1582471800, timestamp2: 1582466400, distanceMiles: 101.1},
			want: 67,
		},
		{
			name: "boston to austin",
			args: args{timestamp1: 1582466400, timestamp2: 1582471800, distanceMiles: 1696},
			want: 1130,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got, _ := SpeedMilesPerHour(tt.args.timestamp1, tt.args.timestamp2, tt.args.distanceMiles); got != tt.want {
				t.Errorf("SpeedMilesPerHour() = %v, want %v", got, tt.want)
			}
		})
	}
}
