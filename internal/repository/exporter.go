package repository

import (
	"encoding/csv"
	"fmt"
	"log"
	"loquegasto-telegram/internal/domain"
	"os"
	"path/filepath"
	"reflect"
	"strconv"
	"time"
)

type ExporterRepository interface {
	Create(userID int64) error
	AddRow(row domain.ExporterRow, userID int64) error
	GetFilePath(userID int64) string
	Delete(userID int64) error
}
type exporterRepository struct {
	filePath string
	header   []string
}

func NewExporterRepository(filePath string) ExporterRepository {
	// Create directory if not exists
	filePath = filepath.Join(os.TempDir(), filePath)
	err := os.MkdirAll(filePath, os.ModeDir|os.ModePerm)
	if err != nil {
		log.Fatalln(err)
	}

	// Store header
	var header []string
	row := &domain.ExporterRow{}
	typeOf := reflect.TypeOf(row)
	refElem := typeOf.Elem()
	for i := 0; i < refElem.NumField(); i++ {
		header = append(header, refElem.Field(i).Tag.Get("csv"))
	}

	return &exporterRepository{
		filePath: filePath,
		header:   header,
	}
}

func (r *exporterRepository) Create(userID int64) error {
	f, err := os.Create(r.GetFilePath(userID))
	defer f.Close()
	if err != nil {
		return err
	}

	csvw := csv.NewWriter(f)
	err = csvw.Write(r.header)
	if err != nil {
		return err
	}
	csvw.Flush()

	return csvw.Error()
}
func (r *exporterRepository) AddRow(row domain.ExporterRow, userID int64) error {
	f, err := os.OpenFile(r.GetFilePath(userID), os.O_APPEND|os.O_WRONLY, os.ModeAppend)
	defer f.Close()
	if err != nil {
		return err
	}

	csvw := csv.NewWriter(f)

	err = csvw.Write([]string{
		row.ID,
		fmt.Sprintf("%.2f", row.Amount),
		row.Description,
		row.CategoryName,
		row.WalletName,
		row.CreatedAt.Format(time.RFC3339),
	})
	if err != nil {
		return err
	}
	csvw.Flush()

	err = csvw.Error()
	return err
}
func (r *exporterRepository) GetFilePath(userID int64) string {
	return filepath.Join(r.filePath, strconv.FormatInt(userID, 10)+".csv")
}
func (r *exporterRepository) Delete(userID int64) error {
	return os.RemoveAll(r.GetFilePath(userID))
}
