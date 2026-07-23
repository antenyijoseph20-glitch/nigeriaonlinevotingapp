package repositories

import (
	"encoding/json"
	"os"

	"nigeriaonlinevoting/models"
)

// load reads all elections from the JSON file.
func (r *ElectionJSONRepository) load() ([]models.Election, error) {

	var elections []models.Election

	file, err := os.Open(r.filePath)

	if os.IsNotExist(err) {

		return elections, nil

	}

	if err != nil {

		return nil, err

	}

	defer file.Close()

	err = json.NewDecoder(file).Decode(&elections)

	if err != nil {

		return nil, err

	}

	return elections, nil
}

// save writes all elections to the JSON file.
func (r *ElectionJSONRepository) save(
	elections []models.Election,
) error {

	file, err := os.Create(r.filePath)

	if err != nil {

		return err

	}

	defer file.Close()

	encoder := json.NewEncoder(file)

	encoder.SetIndent("", "    ")

	return encoder.Encode(elections)

}
