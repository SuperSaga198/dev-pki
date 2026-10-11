//go:build windows

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
	certPath := filepath.Join(dir, "rootCA.pem")

	fmt.Println("Adding Root CA to Windows LocalMachine Root store (requires Admin)...")
	cmd := exec.Command("powershell", "-Command", fmt.Sprintf("Import-Certificate -FilePath '%s' -CertStoreLocation Cert:\\LocalMachine\\Root", certPath))
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}