package conf

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"reflect"
)

type config struct {
	CSRF_KEY_32 string
}

var C *config

// Perhaps the best would be for the application not to relly
// on env variables. But that's not always viable.
func UpdateEnv(c *config) {
	fmt.Printf("that! %+v\n", c)

	v := reflect.ValueOf(c).Elem()
	for i := 0; i < v.NumField(); i += 1 {
		os.Setenv(v.Type().Field(i).Name, v.Field(i).String())
	}
}

func ReadFromJSON(path string, c *config) {
	body, err := ioutil.ReadFile(path)
	if err != nil {
		fmt.Printf("File error: %v\n", err)
		os.Exit(1)
	}

	if err := json.Unmarshal(body, &c); err != nil {
		panic("Failed to read config JSON.\n" + err.Error())
	}
}

func Setup() {
	C = new(config)
	ReadFromJSON("./env.json", C)
	UpdateEnv(C)
}
