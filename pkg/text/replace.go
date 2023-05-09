package text

import (
	"fmt"
	"strings"

	docx "baliance.com/gooxml/document"
	"github.com/sirupsen/logrus"
)

func FindAndReplace(inputFile string, outputFile string) error {

	// read and parse the template docx
	doc, err := docx.Open(inputFile)
	if err != nil {
		logrus.WithError(err).Errorf("input file error: %s", inputFile)
		return err
	}

	parseTimeCode(doc, outputFile)
	//findAndReplaceNamesInTables(doc, outputFile)

	return nil
}

func parseTimeCode(doc *docx.Document, outputFile string) {
	newDocx := docx.New()
	//paragraphs := []docx.Paragraph{}
	tables := doc.Tables()
	for i := 0; i < len(tables); i++ {
		rows := tables[i].Rows()
		for i := 0; i < len(rows); i++ {
			if i == 0 {
				oriStr := ""
				cells := rows[0].Cells()
				cell0Paras := cells[0].Paragraphs()
				for j := 0; j < len(cell0Paras); j++ {
					cell0Runs := cell0Paras[j].Runs()
					for k := 0; k < len(cell0Runs); k++ {
						if cell0Runs[k].Text() == "" {
							continue
						}
						oriStr = attachRuns(cell0Runs)
					}
				}
				if !strings.Contains(oriStr, "TIME CODE") {
					break
				}
				continue
			}
			cells := rows[i].Cells()
			cell0Paras := cells[0].Paragraphs()
			if len(cells) < 3 {
				break
			}
			//cell2Paras := cells[2].Paragraphs()
			for j := 0; j < len(cell0Paras); j++ {
				cell0Runs := cell0Paras[j].Runs()
				//cell2Runs := cell2Paras[j].Runs()
				oriTimeCode := attachRuns(cell0Runs)
				if oriTimeCode == "" /* || cell2Runs[k].Text() == "" */ {
					continue
				}
				newTimeCodes := strings.Split(oriTimeCode, ":")
				newTimeCode := newTimeCodes[1] + newTimeCodes[2]
				fmt.Println(newTimeCode)
				para := newDocx.AddParagraph()
				run := para.AddRun()
				run.AddText(newTimeCode)
			}
		}
	}
	newDocx.SaveToFile(outputFile)
}

func attachRuns(cellRuns []docx.Run) string {
	oriStr := ""
	for l := 0; l < len(cellRuns); l++ {
		oriStr += cellRuns[l].Text()
	}
	return oriStr
}

func findAndReplaceNamesInTables(doc *docx.Document, outputFile string) {
	var nameKV map[string]string

	paragraphs := []docx.Paragraph{}
	tables := doc.Tables()
	for i := 0; i < len(tables); i++ {
		if i == 1 {
			nameKV = getNameKV(tables[i])
			//fmt.Println(nameKV)
			continue
		}
		rows := tables[i].Rows()
		for j := 0; j < len(rows); j++ {
			cells := rows[j].Cells()
			paragraphs = append(paragraphs, cells[0].Paragraphs()...)
		}
	}
	for _, p := range paragraphs {
		characterName := ""
		runs := p.Runs()
		for i := 0; i < len(runs); i++ {
			characterName += runs[i].Text()
			/*if i > 0 {
				runs[i].ClearContent()
			}*/
		}
		nameV, ok := nameKV[characterName]
		if ok {
			runs[0].ClearContent()
			runs[0].AddText(nameV)
		}
	}

	doc.SaveToFile(outputFile)
}

func getNameKV(table docx.Table) map[string]string {
	nameKV := make(map[string]string)
	rows := table.Rows()
	for i := 1; i < len(rows); i++ {
		cells := rows[i].Cells()
		cell0Paras := cells[0].Paragraphs()
		cell1Paras := cells[1].Paragraphs()
		for j := 0; j < len(cell0Paras) && j < len(cell1Paras); j++ {
			cell0Runs := cell0Paras[j].Runs()
			cell1Runs := cell1Paras[j].Runs()
			for k := 0; k < len(cell0Runs); k++ {
				if cell0Runs[k].Text() == "" {
					continue
				}
				oriStr := attachRuns(cell1Runs)
				_, ok := nameKV[oriStr]
				if !ok {
					nameKV[oriStr] = cell0Runs[k].Text()
				}
			}
		}
	}

	fmt.Println(nameKV)

	return nameKV
}
