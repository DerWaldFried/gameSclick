package main

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"

	"github.com/pterm/pterm"
)

func clearTerminal() {
	// We clear the terminal screen before showing the main menu to provide a clean and organized interface for the user.
	// This enhances readability and user experience by removing any previous output or clutter from the terminal.
	var cmd *exec.Cmd
	if runtime.GOOS == "windows" {
		cmd = exec.Command("cls")
	} else {
		cmd = exec.Command("clear")
	}
	cmd.Stdout = os.Stdout
	cmd.Run()
}

func showMainMenu(cfg Config) {
	for {
		clearTerminal()
		header()

		pterm.DefaultSection.Println("Main Menu")

		// Status Box with User and Steam Support Info
		pterm.DefaultBox.WithTitle("System Status").Printfln(
			"User: %s | Steam Support: %v",
			pterm.LightCyan(cfg.Username),
			pterm.LightGreen(cfg.UseSteam),
		)

		// Menü-Options
		options := []string{
			"📦 Install New Game Server",
			"📋 List Installed Servers",
			"⚙️  Edit Settings",
			"🛠️  Run System Diagnosis",
			"🚪 Exit",
		}

		selectedOption, _ := pterm.DefaultInteractiveSelect.
			WithDefaultText("Please choose an action").
			WithOptions(options).
			Show()

		switch selectedOption {
		case "📦 Install New Game Server":
			// Added later: Entry Point for Game Selection and Installation Process
			pterm.Info.Println("Opening Game Selection...")
			// Later Implement: showGameSelection()

		case "📋 List Installed Servers":
			pterm.Info.Println("Fetching installed servers...")
			// Later Implement: listServers()

		case "⚙️  Edit Settings":
			pterm.Info.Println("Reloading Config Wizard...")
			fillConfig() // We call fillConfig again to allow the user to edit their settings. This will overwrite the existing config with the new values.

		case "🛠️  Run System Diagnosis":
			runDiagnosis()

		case "🚪 Exit":
			pterm.Info.Println("Goodbye!")
			os.Exit(0)
		}

		// Break. To make it human readable
		fmt.Scanln()
	}
}

func runDiagnosis() {
	spinner, _ := pterm.DefaultSpinner.Start("Checking System...")

	// First Example Check: Get Distro and SteamCMD Installation Status
	distro := getDistroID()
	hasSteam := checkSteamInstallation()

	spinner.Success("Diagnosis complete!")

	pterm.DefaultTable.WithData(pterm.TableData{
		{"Check", "Status"},
		{"OS Distro", distro},
		{"SteamCMD", fmt.Sprintf("%v", hasSteam)},
		{"Root Privileges", fmt.Sprintf("%v", os.Geteuid() == 0)},
	}).Render()

	pterm.Println("\nPress Enter to return to menu...")
}
