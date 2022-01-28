package ilias

import (
	"net/url"
	"reflect"

	"github.com/PuerkitoBio/goquery"
)

func addQueryParams(path string, params interface{}) (string, error) {
	v := reflect.ValueOf(params)
	if v.Kind() == reflect.Ptr && v.IsNil() {
		return path, nil
	}

	url, err := url.Parse(path)
	if err != nil {
		return path, err
	}

	values := url.Query()
	if err := encoder.Encode(params, values); err != nil {
		return path, err
	}

	url.RawQuery = values.Encode()
	return url.String(), nil
}

func containsSuccessAlert(doc *goquery.Document) bool {
	return containsElement(doc, "div .alert-success")
}

func containsDangerAlert(doc *goquery.Document) bool {
	return containsElement(doc, "div .alert-danger")
}

func containsElement(doc *goquery.Document, selector string) bool {
	return doc.Find(selector).Length() > 0
}
