package gameserver

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/DerWaldFried/gameSclick/src/gameSclick/system"

	"github.com/pterm/pterm"
)

type MinecraftServer struct {
	Type    string
	Version string
	BaseDir string // The root directory for all game servers
}

// API Structure for fetching PaperMC versions and builds
type PaperVersions struct {
	Versions []string `json:"versions"`
}

type PaperBuilds struct {
	Builds []int `json:"builds"`
}

type PaperBuildInfo struct {
	Downloads struct {
		Application struct {
			Name string `json:"name"`
		} `json:"application"`
	} `json:"downloads"`
}

type FabricVersion struct {
	Version string `json:"version"`
	Stable  bool   `json:"stable"`
}

func (m *MinecraftServer) GetName() string {
	return "Minecraft"
}

func (m *MinecraftServer) GetPath() string {
	return m.BaseDir
}

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

// Install now takes the server username to build the correct path
func (m *MinecraftServer) Install(serverUsername string) error {
	pterm.DefaultSection.Println("Minecraft Installation Setup")

	if !system.CheckJavaInstallation() {
		pterm.Warning.Println("Java is not installed. Minecraft requires Java to run.")
		pterm.Info.Println("Attempting to install Java...")
		if err := system.InstallJava(); err != nil {
			return fmt.Errorf("failed to install Java: %v", err)
		}
		return fmt.Errorf("java is installed. And Operation will now continue...")
	} else {
		if !system.IsJavaFunctional() {
			return fmt.Errorf("java is installed but not functional")
		}
	}

	clearTerminal()
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

// Helperfunction to download a file from a URL and save it to a specified path
func (m *MinecraftServer) downloadFile(filepath string, url string) error {
	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()

	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("bad status: %s", resp.Status)
	}

	_, err = io.Copy(out, resp.Body)
	return err
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
	pterm.Info.Println("Abfrage der neuesten PaperMC Informationen...")
	apiBase := "https://api.papermc.io/v2/projects/paper"

	// 1. Get latest version
	resp, err := http.Get(apiBase)
	if err != nil {
		return err
	}
	var vData PaperVersions
	json.NewDecoder(resp.Body).Decode(&vData)
	latestVersion := vData.Versions[len(vData.Versions)-1]

	// 2. Get Latest Build
	resp, err = http.Get(fmt.Sprintf("%s/versions/%s", apiBase, latestVersion))
	if err != nil {
		return err
	}
	var bData PaperBuilds
	json.NewDecoder(resp.Body).Decode(&bData)
	latestBuild := bData.Builds[len(bData.Builds)-1]

	// 3. Call Build Info to get file name
	resp, err = http.Get(fmt.Sprintf("%s/versions/%s/builds/%d", apiBase, latestVersion, latestBuild))
	if err != nil {
		return err
	}
	var biData PaperBuildInfo
	json.NewDecoder(resp.Body).Decode(&biData)
	fileName := biData.Downloads.Application.Name

	// 4. Download URL and save to disk
	downloadUrl := fmt.Sprintf("%s/versions/%s/builds/%d/downloads/%s", apiBase, latestVersion, latestBuild, fileName)
	jarPath := filepath.Join(m.BaseDir, "paper.jar") // Wir nennen es konsistent paper.jar

	pterm.Info.Printfln("Downloade Paper %s (Build %d)...", latestVersion, latestBuild)

	return m.downloadFile(jarPath, downloadUrl)
}

func (m *MinecraftServer) installSpigot() error {
	pterm.Info.Println("Vorbereitung: Spigot BuildTools...")
	btUrl := "https://hub.spigotmc.org/jenkins/job/BuildTools/lastSuccessfulBuild/artifact/target/BuildTools.jar"
	btPath := filepath.Join(m.BaseDir, "BuildTools.jar")

	if err := m.downloadFile(btPath, btUrl); err != nil {
		return err
	}

	pterm.Info.Println("Starte Kompilierung (Spigot 'latest'). Dies kann 5-10 Minuten dauern...")

	// Wir führen BuildTools aus. --rev latest baut die aktuellste stabile Version.
	cmd := exec.Command("java", "-jar", "BuildTools.jar", "--rev", "latest")
	cmd.Dir = m.BaseDir
	cmd.Stdout = os.Stdout // Damit der User den Fortschritt sieht
	cmd.Stderr = os.Stderr

	return cmd.Run()
}

func (m *MinecraftServer) installFabric() error {
	pterm.Info.Println("Ermittle neueste Fabric-Komponenten...")

	// 1. Neueste stabile MC-Version holen
	var mcVersions []FabricVersion
	resp, _ := http.Get("https://meta.fabricmc.net/v2/versions/game")
	json.NewDecoder(resp.Body).Decode(&mcVersions)

	var latestStableMC string
	for _, v := range mcVersions {
		if v.Stable {
			latestStableMC = v.Version
			break
		}
	}

	// 2. Neuesten Loader holen
	var loaders []FabricVersion
	resp, _ = http.Get("https://meta.fabricmc.net/v2/versions/loader")
	json.NewDecoder(resp.Body).Decode(&loaders)
	latestLoader := loaders[0].Version

	// 3. Download der "Server-Launch-Jar"
	// URL-Schema: https://meta.fabricmc.net/v2/versions/loader/<game>/<loader>/<installer>/server/jar
	// Wir nutzen hier den Standard-Installer 1.0.1
	downloadUrl := fmt.Sprintf("https://meta.fabricmc.net/v2/versions/loader/%s/%s/1.0.1/server/jar", latestStableMC, latestLoader)
	jarPath := filepath.Join(m.BaseDir, "fabric-server.jar")

	pterm.Info.Printfln("Downloade Fabric für MC %s (Loader %s)...", latestStableMC, latestLoader)
	return m.downloadFile(jarPath, downloadUrl)
}

func (m *MinecraftServer) installForge() error {
	// Forge ist leider schwer voll-automatisch zu finden, daher nutzen wir hier
	// eine gängige stabile Version oder du müsstest eine Forge-API-Library nutzen.
	mcVersion := "1.20.1"
	forgeVersion := "47.2.0"

	pterm.Info.Printfln("Downloade Forge Installer für %s...", mcVersion)

	downloadUrl := fmt.Sprintf("https://maven.minecraftforge.net/net/minecraftforge/forge/%s-%s/forge-%s-%s-installer.jar",
		mcVersion, forgeVersion, mcVersion, forgeVersion)
	installerPath := filepath.Join(m.BaseDir, "forge-installer.jar")

	if err := m.downloadFile(installerPath, downloadUrl); err != nil {
		return err
	}

	pterm.Info.Println("Installiere Forge Server (headless)...")

	// Forge muss installiert werden, um die eigentliche Server-Jar und Libraries zu generieren
	cmd := exec.Command("java", "-jar", "forge-installer.jar", "--installServer")
	cmd.Dir = m.BaseDir

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("forge installation fehlgeschlagen: %v", err)
	}

	// Installer aufräumen
	os.Remove(installerPath)
	return nil
}
