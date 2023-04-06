package goscript

import (
	"testing"
)

func TestEvaluate(t *testing.T) {
	type args struct {
		states States
		eval   string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "float",
			args: args{
				states: States{
					"sensor.test": &State{
						DomainEntity: "sensor.test",
						Domain:       "humidity",
						Entity:       "test",
						State:        "31.1",
						Attributes:   nil,
					},
				},
				eval: "float(state) > 10.000000",
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Evaluate(tt.args.states, tt.args.eval); got != tt.want {
				t.Errorf("Evaluate() = %v, want %v", got, tt.want)
			}
		})
	}
}
