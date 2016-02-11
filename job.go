package main

const (
	jobFolder        = mompriv + "jobs/"
	scriptFileEnding = ".SC"
)

// Job represents a job
// Example JSON of job
//{
//	"jobid": "36.master",
//	"user": "vagrant",
//	"group": "vagrant",
//	"jobname": "STDIN",
//	"resourcelimits": "neednodes=1,nodes=1,walltime=01:00:00",
//	"jobqueue": "batch",
//	"jobaccount": ""
// }
type Job struct {
	ID             string `json:"jobid"`
	User           string `json:"user"`
	Group          string `json:"group"`
	Name           string `json:"jobname"`
	ResourceLimits string `json:"resourcelimits"`
	Queue          string `json:"jobqueue"`
	Account        string `json:"jobaccount"`
	Cmd            string `json:"cmd"`
}

// GetScript gets the Script of the job
func (j *Job) GetScript() *Script {
	return &Script{jobFolder + j.ID + scriptFileEnding, nil}
}
