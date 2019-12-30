package handlers

import (
	"fmt"
	"net/http/httptest"
	"net/url"
	"os"
	"testing"
	"time"

	"github.com/jasonholmberg/slashspot/internal/store"
	"github.com/joho/godotenv"
	"github.com/nlopes/slack"
	"gotest.tools/v3/assert"
)

func init() {
	godotenv.Load("../../config/test.env")
}

func formValsHelper(in map[string]string) url.Values {
	values := make(url.Values)
	for k, v := range in {
		values.Set(k, v)
	}
	return values
}

func Test_spotCommandHandler(t *testing.T) {
	defer cleanup()
	type args struct {
		cmd *slack.SlashCommand
		rr  *httptest.ResponseRecorder
	}
	tests := []struct {
		name             string
		args             args
		expectedResponse string
	}{
		{
			name: "Test no spot command",
			args: args{
				cmd: &slack.SlashCommand{
					Text:   "",
					UserID: "scooby",
				},
				rr: httptest.NewRecorder(),
			},
			expectedResponse: IDKBlank,
		},
		{
			name: "Test help command",
			args: args{
				cmd: &slack.SlashCommand{
					Text:   "help",
					UserID: "scooby",
				},
				rr: httptest.NewRecorder(),
			},
			expectedResponse: HelpText,
		},
		{
			name: "Test find command",
			args: args{
				cmd: &slack.SlashCommand{
					Text:   "find",
					UserID: "scooby",
				},
				rr: httptest.NewRecorder(),
			},
			expectedResponse: NoSpotsAvailable,
		},
		{
			name: "Test unknown command",
			args: args{
				cmd: &slack.SlashCommand{
					Text:   "bacon",
					UserID: "scooby",
				},
				rr: httptest.NewRecorder(),
			},
			expectedResponse: fmt.Sprintf(IDKTemplate, "bacon"),
		},
		{
			name: "Test reg command",
			args: args{
				cmd: &slack.SlashCommand{
					Text:   "reg 12",
					UserID: "scooby",
				},
				rr: httptest.NewRecorder(),
			},
			expectedResponse: fmt.Sprintf(SpotRegisteredTemplate, "12"),
		},
		{
			name: "Test claim command",
			args: args{
				cmd: &slack.SlashCommand{
					Text:   "take 12",
					UserID: "scooby",
				},
				rr: httptest.NewRecorder(),
			},
			expectedResponse: fmt.Sprintf(SpotClaimedTemplate, "12"),
		},
		{
			name: "Test reg command",
			args: args{
				cmd: &slack.SlashCommand{
					Text:   "reg 13",
					UserID: "scooby",
				},
				rr: httptest.NewRecorder(),
			},
			expectedResponse: fmt.Sprintf(SpotRegisteredTemplate, "13"),
		},
		{
			name: "Test drop command",
			args: args{
				cmd: &slack.SlashCommand{
					Text:   "drop 13",
					UserID: "scooby",
				},
				rr: httptest.NewRecorder(),
			},
			expectedResponse: fmt.Sprintf(SpotDropRegTemplate, "13"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			store.Open()
			spotCommandHandler(tt.args.cmd, tt.args.rr)
			if tt.args.rr.Code >= 300 {
				t.Errorf("Spot call return a non 200 response: %v", tt.args.rr.Code)
			}
			assert.Equal(t, tt.args.rr.Body.String(), tt.expectedResponse)
		})
	}
}

func cleanup() {
	os.Remove(store.DataFilePath())
}

func testSpots() []store.Spot {
	return []store.Spot{
		{
			ID:           "B0",
			OpenDate:     time.Now().AddDate(0, 0, -1).Format(store.SpotDateFormat),
			RegDate:      time.Now().Format(store.SpotDateFormat),
			RegisteredBy: "Fred",
		},
		{
			ID:           "B1",
			OpenDate:     time.Now().Format(store.SpotDateFormat),
			RegDate:      time.Now().Format(store.SpotDateFormat),
			RegisteredBy: "YourMom",
		},
		{
			ID:           "B2",
			OpenDate:     time.Now().Format(store.SpotDateFormat),
			RegDate:      time.Now().Format(store.SpotDateFormat),
			RegisteredBy: "YourMom",
		},
		{
			ID:           "B3",
			OpenDate:     time.Now().AddDate(0, 0, 1).Format(store.SpotDateFormat),
			RegDate:      time.Now().Format(store.SpotDateFormat),
			RegisteredBy: "FredsMom",
		},
		{
			ID:           "B4",
			OpenDate:     time.Now().Format(store.SpotDateFormat),
			RegDate:      time.Now().Format(store.SpotDateFormat),
			RegisteredBy: "BarneysMom",
		},
	}
}

