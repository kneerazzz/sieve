package main

import (
	"os"

	"github.com/kneerazzz/sieve/internal/pacman"
	"github.com/spf13/cobra"
)

func main() {
	root := &cobra.Command{
		Use:   "sieve",
		Short: "Audit AUR packages for suspicious changes",
	}

	root.AddCommand(&cobra.Command{
		Use:   "audit",
		Short: "List Installed packages",
		Run: func(cmd *cobra.Command, args []string) {
			pkgs, err := pacman.ListPackages()

			if err != nil {
				cmd.PrintErrln("error: ", err)
				os.Exit(1)
			}
			for _, p := range pkgs {
				cmd.Printf("%-40s %s\n", p.Name, p.Version)
			}
		},
	})
	root.Execute()
}
