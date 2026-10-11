//go:build darwin

package trust

import (
	"fmt"
	"os/exec"
	"path/filepath"
	"github.com/SuperSaga198/dev-pki/internal/pki"
)

func InstallTrustStore() error {
	dir, err := pki.GetConfigDir()
	if err != nil {
		return err
	}
	certPath := filepath.Join(dir, "rootCA.pem")

	fmt.Println("Adding Root CA to macOS Keychain (requires sudo)...")
	cmd := exec.Command("sudo", "security", "add-trusted-cert", "-d", "-r", "trustRoot", "-k", "/Library/Keychains/System.keychain", certPath)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}