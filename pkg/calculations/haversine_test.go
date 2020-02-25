package calculations

import "testing"

func TestCoordinatesDistance(t *testing.T) {
	type args struct {
		lat1 float64
		lon1 float64
		lat2 float64
		lon2 float64
	}
	tests := []struct {
		name    string
		args    args
		want    float64
		wantErr bool
	}{
		{
			name: "boston to harford",
			args: args{lat1: 42.3601, lon1: 71.0589, lat2: 41.7658, lon2: 72.6734},
			want: 92.4176144506052,
		},
		{
			name: "hartford to boston",
			args: args{lat1: 41.7658, lon1: 72.6734, lat2: 42.3601, lon2: 71.0589},
			want: 92.4176144506052,
		},
		{
			name:    "boston to invalid latitude",
			args:    args{lat1: 42.3601, lon1: 71.0589, lat2: 97.1254, lon2: 72.6734},
			want:    0,
			wantErr: true,
		},
		{
			name:    "boston to invalid longitude",
			args:    args{lat1: 42.3601, lon1: 71.0589, lat2: 97.1254, lon2: 192.6734},
			want:    0,
			wantErr: true,
		},
		{
			name:    "invalid to invalid ",
			args:    args{lat1: 95.3601, lon1: 144.0589, lat2: 97.1254, lon2: 192.6734},
			want:    0,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := CoordinatesDistance(tt.args.lat1, tt.args.lon1, tt.args.lat2, tt.args.lon2)
			if (err != nil) != tt.wantErr {
				t.Errorf("CoordinatesDistance() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("CoordinatesDistance() got = %v, want %v", got, tt.want)
			}
		})
	}
}
