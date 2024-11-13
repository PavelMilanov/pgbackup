package system

import (
	"testing"
)

func TestGetStorageInfo(t *testing.T) {
	data := GetStorageInfo()
	t.Log(data)
}
