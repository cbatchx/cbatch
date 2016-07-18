package main

import (
	"os"

	log "github.com/Sirupsen/logrus"
)

// This file parses all environment variables coming from PBS.

// PBSJob Represents a PBS job. This is just all known Environment variables.
type PBSJob struct {
	Version     string // PBS_VERSION
	Jobname     string // PBS_JOBNAME
	Environment string // PBS_ENVIRONMENT
	TaskNum     string // PBS_TASKNUM
	Walltime    string // PBS_WALLTIME
	MOMPort     string // PBS_MOMPORT
	GpuFile     string // PBS_GPUFILE
	JobCookie   string // PBS_JOBCOOKIE
	NodeNum     string // PBS_NODENUM
	NumNodes    string // PBS_NUM_NODES
	JobID       string // PBS_JOBID
	VNodeNum    string // PBS_VNODENUM
	Queue       string // PBS_QUEUE
	MicFile     string // PBS_MICFILE
	NP          string // PBS_NP
	NumPPn      string // PBS_NUM_PPN
	NodeFile    string // PBS_NODEFILE
	Origin      OriginInfo
}

// OriginInfo Environment variables prefixed with PBS_O (assuming it means origin.)
type OriginInfo struct {
	Queue   string // PBS_O_QUEUE
	Logname string // PBS_O_LOGNAME
	Lang    string // PBS_O_LANG
	Host    string // PBS_O_HOST
	Shell   string // PBS_O_SHELL
	Home    string // PBS_O_HOME
	Workdir string // PBS_O_WORKDIR
	Mail    string // PBS_O_MAIL
	Server  string // PBS_O_SERVER
	Path    string // PBS_O_PATH
}

func getNewJobFromEnv() *NewJob {
	origin := OriginInfo{
		Queue:   os.Getenv("PBS_O_QUEUE"),
		Logname: os.Getenv("PBS_O_LOGNAME"),
		Lang:    os.Getenv("PBS_O_LANG"),
		Host:    os.Getenv("PBS_O_HOST"),
		Shell:   os.Getenv("PBS_O_SHELL"),
		Home:    os.Getenv("PBS_O_HOME"),
		Workdir: os.Getenv("PBS_O_WORKDIR"),
		Mail:    os.Getenv("PBS_O_MAIL"),
		Server:  os.Getenv("PBS_O_SERVER"),
		Path:    os.Getenv("PBS_O_PATH"),
	}

	job := PBSJob{
		Version:     os.Getenv("PBS_VERSION"),
		Jobname:     os.Getenv("PBS_JOBNAME"),
		Environment: os.Getenv("PBS_ENVIRONMENT"),
		TaskNum:     os.Getenv("PBS_TASKNUM"),
		Walltime:    os.Getenv("PBS_WALLTIME"),
		MOMPort:     os.Getenv("PBS_MOMPORT"),
		GpuFile:     os.Getenv("PBS_GPUFILE"),
		JobCookie:   os.Getenv("PBS_JOBCOOKIE"),
		NodeNum:     os.Getenv("PBS_NODENUM"),
		NumNodes:    os.Getenv("PBS_NUM_NODES"),
		JobID:       os.Getenv("PBS_JOBID"),
		VNodeNum:    os.Getenv("PBS_VNODENUM"),
		Queue:       os.Getenv("PBS_QUEUE"),
		MicFile:     os.Getenv("PBS_MICFILE"),
		NP:          os.Getenv("PBS_NP"),
		NumPPn:      os.Getenv("PBS_NUM_PPN"),
		NodeFile:    os.Getenv("PBS_NODEFILE"),
		Origin:      origin,
	}

	args := os.Args[1:]

	log.WithFields(log.Fields{
		"pbs_env": job,
		"args":    args,
	}).Info("New job")

	return &NewJob{
		PBSJob: job,
		Args:   args,
	}
}
