package slot

import (
	"reflect"
	"testing"
)

func TestParsePayoutMysteryJP(t *testing.T) {
	type args struct {
		payout string
	}
	tests := []struct {
		name            string
		args            args
		wantJackpotID   int32
		wantProp        int64
		wantDenominator int64
		wantErr         bool
	}{
		{
			name:            "Error Format",
			args:            args{payout: "23:43|1|3"},
			wantJackpotID:   0,
			wantProp:        0,
			wantDenominator: 0,
			wantErr:         true,
		},
		{
			name:            "Normal Number",
			args:            args{payout: "7|3123:10000000000"},
			wantJackpotID:   7,
			wantProp:        3123,
			wantDenominator: 1e10,
			wantErr:         false,
		},
		{
			name:            "Science Number",
			args:            args{payout: "7|3.123e3:1e10"},
			wantJackpotID:   7,
			wantProp:        3123,
			wantDenominator: 1e10,
			wantErr:         false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotJackpotID, gotProp, gotDenominator, err := ParsePayoutMysteryJP(tt.args.payout)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParsePayoutMysteryJP() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotJackpotID != tt.wantJackpotID {
				t.Errorf("ParsePayoutMysteryJP() gotJackpotID = %v, want %v", gotJackpotID, tt.wantJackpotID)
			}
			if gotProp != tt.wantProp {
				t.Errorf("ParsePayoutMysteryJP() gotProp = %v, want %v", gotProp, tt.wantProp)
			}
			if gotDenominator != tt.wantDenominator {
				t.Errorf("ParsePayoutMysteryJP() gotDenominator = %v, want %v", gotDenominator, tt.wantDenominator)
			}
		})
	}
}

func TestParsePayoutMysteryJPR(t *testing.T) {
	type args struct {
		payout  string
		betRate int64
	}
	tests := []struct {
		name            string
		args            args
		wantJackpotID   int32
		wantProp        int64
		wantDenominator int64
		wantErr         bool
	}{
		{
			name:            "Error Format",
			args:            args{payout: "23:43|1|3", betRate: 8},
			wantJackpotID:   0,
			wantProp:        0,
			wantDenominator: 0,
			wantErr:         true,
		},
		{
			name:            "Normal Number 8",
			args:            args{payout: "7|8:100:10000,18:200:10000,38:300:10000,68:400:10000,88:500:10000", betRate: 8},
			wantJackpotID:   7,
			wantProp:        100,
			wantDenominator: 10000,
			wantErr:         false,
		},
		{
			name:            "Normal Number 18",
			args:            args{payout: "7|8:100:10000,18:200:10000,38:300:10000,68:400:10000,88:500:10000", betRate: 18},
			wantJackpotID:   7,
			wantProp:        200,
			wantDenominator: 10000,
			wantErr:         false,
		},
		{
			name:            "Normal Number 38",
			args:            args{payout: "7|8:100:10000,18:200:10000,38:300:10000,68:400:10000,88:500:10000", betRate: 38},
			wantJackpotID:   7,
			wantProp:        300,
			wantDenominator: 10000,
			wantErr:         false,
		},
		{
			name:            "Normal Number 68",
			args:            args{payout: "7|8:100:10000,18:200:10000,38:300:10000,68:400:10000,88:500:10000", betRate: 68},
			wantJackpotID:   7,
			wantProp:        400,
			wantDenominator: 10000,
			wantErr:         false,
		},
		{
			name:            "Normal Number 88",
			args:            args{payout: "7|8:100:10000,18:200:10000,38:300:10000,68:400:10000,88:500:10000", betRate: 88},
			wantJackpotID:   7,
			wantProp:        500,
			wantDenominator: 10000,
			wantErr:         false,
		},
		{
			name:            "Science Number 8",
			args:            args{payout: "7|8:1.11e2:1e5,18:2.11e2:1e5,38:3.11e2:1e5,68:4.11e2:1e5,88:5.11e2:1e5", betRate: 8},
			wantJackpotID:   7,
			wantProp:        111,
			wantDenominator: 100000,
			wantErr:         false,
		},
		{
			name:            "Science Number 18",
			args:            args{payout: "7|8:1.11e2:1e5,18:2.11e2:1e5,38:3.11e2:1e5,68:4.11e2:1e5,88:5.11e2:1e5", betRate: 18},
			wantJackpotID:   7,
			wantProp:        211,
			wantDenominator: 100000,
			wantErr:         false,
		},
		{
			name:            "Science Number 38",
			args:            args{payout: "7|8:1.11e2:1e5,18:2.11e2:1e5,38:3.11e2:1e5,68:4.11e2:1e5,88:5.11e2:1e5", betRate: 38},
			wantJackpotID:   7,
			wantProp:        311,
			wantDenominator: 100000,
			wantErr:         false,
		},
		{
			name:            "Science Number 68",
			args:            args{payout: "7|8:1.11e2:1e5,18:2.11e2:1e5,38:3.11e2:1e5,68:4.11e2:1e5,88:5.11e2:1e5", betRate: 68},
			wantJackpotID:   7,
			wantProp:        411,
			wantDenominator: 100000,
			wantErr:         false,
		},
		{
			name:            "Science Number 88",
			args:            args{payout: "7|8:1.11e2:1e5,18:2.11e2:1e5,38:3.11e2:1e5,68:4.11e2:1e5,88:5.11e2:1e5", betRate: 88},
			wantJackpotID:   7,
			wantProp:        511,
			wantDenominator: 100000,
			wantErr:         false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotJackpotID, gotProp, gotDenominator, err := ParsePayoutMysteryJPR(tt.args.payout, tt.args.betRate)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParsePayoutMysteryJPR() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotJackpotID != tt.wantJackpotID {
				t.Errorf("ParsePayoutMysteryJPR() gotJackpotID = %v, want %v", gotJackpotID, tt.wantJackpotID)
			}
			if gotProp != tt.wantProp {
				t.Errorf("ParsePayoutMysteryJPR() gotProp = %v, want %v", gotProp, tt.wantProp)
			}
			if gotDenominator != tt.wantDenominator {
				t.Errorf("ParsePayoutMysteryJPR() gotDenominator = %v, want %v", gotDenominator, tt.wantDenominator)
			}
		})
	}
}

