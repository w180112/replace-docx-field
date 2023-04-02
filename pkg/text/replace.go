package text

import (
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

	/*paragraphs := []docx.Paragraph{}
	for _, p := range doc.Paragraphs() {
		paragraphs = append(paragraphs, p)
	}

	for _, p := range paragraphs {
		for _, r := range p.Runs() {
			fmt.Println(r.Text())
		}
	}*/
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
			for _, p := range cells[0].Paragraphs() {
				paragraphs = append(paragraphs, p)
			}
		}
	}
	for _, p := range paragraphs {
		for _, r := range p.Runs() {
			nameV, ok := nameKV[r.Text()]
			if ok {
				r.ClearContent()
				r.AddText(nameV)
			}
		}
	}

	doc.SaveToFile("/home/the/Downloads/edit-test.docx")

	return nil
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
			for k := 0; k < len(cell0Runs) && k < len(cell1Runs); k++ {
				if cell0Runs[k].Text() == "" {
					continue
				}
				nameKV[cell1Runs[k].Text()] = cell0Runs[k].Text()
			}
		}
	}

	return nameKV
}
