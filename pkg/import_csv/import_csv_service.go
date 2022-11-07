package importcsv

import (
	"os"

	"github.com/gocarina/gocsv"
)

type ImportCsv interface {
	UnmarshalCsv(filepath string, model interface{}) error
}

type importCsvImpl struct{}

// UnmarshalCsv implements ImportCsv
func (*importCsvImpl) UnmarshalCsv(filepath string, model interface{}) error {
	csv, err := os.Open(filepath)
	if err != nil {
		return err
	}
	defer csv.Close()
	if err := gocsv.UnmarshalFile(csv, model); err != nil {
		return err
	}
	return nil
}

func NewImportCsv() ImportCsv {
	return &importCsvImpl{}
}
