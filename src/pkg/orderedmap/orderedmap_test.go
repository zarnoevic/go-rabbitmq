package orderedmap

import (
	"testing"
)

func TestOrderedMap_Add(t *testing.T) {
	om := NewOrderedMap()
	om.Add("key1", "value1")

	if val, exists := om.Get("key1"); !exists || val != "value1" {
		t.Errorf("Add failed to add key1 with value1, got: %s, want: value1", val)
	}
}

func TestOrderedMap_Delete(t *testing.T) {
	om := NewOrderedMap()
	om.Add("key1", "value1")
	om.Delete("key1")

	if _, exists := om.Get("key1"); exists {
		t.Error("Delete failed to delete key1")
	}
}

func TestOrderedMap_Get(t *testing.T) {
	om := NewOrderedMap()
	om.Add("key1", "value1")

	if val, exists := om.Get("key1"); !exists || val != "value1" {
		t.Errorf("GetItem failed to retrieve the correct value for key1, got: %s, want: value1", val)
	}
}

func TestOrderedMap_GetAll(t *testing.T) {
	om := NewOrderedMap()
	om.Add("key1", "value1")
	om.Add("key2", "value2")

	items := om.GetAll()
	if len(items) != 2 || items[0].value != "value1" || items[1].value != "value2" {
		t.Errorf("GetAll failed to retrieve all items in order, got: %v", items)
	}
}
