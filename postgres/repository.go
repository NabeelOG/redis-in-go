package postgres

import (
	"fmt"
	"redis-learn/models"

	"gorm.io/gorm"
)

// CREATE
func CreatePerson(db *gorm.DB, person models.Person) error {
	result := db.Create(&person)
	if result.Error != nil {
		return fmt.Errorf("failed to create a person: %v", result.Error)
	}
	return nil
}

// READ
func GetPerson(db *gorm.DB, id string) (*models.Person, error) {
	var person models.Person
	result := db.First(&person, "id = ?", id)
	if result.Error == gorm.ErrRecordNotFound {
		return nil, gorm.ErrRecordNotFound
	}
	if result.Error != nil {
		return nil, fmt.Errorf("failed to get person: %v", result.Error)
	}
	return &person, nil
}

// UPDATE
func UpdatePerson(db *gorm.DB, person *models.Person) error {
	result := db.Save(person)
	if result.Error != nil {
		return fmt.Errorf("failed to update person: %v", result.Error)
	}
	return nil
}

// DELETE
func DeletePerson(db *gorm.DB, id string) error {
	result := db.Delete(&models.Person{}, "id = ?", id)
	if result.Error != nil {
		return fmt.Errorf("failed to delete person: %v", result.Error)
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("person with ID %s does not exist", id)
	}
	return nil
}

// LIST ALL
func ListAllPeople(db *gorm.DB) ([]*models.Person, error) {
	var people []*models.Person
	result := db.Find(&people)
	if result.Error != nil {
		return nil, fmt.Errorf("failed to list people: %v", result.Error)
	}
	return people, nil
}