func TestParsePayoutMysteryJPG(t *testing.T) {
	type args struct {
		payout string
	}
	tests := []struct {
		name            string
		args            args
		wantLinkProp    int64
		wantNolinkProp  int64
		wantDenominator int64
		wantErr         bool
	}{
		{
			name:            "Format Error",
			args:            args{payout: "1:3200:10000:100000000"},
			wantLinkProp:    0,
			wantNolinkProp:  0,
			wantDenominator: 0,
			wantErr:         true,
		},
		{
			name:            "Normal Number",
			args:            args{payout: "200:10000:100000000"},
			wantLinkProp:    200,
			wantNolinkProp:  10000,
			wantDenominator: 100000000,
			wantErr:         false,
		},
		{
			name:            "Science Number",
			args:            args{payout: "2e2:1e4:1e8"},
			wantLinkProp:    200,
			wantNolinkProp:  10000,
			wantDenominator: 100000000,
			wantErr:         false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotLinkProp, gotNolinkProp, gotDenominator, err := ParsePayoutMysteryJPG(tt.args.payout)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParsePayoutMysteryJPG() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotLinkProp != tt.wantLinkProp {
				t.Errorf("ParsePayoutMysteryJPG() gotLinkProp = %v, want %v", gotLinkProp, tt.wantLinkProp)
			}
			if gotNolinkProp != tt.wantNolinkProp {
				t.Errorf("ParsePayoutMysteryJPG() gotNolinkProp = %v, want %v", gotNolinkProp, tt.wantNolinkProp)
			}
			if gotDenominator != tt.wantDenominator {
				t.Errorf("ParsePayoutMysteryJPG() gotDenominator = %v, want %v", gotDenominator, tt.wantDenominator)
			}
		})
	}
}

