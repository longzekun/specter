package models

import (
	"time"

	"github.com/gofrs/uuid"
	"gorm.io/gorm"
)

type Certificate struct {
	ID             uuid.UUID `gorm:"primaryKey;->;<-;type:uuid;"`
	CreatedAt      time.Time `gorm:"<-:create;"`
	CommonName     string
	CAType         string
	KeyType        string
	CertificatePEM string
	PrivateKeyPEM  string
}

func (c *Certificate) BeforeCreate(tx *gorm.DB) (err error) {
	c.ID, err = uuid.NewV4()
	if err != nil {
		return err
	}
	c.CreatedAt = time.Now()
	return nil
}
