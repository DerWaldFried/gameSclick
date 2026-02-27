package system

import "os/exec"

func CheckJavaInstallation() bool {
	// LookPath searching for Java executable.
	_, err := exec.LookPath("java")
	return err == nil
}

func IsJavaFunctional() bool {
	// Führt 'java -version' aus. Wir ignorieren den Output und prüfen nur den Exit-Code.
	cmd := exec.Command("java", "-version")
	if err := cmd.Run(); err != nil {
		return false
	}
	return true
}