func TestParsePayoutMysteryJPGR(t *testing.T) {
	type args struct {
		payout  string
		betRate int64
	}
	tests := []struct {
		name            string
		args            args
		wantLinkProp    int64
		wantNolinkProp  int64
		wantDenominator int64
		wantErr         bool
	}{
		{
			name:            "Format Error",
			args:            args{payout: "1:3|200:10000:100000000"},
			wantLinkProp:    0,
			wantNolinkProp:  0,
			wantDenominator: 0,
			wantErr:         true,
		},
		{
			name:            "Normal Number 8",
			args:            args{payout: "8:10:1000:100000,18:20:1000:100000,38:30:1000:100000,68:40:1000:100000,88:50:1000:100000", betRate: 8},
			wantLinkProp:    10,
			wantNolinkProp:  1000,
			wantDenominator: 100000,
			wantErr:         false,
		},
		{
			name:            "Normal Number 18",
			args:            args{payout: "8:10:1000:100000,18:20:1000:100000,38:30:1000:100000,68:40:1000:100000,88:50:1000:100000", betRate: 18},
			wantLinkProp:    20,
			wantNolinkProp:  1000,
			wantDenominator: 100000,
			wantErr:         false,
		},
		{
			name:            "Normal Number 38",
			args:            args{payout: "8:10:1000:100000,18:20:1000:100000,38:30:1000:100000,68:40:1000:100000,88:50:1000:100000", betRate: 38},
			wantLinkProp:    30,
			wantNolinkProp:  1000,
			wantDenominator: 100000,
			wantErr:         false,
		},
		{
			name:            "Normal Number 68",
			args:            args{payout: "8:10:1000:100000,18:20:1000:100000,38:30:1000:100000,68:40:1000:100000,88:50:1000:100000", betRate: 68},
			wantLinkProp:    40,
			wantNolinkProp:  1000,
			wantDenominator: 100000,
			wantErr:         false,
		},
		{
			name:            "Normal Number 88",
			args:            args{payout: "8:10:1000:100000,18:20:1000:100000,38:30:1000:100000,68:40:1000:100000,88:50:1000:100000", betRate: 88},
			wantLinkProp:    50,
			wantNolinkProp:  1000,
			wantDenominator: 100000,
			wantErr:         false,
		},
		{
			name:            "Science Number 8",
			args:            args{payout: "8:1e1:1e3:1e5,18:2e1:1e3:1e5,38:3e1:1e3:1e5,68:4e1:1e3:1e5,88:5e1:1e3:1e5", betRate: 8},
			wantLinkProp:    10,
			wantNolinkProp:  1000,
			wantDenominator: 100000,
			wantErr:         false,
		},
		{
			name:            "Science Number 18",
			args:            args{payout: "8:1e1:1e3:1e5,18:2e1:1e3:1e5,38:3e1:1e3:1e5,68:4e1:1e3:1e5,88:5e1:1e3:1e5", betRate: 18},
			wantLinkProp:    20,
			wantNolinkProp:  1000,
			wantDenominator: 100000,
			wantErr:         false,
		},
		{
			name:            "Science Number 38",
			args:            args{payout: "8:1e1:1e3:1e5,18:2e1:1e3:1e5,38:3e1:1e3:1e5,68:4e1:1e3:1e5,88:5e1:1e3:1e5", betRate: 38},
			wantLinkProp:    30,
			wantNolinkProp:  1000,
			wantDenominator: 100000,
			wantErr:         false,
		},
		{
			name:            "Science Number 68",
			args:            args{payout: "8:1e1:1e3:1e5,18:2e1:1e3:1e5,38:3e1:1e3:1e5,68:4e1:1e3:1e5,88:5e1:1e3:1e5", betRate: 68},
			wantLinkProp:    40,
			wantNolinkProp:  1000,
			wantDenominator: 100000,
			wantErr:         false,
		},
		{
			name:            "Science Number 88",
			args:            args{payout: "8:1e1:1e3:1e5,18:2e1:1e3:1e5,38:3e1:1e3:1e5,68:4e1:1e3:1e5,88:5e1:1e3:1e5", betRate: 88},
			wantLinkProp:    50,
			wantNolinkProp:  1000,
			wantDenominator: 100000,
			wantErr:         false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotLinkProp, gotNolinkProp, gotDenominator, err := ParsePayoutMysteryJPGR(tt.args.payout, tt.args.betRate)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParsePayoutMysteryJPGR() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotLinkProp != tt.wantLinkProp {
				t.Errorf("ParsePayoutMysteryJPGR() gotLinkProp = %v, want %v", gotLinkProp, tt.wantLinkProp)
			}
			if gotNolinkProp != tt.wantNolinkProp {
				t.Errorf("ParsePayoutMysteryJPGR() gotNolinkProp = %v, want %v", gotNolinkProp, tt.wantNolinkProp)
			}
			if gotDenominator != tt.wantDenominator {
				t.Errorf("ParsePayoutMysteryJPGR() gotDenominator = %v, want %v", gotDenominator, tt.wantDenominator)
			}
		})
	}
}

