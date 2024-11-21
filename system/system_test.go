package system

import (
	"testing"
)

func TestGetStorageInfo(t *testing.T) {
	data := GetStorageInfo()
	t.Log(data)
}

func TestParseOldFiles(t *testing.T) {
	files := ParseOldFiles(1.0)
	t.Log(files)
}
