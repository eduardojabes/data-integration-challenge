package company

import (
	"context"
	"path/filepath"
	"testing"
)

func TestGetCompany(t *testing.T) {
	repository := NewCompanyCSVRepository()
	t.Run("cannot open the archive", func(t *testing.T) {
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
