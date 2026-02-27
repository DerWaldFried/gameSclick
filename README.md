# 🎮 gameSclick

**gameSclick** is a lightweight, interactive CLI tool written in Go, designed to automate the deployment of game servers on Linux. It handles the heavy lifting of system configuration, user creation, and dependency management with a single click.
### ATTENTION: Not a ready to use Build!
---

## 🚀 Features

* **Smart OS Detection:** Automatically identifies your Linux distribution (optimized for Ubuntu/Debian).
* **Automated User Management:** Creates a dedicated system user with restricted privileges to ensure your server stays secure.
* **Seamless SteamCMD Setup:** Installs SteamCMD, handles 32-bit architecture requirements, and auto-accepts EULAs—no manual intervention required.
* **Interactive Configuration:** User-friendly setup wizard powered by [pterm](https://github.com/pterm/pterm).
* **Persistent TOML Settings:** Saves your preferences according to XDG standards in a human-readable format.

---

## 📋 Prerequisites

* **OS:** Ubuntu or Debian.
* **Privileges:** Root or Sudo access (the tool verifies this on startup).
* **Go:** Version 1.21+ (if building from source).

---

## 🛠️ Installation

### Build from Source
```bash
git clone [https://github.com/YOUR_USERNAME/gameSclick.git](https://github.com/YOUR_USERNAME/gameSclick.git)
cd gameSclick
go build -o gamesclick
