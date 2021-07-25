package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMiqChecker_prepareSlackMessage(t *testing.T) {

	type args struct {
		availableDates []string
	}
	tests := []struct {
		name               string
		args               args
		wantMessage        string
		wantHasTargetDates bool
		wantErr            bool
	}{
		{
			name: "has Tuesday",
			args: args{
				availableDates: []string{"2021-09-14"},
			},
			wantMessage:        ":white_check_mark: 14/09 (Tue)",
			wantHasTargetDates: true,
			wantErr:            false,
		},
		{
			name: "no Tuesday",
			args: args{
				availableDates: []string{"2021-08-11"},
			},
			wantMessage:        ":eyes: 11/08 (Wed)",
			wantHasTargetDates: false,
			wantErr:            false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := NewCheckerTask(nil, nil, "", "")
			gotMessage, gotHasTargetDates, err := m.prepareSlackMessage(tt.args.availableDates)
			if (err != nil) != tt.wantErr {
				t.Errorf("prepareSlackMessage() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			assert.Equal(t, gotMessage, tt.wantMessage, "wrong message")
			assert.Equal(t, gotHasTargetDates, tt.wantHasTargetDates, "wrong hasTargetDates")
		})
	}
}
