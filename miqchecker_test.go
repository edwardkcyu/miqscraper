package main

import "testing"

func TestMiqChecker_prepareSlackMessage(t *testing.T) {

	type args struct {
		availableDates []string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "has Tuesday",
			args: args{
				availableDates: []string{"2021-09-14"},
			},
			want:    ":white_check_mark: 14/09 (Tue)",
			wantErr: false,
		},
		{
			name: "no Tuesday",
			args: args{
				availableDates: []string{"2021-08-11"},
			},
			want:    ":eyes: 11/08 (Wed)",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := NewMiqChecker(nil, nil)
			got, err := m.prepareSlackMessage(tt.args.availableDates)
			if (err != nil) != tt.wantErr {
				t.Errorf("prepareSlackMessage() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("prepareSlackMessage() got = %v, want %v", got, tt.want)
			}
		})
	}
}
