package database

import(
    "gorm.io/gorm"
)

type ExampleSentence struct {
	gorm.Model
	Sentence      string `gorm:"not null"`
	TranslationID uint   `gorm:"not null"`
}
type Translation struct {
	gorm.Model
	Translation      string            `gorm:"not null"`
	WordID           uint              `gorm:"not null"`
	ExampleSentences []ExampleSentence `gorm:"foreignKey:TranslationID"`
}
type Word struct {
	gorm.Model
	Word         string        `gorm:"unique;not null"`
	Translations []Translation `gorm:"foreignKey:WordID"`
}

