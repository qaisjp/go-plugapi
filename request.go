package plugapi

import (
	"bufio"
	"bytes"
	"encoding/json"
	"errors"
	// log "github.com/Sirupsen/logrus"
	"io"
	"net/http"
	"reflect"
	"strings"
	// "net/url"
)

func (plug *PlugDJ) authenticateUser() error {
	// NOTE: We don't use plug.Get because we
	// are not accessing plugdj.com/_/, we actually
	// want plugdj.com/ (so we don't want to use the API)
	resp, err := plug.web.Get(plug.config.BaseURL + "/")
	if err != nil {
		return err
	}

	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return ErrUnknownResponse{resp, "/"}
	}

	// non authenticated requests contain a little js snippet
	// beginning with: var _csrf="<60 char token>""
	// we need to use this for all future requests
	// we use scanPrefixes so ve can deal with vis
	csrfPrefix := "var _csrf"
	results, err := scanPrefixes(resp.Body, csrfPrefix)
	if err != nil {
		return err
	}
	csrf := results[csrfPrefix] // get it outta the map

	// check token length for some validity
	if len(csrf) != 60 {
		plug.Log.WithField("_csrf", csrf).Error("csrf token malformed")
		return errors.New("dubapi: csrf token is malformed")
	}
	plug.Log.WithField("_csrf", csrf).Debugln("found csrf token")

	plug.Log.Info("Attempting to log in...")

	data := map[string]string{
		"csrf":     csrf,
		"email":    plug.config.Email,
		"password": plug.config.Password,
	}

	// try to log in
	resp, err = plug.Post(AuthLoginEndpoint, data)
	if err != nil {
		return err
	}

	quickread(resp.Body)
	resp.Body.Close()

	if resp.StatusCode == 401 {
		return ErrAuthentication
	} else if resp.StatusCode != 200 {
		return ErrUnknownResponse{resp, "/"}
	}

	return nil
}

// scanPrefixes searches a Reader for variables (well, prefixes)
// and returns a map with the variable name (prefix name) as key
// and the result as the value (result must be enclosed in quotes in the reader)
func scanPrefixes(reader io.Reader, variables ...string) (map[string]string, error) {
	variablesLeft := len(variables)
	variablesFound := map[string]string{}

	// we use a scanner so we don't have to read it all at once
	scanner := bufio.NewScanner(reader)
	for scanner.Scan() {
		text := scanner.Text()

		// We go through all the variables, checking if it is in this line
		for _, varName := range variables {
			startPos := strings.Index(text, varName+`="`)
			if startPos == -1 {
				continue
			}

			// add the variable length and the `="` to the start
			startPos += len(varName) + 2

			// awesome. so we found the start point!
			// now we find the endPoint...
			endPos := strings.Index(text[startPos+1:], `"`)
			if endPos == -1 {
				// we should always find an end-quote on the same line...
				return nil, errors.New("plugapi: could not find end token for variable " + varName)
			}

			// add it to the map
			variablesFound[varName] = text[startPos : startPos+endPos+1]
			variablesLeft-- // decrement our counter

			if variablesLeft == 0 {
				// break out if we're done finding stuff
				// we use a goto to avoid two checks and breaks (ugly)
				goto FinishedSearching
			}
		}
	}

FinishedSearching:
	// scans stop on error or file finish. let's check for an error
	if err := scanner.Err(); err != nil {
		return nil, err
	}

	// goto is useless because we need an extra check here anyway
	// this returns the actual map if everything was ok
	if variablesLeft == 0 {
		return variablesFound, nil
	}

	return nil, errors.New("plugapi: could not find all variables")
}

// Get makes a get request to the plug API
func (plug *PlugDJ) Get(endpoint string) (*http.Response, error) {
	resp, err := plug.web.Get(plug.config.BaseURL + "/_" + endpoint)
	if err != nil {
		return nil, err
	}

	// If the status code is not 200, error away
	if resp.StatusCode != 200 {
		resp.Body.Close()
		return nil, ErrUnknownResponse{resp, endpoint}
	}

	return resp, nil
}

// apiResponse is a struct for
// data sent by the plug.dj API
type apiEnvelope struct {
	// Note: why can't data be a []interface{} ??
	// Read https://github.com/golang/go/wiki/InterfaceSlice
	Data   []json.RawMessage `json:"data"`
	Meta   interface{}       `json:"meta"`
	Status string            `json:"status"`
	Time   float32           `json:"time"`
}

// GetData allows you to receive info as a struct
func (plug *PlugDJ) GetData(endpoint string, expected interface{}) (interface{}, error) {
	resp, err := plug.Get(endpoint)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	envelope := &apiEnvelope{}

	// err = json.Unmarshal(quickread(resp.Body), envelope)
	err = json.NewDecoder(resp.Body).Decode(envelope)
	if err != nil {
		return nil, err
	}

	if envelope.Status != "ok" {
		return nil, &ErrDataRequestError{envelope, endpoint}
	}

	outType := reflect.TypeOf(expected)         // a "[]Room", for instance
	objects := reflect.MakeSlice(outType, 0, 0) // make the actual slice
	for _, obj := range envelope.Data {
		// we need to Unmarshal our individual objects (each obj is a json.Message)
		var data interface{}

		// init the object so that Unmarshal knows how to treat it
		switch expected.(type) {
		case []*Room:
			data = &Room{}
		default:
			panic("plugapi: don't know how to handle interface provided")
		}

		err := json.Unmarshal([]byte(obj), data)
		if err != nil {
			return nil, err
		}

		// interface{} -> reflect.Value --converted_to-> Room
		newObj := reflect.ValueOf(data).Convert(outType.Elem())
		objects = reflect.Append(objects, newObj)
	}

	return objects.Interface(), nil
}

// Post makes a post request with the map provided as json to the plug API
func (plug *PlugDJ) Post(endpoint string, data map[string]string) (*http.Response, error) {
	body, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	resp, err := plug.web.Post(plug.config.BaseURL+"/_"+endpoint, "application/json", bytes.NewReader(body))
	if err != nil {
		return nil, err
	}

	// If the status code is not 200, error away
	if resp.StatusCode != 200 {
		resp.Body.Close()
		return nil, ErrUnknownResponse{resp, endpoint}
	}

	return resp, nil
}