func TestParseExtraRuleConditionAfterSpin(t *testing.T) {
	type args struct {
		value string
	}
	tests := []struct {
		name         string
		args         args
		wantLinkType string
		wantSymbol   string
		wantCount    int
		wantErr      bool
	}{
		{
			name:         "Format Error",
			args:         args{"1|2:3|4"},
			wantLinkType: "",
			wantSymbol:   "",
			wantCount:    0,
			wantErr:      true,
		},
		{
			name:         "Normal 1",
			args:         args{"any:g1:15"},
			wantLinkType: "any",
			wantSymbol:   "g1",
			wantCount:    15,
			wantErr:      false,
		},
		{
			name:         "Normal 2",
			args:         args{"any:g1:0"},
			wantLinkType: "any",
			wantSymbol:   "g1",
			wantCount:    0,
			wantErr:      false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotLinkType, gotSymbol, gotCount, err := ParseExtraRuleConditionAfterSpin(tt.args.value)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseExtraRuleConditionAfterSpin() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotLinkType != tt.wantLinkType {
				t.Errorf("ParseExtraRuleConditionAfterSpin() gotLinkType = %v, want %v", gotLinkType, tt.wantLinkType)
			}
			if gotSymbol != tt.wantSymbol {
				t.Errorf("ParseExtraRuleConditionAfterSpin() gotSymbol = %v, want %v", gotSymbol, tt.wantSymbol)
			}
			if gotCount != tt.wantCount {
				t.Errorf("ParseExtraRuleConditionAfterSpin() gotCount = %v, want %v", gotCount, tt.wantCount)
			}
		})
	}
}

func TestParseExtraRuleActionRespin(t *testing.T) {
	type args struct {
		value string
	}
	tests := []struct {
		name            string
		args            args
		wantProp        int64
		wantDenominator int64
		wantErr         bool
	}{
		{
			name:            "Format Error",
			args:            args{"1|3:3|23|3"},
			wantProp:        0,
			wantDenominator: 0,
			wantErr:         true,
		},
		{
			name:            "Normal Number",
			args:            args{"1:10"},
			wantProp:        1,
			wantDenominator: 10,
			wantErr:         false,
		},
		{
			name:            "Science Number",
			args:            args{"1e1:1e5"},
			wantProp:        10,
			wantDenominator: 100000,
			wantErr:         false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotProp, gotDenominator, err := ParseExtraRuleActionRespin(tt.args.value)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseExtraRuleActionRespin() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotProp != tt.wantProp {
				t.Errorf("ParseExtraRuleActionRespin() gotProp = %v, want %v", gotProp, tt.wantProp)
			}
			if gotDenominator != tt.wantDenominator {
				t.Errorf("ParseExtraRuleActionRespin() gotDenominator = %v, want %v", gotDenominator, tt.wantDenominator)
			}
		})
	}
}

