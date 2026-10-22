package audit

import (
	"time"

	"github.com/kneerazzz/sieve/internal/aurapi"
)

func Flags(info aurapi.Info, found bool, now time.Time) []string {
	if !found {
		return []string{"NOT IN AUR"}
	}
	var flags []string

	if info.Maintainer == nil {
		flags = append(flags, "ORPHANED")
	}
	if info.OutOfDate != nil {
		flags = append(flags, "OUT-OF-DATE")
	}
	if now.Sub(time.Unix(info.LastModified, 0)) < 30*24*time.Hour {
		flags = append(flags, "RECENTLY-MODIFIED")
	}
	return flags
}
