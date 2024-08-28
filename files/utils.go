package files

import (
	"crypto/sha256"
	"encoding/hex"
	"io"
	"mime"
	"mime/multipart"
	"path/filepath"
	"strings"
)

// getFileExtensions returns the designated folder, media type id, and the mime type based on the file extension
func getFileExtension(filename string) (string, uint, string) {
	// Get file extension
	fileExt := filepath.Ext(filename)
	mimeType := mime.TypeByExtension(fileExt)

	var designatedFolder string
	var typeId uint
	switch strings.Split(mimeType, "/")[0] {
	case "image":
		designatedFolder = "photos"
		typeId = 1
	case "audio":
		designatedFolder = "audios"
		typeId = 2
	case "video":
		designatedFolder = "videos"
		typeId = 3
	case "text":
		designatedFolder = "texts" // Including .txt and other plain text files
		typeId = 4
	case "application":
		designatedFolder = "applications" // Including PPTs, PDFs, Docs, etc
		typeId = 5
	default:
		designatedFolder = "others" // This folder is for files that types are not filtered by the switch
		typeId = 6
	}

	return designatedFolder, typeId, mimeType
}

func calculateFileHash(file *multipart.FileHeader) (string, error) {
	fileReader, err := file.Open()
	if err != nil {
		return "", err
	}
	defer fileReader.Close()

	// Calculate the file hash
	hasher := sha256.New()
	if _, err = io.Copy(hasher, fileReader); err != nil {
		return "", err
	}

	return hex.EncodeToString(hasher.Sum(nil)), nil
}

func getFileData(file *multipart.FileHeader) ([]byte, error) {
	fileReader, err := file.Open()
	if err != nil {
		return nil, err
	}
	defer fileReader.Close()

	fileData, err := io.ReadAll(fileReader)
	if err != nil {
		return nil, err
	}

	return fileData, nil
}
