package cmd

import (
	"testing"
	"time"
)

func TestHumanDuration(t *testing.T) {
	t.Parallel()

	type args struct {
		start time.Time
		now   time.Time
	}

	cases := map[string]struct {
		args args
		want string
	}{
		"not yet": {
			args{
				start: time.Date(2022, 1, 1, 9, 0, 0, 0, time.UTC),
				now:   time.Date(2020, 1, 1, 9, 0, 0, 0, time.UTC),
			},
			"not yet",
		},
		"just arrived": {
			args{
				start: time.Date(2022, 1, 1, 9, 0, 0, 0, time.UTC),
				now:   time.Date(2022, 1, 1, 9, 0, 0, 0, time.UTC),
			},
			"just arrived",
		},
		"simple": {
			args{
				start: time.Date(2020, 1, 1, 9, 0, 0, 0, time.UTC),
				now:   time.Date(2021, 1, 1, 9, 0, 0, 0, time.UTC),
			},
			"1 year",
		},
		"plural": {
			args{
				start: time.Date(2020, 1, 1, 9, 0, 0, 0, time.UTC),
				now:   time.Date(2022, 1, 1, 9, 0, 0, 0, time.UTC),
			},
			"2 years",
		},
		"composed": {
			args{
				start: time.Date(2020, 1, 1, 9, 0, 0, 0, time.UTC),
				now:   time.Date(2022, 6, 1, 9, 0, 0, 0, time.UTC),
			},
			"2 years, 5 months",
		},
		"lower now": {
			args{
				start: time.Date(2021, 12, 1, 9, 0, 0, 0, time.UTC),
				now:   time.Date(2022, 1, 1, 9, 0, 0, 0, time.UTC),
			},
			"1 month",
		},
		"lower year": {
			args{
				start: time.Date(2020, 10, 1, 9, 0, 0, 0, time.UTC),
				now:   time.Date(2022, 2, 1, 9, 0, 0, 0, time.UTC),
			},
			"1 year, 4 months",
		},
	}

	for intention, testCase := range cases {
		intention, testCase := intention, testCase

		t.Run(intention, func(t *testing.T) {
			t.Parallel()

			if got := humanDuration(testCase.args.start, testCase.args.now); got != testCase.want {
				t.Errorf("humanDuration() = `%s`, want `%s`", got, testCase.want)
			}
		})
	}
}
