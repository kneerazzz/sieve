package pacman

import (
	"reflect"
	"testing"
)


func TestParse(t *testing.T) {
	tests := []struct {
		name string
		input string
		want []Package
		wantErr bool
	} {
		{
			name: "two packages",
			input: "yay 12.3.5\n brave-browser 0.3.5\n",
			want: []Package {
				{Name: "yay", Version: "12.3.5"},
				{Name: "brave-browser", Version: "0.3.5"}
			},
		},
		{
			name: "empty",
			input: "",
			want: nil
		},
		{
			name: "malformed line",
			input: "justoneline\n",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.run(tt.name, func(t *testing.T) {
			got, err := parse(tt.input)
			if (err != nil) != tt.wantErr {
				t.Fatalf("err: %v, wantErr: %v", err, tt.wantErr)
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("got: %v, want: %v", got, tt.want)
			}
		})
	}
}
