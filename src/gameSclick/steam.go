package main

import (
	"fmt"
	"os/exec"

	"github.com/pterm/pterm"
)

func checkSteamInstallation() bool {
	// Check, if steamcmd is already installed by looking for its executable in the system's PATH. If it is found, we assume it's installed and return true. If not, we return false, indicating that the installation process needs to be initiated.
	_, err := exec.LookPath("steamcmd")
	return err == nil
}

func startSteamInstallation() {
	pterm.DefaultHeader.WithFullWidth().Println("gameSclick - Steam Installation")

	if checkSteamInstallation() {
		pterm.Success.Println("SteamCMD is already installed.")
		return
	}

	distro := getDistroID()
	if distro != "ubuntu" && distro != "debian" {
		pterm.Warning.Printfln("Unsupported Distro: %s. This tool is optimized for Ubuntu/Debian.", distro)
		return
	}

	// 1. EULA automatisation (Debconf)
	pterm.Info.Println("Accepting SteamCMD EULA...")

	// We build a shell command that uses debconf-set-selections to pre-accept the Steam license agreement. This prevents the installation process from halting to ask for user input.
	licenseCmd := `echo steam steam/question select I AGREE | sudo debconf-set-selections && 
                   echo steam steam/license note '' | sudo debconf-set-selections`

	exec.Command("bash", "-c", licenseCmd).Run()

	// 2. Installation prepare (Multiarch + Update + Install)
	pterm.Info.Println("Enabling i386 architecture and installing steamcmd...")

	// 'DEBIAN_FRONTEND=noninteractive' stops apt-get from asking for user input during installation, which is crucial for automation.
	fullCmd := `sudo dpkg --add-architecture i386 && 
                sudo apt-get update && 
                export DEBIAN_FRONTEND=noninteractive && 
                sudo apt-get install -y steamcmd`

	// Spinner for Operation Feedback
	spinner, _ := pterm.DefaultSpinner.Start("Downloading and installing SteamCMD...")

	cmd := exec.Command("bash", "-c", fullCmd)

	// When we want see Errors, we can do cmd.Stderr = os.Stderr
	err := cmd.Run()

	if err != nil {
		spinner.Fail(fmt.Sprintf("Failed to install SteamCMD: %v", err))
		return
	}

	spinner.Success("SteamCMD installed successfully! Going back to main menu...")
}
