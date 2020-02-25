package geoip

import (
	"net"
	"reflect"
	"testing"

	"github.com/jrxfive/superman-detector/internal/pkg/settings"
	"github.com/stretchr/testify/assert"
)

func TestDefaultLocator_Close(t *testing.T) {
	type fields struct {
		s settings.Specification
	}

	s := settings.NewSettings()
	s.GeoIPDatabaseFileLocation = "../../GeoLite2-City.mmdb"

	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name:    "Close",
			fields:  fields{s: s},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d, err := NewDefaultLocator(s)
			assert.NoError(t, err)
			if err := d.Close(); (err != nil) != tt.wantErr {
				t.Errorf("Close() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestDefaultLocator_Locate(t *testing.T) {
	type fields struct {
		s settings.Specification
	}

	type args struct {
		ip net.IP
	}

	s := settings.NewSettings()
	s.GeoIPDatabaseFileLocation = "../../GeoLite2-City.mmdb"

	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *Location
		wantErr bool
	}{
		{
			name:   "New York",
			fields: fields{s: s},
			args:   args{ip: net.ParseIP("206.81.252.6")},
			want: &Location{
				Latitude:  39.211,
				Longitude: -76.8362,
				Radius:    5,
			},
			wantErr: false,
		},
		{
			name:   "Hong Kong",
			fields: fields{s: s},
			args:   args{ip: net.ParseIP("119.28.48.231")},
			want: &Location{
				Latitude:  22.25,
				Longitude: 114.1667,
				Radius:    50,
			},
			wantErr: false,
		},
		{
			name:   "Moscow",
			fields: fields{s: s},
			args:   args{ip: net.ParseIP("31.173.221.5")},
			want: &Location{
				Latitude:  42.9753,
				Longitude: 47.5022,
				Radius:    1000,
			},
			wantErr: false,
		},
		{
			name:   "Sydney",
			fields: fields{s: s},
			args:   args{ip: net.ParseIP("203.2.218.214")},
			want: &Location{
				Latitude:  -33.8919,
				Longitude: 151.1554,
				Radius:    500,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d, err := NewDefaultLocator(s)
			assert.NoError(t, err)
			got, err := d.Locate(tt.args.ip)
			if (err != nil) != tt.wantErr {
				t.Errorf("Locate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Locate() got = %v, want %v", got, tt.want)
			}
		})
	}
}
