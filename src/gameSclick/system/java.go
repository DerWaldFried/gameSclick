package system

import (
	"os/exec"

	"github.com/pterm/pterm"
)

func CheckJavaInstallation() bool {
	// LookPath searching for Java executable.
	_, err := exec.LookPath("java")
	return err == nil
}

func IsJavaFunctional() bool {
	// Usage java --version to check if Java is functional. If the command executes successfully, we consider Java to be functional.
	// If it returns an error, we assume there is an issue with the Java installation.
	cmd := exec.Command("java", "-version")
	if err := cmd.Run(); err != nil {
		return false
	}
	return true
}

func InstallJava() error {
	// Installing Java using the system's package manager. This example assumes a Debian-based system with apt-get.
	// For other distributions, the command would need to be adjusted accordingly. We only Support Debian/Ubuntu for now, as we check the Distro before.
	cmd := exec.Command("sudo", "apt-get", "update")
	if err := cmd.Run(); err != nil {
		return err
	}

	cmd = exec.Command("sudo", "apt-get", "install", "-y", "default-jre")
	if err := cmd.Run(); err != nil {
		return err
	}

	pterm.Success.Printfln("Java installed successfully.")
	return nil
}
