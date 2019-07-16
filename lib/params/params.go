package params

import (
	"errors"
	"log"
	"os"
	"regexp"
	"strings"

	"github.com/iancoleman/strcase"
)

type Param struct {
	Value interface{}
}

func GetPluginSetting(key string) (param *Param) {
	log.Println("PLUGIN_" + strcase.ToScreamingSnake(key))
	return Parameter(os.Getenv("PLUGIN_" + strcase.ToScreamingSnake(key)))
}

func Parameter(value string) (param *Param) {
	return &Param{Value: string(value)}
}

func (param *Param) Required(message string) *Param {
	switch param.Value.(type) {

	case string:
		if param.Value == "" {
			log.Fatal(errors.New(message))
		}

	default:
		if param.Value == nil {
			log.Fatal(errors.New(message))
		}

	}

	return param
}

func (param *Param) Lower() *Param {
	param.Value = string(strings.ToLower(param.Value.(string)))
	return param
}

func (param *Param) Snake() *Param {
	param.Value = string(strcase.ToSnake(param.Value.(string)))
	return param
}

func (param *Param) MustMatch(regex string, message string) *Param {
	log.Println(param.Value.(string))
	re := regexp.MustCompile(regex)
	if re.MatchString(param.Value.(string)) != true {
		log.Fatal(errors.New(message))
	}
	return param
}

func (param *Param) Array() *Param {
	param.Value = strings.Split(param.Value.(string), ",")
	return param
}
