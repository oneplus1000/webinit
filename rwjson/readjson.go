package rwjson

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/oneplus1000/webinit/sdlog"
)

func StringUnmarshalBody(r *http.Request) (string, error) {
	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

func JsonUnmarshalBody(r *http.Request, obj interface{}) error {

	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return err
	}

	err = json.Unmarshal(data, &obj)
	if err != nil {
		sdlog.Err(err, "Unmarshal json data "+string(data))
		return err
	}
	return nil
}
