package udict

import (
	"fmt"
	"os"
	"os/user"
	"path/filepath"
	"runtime"
)

const AppName = "urban-dict-cli"

// createConfigDir creates a platform-specific directory for storing app configuration files.
func createConfigDir() (string, error) {
	var configDir string

	currentUser, err := user.Current()
	if err != nil {
		return "", fmt.Errorf("failed to get current user: %w", err)
	}

	switch runtime.GOOS {
	case "windows":
		appData := os.Getenv("APPDATA")
		if appData == "" {
			appData = filepath.Join(currentUser.HomeDir, "AppData", "Roaming")
		}
		configDir = filepath.Join(appData, AppName)
	case "linux":
		xdgConfigHome := os.Getenv("XDG_CONFIG_HOME")
		if xdgConfigHome == "" {
			configDir = filepath.Join(currentUser.HomeDir, ".config", AppName)
		} else {
			configDir = filepath.Join(xdgConfigHome, AppName)
		}
	case "darwin":
		configDir = filepath.Join(currentUser.HomeDir, "Library", "Application Support", AppName)
	default:
		return "", fmt.Errorf("unsupported OS: %s", runtime.GOOS)
	}

	if err = os.MkdirAll(configDir, 0755); err != nil {
		return "", fmt.Errorf("failed to create config directory: %w", err)
	}

	return configDir, nil
}

// GetDatabasePath generates the full path for the SQLite database file.
func GetDatabasePath() (string, error) {
	configDir, err := createConfigDir()
	if err != nil {
		return "", err
	}

	// Define the database file path inside the config directory.
	dbPath := filepath.Join(configDir, "store.db")

	return dbPath, nil
}