// registers spots for test and ignores errors
func registerSpotsForTest(spots []store.Spot) {
	for _, spot := range spots {
		od, _ := time.Parse(store.SpotDateFormat, spot.OpenDate)
		store.Register(spot.ID, spot.RegisteredBy, od)
	}
}

func Test_handleFind(t *testing.T) {
	defer cleanup()
	type args struct {
		params []string
		spots  []store.Spot
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "Should find some spots",
			args: args{
				params: []string{"find"},
				spots:  testSpots(),
			},
			want: fmt.Sprintf(OpenSpotsTemplate, "B1,B2,B4"),
		},
		{
			name: "Should not find some spots",
			args: args{
				params: []string{"find"},
				spots:  []store.Spot{},
			},
			want: NoSpotsAvailable,
		},
	}
	for _, tt := range tests {
		cleanup()
		store.Open()
		registerSpotsForTest(tt.args.spots)
		t.Run(tt.name, func(t *testing.T) {
			if got := handleFind(tt.args.params); got != tt.want {
				t.Errorf("handleFind() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_handleRegister(t *testing.T) {
	defer cleanup()
	type args struct {
		params []string
		cmd    *slack.SlashCommand
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "Should register a spot for today",
			args: args{
				params: []string{"reg", "A1"},
				cmd: &slack.SlashCommand{
					UserID: "yourmom",
				},
			},
			want: fmt.Sprintf(SpotRegisteredTemplate, "A1"),
		},
		{
			name: "Should register a spot for one day in the future",
			args: args{
				params: []string{"reg", "A2", time.Now().AddDate(0, 0, 1).Format(store.SpotDateFormat)},
				cmd: &slack.SlashCommand{
					UserID: "yourmom",
				},
			},
			want: fmt.Sprintf(SpotRegisteredTemplate, "A2"),
		},
		{
			name: "Should not register a spot for day in the past",
			args: args{
				params: []string{"reg", "A2", time.Now().AddDate(0, 0, -1).Format(store.SpotDateFormat)},
				cmd: &slack.SlashCommand{
					UserID: "yourmom",
				},
			},
			want: fmt.Sprintf(SpotPastDateRegistrationErrorTemplate, time.Now().AddDate(0, 0, -1).Format(store.SpotDateFormat)),
		},
	}
	for _, tt := range tests {
		store.Open()
		t.Run(tt.name, func(t *testing.T) {
			if got := handleRegister(tt.args.cmd, tt.args.params); got != tt.want {
				t.Errorf("handleRegister() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_handleClaim(t *testing.T) {
	defer cleanup()
	type args struct {
		params []string
		cmd    *slack.SlashCommand
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "should claim a spot",
			args: args{
				params: []string{"take", "B4"},
				cmd: &slack.SlashCommand{
					UserID: "ponyboy",
				},
			},
			want: fmt.Sprintf(SpotClaimedTemplate, "B4"),
		},
		{
			name: "should not claim a spot and get an error",
			args: args{
				params: []string{"take"},
				cmd: &slack.SlashCommand{
					UserID: "ponyboy",
				},
			},
			want: IDKBlank,
		},
		{
			name: "should not claim a spot with wrong ID",
			args: args{
				params: []string{"take", "X11"},
				cmd: &slack.SlashCommand{
					UserID: "ponyboy",
				},
			},
			want: fmt.Sprintf(SpotClaimErrorTemplate, "spot X11 not available"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cleanup()
			store.Open()
			registerSpotsForTest(testSpots())
			if got := handleClaim(tt.args.cmd, tt.args.params); got != tt.want {
				t.Errorf("handleReserve() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_handleHelp(t *testing.T) {
	tests := []struct {
		name string
		want string
	}{
		{
			name: "should get help text",
			want: HelpText,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := handleHelp(); got != tt.want {
				t.Errorf("handleHelp() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_handleUnknown(t *testing.T) {
	type args struct {
		action string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "should handle unknown command",
			args: args{
				action: "scooby-snack",
			},
			want: fmt.Sprintf(IDKTemplate, "scooby-snack"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := handleUnknown(tt.args.action); got != tt.want {
				t.Errorf("handleUnknown() = %v, want %v", got, tt.want)
			}
		})
	}
}
