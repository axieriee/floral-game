package save

import (
	"encoding/json"
	"os"
	"path/filepath"

	"github.com/floral-game/floral-realms/internal/game"
)

func savePath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	dir := filepath.Join(home, ".floral-realms")
	if err := os.MkdirAll(dir, 0755); err != nil {
		return "", err
	}
	return filepath.Join(dir, "save.json"), nil
}

// Save writes the game state to disk.
func Save(state *game.GameState) error {
	path, err := savePath()
	if err != nil {
		return err
	}
	data, err := json.MarshalIndent(state, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0644)
}

// Load reads the game state from disk. Returns nil if no save exists.
func Load() (*game.GameState, error) {
	path, err := savePath()
	if err != nil {
		return nil, err
	}
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}
	var state game.GameState
	if err := json.Unmarshal(data, &state); err != nil {
		return nil, err
	}
	return &state, nil
}
