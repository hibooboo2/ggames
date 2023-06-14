package pollen

import (
	"fmt"
	"html/template"
	"os"
	"time"

	"github.com/hibooboo2/ggames/allplay/logger"
)

var (
	templates        = map[string]*template.Template{}
	templateModTimes = map[string]time.Time{}
)

func LoadTemplate(filename string, args ...interface{}) (*template.Template, error) {
	var funcs template.FuncMap
	for _, arg := range args {
		switch a := arg.(type) {
		case template.FuncMap:
			funcs = a
		default:
			logger.Warnf("Unknown template argument type: %T value: %v", arg, arg)
		}
	}
	info, err := os.Lstat(filename)
	if os.IsNotExist(err) {
		logger.Initf("%q not found", filename)
		return nil, fmt.Errorf("template %q not found", filename)
	}

	t, templateExists := templates[filename]
	if templateModTimes[filename] == info.ModTime() && templateExists {
		logger.Boardf("%q is up-to-date", filename)
		return t, nil
	}

	logger.Initf("Lstat %q: %v", filename, err)
	logger.Initf("Parsing %q", filename)
	t, err = template.New("").Funcs(funcs).ParseFiles(filename)
	if err != nil {
		return nil, err
	}

	templates[filename] = t
	templateModTimes[filename] = info.ModTime()

	return t, nil
}
