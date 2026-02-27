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
	// Check again the System for the user, because if we do it before, we can have a race condition where the user is created after the
	// check but before the creation, which would lead to an error. By checking again right before the creation, we can avoid this issue.
	_, err := exec.LookPath("id")
	if err == nil {
		checkUser := exec.Command("id", "-u", cfg.Username)
		if err := checkUser.Run(); err == nil {
			return fmt.Errorf("user '%s' already exists, skipping creation", cfg.Username)
		}
	}

	pterm.Info.Printfln("Creating system user '%s'...", cfg.Username)

	// Adding User with no password and no shell access (disabled login)
	cmd := exec.Command("sudo", "adduser", "--disabled-password", "--gecos", "", cfg.Username)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("could not create user: %v", err)
	}

	// Set Passwort (Send with Pipe on chpasswd)
	passCmd := exec.Command("sudo", "chpasswd")
	passCmd.Stdin = strings.NewReader(fmt.Sprintf("%s:%s", cfg.Username, cfg.Userpasswd))
	if err := passCmd.Run(); err != nil {
		return fmt.Errorf("could not set password for user: %v", err)
	}

	pterm.Success.Printfln("User '%s' created and password set.", cfg.Username)
	return nil
}
