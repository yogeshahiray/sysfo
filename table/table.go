// This file will have functions related to table output handling
package table

import (
	"encoding/csv"
	"fmt"
	"github.com/bndr/gotabulate"
	"github.com/yogeshahiray/sysfo/common"
	"os"
)

// Structure which will store the table headers and rows
type TableData struct {
	Headers    []string
	Rows       [][]interface{}
	RowIndex   int
	RowLength  int
	TableWidth int // Number of header elements will decide the number of row elements
}

// Create the instance of the TableData structure

var TableContent = TableData{
	RowIndex:   0,
	TableWidth: 0,
	RowLength:  0,
}

func AddSingleHeader(header string) {
	TableContent.Headers = append(TableContent.Headers, header)
	TableContent.TableWidth++
}

func AddHeaders(headers []string) {
	for _, d := range headers {
		TableContent.Headers = append(TableContent.Headers, d)
		TableContent.TableWidth++
	}
}

func AddRow(row []interface{}) {
	TableContent.Rows[TableContent.RowIndex] = append(TableContent.Rows[TableContent.RowIndex], row)
	TableContent.RowLength++
	if TableContent.RowLength >= TableContent.TableWidth {
		// Filled out complete row
		TableContent.RowIndex++
	}

}

// Write in CSV file
func WriteTable(rows [][]interface{}, headers []string, filename string) {

	f, err := os.Create(filename)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer f.Close()
	w := csv.NewWriter(f)
	defer w.Flush()
	w.Write(headers)

	for _, r := range rows {
		s := make([]string, len(r))
		for j, v := range r {
			s[j] = v.(string)
		}
		err = w.Write(s)
		if err != nil {
			fmt.Println(err)
		}
	}
}

// Function to print the table
func PrintTable(rows [][]interface{}, headers []string) {
	t := gotabulate.Create(rows)
	t.SetHeaders(headers)
	// Set Max Cell Size
	t.SetMaxCellSize(14)

	// Turn On String Wrapping
	t.SetWrapStrings(true)
	t.SetAlign("center")
	fmt.Println(t.Render("grid"))
}

func WriteCommandOutput(rows [][]interface{}, headers []string) {
	if len(common.InputParams.OutFile) > 0 {
		WriteTable(rows, headers, common.InputParams.OutFile)
	} else {
		PrintTable(rows, headers)
	}
}
