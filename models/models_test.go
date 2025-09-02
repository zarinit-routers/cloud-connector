package models

import "testing"

func TestParseUUID(t *testing.T) {
	uuid, err := ParseUUID("12345678-1234-5678-1234-567812345678")
	if err != nil {
		t.Error(err)
	}
	t.Log(uuid)
	uuid, err = ParseUUID("1224887889")
	if err == nil {
		t.Errorf("Invalid UUID string parsing not returning error")
	}
	t.Log(uuid)
}
