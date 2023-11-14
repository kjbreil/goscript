package goscript

import "testing"

func TestStateText_Equals(t *testing.T) {
	type args struct {
		cmp string
	}
	tests := []struct {
		name string
		st   StateText
		args args
		want bool
	}{
		{
			name: "Different Casing",
			st:   StateText("on"),
			args: args{
				cmp: "On",
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.st.Equals(tt.args.cmp); got != tt.want {
				t.Errorf("Equals() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStateText_Float(t *testing.T) {
	tests := []struct {
		name string
		st   StateText
		want float64
	}{
		{
			name: "0",
			st:   StateText("0"),
			want: 0,
		},
		{
			name: "1.32",
			st:   StateText("1.32"),
			want: 1.32,
		},
		{
			name: "0.00001",
			st:   StateText("0.00001"),
			want: 0.00001,
		},
		{
			name: "Contains String",
			st:   StateText("On"),
			want: 0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.st.Float(); got != tt.want {
				t.Errorf("Float() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStateText_Int(t *testing.T) {
	tests := []struct {
		name string
		st   StateText
		want int
	}{
		{
			name: "0",
			st:   StateText("0"),
			want: 0,
		},
		{
			name: "1.32",
			st:   StateText("1.32"),
			want: 1,
		},
		{
			name: "0.00001",
			st:   StateText("0.00001"),
			want: 0,
		},
		{
			name: "Contains String",
			st:   StateText("On"),
			want: 0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.st.Int(); got != tt.want {
				t.Errorf("Int() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStateText_String(t *testing.T) {
	tests := []struct {
		name string
		st   StateText
		want string
	}{
		{
			name: "0",
			st:   StateText("0"),
			want: "0",
		},
		{
			name: "On",
			st:   StateText("On"),
			want: "On",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.st.String(); got != tt.want {
				t.Errorf("String() = %v, want %v", got, tt.want)
			}
		})
	}
}