func TestParseExtraRuleActionRespinR(t *testing.T) {
	type args struct {
		pattern string
		rate    int64
	}
	tests := []struct {
		name            string
		args            args
		wantProp        int64
		wantDenominator int64
		wantErr         bool
	}{
		{
			name:            "Format Error",
			args:            args{"1|3:3|23|3", 8},
			wantProp:        0,
			wantDenominator: 0,
			wantErr:         true,
		},
		{
			name:            "Normal Number",
			args:            args{"8:1:10,18:5:15,38:3:20", 8},
			wantProp:        1,
			wantDenominator: 10,
			wantErr:         false,
		},
		{
			name:            "Science Number",
			args:            args{"8:1e1:1e5,18:1e2:1e5,28:1e3:1e5,38:1e4:1e5", 18},
			wantProp:        100,
			wantDenominator: 1e5,
			wantErr:         false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotProp, gotDenominator, err := ParseExtraRuleActionRespinR(tt.args.pattern, tt.args.rate)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseExtraRuleActionRespinR() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotProp != tt.wantProp {
				t.Errorf("ParseExtraRuleActionRespinR() gotProp = %v, want %v", gotProp, tt.wantProp)
			}
			if gotDenominator != tt.wantDenominator {
				t.Errorf("ParseExtraRuleActionRespinR() gotDenominator = %v, want %v", gotDenominator, tt.wantDenominator)
			}
		})
	}
}

func TestParseExtraRuleActionUniGoldBall(t *testing.T) {
	type args struct {
		value string
	}
	tests := []struct {
		name            string
		args            args
		wantSymbol      string
		wantProp        int64
		wantDenominator int64
		wantErr         bool
	}{
		{
			name:            "Format Error",
			args:            args{"1|3:3|23|3"},
			wantSymbol:      "",
			wantProp:        0,
			wantDenominator: 0,
			wantErr:         true,
		},
		{
			name:            "Normal Number",
			args:            args{"g19:1:10"},
			wantSymbol:      "g19",
			wantProp:        1,
			wantDenominator: 10,
			wantErr:         false,
		},
		{
			name:            "Science Number",
			args:            args{"g8:1e1:1e5"},
			wantSymbol:      "g8",
			wantProp:        10,
			wantDenominator: 100000,
			wantErr:         false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotSymbol, gotProp, gotDenominator, err := ParseExtraRuleActionUniGoldBall(tt.args.value)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseExtraRuleActionUniGoldBall() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotSymbol != tt.wantSymbol {
				t.Errorf("ParseExtraRuleActionUniGoldBall() gotSymbol = %v, want %v", gotSymbol, tt.wantSymbol)
			}
			if gotProp != tt.wantProp {
				t.Errorf("ParseExtraRuleActionUniGoldBall() gotProp = %v, want %v", gotProp, tt.wantProp)
			}
			if gotDenominator != tt.wantDenominator {
				t.Errorf("ParseExtraRuleActionUniGoldBall() gotDenominator = %v, want %v", gotDenominator, tt.wantDenominator)
			}
		})
	}
}

func TestParseExtraRuleConditionAppearSymbol(t *testing.T) {
	type args struct {
		value string
	}
	tests := []struct {
		name       string
		args       args
		wantSymbol string
		wantRor    string
		wantVals   []int
		wantErr    bool
	}{
		{
			name:       "Format Error",
			args:       args{"3:3ww|2ee3|qq3"},
			wantSymbol: "",
			wantRor:    "",
			wantVals:   []int(nil),
			wantErr:    true,
		},
		{
			name:       "Normal Number",
			args:       args{"w14|Reel|1,2,3,4"},
			wantSymbol: "w14",
			wantRor:    "Reel",
			wantVals:   []int{1, 2, 3, 4},
			wantErr:    false,
		},
		{
			name:       "Normal Number",
			args:       args{"w14|Row|1,2,3,4"},
			wantSymbol: "w14",
			wantRor:    "Row",
			wantVals:   []int{1, 2, 3, 4},
			wantErr:    false,
		},
		{
			name:       "Normal Number",
			args:       args{"w14|Any"},
			wantSymbol: "w14",
			wantRor:    "Any",
			wantVals:   []int(nil),
			wantErr:    false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotSymbol, gotRor, gotVals, err := ParseExtraRuleConditionAppearSymbol(tt.args.value)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseExtraRuleConditionAppearSymbol() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotSymbol != tt.wantSymbol {
				t.Errorf("ParseExtraRuleConditionAppearSymbol() gotSymbol = %v, want %v", gotSymbol, tt.wantSymbol)
			}
			if gotRor != tt.wantRor {
				t.Errorf("ParseExtraRuleConditionAppearSymbol() gotRor = %v, want %v", gotRor, tt.wantRor)
			}
			if !reflect.DeepEqual(gotVals, tt.wantVals) {
				t.Errorf("ParseExtraRuleConditionAppearSymbol() gotVals = %v, want %v", gotVals, tt.wantVals)
			}
		})
	}
}

