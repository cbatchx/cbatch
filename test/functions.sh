# Submit a job.
# Job is passed as a string.
function job_submit() {    
    echo $($1 | qsub)
}

# Takes a job an prints the status of it.
function job_status() {
    echo $(qstat -f $1)
}