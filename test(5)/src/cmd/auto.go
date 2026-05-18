package cmd

import (
	"os"
	"runtime"

	backupstatus "github.com/crtwheel/cartify/src/status/backup"
	spotifystatus "github.com/crtwheel/cartify/src/status/spotify"
)

// Auto checks Spotify state, re-backup and apply if needed, then launch
// Spotify client normally. Blocks Spotify updates and auto-updates Cartify.
func Auto(CartifyVersion string) {
	if runtime.GOOS == "windows" || runtime.GOOS == "darwin" {
		BlockSpotifyUpdates(true)
	}

	if CartifyVersion != "Dev" {
		AutoUpdateCheck(CartifyVersion)
	}

	backupVersion := backupSection.Key("version").MustString("")
	spotStat := spotifystatus.Get(appPath)
	backStat := backupstatus.Get(prefsPath, backupFolder, backupVersion)

	if spotStat.IsBackupable() && (backStat.IsEmpty() || backStat.IsOutdated()) {
		Backup(CartifyVersion, true)
		backupVersion := backupSection.Key("version").MustString("")
		backStat = backupstatus.Get(prefsPath, backupFolder, backupVersion)
	}

	if !backStat.IsBackuped() {
		os.Exit(1)
	}

	if isAppX {
		spotStat = spotifystatus.Get(appDestPath)
	}

	if !spotStat.IsApplied() && backStat.IsBackuped() {
		CheckStates()
		InitSetting()
		Apply(CartifyVersion)
	}
}


