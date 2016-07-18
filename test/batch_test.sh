#!/bin/bash

source functions.sh

NUM_JOBS=$1

COUNTER=$NUM_JOBS

until [  $COUNTER -lt 1 ]; do
    
    
    job_submit "echo hello"
    echo Submitted job $(($NUM_JOBS-$COUNTER+1)) of $NUM_JOBS
    
    
    let COUNTER-=1
done

