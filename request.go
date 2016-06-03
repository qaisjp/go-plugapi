package plugapi

import (
	"bufio"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	// log "github.com/Sirupsen/logrus"
	"io"
	"io/ioutil"
	"strings"
	"time"
	// "net/url"
	"net/http"
)

func (plug *PlugDJ) authenticateUser() error {
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
	plug.Log.WithField("_csrf", csrf).Info("found csrf token")

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

func quickread(reader io.ReadCloser) {
	b, _ := ioutil.ReadAll(reader)
	fmt.Printf("%s\n", b)
}

func (plug *PlugDJ) connectSocket() error {
	// Socket connections depend on a few things from plug:
	// - the actual socket url (_gws)
	// - the server time (_st)
	// - the ??some??sort??of??authorization??code?? (_jm)
	// This works in the same way as getting our csrf token
	// but with different prefixes

	// Let's go ahead and grab that!
	resp, err := plug.web.Get(plug.config.BaseURL)
	if err != nil {
		return err
	}

	defer resp.Body.Close()
	variables, err := scanPrefixes(resp.Body, ",_gws", ",_jm", ",_st")
	if err != nil {
		return err
	}

	// We'll use this later
	plug.authCode = variables[",jm"]

	// We don't want to override any forced socket urls...
	if plug.config.SocketURL == "" {
		plug.config.SocketURL = variables[",gws"]
	}

	// Now we want to also store a time offset because
	// one of us (bot or plug) are ahead of each other

	// Format of the time used (manually composed)
	// see: https://golang.org/src/time/format.go
	format := "2006-01-02 15:04:05.000000"
	theirTime, err := time.Parse(format, variables[",_st"])
	if err != nil {
		plug.Log.WithField("_st", variables[",_st"]).Warnln("could not parse correctly")
		return errors.New("plugapi: could not parse _st correctly")
	}

	// offset the time, store it in seconds
	// note: Seconds() returns a float, and int() truncates it
	offset := int(time.Now().Sub(theirTime).Seconds())

	// plugdj runs in their own timezone... valve time
	// now we can use this Location to handle times
	// correctly everywhere on our bot
	plug.location = time.FixedZone("plugdj", offset)

	return errors.New("not implemented yet!")
}

// Get information about ourselves
// func (c *apiClient) loadSession() error {
// 	endpoint := ""

// 	userdata := &User{}
// 	err := c.GetData(endpoint, userdata)
// 	if err != nil {
// 		return err
// 	}
// 	c.plug.User = userdata
// 	return nil
// }

// Get makes a connection and returns the response,
// handling any reauthentications required
// func (c *apiClient) Get(endpoint string) (*http.Response, error) {
// 	target := c.plug.config.BaseURL + endpoint

// 	resp, err := c.client.Get(target)

// 	// If we need authentication
// 	if (err != nil) && (err.(*url.Error).Err == ErrAuthenticationRequired) {
// 		// Try to authenticate
// 		if err := c.authenticate(); err != nil {
// 			return nil, err
// 		}

// 		// Resubmit the request
// 		resp, err = c.client.Get(target)

// 		// Do we still need authentication?
// 		if (err != nil) && (err.(*url.Error).Err == ErrAuthenticationRequired) {
// 			// Don't try again. Return "auth failed"
// 			return nil, ErrAuthentication
// 		}
// 	}

// 	if err != nil {
// 		return nil, err
// 	}

// 	return resp, err
// }

// Post makes a post request with the map provided as json
func (plug *PlugDJ) Post(endpoint string, data map[string]string) (*http.Response, error) {
	body, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	resp, err := plug.web.Post(plug.config.BaseURL+"/_"+endpoint, "application/json", bytes.NewReader(body))
	if err != nil {
		return nil, err
	}

	return resp, nil
}

// Responses by plugtrack are like:
// {"code":200,"message":"OK","data":{...}}
// So we need to parse the responses and pick out the data
type apiData struct {
	Data interface{} `json:"data"`
	Code int         `json:"code"`
}

// GetData allows you to receive info as a struct
// func (c *http.Client) GetData(endpoint string, v interface{}) error {
// 	resp, err := c.Get(endpoint)
// 	if err != nil {
// 		return err
// 	}
// 	defer resp.Body.Close()

// 	data := &apiData{}

// 	switch v.(type) {
// 	case *User:
// 		data.Data = v.(*User)
// 	case *Room:
// 		data.Data = v.(*Room)
// 	default:
// 		return ErrUnknownData
// 	}

// 	// str, err := ioutil.ReadAll(resp.Body)
// 	// fmt.Println(string(str))
// 	// if true {
// 	// 	return nil
// 	// }

// 	err = json.NewDecoder(resp.Body).Decode(data)
// 	if err != nil {
// 		return err
// 	}

// 	if data.Code != 200 {
// 		return &ErrDataRequestError{data, endpoint}
// 	}

// 	return nil
// }
