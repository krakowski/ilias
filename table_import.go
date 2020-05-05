package ilias

import (
	"github.com/360EntSecGroup-Skylar/excelize"
	"github.com/PuerkitoBio/goquery"
	"net/http"
	"net/textproto"
	"net/url"
)

const (
	importPath string = "ilias.php?mode=1&cmd=post&cmdClass=ildclrecordlistgui&cmdNode=6k:yf:yw&baseClass=ilrepositorygui&fallbackCmd=showImportExcel"
)

type ImportParams struct {
	Reference		string	`schema:"ref_id"`
	Table			string  `schema:"table_id"`
	Token		    string  `schema:"rtoken"`
}

func (table *TableService) Import(params *ImportParams, sheet *excelize.File) error {

	buffer, err := sheet.WriteToBuffer()
	if err != nil {
		return err
	}

	// Prepare request url
	path, err := addQueryParams(importPath, params)
	if err != nil {
		return err
	}

	hash, err := table.getFileHash(params)
	if err != nil {
		return err
	}

	values := url.Values{
		"ilfilehash" : { hash },
		"cmd[importExcel]" : { "Importieren" },
	}

	req, err := table.client.NewMultipartRequest(http.MethodPost, path, values, &UploadFile{
		Content:   buffer,
		Header: textproto.MIMEHeader{
			"Content-Disposition": { "form-data; name=\"import_file\"; filename=\"import.xlsx\" " },
			"Content-Type": { "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet" },
		},
	})

	// Create request
	if err != nil {
		return  err
	}

	// Execute request
	_, err = table.client.Do(req)
	if err != nil {
		return err
	}

	return nil
}

func (table *TableService) getFileHash(params *ImportParams) (string, error) {

	// Prepare request url
	path, err := addQueryParams(importPath, params)
	if err != nil {
		return "", err
	}

	req, err := table.client.NewRequest(http.MethodGet, path, nil)
	if err != nil {
		return "", err
	}

	resp, err := table.client.Do(req)
	if err != nil {
		return "", err
	}

	defer resp.Body.Close()
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return "", err
	}

	value, exists := doc.Find("#ilfilehash").Eq(0).Attr("value")
	if !exists {
		return "", ErrFileHash
	}

	return value, nil
}
