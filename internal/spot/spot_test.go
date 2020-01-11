package spot

import (
	"fmt"
	"os"
	"reflect"
	"testing"
	"time"

	"github.com/jasonholmberg/slashspot/internal/data"
	"github.com/jasonholmberg/slashspot/internal/util"
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
			data.Open()
			_, err := os.Open(data.FilePath())
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
		want data.Spot
	}{
		{
			name: "Should create a new spot instance",
			args: args{
				ID:           "11",
				registeredBy: "jjrambo",
				openDate:     localTime(),
			},
			want: data.Spot{
				ID:           "11",
				OpenDate:     localTime().Format(util.SpotDateFormat),
				RegDate:      localTime().Format(util.SpotDateFormat),
				RegisteredBy: "jjrambo",
			},
		},
		{
			name: "Should create a new spot instance open tomorrow",
			args: args{
				ID:           "11",
				registeredBy: "jjrambo",
				openDate:     localTime().AddDate(0, 0, 1),
			},
			want: data.Spot{
				ID:           "11",
				OpenDate:     localTime().AddDate(0, 0, 1).Format(util.SpotDateFormat),
				RegDate:      localTime().Format(util.SpotDateFormat),
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
				OpenDate:     localTime(),
				RegDate:      localTime(),
				RegisteredBy: "slackuser",
			},
			want: fmt.Sprintf("%s-%s", "44", localTime().Format(util.SpotDateFormat)),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			spot := data.Spot{
				ID:           tt.fields.ID,
				OpenDate:     tt.fields.OpenDate.Format(util.SpotDateFormat),
				RegDate:      tt.fields.RegDate.Format(util.SpotDateFormat),
				RegisteredBy: tt.fields.RegisteredBy,
			}
			if got := spot.Key(); got != tt.want {
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
				date: localTime(),
			},
			want: fmt.Sprintf("%s-%s", "B11", localTime().Format(util.SpotDateFormat)),
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
	os.Remove(data.FilePath())
}

func localTime() time.Time {
	return time.Now().In(time.Local)
}

func testSpots() []data.Spot {
	return []data.Spot{
		{
			ID:           "B0",
			OpenDate:     localTime().AddDate(0, 0, -1).Format(util.SpotDateFormat),
			RegDate:      localTime().Format(util.SpotDateFormat),
			RegisteredBy: "Fred",
		},
		{
			ID:           "B1",
			OpenDate:     localTime().Format(util.SpotDateFormat),
			RegDate:      localTime().Format(util.SpotDateFormat),
			RegisteredBy: "slackuser",
		},
		{
			ID:           "B2",
			OpenDate:     localTime().Format(util.SpotDateFormat),
			RegDate:      localTime().Format(util.SpotDateFormat),
			RegisteredBy: "slackuser",
		},
		{
			ID:           "B3",
			OpenDate:     localTime().AddDate(0, 0, 1).Format(util.SpotDateFormat),
			RegDate:      localTime().Format(util.SpotDateFormat),
			RegisteredBy: "FredsMom",
		},
		{
			ID:           "B4",
			OpenDate:     localTime().Format(util.SpotDateFormat),
			RegDate:      localTime().Format(util.SpotDateFormat),
			RegisteredBy: "BarneysMom",
		},
	}
}

// registers spots for test and ignores errors
func registerSpotsForTest(spots []data.Spot) {
	for _, spot := range spots {
		od, _ := time.Parse(util.SpotDateFormat, spot.OpenDate)
		Register(spot.ID, spot.RegisteredBy, od)
	}
}

func TestSpotBase_Find(t *testing.T) {
	defer cleanup()
	type fields struct {
		spots []data.Spot
	}
	tests := []struct {
		name    string
		fields  fields
		want    map[string]data.Spot
		wantErr bool
	}{
		{
			name: "Should find nothing",
			fields: fields{
				spots: []data.Spot{},
			},
			want:    make(map[string]data.Spot),
			wantErr: true,
		},
		{
			name: "Should find three available",
			fields: fields{
				spots: testSpots(),
			},
			want: map[string]data.Spot{
				formatKey("B1", localTime()): data.Spot{
					ID:           "B1",
					OpenDate:     localTime().Format(util.SpotDateFormat),
					RegDate:      localTime().Format(util.SpotDateFormat),
					RegisteredBy: "slackuser",
				},
				formatKey("B2", localTime()): data.Spot{
					ID:           "B2",
					OpenDate:     localTime().Format(util.SpotDateFormat),
					RegDate:      localTime().Format(util.SpotDateFormat),
					RegisteredBy: "slackuser",
				},
				formatKey("B4", localTime()): data.Spot{
					ID:           "B4",
					OpenDate:     localTime().Format(util.SpotDateFormat),
					RegDate:      localTime().Format(util.SpotDateFormat),
					RegisteredBy: "BarneysMom",
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cleanup()
			data.Open()
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
	defer cleanup()
	type fields struct {
		spots []data.Spot
	}
	type args struct {
		id   string
		user string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    data.Spot
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
			want: data.Spot{
				ID:           "B1",
				OpenDate:     localTime().Format(util.SpotDateFormat),
				RegDate:      localTime().Format(util.SpotDateFormat),
				RegisteredBy: "slackuser",
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
			want: data.Spot{
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
			want: data.Spot{
				ID: NotAvailable,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data.Open()
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
		spots []data.Spot
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
		want    data.Spot
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
				openDate: localTime(),
			},
			want: data.Spot{
				ID:           "B11",
				RegisteredBy: "pparker",
				OpenDate:     localTime().Format(util.SpotDateFormat),
				RegDate:      localTime().Format(util.SpotDateFormat),
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
				openDate: localTime().AddDate(0, 0, 1),
			},
			want: data.Spot{
				ID:           "B12",
				RegisteredBy: "pparker",
				OpenDate:     localTime().AddDate(0, 0, 1).Format(util.SpotDateFormat),
				RegDate:      localTime().Format(util.SpotDateFormat),
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
				openDate: localTime(),
			},
			want: data.Spot{
				ID:           "B1",
				RegisteredBy: "slackuser",
				OpenDate:     localTime().Format(util.SpotDateFormat),
				RegDate:      localTime().Format(util.SpotDateFormat),
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data.Open()
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
		spots []data.Spot
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
				user: "slackuser",
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
				user: "Notslackusermy",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data.Open()
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
			name: "should drop all registrations for slackuser",
			args: args{
				user: "slackuser",
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
		data.Open()
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
