//go:build linux

package trust

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"github.com/SuperSaga198/dev-pki/internal/pki"
)

func InstallTrustStore() error {
	dir, err := pki.GetConfigDir()
	if err != nil {
		return err
	}
	srcPath := filepath.Join(dir, "rootCA.pem")
	destPath := "/usr/local/share/ca-certificates/dev-pki-root.crt"

	fmt.Println("Copying Root CA and updating certificates on Linux (requires sudo)...")
	
	// Copy file via bash/sudo
	cpCmd := exec.Command("sudo", "cp", srcPath, destPath)
	cpCmd.Stdout = os.Stdout
	cpCmd.Stderr = os.Stderr
	if err := cpCmd.Run(); err != nil {
		return err
	}

	updateCmd := exec.Command("sudo", "update-ca-certificates")
	updateCmd.Stdout = os.Stdout
	updateCmd.Stderr = os.Stderr
	return updateCmd.Run()
}