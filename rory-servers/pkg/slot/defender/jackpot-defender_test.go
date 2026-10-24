package defender

import (
	"context"
	"log"
	"testing"
	"time"

	"example.invalid/golang/protobuf/proto"
	"example.invalid/ServerTeam/Servers/pkg/slot"
	"example.invalid/ServerTeam/Servers/pkg/slot/protodata"
)

func TestJackpotDefender_Pool(t *testing.T) {
	type fields struct {
		contribution int64
		award        int64
		increment    int64
		target       []slot.JackpotLevel
	}
	tests := []struct {
		name   string
		fields fields
		want   int64
	}{
		{
			name: "Normal",
			fields: fields{
				contribution: 123456789,
				award:        12,
				increment:    100,
				target:       []slot.JackpotLevel{slot.JackpotGrand},
			},
			want: 12333,
		},
		{
			name: "Negative",
			fields: fields{
				contribution: 12345,
				award:        12345,
				increment:    100,
				target:       []slot.JackpotLevel{slot.JackpotGrand},
			},
			want: -12344,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			def := &JackpotDefender{
				contribution: tt.fields.contribution,
				award:        tt.fields.award,
				increment:    tt.fields.increment,
				target:       tt.fields.target,
			}
			if got := def.Pool(); got != tt.want {
				t.Errorf("JackpotDefender.Pool() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestJackpotDefender_Increase(t *testing.T) {
	type fields struct {
		contribution int64
		award        int64
		increment    int64
		target       []slot.JackpotLevel
	}
	type args struct {
		pay int64
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   int64
	}{
		{
			name: "Increase Zero",
			fields: fields{
				contribution: 123456789,
				award:        12,
				increment:    100,
				target:       []slot.JackpotLevel{slot.JackpotGrand},
			},
			args: args{pay: 0},
			want: 0,
		},
		{
			name: "Increase Normal",
			fields: fields{
				contribution: 123456789,
				award:        12,
				increment:    100,
				target:       []slot.JackpotLevel{slot.JackpotGrand},
			},
			args: args{pay: 88},
			want: 8800,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			def := &JackpotDefender{
				contribution: tt.fields.contribution,
				award:        tt.fields.award,
				increment:    tt.fields.increment,
				target:       tt.fields.target,
			}
			if got := def.Increase(tt.args.pay); got != tt.want {
				t.Errorf("JackpotDefender.Increase() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestJackpotDefender_Award(t *testing.T) {
	type fields struct {
		contribution int64
		award        int64
		increment    int64
		target       []slot.JackpotLevel
	}
	type args struct {
		award int64
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{
			name: "Award Zero",
			fields: fields{
				contribution: 123456789,
				award:        12,
				increment:    100,
				target:       []slot.JackpotLevel{slot.JackpotGrand},
			},
			args: args{award: 0},
		},
		{
			name: "Award Normal",
			fields: fields{
				contribution: 123456789,
				award:        12,
				increment:    100,
				target:       []slot.JackpotLevel{slot.JackpotGrand},
			},
			args: args{award: 6001234},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			def := &JackpotDefender{
				contribution: tt.fields.contribution,
				award:        tt.fields.award,
				increment:    tt.fields.increment,
				target:       tt.fields.target,
			}
			def.Award(tt.args.award)
		})
	}
}

func TestJackpotDefender_Eligible(t *testing.T) {
	type fields struct {
		contribution int64
		award        int64
		increment    int64
		target       []slot.JackpotLevel
	}
	type args struct {
		jackpot *slot.Jackpot
		poolFn  func() int64
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
	}{
		{
			name: "Not In Target",
			fields: fields{
				contribution: 500 * 1e4,
				award:        0,
				increment:    100,
				target:       []slot.JackpotLevel{slot.JackpotGrand},
			},
			args: args{
				jackpot: &slot.Jackpot{
					Level: slot.JackpotMinor,
				},
				poolFn: func() int64 { return 20000 },
			},
			want: true,
		},
		{
			name: "Not Enough",
			fields: fields{
				contribution: 500 * 1e4,
				award:        0,
				increment:    100,
				target:       []slot.JackpotLevel{slot.JackpotGrand},
			},
			args: args{
				jackpot: &slot.Jackpot{
					Level: slot.JackpotGrand,
				},
				poolFn: func() int64 { return 20000 },
			},
			want: false,
		},
		{
			name: "Enough",
			fields: fields{
				contribution: 50000 * 1e4,
				award:        0,
				increment:    100,
				target:       []slot.JackpotLevel{slot.JackpotGrand},
			},
			args: args{
				jackpot: &slot.Jackpot{
					Level: slot.JackpotGrand,
				},
				poolFn: func() int64 { return 20000 },
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			def := &JackpotDefender{
				contribution: tt.fields.contribution,
				award:        tt.fields.award,
				increment:    tt.fields.increment,
				target:       tt.fields.target,
			}
			if got := def.Eligible(tt.args.jackpot, tt.args.poolFn); got != tt.want {
				t.Errorf("JackpotDefender.Eligible() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestJackpotDefender_RealCase_Eligible(t *testing.T) {

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	def := NewJackpotDefender(100, slot.JackpotGrand)
	def.Increase(1e9)

	q := make(chan slot.JackpotLevel)

	go func() {
		jp := new(slot.Jackpot)
		jp.Origin = &protodata.JackpotTable{
			IsLinked: proto.Int32(0),
			StartUp:  proto.Int64(60000),
			Level:    proto.String("Mini"),
		}
		jp.Level = slot.JackpotMini

		f := func() int64 {
			return jp.GetCurrentPoolByBet(88, 8800)
		}

		for {
			time.Sleep(500 * time.Millisecond)
			if def.Eligible(jp, f) {
				log.Printf("Win Jackpot %v %d, Pool: %d\n", jp.Level, f(), def.Pool())
				def.Award(f())
			} else {
				log.Printf("Jackpot %v can not Win, want %d, Pool: %d\n", jp.Level, f(), def.Pool())
				q <- jp.Level
			}
		}
	}()

	go func() {
		jp := new(slot.Jackpot)
		jp.Origin = &protodata.JackpotTable{
			IsLinked: proto.Int32(1),
			StartUp:  proto.Int64(2250000),
			Level:    proto.String("Major"),
		}
		jp.Level = slot.JackpotMajor

		f := func() int64 {
			return jp.GetCurrentPoolByBet(8, 8800)
		}

		for {
			time.Sleep(500 * time.Millisecond)
			if def.Eligible(jp, f) {
				log.Printf("Win Jackpot %v %d, Pool: %d\n", jp.Level, f(), def.Pool())
				def.Award(f())
			} else {
				log.Printf("Jackpot %v can not Win, want %d, Pool: %d\n", jp.Level, f(), def.Pool())
				q <- jp.Level
			}
		}
	}()

	go func() {
		jp := new(slot.Jackpot)
		jp.Origin = &protodata.JackpotTable{
			IsLinked: proto.Int32(1),
			StartUp:  proto.Int64(2250000),
			Level:    proto.String("Grand"),
		}
		jp.Level = slot.JackpotGrand

		f := func() int64 {
			return jp.GetCurrentPoolByBet(8, 8800)
		}

		for {
			time.Sleep(1000 * time.Millisecond)
			if def.Eligible(jp, f) {
				log.Printf("Win Jackpot %v %d, Pool: %d\n", jp.Level, f(), def.Pool())
				def.Award(f())
			} else {
				log.Printf("Jackpot %v can not Win, want %d, Pool: %d\n", jp.Level, f(), def.Pool())
				q <- jp.Level
			}
		}
	}()

	select {
	case lv := <-q:
		log.Printf("JackpotDefender defensed %v \n", lv)
	case <-ctx.Done():
		log.Println("timeout of 5 seconds.")
	}

}
