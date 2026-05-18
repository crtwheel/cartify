package cmd

import (
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"

	"github.com/crtwheel/cartify/src/utils"
)

func Update(currentVersion string) bool {
	tagName, err := utils.FetchLatestTag()
	if err != nil {
		utils.PrintError("Cannot fetch latest release info")
		utils.PrintError(err.Error())
		return false
	}
	if currentVersion == tagName {
		utils.PrintSuccess("cartify is up-to-date.")
		return false
	}

	utils.PrintInfo("Latest release: " + tagName)
	if currentVersion != "Dev" {
		updateFromTag(currentVersion, tagName)
	}
	return true
}

func AutoUpdateCheck(currentVersion string) {
	tagName, err := utils.FetchLatestTag()
	if err != nil || tagName == currentVersion || currentVersion == "Dev" {
		return
	}

	utils.PrintInfo("Auto-updating Cartify: v" + currentVersion + " -> v" + tagName)
	updateFromTag(currentVersion, tagName)
}

func updateFromTag(currentVersion, tagName string) {
	var assetURL string = "https://github.com/crtwheel/cartify/releases/download/v" + tagName + "/Cartify-" + tagName + "-" + runtime.GOOS + "-"
	var location string = os.TempDir() + "/Cartify-" + tagName

	if runtime.GOARCH == "386" && runtime.GOOS == "windows" {
		assetURL += "x32"
	} else if runtime.GOARCH == "arm64" {
		assetURL += "arm64"
	} else if runtime.GOOS == "windows" {
		assetURL += "x64"
	} else {
		assetURL += "amd64"
	}

	if runtime.GOOS == "windows" {
		assetURL += ".zip"
		location += ".zip"
	} else {
		assetURL += ".tar.gz"
		location += ".tar.gz"
	}

	spinner, _ := utils.Spinner.Start("Downloading Cartify update")
	out, err := os.Create(location)
	if err != nil {
		spinner.Fail("Auto-update failed")
		return
	}
	defer out.Close()

	resp, err := http.Get(assetURL)
	if err != nil {
		spinner.Fail("Auto-update failed")
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		spinner.Fail("Auto-update failed")
		return
	}

	_, err = io.Copy(out, resp.Body)
	if err != nil {
		spinner.Fail("Auto-update failed")
		return
	}
	spinner.Success("Downloaded Cartify update")

	exe, err := os.Executable()
	if err != nil {
		return
	}
	if exe, err = filepath.EvalSymlinks(exe); err != nil {
		return
	}

	exeOld := exe + ".old"
	utils.CheckExistAndDelete(exeOld)

	if err = os.Rename(exe, exeOld); err != nil {
		return
	}

	switch runtime.GOOS {
	case "windows":
		err = utils.Unzip(location, utils.GetExecutableDir())
	case "linux", "darwin":
		err = exec.Command("tar", "-xzf", location, "-C", utils.GetExecutableDir()).Run()
	}

	if err != nil {
		os.Rename(exeOld, exe)
		return
	}

	utils.CheckExistAndDelete(exeOld)
	utils.PrintSuccess("Cartify updated to v" + tagName)
}

func permissionError(err error) {
	utils.PrintInfo("If fatal error is \"Permission denied\", please check read/write permission of Cartify executable directory.")
	utils.PrintInfo("However, if you used a package manager to install Cartify, please upgrade by using the same package manager.")
	utils.Fatal(err)
}


