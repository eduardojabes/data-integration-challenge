package company

import (
	"bytes"
	"context"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
)

func TestReadFile(t *testing.T) {
	t.Run("Correct data read", func(t *testing.T) {
		repository := NewCompanyCSVRepository()
		writeData := "test; test"

		want := strings.Split(writeData, ";")

		var buffer bytes.Buffer
		buffer.WriteString(writeData)
		data, err := repository.Read_File(&buffer)

		if err != nil {
			t.Errorf("got %v ,but it should be nil", err)
		}
		if !reflect.DeepEqual(data[0], want) {
			t.Errorf("got %s ,but it should be %s", data[0], want)
		}

	})

}

func TestGetCompany(t *testing.T) {
	repository := NewCompanyCSVRepository()
	t.Run("Cannot open the archive", func(t *testing.T) {
		_, err := repository.GetCompany(context.Background(), "...")

		if err == nil {
			t.Errorf("got %v ,but it should be an error", err)
		}

	})

	t.Run("open the archive", func(t *testing.T) {

		path, _ := filepath.Abs("test_CSV.csv")

		_, err := repository.GetCompany(context.Background(), path)

		if err != nil {
			t.Errorf("got %v ,but it should be nil", path)
		}

	})

}
