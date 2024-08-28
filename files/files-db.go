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
	gormTx = gConn.WithContext(ctx).Session(&gorm.Session{NewDB: true}).Begin()

	if tx != nil {
		// If tx is provided then load the said transaction
		gormTx = gormTx.Set("gorm:db", tx)
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

	//Calculate the file hash
	hasher := sha256.New()
	if _, err = io.Copy(hasher, fileReader); err != nil {
		if tx == nil {
			gormTx.Rollback()
		}
		return err
	}

	fileInfo.HashValue = hex.EncodeToString(hasher.Sum(nil))
	designatedFolder, fileInfo.MIMEType = getFileExtension(file.Filename)

	designatedFolder = fmt.Sprintf("%s/%s/%s", config.FileRootPath, designatedFolder, file.Filename)

	media := Media{
		RefID:       refId,
		SourceTable: refTable,
		File:        fileInfo,
	}

	result := gormTx.Create(&media)

	if result.Error != nil {
		if tx == nil {
			gormTx.Rollback()
		}
		return result.Error
	}

	fileData, err := io.ReadAll(fileReader)
	if err != nil {
		if tx == nil {
			gormTx.Rollback()
		}
		return err
	}

	// Save to the designated folder
	if err = SaveFile(designatedFolder, fileData); err != nil {
		return err
	}

	if tx == nil {
		gormTx.Commit()
	}

	return nil
}

func DeleteFromDB() {

}

func GetFromDB() {

}

func UpdateInDB() {

}

// getFileExtensions returns the designated folder and mime type based on the file extension
func getFileExtension(filename string) (string, string) {
	// Get file extension
	fileExt := filepath.Ext(filename)
	mimeType := mime.TypeByExtension(fileExt)

	var designatedFolder string
	switch strings.Split(mimeType, "/")[0] {
	case "image":
		designatedFolder = "photos"
	case "application":
		designatedFolder = "applications" // Including PPTs, PDFs, Docs, etc
	case "video":
		designatedFolder = "videos"
	case "audio":
		designatedFolder = "audios"
	case "text":
		designatedFolder = "plaintexts" // Including .txt and other plain text files
	default:
		designatedFolder = "others" // This folder is for files that types are not filtered by the switch
	}

	return designatedFolder, mimeType
}
