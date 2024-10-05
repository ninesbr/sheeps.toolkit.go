package jsonstorage

import (
	"encoding/json"
	"fmt"
	"reflect"
)

type Document map[string]any

func (d Document) getFirstNotNull(keys ...string) any {
	for _, key := range keys {
		if value, ok := d[key]; ok {
			return value
		}
	}
	return nil
}

func (d Document) Marshal() ([]byte, error) {
	return json.Marshal(d)
}

func (d Document) GetID() any {
	return d.getFirstNotNull("id", "ID", "Id", "_id")
}

func (d Document) Validate() error {
	if d.GetID() == nil {
		return fmt.Errorf("document must have an id")
	}
	return nil
}

func (d Document) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]any(d))
}

func structToDoc(obj any) (Document, error) {
	immutable := reflect.ValueOf(obj)
	if immutable.Kind() == reflect.Ptr {
		immutable = immutable.Elem()
	}

	if immutable.Kind() != reflect.Struct {
		return nil, fmt.Errorf("obj must be a struct")
	}
	var doc Document
	buff, err := json.Marshal(obj)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(buff, &doc)
	if err != nil {
		return nil, err
	}
	doc["id"] = doc.GetID()

	return doc, nil
}
