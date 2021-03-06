package util

import (
	"testing"
	"time"
)

func TestBeforeNow(t *testing.T) {
	type args struct {
		test string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "should be before",
			args: args{
				test: time.Now().AddDate(0, 0, -1).Format(SpotDateFormat),
			},
			want: true,
		},
		{
			name: "should not be before",
			args: args{
				test: time.Now().AddDate(0, 0, 1).Format(SpotDateFormat),
			},
			want: false,
		},
		{
			name: "should not be before",
			args: args{
				test: time.Now().Format(SpotDateFormat),
			},
			want: false,
		},
		{
			name: "should not be before",
			args: args{
				test: time.Now().Format(SpotDateFormat),
			},
			want: false,
		},
		{
			name: "should not be before",
			args: args{
				test: time.Now().AddDate(0, 1, 1).Format(SpotDateFormat),
			},
			want: false,
		},
		{
			name: "should not be before",
			args: args{
				test: time.Now().AddDate(1, 1, 1).Format(SpotDateFormat),
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := BeforeNow(tt.args.test); got != tt.want {
				t.Errorf("BeforeNow() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAfterNow(t *testing.T) {
	type args struct {
		test string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "should be after 0",
			args: args{
				test: time.Now().AddDate(0, 0, 1).Format(SpotDateFormat),
			},
			want: true,
		},
		{
			name: "should be after 1",
			args: args{
				test: time.Now().AddDate(0, 1, 1).Format(SpotDateFormat),
			},
			want: true,
		},
		{
			name: "should be after 2",
			args: args{
				test: time.Now().AddDate(1, 1, 1).Format(SpotDateFormat),
			},
			want: true,
		},
		{
			name: "should not be after 0",
			args: args{
				test: time.Now().AddDate(0, 0, -1).Format(SpotDateFormat),
			},
			want: false,
		},
		{
			name: "should not be after 1",
			args: args{
				test: time.Now().Format(SpotDateFormat),
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := AfterNow(tt.args.test); got != tt.want {
				t.Errorf("AfterNow() = %v, want %v", got, tt.want)
			}
		})
	}
}
