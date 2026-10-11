package hosts

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"
)

func AddDomain(domain string) error {
	var hostsPath string
	if runtime.GOOS == "windows" {
		hostsPath = "C:\\Windows\\System32\\drivers\\etc\\hosts"
	} else {
		hostsPath = "/etc/hosts"
	}

	content, err := os.ReadFile(hostsPath)
	if err != nil {
		return fmt.Errorf("failed to read hosts file: %v", err)
	}

	entry := fmt.Sprintf("127.0.0.1 %s", domain)
	if strings.Contains(string(content), domain) {
		fmt.Printf("Domain %s already exists in hosts file.\n", domain)
		return nil
	}

	fmt.Printf("Adding %s to hosts file (requires sudo/admin)...\n", domain)
	
	if runtime.GOOS == "windows" {
		cmd := exec.Command("powershell", "-Command", fmt.Sprintf("Add-Content -Path '%s' -Value '`n%s'", hostsPath, entry))
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		return cmd.Run()
	} else {
		// Use echo and sudo tee for Unix systems
		cmdStr := fmt.Sprintf("echo '%s' | sudo tee -a %s > /dev/null", entry, hostsPath)
		cmd := exec.Command("sh", "-c", cmdStr)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		return cmd.Run()
	}
}