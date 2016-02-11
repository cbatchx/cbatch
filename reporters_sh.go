package main

// http://docs.adaptivecomputing.com/torque/6-0-0/help.htm#topics/torque/13-appendices/scriptEnvironment.htm

const prologuesh = `#!/bin/sh
curl -s -H "Content-Type: application/json" -X POST -d "{ \"jobid\":\"$1\", \"user\":\"$2\", \"group\":\"$3\", \"jobname\":\"$4\", \"resourcelimits\":\"$5\", \"jobqueue\":\"$6\", \"jobaccount\":\"$7\" }" http://localhost:8080/new
`

const epiloguesh = `#!/bin/sh
curl -s -H "Content-Type: application/json" -X POST -d "{ \"jobid\":\"$1\", \"user\":\"$2\", \"group\":\"$3\", \"jobname\":\"$4\", \"sessionid\":\"$5\", \"resourcelimits\":\"$6\", \"resourcesused\":\"$7\", \"jobqueue\":\"$8\", \"jobaccount\":\"$9\", \"jobexitcode\":\"$10\" }" http://localhost:8080/done
`

// TODO: need to actually start docker :D
const jobstartersh = `#!/bin/sh
resp=$(curl -s -H "Content-Type: application/json" -X POST -d "{ \"cmd\":\"$*\", \"jobid\":\"$PBS_JOBID\" }" http://localhost:8080/exec)
eval $resp
# Check if we need to exit with error
if [ $DOC_ERR != "nil" ]
  then
    export info="Encountered some error"
    $*
    exit 1
fi
# Run in container
if [ -n $DOC_CONTAINER ]
  then
    export info="Would have started docker"
    $*
fi
# Run on host
if [ -z $DOC_CONTAINER ]
  then
    export info="Normal job"
    $*
fi
`
