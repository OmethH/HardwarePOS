package reports

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"strconv"

	"github.com/go-pdf/fpdf"
)

// ToCSV renders any report's rows into a CSV byte buffer given the column
// headers and a row-extraction function.
func ToCSV(headers []string, rowCount int, cell func(row, col int) string) ([]byte, error) {
	var buf bytes.Buffer
	w := csv.NewWriter(&buf)

	if err := w.Write(headers); err != nil {
		return nil, err
	}
	for r := 0; r < rowCount; r++ {
		record := make([]string, len(headers))
		for c := range headers {
			record[c] = cell(r, c)
		}
		if err := w.Write(record); err != nil {
			return nil, err
		}
	}
	w.Flush()
	return buf.Bytes(), w.Error()
}

// ToPDF renders a simple tabular PDF report.
func ToPDF(title string, headers []string, rowCount int, cell func(row, col int) string) ([]byte, error) {
	pdf := fpdf.New("P", "mm", "A4", "")
	pdf.AddPage()
	pdf.SetFont("Arial", "B", 14)
	pdf.CellFormat(0, 10, title, "", 1, "L", false, 0, "")
	pdf.Ln(2)

	colWidth := 190.0 / float64(len(headers))

	pdf.SetFont("Arial", "B", 10)
	for _, header := range headers {
		pdf.CellFormat(colWidth, 8, header, "1", 0, "L", false, 0, "")
	}
	pdf.Ln(-1)

	pdf.SetFont("Arial", "", 9)
	for r := 0; r < rowCount; r++ {
		for c := range headers {
			pdf.CellFormat(colWidth, 7, cell(r, c), "1", 0, "L", false, 0, "")
		}
		pdf.Ln(-1)
	}

	var buf bytes.Buffer
	if err := pdf.Output(&buf); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func fmtMoney(v float64) string {
	return strconv.FormatFloat(v, 'f', 2, 64)
}

func fmtQty(v float64) string {
	if v == float64(int64(v)) {
		return fmt.Sprintf("%d", int64(v))
	}
	return strconv.FormatFloat(v, 'f', 2, 64)
}
