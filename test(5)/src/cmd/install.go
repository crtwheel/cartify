package cmd

import (
	"os"
	"os/exec"
	"path/filepath"
	"runtime"

	"github.com/crtwheel/cartify/src/utils"
)

func Install(CartifyVersion string) {
	utils.PrintBold("Cartify Installer")
	utils.PrintInfo("Setting up Cartify...")

	InitPaths()
	CheckStates()
	InitSetting()
	Apply(CartifyVersion)

	if runtime.GOOS == "windows" {
		setupScheduledTask()
		utils.PrintInfo("Created scheduled task: Cartify will auto-apply on login.")
	}

	utils.PrintSuccess("Cartify installed! It will automatically re-apply after Spotify updates.")
	utils.PrintInfo("Ad blocking is active. Look for the green shield icon in the top-right of Spotify.")
}

func setupScheduledTask() {
	exe, err := os.Executable()
	if err != nil {
		utils.PrintError("Cannot find cartify executable path")
		return
	}

	taskName := "CartifyAutoApply"
	cmd := exec.Command("schtasks", "/Create",
		"/TN", taskName,
		"/TR", exe+" auto",
		"/SC", "ONLOGON",
		"/DELAY", "0000:00:30",
		"/F",
		"/RL", "HIGHEST",
	)

	if output, err := cmd.CombinedOutput(); err != nil {
		utils.PrintWarning("Could not create scheduled task: " + string(output))
		utils.PrintInfo("You can manually run 'cartify auto' after each Spotify update.")
	} else {
		utils.PrintSuccess("Scheduled task created: cartify will auto-apply at login.")
	}

	startupDir := filepath.Join(os.Getenv("APPDATA"), "Microsoft", "Windows", "Start Menu", "Programs", "Startup")
	shortcutPath := filepath.Join(startupDir, "CartifyAutoApply.bat")
	batchContent := "@echo off\n\"" + exe + "\" auto\n"
	os.WriteFile(shortcutPath, []byte(batchContent), 0644)
	utils.PrintInfo("Startup script created.")
}
