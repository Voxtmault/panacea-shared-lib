package files

import (
	"context"
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"fmt"
	"io"
	"mime/multipart"
	"strings"

	"github.com/voxtmault/panacea-shared-lib/config"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func SaveToDB(ctx context.Context, tx *sql.Tx, refId uint, refTable, preferedName string, file *multipart.FileHeader) (string, error) {
	gormTx, err := gorm.Open(mysql.New(mysql.Config{
		Conn: tx,
	}), &gorm.Config{})
	if err != nil {
		tx.Rollback()
		return "", err
	}

	config := config.GetConfig().FileHandlingConfig

	// Save the file entry to the database
	var designatedFolder string

	fileInfo := File{}
	media := Media{}

	media.RefID = refId
	media.SourceTable = refTable

	fileInfo.Filename = file.Filename
	fileInfo.Size = uint(file.Size)
	fileInfo.HashValue, err = calculateFileHash(file)
	if err != nil {
		if tx == nil {
			gormTx.Rollback()
		}
		return "", fmt.Errorf("calculate file hash: %e", err)
	}

	designatedFolder, media.IDMediaType, fileInfo.MIMEType = getFileExtension(file.Filename)
	fileInfo.FilePath = fmt.Sprintf("%s/%s/%d%s-%s-%s", config.FileRootPath, designatedFolder, refId, refTable, preferedName, strings.Replace(file.Filename, " ", "_", -1))

	media.File = fileInfo

	result := gormTx.Create(&media)
	if result.Error != nil {
		if tx == nil {
			gormTx.Rollback()
		}
		return "", result.Error
	}

	fileData, err := getFileData(file)
	if err != nil {
		if tx == nil {
			gormTx.Rollback()
		}
		return "", fmt.Errorf("get file data: %e", err)
	}

	// Save to the designated folder
	if err = SaveFile(fileInfo.FilePath, fileData); err != nil {
		if tx == nil {
			gormTx.Rollback()
		}
		return "", err
	}

	if tx == nil {
		// Commit the GORM transaction if no external transaction is provided
		if err := gormTx.Commit().Error; err != nil {
			gormTx.Rollback()
			return "", err
		}
	}

	return fileInfo.FilePath, nil
}

func UpdateInDB(ctx context.Context, tx *sql.Tx, mediaId, refId uint, refTable string, file *multipart.FileHeader) (string, error) {
	gormTx, err := gorm.Open(mysql.New(mysql.Config{
		Conn: tx,
	}), &gorm.Config{})
	if err != nil {
		tx.Rollback()
		return "", err
	}

	config := config.GetConfig().FileHandlingConfig

	// Save the file entry to the database
	fileInfo := File{}
	media := Media{}
	var designatedFolder string

	// Fetch the existing Media record
	result := gormTx.First(&media, mediaId)
	if result.Error != nil {
		if tx == nil {
			gormTx.Rollback()
		}
		return "", result.Error
	}

	fileReader, err := file.Open()
	if err != nil {
		if tx == nil {
			gormTx.Rollback()
		}
		return "", err
	}
	defer fileReader.Close()

	// Calculate the file hash
	hasher := sha256.New()
	if _, err = io.Copy(hasher, fileReader); err != nil {
		if tx == nil {
			gormTx.Rollback()
		}
		return "", err
	}

	fileInfo.HashValue = hex.EncodeToString(hasher.Sum(nil))
	fileInfo.Filename = file.Filename
	fileInfo.Size = uint(file.Size)

	designatedFolder, media.IDMediaType, fileInfo.MIMEType = getFileExtension(file.Filename)
	designatedFolder = fmt.Sprintf("%s/%s/%d%s-%s", config.FileRootPath, designatedFolder, refId, refTable, strings.Replace(file.Filename, " ", "_", -1))

	fileInfo.FilePath = designatedFolder
	media.File = fileInfo

	result = gormTx.Save(&media)
	if result.Error != nil {
		if tx == nil {
			gormTx.Rollback()
		}
		return "", result.Error
	}

	// Re-open the file reader to read the file data
	fileReader, err = file.Open()
	if err != nil {
		if tx == nil {
			gormTx.Rollback()
		}
		return "", err
	}
	defer fileReader.Close()

	fileData, err := io.ReadAll(fileReader)
	if err != nil {
		if tx == nil {
			gormTx.Rollback()
		}
		return "", err
	}

	// Save to the designated folder
	if err = SaveFile(designatedFolder, fileData); err != nil {
		if tx == nil {
			gormTx.Rollback()
		}
		return "", err
	}

	if tx == nil {
		// Commit the GORM transaction if no external transaction is provided
		if err := gormTx.Commit().Error; err != nil {
			gormTx.Rollback()
			return "", err
		}
	}

	return fileInfo.FilePath, nil
}

// DeleteFromDB is not implemented yet
func DeleteFromDB() {
}

// GetFromDB is not implemented yet
func GetFromDB() {
}
