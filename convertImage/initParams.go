package convertImage

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

func initParams(arg string, inForm string, outForm string) (Params, error) {
	var err error
	var params Params

	params.Inform, params.Outform, err = setFileFormat(inForm, outForm)
	if err != nil {
		return Params{}, err
	}
	files, err := walkImageDir(arg, params.Inform)
	if err != nil {
		return Params{}, err
	}
	params.Infile = make([]*os.File, len(files))
	params.Outfile = make([]*os.File, len(files))
	for i, file := range files {
		params.Infile[i], err = os.Open(file)
		if err != nil {
			return Params{}, getError(err)
		}
		params.Outfile[i], err = os.Create(file[:len(file)-len(filepath.Ext(file))] + params.Outform)
		if err != nil {
			return Params{}, getError(err)
		}
	}
	return params, nil
}

func setFileFormat(inForm, outForm string) (string, string, error) {
	var newForm [2]string

	for i, form := range []string{inForm, outForm} {
		form = "." + form
		switch {
		case form == PNG:
			newForm[i] = PNG
		case form == JPEG:
			newForm[i] = JPEG
		case form == PGM && i == 1:
			newForm[i] = PGM
		case form == GIF:
			newForm[i] = GIF
		default:
			return "", "", fmt.Errorf("error: %s: invalid format", form[1:])
		}
	}
	if inForm == outForm {
		return "", "", fmt.Errorf("error: %s: input and output formats are same", inForm)
	}
	return newForm[0], newForm[1], nil
}

func walkImageDir(dir string, form string) ([]string, error) {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return confirmFileCondition(dir, form, err)
	}

	var ImageFileNames []string
	for _, file := range files {
		searchPath := filepath.Join(dir, file.Name())
		if file.IsDir() {
			subDirFiles, err := walkImageDir(searchPath, form)
			if err != nil {
				return nil, getError(err)
			}
			ImageFileNames = append(ImageFileNames, subDirFiles...)
			continue
		}
		switch filepath.Ext(file.Name()) {
		case form:
			ImageFileNames = append(ImageFileNames, searchPath)
		case JPEG, PNG, PGM:
			continue
		default:
			return nil, fmt.Errorf("error: %s is not a valid file", searchPath)
		}
	}
	return ImageFileNames, nil
}

func confirmFileCondition(fileName string, format string, dirErr error) ([]string, error) {
	stat, err := os.Stat(fileName)

	if err != nil {
		return nil, getError(err)
	}
	if stat.IsDir() {
		return nil, getError(dirErr)
	}
	if filepath.Ext(fileName) != format {
		return nil, fmt.Errorf("error: %s is not a valid file", fileName)
	}
	return []string{fileName}, nil
}

func getError(err error) error {
	newErrorMsg := err.Error()
	if strings.Contains(err.Error(), " ") {
		newErrorMsg = "error:" + err.Error()[strings.Index(err.Error(), " "):]
	}
	return fmt.Errorf("%s", newErrorMsg)
}