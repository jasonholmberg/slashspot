package util

import (
	"testing"
	"time"
)

func TestBeforeNow(t *testing.T) {
	type args struct {
		test time.Time
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "should be before",
			args: args{
				test: time.Now().AddDate(0, 0, -1),
			},
			want: true,
		},
		{
			name: "should not be before",
			args: args{
				test: time.Now().AddDate(0, 0, 1),
			},
			want: false,
		},
		{
			name: "should not be before",
			args: args{
				test: time.Now(),
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
		test time.Time
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "should be after",
			args: args{
				test: time.Now().AddDate(0, 0, 1),
			},
			want: true,
		},
		{
			name: "should not be after",
			args: args{
				test: time.Now().AddDate(0, 0, -1),
			},
			want: false,
		},
		{
			name: "should not be after",
			args: args{
				test: time.Now(),
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
