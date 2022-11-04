package main

import (
	"sync"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/sirupsen/logrus"
)

// All metrics should be placed here
var (
	jobStart = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "jams_job_start_time",
		Help: "Unix time of job Start",
	}, []string{"job_name"})
	JobComplete = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "jams_job_completion_time",
		Help: "Unix time in milliseconds of job completion",
	}, []string{"job_name"})
	jobStatusCode = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "jams_job_status_code",
		Help: "Unix time in milliseconds of job Start",
	}, []string{"job_name"})
	jobHistoryCode = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "jams_job_history_id",
		Help: "The historicID of last job run",
	}, []string{"job_name"})
	totalJobsMec = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "jams_jobs_count_total",
		Help: "Total number of jobs in JAMS",
	})
	totalAgentsMec = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "jams_agents_count_total",
		Help: "Total number of agents in JAMS",
	})
	activeAgents = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "jams_agents_active_count_total",
		Help: "Total number of agents *online* in JAMS",
	})
	jobLimit = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "jams_job_limit_count",
		Help: "Job limit for Jams agent",
	}, []string{"agent_name"})
	jobCount = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "jams_agent_job_count",
		Help: "Job count for Jams agent",
	}, []string{"agent_name"})
	totalFoldersMec = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "jams_folders_count_total",
		Help: "Total number of folders in JAMS",
	})
	jamsScrapeDuration = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "jams_scrape_duration",
		Help: "Total number of milliseconds to scrape JAMS",
	})
)

func registerJAMSMetrics(reg *prometheus.Registry, client JAMSClient) {
	start := time.Now()
	var wg sync.WaitGroup

	wg.Add(1)
	go func(client JAMSClient) {
		defer wg.Done()
		agentsMec(reg, &client)
	}(client)
	folders := foldersMec(reg, &client)
	for _, i := range folders {
		wg.Add(1)
		go func(client JAMSClient, folderID int) {
			defer wg.Done()
			jobs, err := client.GetJobsByFolder(folderID)
			if err != nil {
				logrus.Warnln("Error getting jobs for folder", folderID, err)
				return
			}
			totalJobsMec.Set(float64(len(jobs)))

			for _, j := range jobs {
				wg.Add(1)
				go func(client JAMSClient, reg *prometheus.Registry, j Job) {
					defer wg.Done()
					historicJobMec(reg, &client, &j)
				}(client, reg, j)
			}

		}(client, i.ID)
	}
	wg.Wait()
	total := time.Since(start)
	jamsScrapeDuration.Set(float64(total.Milliseconds()))
	reg.MustRegister(jamsScrapeDuration)
	reg.MustRegister(totalAgentsMec)
	reg.MustRegister(activeAgents)
	reg.MustRegister(jobLimit)
	reg.MustRegister(jobCount)
	reg.MustRegister(totalJobsMec)
	reg.MustRegister(totalFoldersMec)
	reg.MustRegister(jobStart)
	reg.MustRegister(jobHistoryCode)
	reg.MustRegister(JobComplete)
	reg.MustRegister(jobStatusCode)

}

func agentsMec(reg *prometheus.Registry, client *JAMSClient) {

	agents, err := client.GetAgents()
	if err != nil {
		logrus.Error(err)
	}
	logrus.Debugf("Found %d agents", len(agents))
	for _, agent := range agents {
		activeAgents.Set(0)
		if agent.Online {
			activeAgents.Add(1)
		}
		jobLimit.With(prometheus.Labels{"agent_name": agent.Name}).Set(float64(agent.JobLimit))
		jobCount.With(prometheus.Labels{"agent_name": agent.Name}).Set(float64(agent.JobCount))
	}
	totalAgentsMec.Set(float64(len(agents)))

}

func foldersMec(reg *prometheus.Registry, client *JAMSClient) []Folder {

	folders, err := client.GetAllSubFolders(1)
	if err != nil {
		logrus.Error(err)
		return nil
	}

	logrus.Debugf("Found %d folders", len(folders))

	totalFoldersMec.Set(float64(len(folders)))
	return folders
}

func historicJobMec(reg *prometheus.Registry, client *JAMSClient, job *Job) {
	historicJob, err := client.GetLastHistorybyJob(job.ID)
	if err != nil {
		logrus.Warnln("Error getting history for job ", job.ID, err)
	}
	jobStart.With(prometheus.Labels{"job_name": historicJob.JobName}).Set(float64(historicJob.JobStart.UnixMilli()))
	JobComplete.With(prometheus.Labels{"job_name": historicJob.JobName}).Set(float64(historicJob.JobComplete.UnixMilli()))
	jobHistoryCode.With(prometheus.Labels{"job_name": historicJob.JobName}).Set(float64(historicJob.JobFinalStatusCode))
	jobStatusCode.With(prometheus.Labels{"job_name": historicJob.JobName}).Set(float64(historicJob.JobFinalStatusCode))

}
