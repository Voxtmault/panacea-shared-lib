package files

import (
	"context"
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"fmt"
	"io"
	"mime"
	"mime/multipart"
	"path/filepath"
	"strings"

	"github.com/voxtmault/panacea-shared-lib/config"
	"github.com/voxtmault/panacea-shared-lib/storage"
	"gorm.io/gorm"
)

func SaveToDB(ctx context.Context, tx *sql.Tx, refId uint, refTable string, file *multipart.FileHeader) error {
	// Get the ORM Connection
	gConn := storage.GetGORMMariaDB()
	config := config.GetConfig().FileHandlingConfig

	var gormTx *gorm.DB
	if tx != nil {
		// If tx is provided, use it to create a GORM transaction
		gormTx = gConn.WithContext(ctx).Session(&gorm.Session{NewDB: true}).Begin()
		gormTx = gormTx.Set("gorm:db", tx)
	} else {
		// Otherwise, start a new GORM transaction
		gormTx = gConn.WithContext(ctx).Begin()
	}

	// Save the file entry to the database
	fileInfo := File{}
	var designatedFolder string

	fileReader, err := file.Open()
	if err != nil {
		if tx == nil {
			gormTx.Rollback()
		}
		return err
	}
	defer fileReader.Close()

	// Calculate the file hash
	hasher := sha256.New()
	if _, err = io.Copy(hasher, fileReader); err != nil {
		if tx == nil {
			gormTx.Rollback()
		}
		return err
	}

	fileInfo.HashValue = hex.EncodeToString(hasher.Sum(nil))

	media := Media{
		RefID:       refId,
		SourceTable: refTable,
		File:        fileInfo,
	}
	designatedFolder, media.IDMediaType, fileInfo.MIMEType = getFileExtension(file.Filename)
	designatedFolder = fmt.Sprintf("%s/%s/%s", config.FileRootPath, designatedFolder, file.Filename)

	result := gormTx.Create(&media)
	if result.Error != nil {
		if tx == nil {
			gormTx.Rollback()
		}
		return result.Error
	}

	// Re-open the file reader to read the file data
	fileReader, err = file.Open()
	if err != nil {
		if tx == nil {
			gormTx.Rollback()
		}
		return err
	}
	defer fileReader.Close()

	fileData, err := io.ReadAll(fileReader)
	if err != nil {
		if tx == nil {
			gormTx.Rollback()
		}
		return err
	}

	// Save to the designated folder
	if err = SaveFile(designatedFolder, fileData); err != nil {
		if tx == nil {
			gormTx.Rollback()
		}
		return err
	}

	if tx == nil {
		// Commit the GORM transaction if no external transaction is provided
		if err := gormTx.Commit().Error; err != nil {
			gormTx.Rollback()
			return err
		}
	}

	return nil
}

func DeleteFromDB() {

}

func GetFromDB() {

}

func UpdateInDB() {

}

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
