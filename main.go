package main

import (
	"time"

	"github.com/kneerazzz/sieve/internal/audit"
	"github.com/kneerazzz/sieve/internal/aurapi"
	"github.com/kneerazzz/sieve/internal/pacman"
	"github.com/spf13/cobra"
)

func main() {
	root := &cobra.Command{
		Use:   "sieve",
		Short: "Audit AUR packages for suspicious changes",
	}
	auditCmd := &cobra.Command{
		Use:   "audit",
		Short: "List installed AUR packages.",
		RunE: func(cmd *cobra.Command, args []string) error {
			pkgs, err := pacman.ListPackages()
			if err != nil {
				return err
			}
			names := make([]string, len(pkgs))
			for i, p := range pkgs {
				names[i] = p.Name
			}
			infos, err := aurapi.New().Info(cmd.Context(), names)

			if err != nil {
				return err
			}

			now := time.Now()

			for _, p := range pkgs {
				info, found := infos[p.Name]
				maint := "-"
				modified := "-"

				if found {
					if info.Maintainer != nil {
						maint = *info.Maintainer
					} else {
						maint = "(none"
					}
					modified = time.Unix(info.LastModified, 0).Format("2026-01-02")
				}
				cmd.Printf("%-35s %-16s %-12s %v\n", p.Name, maint, modified, audit.Flags(info, found, now))
			}
			return nil
		},
	}
	root.AddCommand(auditCmd)
	root.Execute()
}
