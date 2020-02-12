package xhttp

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type HttpClient struct {
	apiUrl string
}

func NewClient(apiserverUrl string) *HttpClient {
	return &HttpClient{
		apiUrl: apiserverUrl + "/apis/",
	}
}

func (client *HttpClient) GetApiResourceLists() ([]*APIResourceList, error) {
	groupList, err := client.GetApiGroupList()
	if err != nil {
		return nil, err
	}

	var resourceList = []*APIResourceList{}

	for _, group := range groupList.Groups {
		subPath := group.PreferredVersion.GroupVersion
		list, err := client.GetApiResourceListForUrl(subPath)
		if err != nil {
			return nil, err
		}
		resourceList = append(resourceList, list)
	}
	return resourceList, nil
}

func (client *HttpClient) GetApiResourceListForUrl(subPath string) (*APIResourceList, error) {
	req, err := http.NewRequest("GET", client.apiUrl + subPath, nil)
	if err != nil {
		return nil, err
	}

	bytes, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	var resourceList = APIResourceList{}
	err = json.Unmarshal(bytes, &resourceList)
	if err != nil {
		return nil, err
	}
	return &resourceList, nil
}

// GetApiGroupList executes and http get request and receives all available apiGroups from the kubernetes apiserver
func (client *HttpClient) GetApiGroupList() (*APIGroupList, error) {
	req, err := http.NewRequest("GET", client.apiUrl, nil)
	if err != nil {
		return nil, err
	}

	bytes, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	var groupList = APIGroupList{}
	err = json.Unmarshal(bytes, &groupList)
	if err != nil {
		return nil, err
	}
	return &groupList, nil
}

// Do executes an HTTP request and returns the response body.
// Any errors or non-200 status code result in an error.
func (client *HttpClient) Do(req *http.Request) ([]byte, error) {
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("StatusCode: %d, Body: %s", resp.StatusCode, body)
	}

	return body, nil
}
