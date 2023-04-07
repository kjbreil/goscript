package helpers

import (
	"reflect"
	"testing"
	"time"
)

func Test_timeToCron(t *testing.T) {
	type args struct {
		t time.Time
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "6:00am",
			args: args{t: time.Date(0, 0, 0, 6, 0, 0, 0, time.Local)},
			want: "0 6 * * *",
		},
		{
			name: "6:00am UTC",
			args: args{t: time.Date(0, 0, 0, 6, 0, 0, 0, time.UTC)},
			want: "0 6 * * *",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := TimeToCron(tt.args.t); got != tt.want {
				t.Errorf("TimeToCron() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSortCronJobs(t *testing.T) {
	type args struct {
		expressions []string
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		{
			name: "empty",
			args: args{
				expressions: []string{"30 1 0 0 0", "0 12 0 0 0", "0 1 0 0 0", "25 1 0 0 0"},
			},
			want: []string{"0 1 0 0 0", "25 1 0 0 0", "30 1 0 0 0", "0 12 0 0 0"},
		},

		{
			name: "empty",
			args: args{
				expressions: []string{"1 1 * * *", "0 12 * * *"},
			},
			want: []string{"1 1 * * *", "0 12 * * *"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := SortCronJobs(tt.args.expressions); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SortCronJobs() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestLastValidCron(t *testing.T) {
	type args struct {
		crons []string
		t     time.Time
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "test 1",
			args: args{
				crons: []string{"0 1 * * *", "0 12 * * *"},
				t:     time.Date(2023, 1, 1, 11, 0, 0, 0, time.Local),
			},
			want:    "0 1 * * *",
			wantErr: false,
		},
		{
			name: "test 1",
			args: args{
				crons: []string{"1 1 * * *", "0 12 * * *"},
				t:     time.Date(2023, 1, 1, 11, 0, 0, 0, time.Local),
			},
			want:    "1 1 * * *",
			wantErr: false,
		},
		{
			name: "test 1",
			args: args{
				crons: []string{"30 1 * * *", "0 12 * * *"},
				t:     time.Date(2023, 1, 1, 11, 0, 0, 0, time.Local),
			},
			want:    "30 1 * * *",
			wantErr: false,
		},
		{
			name: "test 1",
			args: args{
				crons: []string{"0 1 * * *", "0 12 * * *"},
				t:     time.Date(2023, 1, 1, 13, 0, 0, 0, time.Local),
			},
			want:    "0 12 * * *",
			wantErr: false,
		},
		{
			name: "test 1",
			args: args{
				crons: []string{"0 1 * * *", "0 12 * * *"},
				t:     time.Date(2023, 1, 1, 13, 0, 0, 0, time.Local),
			},
			want:    "0 12 * * *",
			wantErr: false,
		},
		{
			name: "test 1",
			args: args{
				crons: []string{"25 1 * * *", "0 1 * * *", "30 1 * * *"},
				t:     time.Date(2023, 1, 1, 1, 27, 0, 0, time.Local),
			},
			want:    "25 1 * * *",
			wantErr: false,
		},
		{
			name: "test 1",
			args: args{
				crons: []string{"0 1 * * *", "30 1 * * *", "25 1 * * *"},
				t:     time.Date(2023, 1, 1, 1, 27, 0, 0, time.Local),
			},
			want:    "25 1 * * *",
			wantErr: false,
		},
		{
			name: "test 1",
			args: args{
				crons: []string{"* 1 * * *", "18 11 * * *", "* 16 * * *"},
				t:     time.Date(2023, 1, 1, 11, 16, 0, 0, time.Local),
			},
			want:    "* 1 * * *",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := LastValidCron(tt.args.crons, tt.args.t)
			if (err != nil) != tt.wantErr {
				t.Errorf("LastValidCron() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("LastValidCron() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNextValidCron(t *testing.T) {
	type args struct {
		crons []string
		t     time.Time
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "test 1",
			args: args{
				crons: []string{"0 1 * * *", "0 12 * * *"},
				t:     time.Date(2023, 1, 1, 11, 0, 0, 0, time.Local),
			},
			want:    "0 12 * * *",
			wantErr: false,
		},
		{
			name: "test big diff",
			args: args{
				crons: []string{"0 1 * * *"},
				t:     time.Date(2023, 1, 1, 23, 0, 0, 0, time.Local),
			},
			want:    "0 1 * * *",
			wantErr: false,
		},
		{
			name: "test big diff",
			args: args{
				crons: []string{"* 1 * * *", "35 11 * * *", "* 21 * * *"},
				t:     time.Date(2023, 1, 1, 11, 30, 0, 0, time.Local),
			},
			want:    "35 11 * * *",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NextValidCron(tt.args.crons, tt.args.t)
			if (err != nil) != tt.wantErr {
				t.Errorf("NextValidCron() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("NextValidCron() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNextTime(t *testing.T) {
	type args struct {
		crons []string
		t     time.Time
	}
	tests := []struct {
		name    string
		args    args
		want    time.Time
		wantErr bool
	}{
		{
			name: "test 1",
			args: args{
				crons: []string{"0 1 * * *", "0 12 * * *"},
				t:     time.Date(2023, 1, 1, 11, 0, 0, 0, time.Local),
			},
			want:    time.Date(2023, 1, 1, 12, 0, 0, 0, time.Local),
			wantErr: false,
		},
		{
			name: "next day",
			args: args{
				crons: []string{"0 1 * * *"},
				t:     time.Date(2023, 1, 1, 23, 0, 0, 0, time.Local),
			},
			want:    time.Date(2023, 1, 2, 1, 0, 0, 0, time.Local),
			wantErr: false,
		},
		{
			name: "minutes not zer",
			args: args{
				crons: []string{"0 1 * * *"},
				t:     time.Date(2023, 1, 1, 23, 12, 0, 0, time.Local),
			},
			want:    time.Date(2023, 1, 2, 1, 0, 0, 0, time.Local),
			wantErr: false,
		},
		{
			name: "minutes cron *",
			args: args{
				crons: []string{"* 1 * * *"},
				t:     time.Date(2023, 1, 1, 23, 12, 0, 0, time.Local),
			},
			want:    time.Date(2023, 1, 2, 1, 0, 0, 0, time.Local),
			wantErr: false,
		},
		{
			name: "minutes cron * in current hour multiple entry",
			args: args{
				crons: []string{"* 1 * * *", "* 4 * * *"},
				t:     time.Date(2023, 1, 1, 1, 12, 0, 0, time.Local),
			},
			want:    time.Date(2023, 1, 1, 1, 13, 0, 0, time.Local),
			wantErr: false,
		},
		{
			name: "minutes cron * in current hour",
			args: args{
				crons: []string{"* 1 * * *"},
				t:     time.Date(2023, 1, 1, 1, 12, 0, 0, time.Local),
			},
			want:    time.Date(2023, 1, 1, 1, 13, 0, 0, time.Local),
			wantErr: false,
		},

		{
			name: "hours cron *",
			args: args{
				crons: []string{"* * * * *"},
				t:     time.Date(2023, 1, 1, 1, 12, 0, 0, time.Local),
			},
			want:    time.Date(2023, 1, 1, 1, 13, 0, 0, time.Local),
			wantErr: false,
		},
		{
			name: "hours cron * on the hour",
			args: args{
				crons: []string{"* * * * *"},
				t:     time.Date(2023, 1, 1, 1, 00, 0, 0, time.Local),
			},
			want:    time.Date(2023, 1, 1, 1, 01, 0, 0, time.Local),
			wantErr: false,
		},
		{
			name: "hours cron * on the 59 minute",
			args: args{
				crons: []string{"* * * * *"},
				t:     time.Date(2023, 1, 1, 1, 59, 0, 0, time.Local),
			},
			want:    time.Date(2023, 1, 1, 2, 00, 0, 0, time.Local),
			wantErr: false,
		},
		{
			name: "hours cron every 5 minutes",
			args: args{
				crons: []string{"*/5 * * * *"},
				t:     time.Date(2023, 1, 1, 1, 02, 0, 0, time.Local),
			},
			want:    time.Date(2023, 1, 1, 1, 05, 0, 0, time.Local),
			wantErr: false,
		},
		{
			name: "every 10 seconds",
			args: args{
				crons: []string{"*/10 * * * * *"},
				t:     time.Date(2023, 1, 1, 1, 2, 2, 0, time.Local),
			},
			want:    time.Date(2023, 1, 1, 1, 2, 10, 0, time.Local),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NextTime(tt.args.crons, tt.args.t)
			if (err != nil) != tt.wantErr {
				t.Errorf("NextTime() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NextTime() got = %v, want %v", got, tt.want)
			}
		})
	}
}
