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

	$ go get bitbucket.org/dizk/cbatch
	$ cd $GOPATH/src/bitbucket.org/dizk/cbatch
    $ glide install
	$ go build

After building you can run it on a torque cluster provided by vagrant.

Install
-------
Move `config/config.toml` to `/etc/cbatch.toml` and edit to fit your system.

    $ sudo cp config.toml /etc/cbatch.toml

Move `cbatch` to some bin folder.

    $ sudo cp cbatch /usr/bin/cbatch

Add `$jobstarter` option to your Torque mom config.

    $ echo '$job_starter /usr/bin/cbatch' >> /var/spool/torque/mom_priv/config
    $ echo '$job_starter_run_privileged true' >> /var/spool/torque/mom_priv/config


Development Requirements
-------------------------
- Vagrant [https://www.vagrantup.com/](https://www.vagrantup.com/)

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


Vagrant setup is built upon crcollins/torquecluster