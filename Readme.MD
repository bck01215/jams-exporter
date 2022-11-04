# Jams Exporter

## Why

The jams monitor is difficult to integrate with other alerting tools. By creating an exporter, we can display data in Grafana or anything compatible with prometheus.

## How to use

Use either the Docker container `bkauffman7/jams-exporter or build from source using`go build`

### Flags

|Flag|Definition|Options|
| ------------- | ------------- |  ------------- |
| --help  | show this  |
| --host  | method, host, and port of jams server (<https://localhost:6371>)||
| --skip-verify  | If present, Jams cert will not be validated  ||
| --port  | The port to bind the application to (Defaults to 8000)  ||
| --password, -p  | password of user authenticating to Jams  ||
| --username, -u  | username of user authenticating to Jams  |
| --log-level, -l  | Log level to run the application (Defaults to error)|debug,info,warn,error|
| --skip-verify  | If present, Jams cert will not be validated ||

## Exported metrics

All job metrics are based on the last historic entry for the job. Some metrics may be incomplete based on the level of access your jams user has. Turn on warn level to see warning messages about invalid response codes (such as 401s) from the server.
|Name|Definition|
| ------------- | -------------  |
|jams_job_start_time|Unix time of when the job started|
|jams_job_completion_time|Unix time of when the job completed|
|jams_job_status_code|Status code of exited job|
|jams_job_history_id| Historic ID for the job. This can be used to count changes and get a total times rune|
|jams_jobs_count_total| used by getting the length of all jobs returned by searching folders |
|jams_agents_count_total| Total number of agents found|
|jams_agents_active_count_total| Agents with the status of online|
|jams_job_limit_count| Job limit for each agent|
|jams_agent_job_count| Jobs assigned to each agent|
|jams_folders_count_total| total folders found |
|jams_scrape_duration| Duration of the scrape|

## Authors

- [Brandon Kauffman](mailto:bck01215@gmail.com)
