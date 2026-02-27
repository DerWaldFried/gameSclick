package gameserver

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/DerWaldFried/gameSclick/src/gameSclick/system"

	"github.com/pterm/pterm"
)

type MinecraftServer struct {
	Type    string
	Version string
	BaseDir string // The root directory for all game servers
}

func (m *MinecraftServer) GetName() string {
	return "Minecraft"
}

func (m *MinecraftServer) GetPath() string {
	return m.BaseDir
}

// Install now takes the server username to build the correct path
func (m *MinecraftServer) Install(serverUsername string) error {
	pterm.DefaultSection.Println("Minecraft Installation Setup")

	if !system.CheckJavaInstallation() {
		return fmt.Errorf("java is not installed")
	} else {
		if !system.IsJavaFunctional() {
			return fmt.Errorf("java is installed but not functional")
		}
	}

	selected, _ := pterm.DefaultInteractiveSelect.
		WithDefaultText("Choose the Minecraft System you want to install").
		WithOptions([]string{"Vanilla", "Paper", "Spigot", "Fabric", "Forge"}).
		Show()

	m.Type = selected
	// In the future, you could add another selection for the Minecraft Version (e.g., 1.20.4)
	m.Version = "latest"

	// Build path: /home/gsuser/server/minecraft/<type>
	m.BaseDir = filepath.Join("/home", serverUsername, "server", "minecraft", strings.ToLower(m.Type))

	pterm.Info.Printfln("Installing %s Minecraft to: %s", m.Type, m.BaseDir)

	// Create directory with correct permissions
	err := os.MkdirAll(m.BaseDir, 0755)
	if err != nil {
		return fmt.Errorf("failed to create directory: %v", err)
	}

	switch m.Type {
	case "Vanilla":
		return m.installVanilla()
	case "Paper":
		return m.installPaper()
	case "Spigot":
		return m.installSpigot()
	case "Fabric":
		return m.installFabric()
	case "Forge":
		return m.installForge()
	default:
		return fmt.Errorf("unsupported Minecraft type: %s", m.Type)
	}
}

func (m *MinecraftServer) installVanilla() error {
	pterm.Info.Println("Starting Vanilla Minecraft Installation...")

	// Path for the server jar, e.g., /home/gsuser/server/minecraft/vanilla/minecraft_server.jar
	jarPath := filepath.Join(m.BaseDir, "minecraft_server.jar")

	// Prepare Download-Command
	downloadCmd := exec.Command("wget", "-O", jarPath, "https://piston-data.mojang.com/v1/objects/64bb6d763bed0a9f1d632ec347938594144943ed/server.jar")

	// Optional: Usage BaseDir as working directory for the command, so that the jar is downloaded directly to the correct location.
	// This can help avoid issues with relative paths and ensure that the file ends up in the intended directory.
	downloadCmd.Dir = m.BaseDir

	if err := downloadCmd.Run(); err != nil {
		return fmt.Errorf("failed to download Minecraft server jar: %v", err)
	}

	return nil
}

func (m *MinecraftServer) installPaper() error {
	pterm.Info.Println("Starting Paper Minecraft Installation...")

	return nil
}

func (m *MinecraftServer) installSpigot() error {
	pterm.Info.Println("Starting Spigot Minecraft Installation...")

	return nil
}

func (m *MinecraftServer) installFabric() error {
	pterm.Info.Println("Starting Fabric Minecraft Installation...")

	return nil
}

func (m *MinecraftServer) installForge() error {
	pterm.Info.Println("Starting Forge Minecraft Installation...")

	return nil
}
