package company

import (
	"bytes"
	"context"
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
)

func CreatTestFile(data string) string {
	file, _ := ioutil.TempFile("./", "test_file_*.csv")

	buffer := &bytes.Buffer{}
	buffer.WriteString(data)

	file.Write(buffer.Bytes())

	if _, err := file.Seek(0, os.SEEK_SET); err != nil {
		panic(err)
	}

	file.Close()

	return file.Name()
}

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
		data := "name;addresszip;website\ntola sales group;78229;http://repsources.com"
		fileName := CreatTestFile(data)
		path, _ := filepath.Abs(fileName)

		got, err := repository.GetCompany(context.Background(), path)

		if err != nil {
			t.Errorf("got %v ,but it should be nil", path)
		}
		if got[0].Name != "TOLA SALES GROUP" || got[0].Zip != "78229" || got[0].Website != "http://repsources.com" {
			t.Errorf("got wrong data: %v", got[0])
		}
		os.Remove(fileName)
	})

}
