package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/pterm/pterm"
)

// Check the Distro to determine the package manager and installation commands
func getDistroID() string {
	file, err := os.Open("/etc/os-release")
	if err != nil {
		return "unknown"
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "ID=") {
			return strings.Trim(strings.TrimPrefix(line, "ID="), "\"")
		}
	}
	return "unknown"
}

func (cfg *Config) addUser() error {
	pterm.Info.Printfln("Creating system user '%s'...", cfg.Username)

	// --disabled-password stops the system from asking for a password during user creation, we will set it later
	// --gecos "" passes empty information for the user (like full name, room number, etc.) to avoid interactive prompts
	cmd := exec.Command("sudo", "adduser", "--disabled-password", "--gecos", "", cfg.Username)

	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("could not create user: %v", err)
	}

	pterm.Success.Printfln("User '%s' created successfully.", cfg.Username)
	return nil
}
