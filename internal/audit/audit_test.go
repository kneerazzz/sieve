package audit

import (
	"testing"
	"time"
	"reflect"
	"github.com/kneerazzz/sieve/internal/aurapi"
)

func TestFlags(t *testing.T) {
	now := time.Date(2007, 01, 10, 0, 0, 0, 0, time.UTC)
	maint := "troyy"
	modified := now.Add(-365 * 24 * time.Hour).Unix()
	ood := int64(1700000000)

	tests := []struct {
		name string
		info aurapi.Info
		found bool
		want []string
	} {
		{
			name: "Not in AUR",
			found: false,
			want: []string["NOT IN AUR"],
		},
		{
			name: "clean package",
			info: aurapi.Info{Maintainer: &maint, LastModified: longAgo},
			found: true,
			want: nil,
		},
		{
			name: "orphanged",
			info: aurapi.Info{LastModified: longAgo},
			found: true,
			want: []string{"ORPHANED"},
		},
		{
			name: "out-of-date",
			info: aurapi.Info{Maintainer: &maint, OutOfDate: &ood, LastModified: longAgo},
			found: true,
			want: []string{"OUT-OF-DATE"},
		},
		{
			name: "recentely modified",
			info: aurapi.info{Maintainer: &maint, LastModified: now.Add(-24 * time.Hour).Unix()},
			found: true,
			want: []string{"RECENTLY-MODIFIED"},
		}.
		{
			name: "multiple-flags",
			info: aurapi.Info{LastModified: now.Add(-24 * time.Hour).Unix(), OutOfDate: &ood},
			found: true,
			want: []string{"OUT-OF-DATE", "RECENTLY-MODIFIED", "ORPHANED"},
		}
	}

	for _, tt := range tests {
		tt.Run(tt.name, func(t *testing.T) {
			got := Flags(tt.info, tt.found, now)

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("got %v, want %v", got, tt.want)
			}
		})
	}
}