func TestParseExtraRuleActionReplaceBG(t *testing.T) {
	type args struct {
		pattern string
	}
	tests := []struct {
		name            string
		args            args
		wantRor         string
		wantTarget      int
		wantSymbol      string
		wantProp        int64
		wantDenominator int64
		wantErr         bool
	}{
		{
			name:            "Error Format",
			args:            args{pattern: "1,2,3,4:342;e"},
			wantRor:         "",
			wantTarget:      0,
			wantSymbol:      "",
			wantProp:        0,
			wantDenominator: 0,
			wantErr:         true,
		},
		{
			name:            "Normal Number",
			args:            args{pattern: "Reel:5:bg4:100:1000000"},
			wantRor:         "Reel",
			wantTarget:      5,
			wantSymbol:      "bg4",
			wantProp:        100,
			wantDenominator: 1000000,
			wantErr:         false,
		},
		{
			name:            "Science Number",
			args:            args{pattern: "Reel:5:bg4:1e2:1e6"},
			wantRor:         "Reel",
			wantTarget:      5,
			wantSymbol:      "bg4",
			wantProp:        100,
			wantDenominator: 1000000,
			wantErr:         false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotRor, gotTarget, gotSymbol, gotProp, gotDenominator, err := ParseExtraRuleActionReplaceBG(tt.args.pattern)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseExtraRuleActionReplaceBG() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotRor != tt.wantRor {
				t.Errorf("ParseExtraRuleActionReplaceBG() gotRor = %v, want %v", gotRor, tt.wantRor)
			}
			if gotTarget != tt.wantTarget {
				t.Errorf("ParseExtraRuleActionReplaceBG() gotTarget = %v, want %v", gotTarget, tt.wantTarget)
			}
			if gotSymbol != tt.wantSymbol {
				t.Errorf("ParseExtraRuleActionReplaceBG() gotSymbol = %v, want %v", gotSymbol, tt.wantSymbol)
			}
			if gotProp != tt.wantProp {
				t.Errorf("ParseExtraRuleActionReplaceBG() gotProp = %v, want %v", gotProp, tt.wantProp)
			}
			if gotDenominator != tt.wantDenominator {
				t.Errorf("ParseExtraRuleActionReplaceBG() gotDenominator = %v, want %v", gotDenominator, tt.wantDenominator)
			}
		})
	}
}

