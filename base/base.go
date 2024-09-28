package base

import (
	"encoding/json"
	"fmt"
	"os"
	"reflect"
)

type BaseService struct {
	FilePath string
}

func (s *BaseService) ReadFromFile(data interface{}) error {
	fileContent, err := os.ReadFile(s.FilePath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil // No file, treat it as an empty state
		}
		return fmt.Errorf("cannot read file: %v", err)
	}

	if len(fileContent) > 0 {
		if err := json.Unmarshal(fileContent, data); err != nil {
			return fmt.Errorf("cannot convert JSON to struct: %v", err)
		}
	}
	return nil
}

func (s *BaseService) WriteToFile(data interface{}) error {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("cannot convert struct to JSON: %v", err)
	}

	if err := os.WriteFile(s.FilePath, jsonData, 0644); err != nil {
		return fmt.Errorf("cannot write to file: %v", err)
	}
	return nil
}

func (s *BaseService) GetNextID(items interface{}) int {
	v := reflect.ValueOf(items)
	if v.Kind() == reflect.Slice && v.Len() > 0 {
		lastItem := v.Index(v.Len() - 1).FieldByName("Id")
		lastID := int(lastItem.Int())
		return lastID + 1
	}
	return 1
}
