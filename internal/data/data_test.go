package data

import (
	"os"
	"reflect"
	"testing"

	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
)

func init() {
	godotenv.Load("../../config/test.env")
	os.MkdirAll(os.Getenv("SPOT_DATA_DIR"), os.ModePerm)
}

var (
	dataStr = `{
		"B1-2020-01-05": {
			"ID": "B1",
			"OpenDate": "2020-01-05",
			"RegDate": "2020-01-05",
			"RegisteredBy": "YourMom"
		},
		"B2-2020-01-05": {
			"ID": "B2",
			"OpenDate": "2020-01-05",
			"RegDate": "2020-01-05",
			"RegisteredBy": "YourMom"
		},
		"B3-2020-01-06": {
			"ID": "B3",
			"OpenDate": "2020-01-06",
			"RegDate": "2020-01-05",
			"RegisteredBy": "FredsMom"
		},
		"B4-2020-01-05": {
			"ID": "B4",
			"OpenDate": "2020-01-05",
			"RegDate": "2020-01-05",
			"RegisteredBy": "BarneysMom"
		}
	}`
)

func setupTestStore() {
	f, err := os.Create(FilePath())
	defer f.Close()
	if err != nil {
		panic(err)
	}
	_, err = f.WriteString(dataStr)
	if err != nil {
		panic(err)
	}
}

func cleanup() {
	os.Remove(FilePath())
	store = nil
}

func TestOpen(t *testing.T) {
	defer cleanup()
	tests := []struct {
		name    string
		preload bool
	}{
		{
			name:    "should open the store",
			preload: false,
		},
		{
			name:    "should open an existing store",
			preload: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.preload {
				setupTestStore()
			}
			Open()
			assert.NotNil(t, store, "should not be nil")
			assert.True(t, IsOpen(), "should be open")
			if tt.preload {
				assert.Equal(t, 4, len(store), "should have 4 spots")
			}
		})
	}
}

func TestIsOpen(t *testing.T) {
	defer cleanup()
	tests := []struct {
		name string
		want bool
	}{
		{
			name: "should be open",
			want: true,
		},
		{
			name: "should not be open",
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.want {
				Open()
			}
			if got := IsOpen(); got != tt.want {
				t.Errorf("IsOpen() = %v, want %v", got, tt.want)
			}
			store = nil
		})
	}
}

func TestSave(t *testing.T) {
	defer cleanup()
	tests := []struct {
		name    string
		wantErr bool
	}{
		{
			name:    "should save",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := Save(); (err != nil) != tt.wantErr {
				t.Errorf("Save() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestLoad(t *testing.T) {
	defer cleanup()
	tests := []struct {
		name    string
		want    map[string]Spot
		wantErr bool
	}{
		{
			name: "should load",
			want: map[string]Spot{
				"B1-2020-01-05": Spot{
					ID:           "B1",
					OpenDate:     "2020-01-05",
					RegDate:      "2020-01-05",
					RegisteredBy: "YourMom",
				},
				"B2-2020-01-05": Spot{
					ID:           "B2",
					OpenDate:     "2020-01-05",
					RegDate:      "2020-01-05",
					RegisteredBy: "YourMom",
				},
				"B3-2020-01-06": Spot{
					ID:           "B3",
					OpenDate:     "2020-01-06",
					RegDate:      "2020-01-05",
					RegisteredBy: "FredsMom",
				},
				"B4-2020-01-05": Spot{
					ID:           "B4",
					OpenDate:     "2020-01-05",
					RegDate:      "2020-01-05",
					RegisteredBy: "BarneysMom",
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		cleanup()
		setupTestStore()
		Open()
		t.Run(tt.name, func(t *testing.T) {
			got, err := Load()
			if (err != nil) != tt.wantErr {
				t.Errorf("Load() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Load() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDelete(t *testing.T) {
	defer cleanup()
	type args struct {
		id string
	}
	tests := []struct {
		name        string
		args        args
		shouldExist bool
	}{
		{
			name: "should delete an id that exists",
			args: args{
				id: "B1-2020-01-05",
			},
			shouldExist: true,
		},
		{
			name: "should not delete an id that doesn't exists",
			args: args{
				id: "B11-2020-01-05",
			},
			shouldExist: false,
		},
	}
	for _, tt := range tests {
		setupTestStore()
		Open()
		t.Run(tt.name, func(t *testing.T) {
			if tt.shouldExist {
				assert.False(t, store[tt.args.id].IsZeroValue())
			}
			Delete(tt.args.id)
			assert.True(t, store[tt.args.id].IsZeroValue())
		})
	}
}

func TestInsert(t *testing.T) {
	defer cleanup()
	type args struct {
		s Spot
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "should insert",
			args: args{
				s: Spot{
					ID:           "TEST",
					OpenDate:     "2020-01-05",
					RegDate:      "2020-01-05",
					RegisteredBy: "BarneysMom",
				},
			},
		},
	}
	for _, tt := range tests {
		Open()
		t.Run(tt.name, func(t *testing.T) {
			Insert(tt.args.s)
			assert.False(t, store[tt.args.s.Key()].IsZeroValue())
		})
	}
}
