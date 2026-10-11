package cmd

import (
	"flag"
	"fmt"
	"os"

	"github.com/SuperSaga198/dev-pki/internal/hosts"
	"github.com/SuperSaga198/dev-pki/internal/pki"
	"github.com/SuperSaga198/dev-pki/internal/trust"
)

func Execute() {
	initCmd := flag.NewFlagSet("init", flag.ExitOnError)
	issueCmd := flag.NewFlagSet("issue", flag.ExitOnError)
	domainFlag := issueCmd.String("domain", "", "Domain name to issue certificate for (e.g., myapp.test)")

	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}

	switch os.Args[1] {
	case "init":
		initCmd.Parse(os.Args[2:])
		fmt.Println("Initializing Local Root CA...")
		_, _, err := pki.GenerateRootCA()
		if err != nil {
			fmt.Printf("Error generating CA: %v\n", err)
			os.Exit(1)
		}

		err = trust.InstallTrustStore()
		if err != nil {
			fmt.Printf("Error installing trust store: %v\n", err)
			os.Exit(1)
		}
		fmt.Println("SUCCESS: Local Root CA generated and trusted successfully!")

	case "issue":
		issueCmd.Parse(os.Args[2:])
		if *domainFlag == "" {
			fmt.Println("Error: --domain flag is required for issue command.")
			issueCmd.Usage()
			os.Exit(1)
		}

		if !pki.CAExists() {
			fmt.Println("Error: Local Root CA not found. Please run 'dev-pki init' first.")
			os.Exit(1)
		}

		fmt.Printf("Issuing certificate for %s...\n", *domainFlag)
		certFile, keyFile, err := pki.IssueCertificate(*domainFlag)
		if err != nil {
			fmt.Printf("Error issuing certificate: %v\n", err)
			os.Exit(1)
		}

		err = hosts.AddDomain(*domainFlag)
		if err != nil {
			fmt.Printf("Warning: Failed to update hosts file: %v\n", err)
		}

		fmt.Printf("SUCCESS:\n  Certificate: %s\n  Private Key: %s\n  URL: https://%s\n", certFile, keyFile, *domainFlag)

	default:
		printUsage()
		os.Exit(1)
	}
}

func printUsage() {
	fmt.Println("dev-pki - Local CA & Identity Provisioner for Devs")
	fmt.Println("\nUsage:")
	fmt.Println("  dev-pki init                Initialize local Root CA and trust it on the OS")
	fmt.Println("  dev-pki issue --domain <d>  Issue a signed wildcard cert for a domain and update hosts")
}