package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"sync"
)

type JAMSClient struct {
	Username string
	Password string
	// e.g. https://localhost:6371
	Host          string
	Client        *http.Client
	Authenticated bool
	AccessToken   string `json:"access_token"`
}

func (jams *JAMSClient) Login() error {
	loginUrl := fmt.Sprintf("%s/jams/api/authentication/login", jams.Host)
	if jams.Username == "" || jams.Password == "" || jams.Client == nil {
		return errors.New("username, password, or client not specified")
	}
	postBody := []byte(fmt.Sprintf(`{"username":"%s", "password":"%s"}`, jams.Username, jams.Password))
	request, _ := http.NewRequest("POST", loginUrl, bytes.NewBuffer(postBody))
	request.Header.Set("Content-Type", "application/json; charset=UTF-8")
	response, err := jams.Client.Do(request)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	if response.StatusCode > 299 {
		return errors.New("Received an invalid response from server: " + response.Status)
	}
	err = json.NewDecoder(response.Body).Decode(&jams)
	if err != nil {
		return err
	}
	jams.Authenticated = true
	return nil
}

func (jams *JAMSClient) GetAgents() ([]Agent, error) {
	var body []Agent
	if err := jams.checkLogin(); err != nil {
		return nil, err
	}
	url := fmt.Sprintf("%s/jams/api/agent", jams.Host)
	request, _ := http.NewRequest("GET", url, nil)
	request.Header.Set("Content-Type", "application/json; charset=UTF-8")
	request.Header.Set("Accept", "application/json; charset=UTF-8")

	request.Header.Set("Authorization", "Bearer "+jams.AccessToken)
	response, err := jams.Client.Do(request)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	if response.StatusCode > 299 {
		return nil, errors.New("Received an invalid response from server: " + response.Status)
	}
	err = json.NewDecoder(response.Body).Decode(&body)
	if err != nil {
		return nil, err
	}
	return body, err
}

func (jams *JAMSClient) checkLogin() error {
	if !jams.Authenticated {
		return jams.Login()
	}
	return nil
}

// Use topID of 1 to get all jams folders
func (jams *JAMSClient) GetAllSubFolders(topID int) ([]Folder, error) {
	var body []Folder
	var wg sync.WaitGroup
	if err := jams.checkLogin(); err != nil {
		return nil, err
	}
	url := fmt.Sprintf("%s/jams/api/folder/children?id=%d", jams.Host, topID)
	request, _ := http.NewRequest("GET", url, nil)
	request.Header.Set("Content-Type", "application/json; charset=UTF-8")
	request.Header.Set("Accept", "application/json; charset=UTF-8")

	request.Header.Set("Authorization", "Bearer "+jams.AccessToken)
	response, err := jams.Client.Do(request)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	if response.StatusCode > 299 {
		return nil, errors.New("Received an invalid response from server: " + response.Status)
	}
	err = json.NewDecoder(response.Body).Decode(&body)
	if err != nil {
		return nil, err
	}
	for _, v := range body {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			folders, err := jams.GetAllSubFolders(id)
			if err == nil {
				body = append(body, folders...)
			}
		}(v.ID)
	}
	wg.Wait()
	return body, err

}

func (jams *JAMSClient) GetJobsByFolder(folder int) ([]Job, error) {
	var body []Job
	if err := jams.checkLogin(); err != nil {
		return nil, err
	}
	url := fmt.Sprintf("%s/jams/api/job/folder/%d", jams.Host, folder)
	request, _ := http.NewRequest("GET", url, nil)
	request.Header.Set("Content-Type", "application/json; charset=UTF-8")
	request.Header.Set("Accept", "application/json; charset=UTF-8")

	request.Header.Set("Authorization", "Bearer "+jams.AccessToken)
	response, err := jams.Client.Do(request)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	if response.StatusCode > 299 {
		return nil, errors.New("Received an invalid response from server: " + response.Status)
	}
	err = json.NewDecoder(response.Body).Decode(&body)
	if err != nil {
		return nil, err
	}

	return body, err

}

func (jams *JAMSClient) GetLastHistorybyJob(job int) (HistoricJob, error) {
	var body []HistoricJob
	if err := jams.checkLogin(); err != nil {
		return HistoricJob{}, err
	}
	url := fmt.Sprintf("%s/jams/api/history/job/%d", jams.Host, job)
	request, _ := http.NewRequest("GET", url, nil)
	request.Header.Set("Content-Type", "application/json; charset=UTF-8")
	request.Header.Set("Accept", "application/json; charset=UTF-8")

	request.Header.Set("Authorization", "Bearer "+jams.AccessToken)
	response, err := jams.Client.Do(request)
	if err != nil {
		return HistoricJob{}, err
	}
	defer response.Body.Close()

	if response.StatusCode > 299 {
		return HistoricJob{}, errors.New("Received an invalid response from server: " + response.Status)
	}
	err = json.NewDecoder(response.Body).Decode(&body)
	if err != nil {
		return HistoricJob{}, err
	}
	if len(body) == 0 {
		return HistoricJob{}, errors.New("No history found for this job")
	}
	return body[0], err

}
