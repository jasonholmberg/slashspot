package data

import (
	"testing"
)

func TestSpot_Key(t *testing.T) {
	type fields struct {
		ID           string
		OpenDate     string
		RegDate      string
		RegisteredBy string
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "should correctly format spot key",
			fields: fields{
				ID: "T1",
				OpenDate: "2020-01-01",
				RegDate: "2020-01-01",
				RegisteredBy: "SuperFuzz",
			},
			want: "T1-2020-01-01",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			spot := Spot{
				ID:           tt.fields.ID,
				OpenDate:     tt.fields.OpenDate,
				RegDate:      tt.fields.RegDate,
				RegisteredBy: tt.fields.RegisteredBy,
			}
			if got := spot.Key(); got != tt.want {
				t.Errorf("Spot.Key() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSpot_IsZeroValue(t *testing.T) {
	type fields struct {
		ID           string
		OpenDate     string
		RegDate      string
		RegisteredBy string
	}
	tests := []struct {
		name   string
		fields fields
		want   bool
	}{
		{
			name: "should be zeroed",
			fields: fields{
				ID: "",
				OpenDate: "",
				RegDate: "",
				RegisteredBy: "",
			},
			want: true,
		},
		{
			name: "should not be zeroed",
			fields: fields{
				ID: "T1",
				OpenDate: "",
				RegDate: "",
				RegisteredBy: "",
			},
			want: false,
		},
		{
			name: "should not be zeroed",
			fields: fields{
				ID: "T1",
				OpenDate: "X",
				RegDate: "X",
				RegisteredBy: "",
			},
			want: false,
		},
		{
			name: "should not be zeroed",
			fields: fields{
				ID: "T1",
				OpenDate: "X",
				RegDate: "",
				RegisteredBy: "",
			},
			want: false,
		},
		{
			name: "should not be zeroed",
			fields: fields{
				ID: "",
				OpenDate: "",
				RegDate: "",
				RegisteredBy: "X",
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := Spot{
				ID:           tt.fields.ID,
				OpenDate:     tt.fields.OpenDate,
				RegDate:      tt.fields.RegDate,
				RegisteredBy: tt.fields.RegisteredBy,
			}
			if got := s.IsZeroValue(); got != tt.want {
				t.Errorf("Spot.IsZeroValue() = %v, want %v", got, tt.want)
			}
		})
	}
}
