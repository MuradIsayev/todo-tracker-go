package base

import (
	"encoding/json"
	"fmt"
	"os"
	"reflect"
	"time"

	"github.com/MuradIsayev/todo-tracker/helpers"
	"github.com/MuradIsayev/todo-tracker/status"
)

type BaseService[T any] struct {
	FilePath string
}

func (s *BaseService[T]) ReadFromFile(data *[]T) error {
	fileContent, err := os.ReadFile(s.FilePath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil // No file, treat it as an empty state
		}
		return fmt.Errorf("cannot read file: %v", err)
	}

	if len(fileContent) > 0 {
		if err := json.Unmarshal(fileContent, &data); err != nil {
			return fmt.Errorf("cannot convert JSON to struct: %v", err)
		}
	}
	return nil
}

func (s *BaseService[T]) WriteToFile(data []T) error {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("cannot convert struct to JSON: %v", err)
	}

	if err := os.WriteFile(s.FilePath, jsonData, 0644); err != nil {
		return fmt.Errorf("cannot write to file: %v", err)
	}
	return nil
}

func (s *BaseService[T]) GetNextID(items []T) int {
	v := reflect.ValueOf(items)
	if v.Kind() == reflect.Slice && v.Len() > 0 {
		lastItem := v.Index(v.Len() - 1).FieldByName("Id")
		lastID := int(lastItem.Int())
		return lastID + 1
	}
	return 1
}

func (s *BaseService[T]) DeleteAllItems() error {
	var items []T

	return s.WriteToFile(items)
}

func (s *BaseService[T]) FindItemById(items []T, id int) (int, *T, error) {
	for i, item := range items {
		v := reflect.ValueOf(item)
		if v.Kind() == reflect.Struct {
			itemID := int(v.FieldByName("Id").Int())
			if itemID == id {
				return i, &item, nil
			}
		}
	}

	return -1, nil, fmt.Errorf("item with ID=%d not found", id)
}

func (s *BaseService[T]) UpdateItemName(id, name string) error {
	itemId, err := helpers.ValidateIdAndConvertToInt(id)
	if err != nil {
		return err
	}

	items := []T{}
	err = s.ReadFromFile(&items)
	if err != nil {
		return err
	}

	index, item, err := s.FindItemById(items, itemId)
	if err != nil {
		return err
	}

	if name != "" {
		v := reflect.ValueOf(item).Elem().FieldByName("Name")
		v.SetString(name)

		v = reflect.ValueOf(item).Elem().FieldByName("UpdatedAt")
		v.Set(reflect.ValueOf(time.Now()))

		items[index] = *item
	}

	return s.WriteToFile(items)
}

func (s *BaseService[T]) DeleteItem(id string) error {
	itemId, err := helpers.ValidateIdAndConvertToInt(id)
	if err != nil {
		return err
	}

	items := []T{}
	err = s.ReadFromFile(&items)
	if err != nil {
		return err
	}

	index, _, err := s.FindItemById(items, itemId)
	if err != nil {
		return err
	}

	items = append(items[:index], items[index+1:]...)

	return s.WriteToFile(items)
}

func (s *BaseService[T]) UpdateItemStatus(id string, itemStatus status.ItemStatus) error {
	taskId, err := helpers.ValidateIdAndConvertToInt(id)
	if err != nil {
		return err
	}

	items := []T{}
	err = s.ReadFromFile(&items)
	if err != nil {
		return err
	}

	index, item, err := s.FindItemById(items, taskId)
	if err != nil {
		return err
	}

	v := reflect.ValueOf(item).Elem().FieldByName("Status")
	v.Set(reflect.ValueOf(itemStatus))

	v = reflect.ValueOf(item).Elem().FieldByName("UpdatedAt")
	v.Set(reflect.ValueOf(time.Now()))

	items[index] = *item

	return s.WriteToFile(items)
}
