package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
)

// RestClient manages the REST interface for a calling user.
type RestClient struct {
	serverURL url.URL
	apiToken  string
}

// MakeRestClient is the factory for constructing a RestClient for a given endpoint
func MakeRestClient(url url.URL, apiToken string) RestClient {
	return RestClient{
		serverURL: url,
		apiToken:  apiToken,
	}
}

// extractError checks if the response signifies an error (for now, StatusCode != 200).
// If so, it returns the error.
// Otherwise, it returns nil.
func extractError(resp *http.Response) error {
	if resp.StatusCode == 200 {
		return nil
	}

	errorBuf, _ := ioutil.ReadAll(resp.Body) // ignore returned error
	return fmt.Errorf("HTTP %v: %s", resp.Status, errorBuf)
}

// submitForm is a helper used for submitting (ex.) GETs and POSTs to the server
func (client RestClient) submitForm(response interface{}, path string, request interface{}, requestMethod string, encodeJSON bool) error {
	var err error
	queryURL := client.serverURL
	queryURL.Path = path

	var req *http.Request
	var body io.Reader

	if request != nil {
		//if rawRequestPaths[path] {
			reqBytes, ok := request.([]byte)
			if !ok {
				return fmt.Errorf("couldn't decode raw request as bytes")
			}
			body = bytes.NewBuffer(reqBytes)
		//} else {
		//	v, err := query.Values(request)
		//	if err != nil {
		//		return err
		//	}
		//
		//	queryURL.RawQuery = v.Encode()
		//	if encodeJSON {
		//		jsonValue, _ := json.Marshal(request)
		//		body = bytes.NewBuffer(jsonValue)
		//	}
		//}
	}

	req, err = http.NewRequest(requestMethod, queryURL.String(), body)
	if err != nil {
		return err
	}

	httpClient := &http.Client{}
	resp, err := httpClient.Do(req)

	if err != nil {
		return err
	}

	defer resp.Body.Close()

	err = extractError(resp)
	if err != nil {
		return err
	}

	dec := json.NewDecoder(resp.Body)
	return dec.Decode(&response)
}

// get performs a GET request to the specific path against the server
func (client RestClient) get(response interface{}, path string, request interface{}) error {
	return client.submitForm(response, path, request, "GET", false /* encodeJSON */)
}

// post sends a POST request to the given path with the given request object.
// No query parameters will be sent if request is nil.
// response must be a pointer to an object as post writes the response there.
func (client RestClient) post(response interface{}, path string, request interface{}) error {
	return client.submitForm(response, path, request, "POST", true /* encodeJSON */)
}

func (client RestClient) FetchVersions(response SolcVersion) error {
	err := client.get(&response, "/solc_versions", nil)
	return err
}

func (client RestClient) FetchVersion(response SolcVersion, version string) (SolcBuild, error) {
	err := client.get(&response, "/solc_versions", nil)
	if err != nil {
		return SolcBuild{}, err
	}
	for _, b := range response.Builds {
		if b.Version == version {
			return b, nil
		}
	}

	return SolcBuild{}, fmt.Errorf("given version not found")
}