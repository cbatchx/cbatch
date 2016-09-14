cbatch - Containerized Torque
================
Run your batch jobs in containers, created dynamically at run time. Without changes to your jobs.

https://asciinema.org/a/e9gfcdnbrt5e6hb1jor5lp8x8
________________________________________________________________________________
Build Requirements
----------------------

- Go [https://golang.org/](https://golang.org/)
- Glide [https://github.com/Masterminds/glide](https://github.com/Masterminds/glide)

Build
------
You can download and compile the project with:

	$ go get github.com/cbatchx/cbatch
	$ cd $GOPATH/src/github.com/cbatchx/cbatch
    $ glide install
	$ go build

After building you can run it on a torque cluster provided by vagrant.

Install
-------
Download .tar file from releases section of github.

    $ sudo mkdir -p /var/lib/cbatch
	$ sudo tar -xvf cbatch-0.0.6-amd64linux.tar -C /var/lib/cbatch

Edit `/var/lib/cbatch/config/config.toml` to fit your system.

Add `$jobstarter` option to your Torque mom config.

    $ echo '$job_starter /var/lib/cbatch/cbatch' >> /var/spool/torque/mom_priv/config
    $ echo '$job_starter_run_privileged true' >> /var/spool/torque/mom_priv/config

Development Requirements without Torque
--------------------------
- Docker installed locally

In this mode you simply give cbatch some environment variables to "mock" a submitted torque job.


Run without Torque
------------------

To run cbatch without Torque some environment variables must be "mocked". 

* `PBS_O_LOGNAME` must be set to a user on the system.
* `PBS_O_HOME` must be set to the home directory of the user.
* `PBS_JOBID` must be set to the name of one of the jobs in `test/jobs` without .SC

The config.toml file must also be configured to make it run. See the `config/config.toml`
	
You must also supply which shell to run the job as and the job as parameters to cbatch as Torque normally supplies cbatch with these.
	
	$ go build
	$ PBS_O_LOGNAME=$USER PBS_O_HOME=$HOME PBS_JOBID=hello ./cbatch /bin/bash `pwd`/test/jobs/hello.SC


You can also make the job interactive by running (you will end up with a shell in the container): 
	$ PBS_O_LOGNAME=$USER PBS_O_HOME=$HOME PBS_JOBID=hello PBS_ENVIRONMENT=PBS_INTERACTIVE ./cbatch /bin/bash


Development Requirements with Torque
-------------------------
- Vagrant [https://www.vagrantup.com/](https://www.vagrantup.com/)

Run with Torque
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


Vagrant setup is built upon crcollins/torquecluster