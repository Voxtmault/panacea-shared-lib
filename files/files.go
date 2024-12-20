package files

import (
	"errors"
	"os"
)

// Saves file with the provided filepath, will truncate the file if it already exists
func SaveFile(filePath string, data []byte) error {
	return os.WriteFile(filePath, data, 0644)
}

func DeleteFile(filePath string) error {
	if err := os.Remove(filePath); err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil
		}
		if errors.Is(err, os.ErrPermission) {
			return errors.New("permission denied")
		}
		if errors.Is(err, os.ErrInvalid) {
			return errors.New("invalid file type")
		}
		return err
	}

	return nil
}
