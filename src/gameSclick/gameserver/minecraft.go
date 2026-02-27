package gameserver

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

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
