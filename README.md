cbatch - Containerized Torque
================
A simple torque cluster running the work inside a docker container. Comes with a Vagrant cluster if you want to test it.

And a daemon implemented in Go which manages the containers. The daemon must run on each of the Torque MOM nodes.

________________________________________________________________________________
Build/Run Requirements
----------------------

- Vagrant
- Virtualbox
- Go

Download
--------
You can download and compile the project with:

	$ go get bitbucket.org/dizk/cbatch
	$ cd $GOPATH/src/bitbucket.org/dizk/cbatch
	$ go build

After building you can run it on a torque cluster provided by vagrant.

Run
-----
The following command builds a three part virtual Torque cluster with 1 master host and 2 slaves.

	$ NODES=2 vagrant up

`NODES` defines the number of slave nodes that will be created for the cluster.

Access
------
	$ vagrant ssh master
	$ echo 'echo "Hello Docker"' | qsub

You can also ssh into `slave1`, `slave2` etc. to look at the logs.  


Suspend
-------
	$ NODES=2 vagrant suspend

Take Down
---------
	$ NODES=2 vagrant destroy

________________________________________________________________________________
How it works:
---------------------
The daemon places 3 scripts in the mom_priv folder.

 - prologue
 - jobstarter
 - epilogue

All of them just feed the cbatch information about the job, so cbatch can manage the container for the job.

Jobstart runs in priviledged mode. It executes a `docker exec` to run the job.  

Vagrant setup is built upon crcollins/torquecluster
