package database

import (
)

type ExampleSentence struct {
	ID           uint           `gorm:"primaryKey"`
	Sentence      string `gorm:"not null"`
	TranslationID uint   `gorm:"not null"`
}

type Translation struct {
	ID           uint           `gorm:"primaryKey"`
	Translation      string            `gorm:"uniqueIndex; not null"`
	WordID           uint              `gorm:"not null"` 
	ExampleSentences []ExampleSentence  `gorm:"foreignKey:TranslationID;constraint:OnDelete:CASCADE;"`
}

type Word struct {
	ID           uint           `gorm:"primaryKey"`
	Word         string        `gorm:"uniqueIndex; not null"`
	Translations []Translation `gorm:"foreignKey:WordID;constraint:OnDelete:CASCADE;"`
}

