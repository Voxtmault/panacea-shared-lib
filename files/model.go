package files

import (
	"time"

	"gorm.io/gorm"
)

type File struct {
	Filename  string `json:"filename" validate:"required" example:"filename.ext"`
	MIMEType  string `json:"mime_type" validate:"required" example:"image/jpg"`
	Size      uint   `json:"size" validate:"omitempty,number,gte=0" example:"1024"`
	FilePath  string `json:"file_path" example:"path/to/your/file"`
	HashValue string `json:"hash_value" example:"hash_value"`
}

type Media struct {
	ID          uint           `json:"id" gorm:"primaryKey"`
	IDMediaType uint           `json:"id_media_type" gorm:"index;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;foreignKey:IDMediaType;references:ID"`
	RefID       uint           `json:"ref_id" gorm:"index"`
	SourceTable string         `json:"source_table" gorm:"index"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"deleted_at" gorm:"index"`
	File
}
