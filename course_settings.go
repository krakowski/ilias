package ilias

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"net/textproto"
	"net/url"
)

const (
	settingsPath = "ilias.php?cmdClass=ilobjcoursegui&cmd=post&cmdNode=v6:ko&baseClass=ilrepositorygui&fallbackCmd=update"
)

type SettingsParams struct {
	Reference	string  `schema:"ref_id"`
	Token		string	`schema:"rtoken"`
}

func toDirection(sortingMode string) string {
	if sortingMode == "ascending" {
		return "0"
	}

	if sortingMode == "descending" {
		return "1"
	}

	return "0"
}

func setMandatorySettings(values url.Values, config *CourseSettings) {
	values["title"] = []string{config.Title}

	if config.EcsExport {
		values["ecs_export"] = []string{"1"}
	} else {
		values["ecs_export"] = []string{"0"}
	}

	direction := toDirection(config.Presentation.Sorting.Mode)
	switch config.Presentation.Sorting.Mode {
	case "alphabetical":
		values["sorting"] = []string{"0"}
		values["title_sorting_direction"] = []string{direction}
		break
	case "creation":
		values["sorting"] = []string{"4"}
		values["creation_sorting_direction"] = []string{direction}
		break
	case "deadline":
		values["sorting"] = []string{"2"}
		values["activation_sorting_direction"] = []string{direction}
		break
	}
}

func setOptionalSettings(values url.Values, config *CourseSettings) {
	if presentationMode := config.Presentation.Mode; presentationMode != nil {
		switch *presentationMode {
		case "tile":
			values["list_presentation"] = []string{"tile"}
			break
		case "list":
			values["list_presentation"] = []string{""}
			break
		}
	}
}

func (course *CourseService) SynchronizeSettings(params *SettingsParams, settings *CourseSettings) error {

	// Prepare request url
	path, err := addQueryParams(settingsPath, params)
	if err != nil {
		return err
	}

	values := url.Values{
		"cmd[Update]": {"Speichern"},
	}

	setMandatorySettings(values, settings)
	setOptionalSettings(values, settings)
	
	content := bytes.Buffer{}
	req, err := course.client.NewMultipartRequest(http.MethodPost, path, values, &UploadFile{
		Content:   &content,
		Header: textproto.MIMEHeader{
			"Content-Disposition": { "form-data; name=\"tile_image\"; filename=\"\" " },
			"Content-Type": { "application/octet-stream" },
		},
	})

	if err != nil {
		return err
	}

	// Retrieve the HTML source
	resp, err := course.client.Do(req)
	if err != nil {
		return err
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	print(string(body))

	//if !strings.HasPrefix(string(body), "{\"result\":true") {
	//	return ErrUpdate
	//}

	return nil
}