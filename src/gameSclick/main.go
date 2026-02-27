package main

import (
	"os"
	"os/exec"
	"path/filepath"

	"github.com/BurntSushi/toml"
	"github.com/pterm/pterm"
)

type Config struct {
	Username      string `toml:"username"`
	Userpasswd    string `toml:"userpasswd"`
	SteamUsername string `toml:"steam_username"`
	UseSteam      bool   `toml:"use_steam"`
	AutoSave      bool   `toml:"autosave"`
}

func checkPrivileges() {
	if os.Geteuid() != 0 {
		_, err := exec.LookPath("sudo")
		if err != nil {
			pterm.Error.Println("This script needs root privileges. Please run as root (or install sudo).")
			os.Exit(1)
		}
	}
}

// Helperfunctions for Config Management
func getConfigPath() (string, error) {
	configDir, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(configDir, "gameSclick", "config.toml"), nil
}

func checkConfigExists() {
	configPath, err := getConfigPath()
	if err != nil {
		pterm.Error.Println("Could not find config path:", err)
		return
	}

	var cfg Config
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		pterm.Warning.Println("No Config found. You started first time? Let's set it up!")
		cfg, err = fillConfig()
		if err != nil {
			pterm.Error.Println("Error saving config:", err)
			return
		}
	} else {
		cfg, err = loadConfig()
		if err != nil {
			pterm.Error.Println("Error loading config:", err)
			return
		}
		pterm.Success.Printf("Welcome back, %s!\n", cfg.Username)
	}
}

func loadConfig() (Config, error) {
	configPath, err := getConfigPath()
	if err != nil {
		return Config{}, err
	}

	var cfg Config
	if _, err := toml.DecodeFile(configPath, &cfg); err != nil {
		return Config{}, err
	}
	return cfg, nil
}

func fillConfig() (Config, error) {
	configPath, err := getConfigPath()
	if err != nil {
		return Config{}, err
	}

	// Add Folder for Config if not exists
	os.MkdirAll(filepath.Dir(configPath), 0755)

	pterm.DefaultSection.Println("Setup Your Settings")

	pterm.Description.Println("Type in the Username for the Server User that will be created.")
	rusernameResult, _ := pterm.DefaultInteractiveTextInput.WithDefaultValue("gsuser").Show()

	pterm.Description.Println("Type in the Password for the Server User that will be created.")
	rpasswordResult, _ := pterm.DefaultInteractiveTextInput.WithDefaultValue("******").Show()

	pterm.Description.Println("You want to use Steam Gameservers? (y/n)")
	useSteamResult, _ := pterm.DefaultInteractiveConfirm.WithDefaultValue(true).Show()

	var susernameResult string

	if useSteamResult {
		pterm.Description.Println("Steam Login (type 'anonymous' or your username):")
		susernameResult, _ = pterm.DefaultInteractiveTextInput.WithDefaultValue("anonymous").Show()
	} else {
		susernameResult = "none" // Value when Steam is not used, can be ignored in the rest of the code
	}

	cfg := Config{
		Username:      rusernameResult,
		Userpasswd:    rpasswordResult,
		SteamUsername: susernameResult,
		UseSteam:      useSteamResult,
		AutoSave:      true,
	}

	// Write File
	f, err := os.Create(configPath)
	if err != nil {
		return Config{}, err
	}
	defer f.Close()

	if err := toml.NewEncoder(f).Encode(cfg); err != nil {
		return Config{}, err
	}

	pterm.Success.Println("Config saved to:", configPath)
	return cfg, nil
}

func header() {
	pterm.DefaultHeader.WithFullWidth().Println("gameSclick")
	pterm.DefaultCenter.Println("Your simple CLI to install GameServer on your Server.")
}

func main() {
	// Load Header
	header()
	// Check if the user has root privileges, if not, exit with an error message.
	checkPrivileges()
	// Check if config exists, if not create it. After that load the config and work with it.
	checkConfigExists()
	cfg, err := loadConfig()
	if err != nil {
		pterm.Error.Println("Error loading config:", err)
		return
	}

	// --- SYSTEM SETUP ---
	// Add SystemUser (using cfg.Username)
	err = cfg.addUser()
	if err != nil {
		pterm.Warning.Println(err) // When the user already exists, we can ignore the error and continue with the installation.
	}

	// --- STEAM SETUP ---
	if cfg.UseSteam {
		startSteamInstallation()
	}

	// On this zone you can work with the loaded config. We go directly to main menu.
	showMainMenu(cfg)
}
