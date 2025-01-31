package repositories

import (
	"encoding/json"
	"fmt"
	"os"
	"sync"
)

func LoadFromFile(fileName string, dest interface{}, mu *sync.Mutex) error {
	mu.Lock()
	defer mu.Unlock()

	file, err := os.Open(fileName)
	if os.IsNotExist(err) {
		return nil
	} else if err != nil {
		return fmt.Errorf("failed to open file %s: %v", fileName, err)
	}
	defer file.Close()

	if err := json.NewDecoder(file).Decode(dest); err != nil {
		return fmt.Errorf("failed to decode json from %s: %v", fileName, err)
	}
	return nil
}

func SaveToFile(fileName string, src interface{}, mu *sync.Mutex) error {
	mu.Lock()
	defer mu.Unlock()

	data, err := json.MarshalIndent(src, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal data to JSON: %v", err)
	}
	if err := os.WriteFile(fileName, data, 0644); err != nil {
		return fmt.Errorf("failed to write to file %s: %v", fileName, err)
	}
	return nil
}