func TestParseExtraRuleActionReplaceBGR(t *testing.T) {
	type args struct {
		pattern string
		betRate int64
	}
	tests := []struct {
		name            string
		args            args
		wantRor         string
		wantTarget      int
		wantSymbol      string
		wantProp        int64
		wantDenominator int64
		wantErr         bool
	}{
		{
			name:            "Error Format",
			args:            args{pattern: "1,2,3,4:342;e", betRate: 1},
			wantRor:         "",
			wantTarget:      0,
			wantSymbol:      "",
			wantProp:        0,
			wantDenominator: 0,
			wantErr:         true,
		},
		{
			name:            "Normal Number 28",
			args:            args{pattern: "28:Reel:5:bg4:200:1000000,38:Reel:5:bg4:300:1000000,68:Reel:5:bg4:600:1000000,88:Reel:5:bg4:800:1000000", betRate: 28},
			wantRor:         "Reel",
			wantTarget:      5,
			wantSymbol:      "bg4",
			wantProp:        200,
			wantDenominator: 1000000,
			wantErr:         false,
		},
		{
			name:            "Normal Number 38",
			args:            args{pattern: "28:Reel:5:bg4:200:1000000,38:Reel:5:bg4:300:1000000,68:Reel:5:bg4:600:1000000,88:Reel:5:bg4:800:1000000", betRate: 38},
			wantRor:         "Reel",
			wantTarget:      5,
			wantSymbol:      "bg4",
			wantProp:        300,
			wantDenominator: 1000000,
			wantErr:         false,
		},
		{
			name:            "Normal Number 68",
			args:            args{pattern: "28:Reel:5:bg4:200:1000000,38:Reel:5:bg4:300:1000000,68:Reel:5:bg4:600:1000000,88:Reel:5:bg4:800:1000000", betRate: 68},
			wantRor:         "Reel",
			wantTarget:      5,
			wantSymbol:      "bg4",
			wantProp:        600,
			wantDenominator: 1000000,
			wantErr:         false,
		},
		{
			name:            "Normal Number 88",
			args:            args{pattern: "28:Reel:5:bg4:200:1000000,38:Reel:5:bg4:300:1000000,68:Reel:5:bg4:600:1000000,88:Reel:5:bg4:800:1000000", betRate: 88},
			wantRor:         "Reel",
			wantTarget:      5,
			wantSymbol:      "bg4",
			wantProp:        800,
			wantDenominator: 1000000,
			wantErr:         false,
		},
		{
			name:            "Not Exist BetRate",
			args:            args{pattern: "28:Reel:5:bg4:200:1000000,38:Reel:5:bg4:300:1000000,68:Reel:5:bg4:600:1000000,88:Reel:5:bg4:800:1000000", betRate: 2},
			wantRor:         "",
			wantTarget:      0,
			wantSymbol:      "",
			wantProp:        0,
			wantDenominator: 0,
			wantErr:         true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotRor, gotTarget, gotSymbol, gotProp, gotDenominator, err := ParseExtraRuleActionReplaceBGR(tt.args.pattern, tt.args.betRate)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseExtraRuleActionReplaceBGR() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotRor != tt.wantRor {
				t.Errorf("ParseExtraRuleActionReplaceBGR() gotRor = %v, want %v", gotRor, tt.wantRor)
			}
			if gotTarget != tt.wantTarget {
				t.Errorf("ParseExtraRuleActionReplaceBGR() gotTarget = %v, want %v", gotTarget, tt.wantTarget)
			}
			if gotSymbol != tt.wantSymbol {
				t.Errorf("ParseExtraRuleActionReplaceBGR() gotSymbol = %v, want %v", gotSymbol, tt.wantSymbol)
			}
			if gotProp != tt.wantProp {
				t.Errorf("ParseExtraRuleActionReplaceBGR() gotProp = %v, want %v", gotProp, tt.wantProp)
			}
			if gotDenominator != tt.wantDenominator {
				t.Errorf("ParseExtraRuleActionReplaceBGR() gotDenominator = %v, want %v", gotDenominator, tt.wantDenominator)
			}
		})
	}
}

func TestParseExtraRuleActionConditionalAction(t *testing.T) {
	type args struct {
		pattern string
	}
	tests := []struct {
		name          string
		args          args
		wantCondition ExtraRuleCondition
		wantAction    ExtraRuleAction
		wantErr       bool
	}{
		{
			name:          "JPDefense",
			args:          args{pattern: "JPDefense:Respin"},
			wantCondition: ExtraRuleConditionJPDefense,
			wantAction:    ExtraRuleActionRespin,
			wantErr:       false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotCondition, gotAction, err := ParseExtraRuleActionConditionalAction(tt.args.pattern)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseExtraRuleActionConditionalAction() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotCondition != tt.wantCondition {
				t.Errorf("ParseExtraRuleActionConditionalAction() gotCondition = %v, want %v", gotCondition, tt.wantCondition)
			}
			if gotAction != tt.wantAction {
				t.Errorf("ParseExtraRuleActionConditionalAction() gotAction = %v, want %v", gotAction, tt.wantAction)
			}
		})
	}
}
