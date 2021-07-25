package main

import (
	"reflect"
	"testing"
)

func TestMiqManager_fetchAvailableDates(t *testing.T) {
	type fields struct {
		url string
	}
	tests := []struct {
		name    string
		fields  fields
		want    []string
		wantErr bool
	}{
		{
			name:   "query miq portal",
			fields: fields{url: "https://allocation.miq.govt.nz/portal/"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := NewMiqManager(tt.fields.url)
			got, err := m.fetchAvailableDates()
			if (err != nil) != tt.wantErr {
				t.Errorf("fetchAvailableDates() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("fetchAvailableDates() got = %v, want %v", got, tt.want)
			}
		})
	}
}
