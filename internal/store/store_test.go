package store

import (
	"fmt"
	"os"
	"reflect"
	"testing"
	"time"

	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
)

func init() {
	godotenv.Load("../../config/test.env")
}

func TestOpen(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			name: "Should open the store",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			Open()
			_, err := os.Open(DataFilePath())
			if err != nil {
				t.Error("Fialed to open data store:", err)
			}
		})
	}
}

func TestNewSpot(t *testing.T) {
	type args struct {
		ID           string
		registeredBy string
		openDate     time.Time
	}
	tests := []struct {
		name string
		args args
		want Spot
	}{
		{
			name: "Should create a new spot instance",
			args: args{
				ID:           "11",
				registeredBy: "jjrambo",
				openDate:     time.Now(),
			},
			want: Spot{
				ID:           "11",
				OpenDate:     time.Now().Format(SpotDateFormat),
				RegDate:      time.Now().Format(SpotDateFormat),
				RegisteredBy: "jjrambo",
			},
		},
		{
			name: "Should create a new spot instance open tomorrow",
			args: args{
				ID:           "11",
				registeredBy: "jjrambo",
				openDate:     time.Now().AddDate(0, 0, 1),
			},
			want: Spot{
				ID:           "11",
				OpenDate:     time.Now().AddDate(0, 0, 1).Format(SpotDateFormat),
				RegDate:      time.Now().Format(SpotDateFormat),
				RegisteredBy: "jjrambo",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewSpot(tt.args.ID, tt.args.registeredBy, tt.args.openDate); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewSpot() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSpot_key(t *testing.T) {
	type fields struct {
		ID           string
		OpenDate     time.Time
		RegDate      time.Time
		RegisteredBy string
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "Should return a key",
			fields: fields{
				ID:           "44",
				OpenDate:     time.Now(),
				RegDate:      time.Now(),
				RegisteredBy: "YourMom",
			},
			want: fmt.Sprintf("%s-%s", "44", time.Now().Format(SpotDateFormat)),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			spot := Spot{
				ID:           tt.fields.ID,
				OpenDate:     tt.fields.OpenDate.Format(SpotDateFormat),
				RegDate:      tt.fields.RegDate.Format(SpotDateFormat),
				RegisteredBy: tt.fields.RegisteredBy,
			}
			if got := spot.key(); got != tt.want {
				t.Errorf("Spot.key() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_formatKey(t *testing.T) {
	type args struct {
		id   string
		date time.Time
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "Should format a correct key",
			args: args{
				id:   "B11",
				date: time.Now(),
			},
			want: fmt.Sprintf("%s-%s", "B11", time.Now().Format(SpotDateFormat)),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := formatKey(tt.args.id, tt.args.date); got != tt.want {
				t.Errorf("formatKey() = %v, want %v", got, tt.want)
			}
		})
	}
}

func cleanup() {
	os.Remove(DataFilePath())
}

func testSpots() []Spot {
	return []Spot{
		{
			ID:           "B0",
			OpenDate:     time.Now().AddDate(0, 0, -1).Format(SpotDateFormat),
			RegDate:      time.Now().Format(SpotDateFormat),
			RegisteredBy: "Fred",
		},
		{
			ID:           "B1",
			OpenDate:     time.Now().Format(SpotDateFormat),
			RegDate:      time.Now().Format(SpotDateFormat),
			RegisteredBy: "YourMom",
		},
		{
			ID:           "B2",
			OpenDate:     time.Now().Format(SpotDateFormat),
			RegDate:      time.Now().Format(SpotDateFormat),
			RegisteredBy: "YourMom",
		},
		{
			ID:           "B3",
			OpenDate:     time.Now().AddDate(0, 0, 1).Format(SpotDateFormat),
			RegDate:      time.Now().Format(SpotDateFormat),
			RegisteredBy: "FredsMom",
		},
		{
			ID:           "B4",
			OpenDate:     time.Now().Format(SpotDateFormat),
			RegDate:      time.Now().Format(SpotDateFormat),
			RegisteredBy: "BarneysMom",
		},
	}
}

// registers spots for test and ignores errors
func registerSpotsForTest(spots []Spot) {
	for _, spot := range spots {
		od, _ := time.Parse(SpotDateFormat, spot.OpenDate)
		Register(spot.ID, spot.RegisteredBy, od)
	}
}

func TestSpotBase_Find(t *testing.T) {
	defer cleanup()
	type fields struct {
		spots []Spot
	}
	tests := []struct {
		name    string
		fields  fields
		want    map[string]Spot
		wantErr bool
	}{
		{
			name: "Should find nothing",
			fields: fields{
				spots: []Spot{},
			},
			want:    make(map[string]Spot),
			wantErr: true,
		},
		{
			name: "Should find two available",
			fields: fields{
				spots: testSpots(),
			},
			want: map[string]Spot{
				formatKey("B1", time.Now()): Spot{
					ID:           "B1",
					OpenDate:     time.Now().Format(SpotDateFormat),
					RegDate:      time.Now().Format(SpotDateFormat),
					RegisteredBy: "YourMom",
				},
				formatKey("B2", time.Now()): Spot{
					ID:           "B2",
					OpenDate:     time.Now().Format(SpotDateFormat),
					RegDate:      time.Now().Format(SpotDateFormat),
					RegisteredBy: "YourMom",
				},
				formatKey("B4", time.Now()): Spot{
					ID:           "B4",
					OpenDate:     time.Now().Format(SpotDateFormat),
					RegDate:      time.Now().Format(SpotDateFormat),
					RegisteredBy: "BarneysMom",
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cleanup()
			Open()
			registerSpotsForTest(tt.fields.spots)
			got, err := Find()
			if (err != nil) != tt.wantErr {
				t.Errorf("SpotBase.Find() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SpotBase.Find() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSpotBase_Claim(t *testing.T) {
	// defer cleanup()
	type fields struct {
		spots []Spot
	}
	type args struct {
		id   string
		user string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    Spot
		wantErr bool
	}{
		{
			name: "Should claim a spot",
			fields: fields{
				spots: testSpots(),
			},
			args: args{
				id:   "B1",
				user: "Captain Fantastic",
			},
			want: Spot{
				ID:           "B1",
				OpenDate:     time.Now().Format(SpotDateFormat),
				RegDate:      time.Now().Format(SpotDateFormat),
				RegisteredBy: "YourMom",
			},
			wantErr: false,
		},
		{
			name: "Should not claim a spot, no spots registered",
			fields: fields{
				spots: testSpots(),
			},
			args: args{
				id:   "B5",
				user: "Captain Fantastic",
			},
			want: Spot{
				ID: NotAvailable,
			},
			wantErr: true,
		},
		{
			name: "Should not claim a registered spot the is not avialable today",
			fields: fields{
				spots: testSpots(),
			},
			args: args{
				id:   "B3",
				user: "Captain Fantastic",
			},
			want: Spot{
				ID: NotAvailable,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			Open()
			cleanup()
			registerSpotsForTest(tt.fields.spots)
			got, err := Claim(tt.args.id, tt.args.user)
			if (err != nil) != tt.wantErr {
				t.Errorf("SpotBase.Claim() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SpotBase.Claim() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSpotBase_Register(t *testing.T) {
	defer cleanup()
	type fields struct {
		spots []Spot
	}
	type args struct {
		id       string
		user     string
		openDate time.Time
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    Spot
		wantErr bool
	}{
		{
			name: "should register",
			fields: fields{
				spots: testSpots(),
			},
			args: args{
				id:       "B11",
				user:     "pparker",
				openDate: time.Now(),
			},
			want: Spot{
				ID:           "B11",
				RegisteredBy: "pparker",
				OpenDate:     time.Now().Format(SpotDateFormat),
				RegDate:      time.Now().Format(SpotDateFormat),
			},
			wantErr: false,
		},
		{
			name: "should register in the future",
			fields: fields{
				spots: testSpots(),
			},
			args: args{
				id:       "B12",
				user:     "pparker",
				openDate: time.Now().AddDate(0, 0, 1),
			},
			want: Spot{
				ID:           "B12",
				RegisteredBy: "pparker",
				OpenDate:     time.Now().AddDate(0, 0, 1).Format(SpotDateFormat),
				RegDate:      time.Now().Format(SpotDateFormat),
			},
			wantErr: false,
		},
		{
			name: "should not register alreasy registered",
			fields: fields{
				spots: testSpots(),
			},
			args: args{
				id:       "B1",
				user:     "pparker",
				openDate: time.Now(),
			},
			want: Spot{
				ID:           "B1",
				RegisteredBy: "YourMom",
				OpenDate:     time.Now().Format(SpotDateFormat),
				RegDate:      time.Now().Format(SpotDateFormat),
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			Open()
			registerSpotsForTest(tt.fields.spots)
			got, err := Register(tt.args.id, tt.args.user, tt.args.openDate)
			if (err != nil) != tt.wantErr {
				t.Errorf("SpotBase.Register() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SpotBase.Register() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSpotBase_DropRegistration(t *testing.T) {
	defer cleanup()
	type fields struct {
		spots []Spot
	}
	type args struct {
		id   string
		user string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "should drop registration",
			fields: fields{
				spots: testSpots(),
			},
			args: args{
				id:   "B1",
				user: "YourMom",
			},
			wantErr: false,
		},
		{
			name: "should not drop registration not original registered-by",
			fields: fields{
				spots: testSpots(),
			},
			args: args{
				id:   "B1",
				user: "NotYourMommy",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			Open()
			registerSpotsForTest(tt.fields.spots)
			if err := DropRegistration(tt.args.id, tt.args.user); (err != nil) != tt.wantErr {
				t.Errorf("SpotBase.DropRegistration() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestDropAllRegistrations(t *testing.T) {
	defer cleanup()
	type args struct {
		user string
	}
	tests := []struct {
		name          string
		args          args
		expectedCount int
	}{
		{
			name: "should drop all registrations for YourMom",
			args: args{
				user: "YourMom",
			},
			expectedCount: 1,
		},
		{
			name: "shouldn't drop anything",
			args: args{
				user: "NothingToDrop",
			},
			expectedCount: 3,
		},
	}
	for _, tt := range tests {
		Open()
		registerSpotsForTest(testSpots())
		t.Run(tt.name, func(t *testing.T) {
			DropAllRegistrations(tt.args.user)
			openspots, _ := Find()
			assert.True(t, len(openspots) == tt.expectedCount)
			for _, spot := range openspots {
				if spot.RegisteredBy == tt.args.user {
					t.Errorf("Expected no reservations for user: %s", tt.args.user)
				}
			}
		})
	}
}
