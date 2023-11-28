package domain

import (
	"testing"
	"time"
)

func TestCalculateCountdown(t *testing.T) {
	type args struct {
		input time.Time
	}
	tests := []struct {
		name string
		args args
		want time.Duration
	}{
		{
			name: "01-01-2023T00:00:00+00:00 (no holidays expected)",
			args: args{
				input: time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC),
			},
			want: time.Date(2023, 1, 9, 0, 0, 0, 0, time.UTC).Sub(time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC)),
		},
		{
			name: "01-01-2023T18:00:00+00:00 (no holidays expected)",
			args: args{
				input: time.Date(2023, 1, 1, 18, 0, 0, 0, time.UTC),
			},
			want: time.Date(2023, 1, 10, 0, 0, 0, 0, time.UTC).Sub(time.Date(2023, 1, 1, 18, 0, 0, 0, time.UTC)),
		},
		{
			name: "20-03-2024T13:00:00+00:00",
			args: args{
				input: time.Date(2024, 3, 20, 13, 0, 0, 0, time.UTC),
			},
			want: time.Date(2024, 4, 4, 0, 0, 0, 0, time.UTC).Sub(time.Date(2024, 3, 20, 13, 0, 0, 0, time.UTC)),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := CalculateCountdown(tt.args.input); *got != tt.want {
				t.Errorf("CalculateCountdown() = %v, want %v", got, tt.want)
			}
		})
	}
}
