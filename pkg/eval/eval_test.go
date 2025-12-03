package eval

import (
	"testing"

	"github.com/kjbreil/goscript/pkg/state"
)

func TestEvaluate(t *testing.T) {
	type args struct {
		states state.States
		eval   string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "on",
			args: args{
				//nolint:exhaustruct // LastChanged and LastUpdated not needed for test
				states: state.NewSingleStates("sensor.test", &state.State{
					DomainEntity: "sensor.test",
					Domain:       "humidity",
					Entity:       "test",
					State:        "on",
					Attributes:   nil,
				}),
				eval: `state == "on"`,
			},
			want: true,
		},
		{
			name: "float",
			args: args{
				//nolint:exhaustruct // LastChanged and LastUpdated not needed for test
				states: state.NewSingleStates("sensor.test", &state.State{
					DomainEntity: "sensor.test",
					Domain:       "sensor",
					Entity:       "test",
					State:        "31.1",
					Attributes:   nil,
				}),
				eval: "float(state) > 10.000000",
			},
			want: true,
		},
		{
			name: "int",
			args: args{
				//nolint:exhaustruct // LastChanged and LastUpdated not needed for test
				states: state.NewSingleStates("sensor.test", &state.State{
					DomainEntity: "sensor.test",
					Domain:       "sensor",
					Entity:       "test",
					State:        "1",
					Attributes:   nil,
				}),
				eval: "int(state) > 0",
			},
			want: true,
		},
		{
			name: "int is 1 want 0",
			args: args{
				//nolint:exhaustruct // LastChanged and LastUpdated not needed for test
				states: state.NewSingleStates("sensor.test", &state.State{
					DomainEntity: "sensor.test",
					Domain:       "sensor",
					Entity:       "test",
					State:        "1",
					Attributes:   nil,
				}),
				eval: "int(state) = 0",
			},
			want: false,
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
