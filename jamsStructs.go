package main

import "time"

type Agent struct {
	Name     string `json:"agentName"`
	ID       int    `json:"agentID"`
	Online   bool   `json:"online"`
	JobCount int    `json:"jobCount"`
	JobLimit int    `json:"jobLimit"`
}

type Folder struct {
	ID int `json:"folderID"`
}
type Job struct {
	ID int `json:"jobID"`
}
type HistoricJob struct {
	JobName            string    `json:"jobName"`
	JobStart           time.Time `json:"startTimeUTC"`
	JobComplete        time.Time `json:"completionTimeUTC"`
	JobFinalStatusCode int       `json:"finalStatusCode"`
	JobHistoryID       int       `json:"historyID"`
}
