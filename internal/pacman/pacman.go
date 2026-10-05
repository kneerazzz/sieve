package pacman

import (
	"bufio"
	"errors"
	"fmt"
	"os/exec"
	"strings"
)

type Package struct {
	Name    string
	Version string
}

func ListPackages() ([]Package, error) {
	out, err := exec.Command("pacman", "-Qm").Output()

	if err != nil {
		//pacman exits with code 1 when there are no foreign packages
		var exitError *exec.ExitError
		if errors.As(err, &exitError) && exitError.ExitCode() == 1 && len(out) == 0 {
			return nil, nil
		}
		return nil, fmt.Errorf("Running pacman -Qm %w", err)
	}
	return parse(string(out))
}

// parse will turn the "name version" into the Packages vaules
func parse(output string) ([]Package, error) {

	var pkgs []Package
	scanner := bufio.NewScanner(strings.NewReader(output))
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		fields := strings.Fields(line)
		if len(fields) != 2 {
			fmt.Errorf("Unexpected pacman output line: %q", line)
		}
		pkgs = append(pkgs, Package{Name: fields[0], Version: fields[1]})
	}
	return pkgs, scanner.Err()
}
