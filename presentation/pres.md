#Integrating Linux Containers into Batch systems

---

##Goals
- Use containers to create homogeneous job environments.
- Give users reproducible and predictable environments independent of host configuration.
- Integrate well with existing infrastructure.

---

##Implementation
- Implemented in Go
- Supports Docker images.
- Has to be installed on each worker node.
- Supports Torque as the resource manager.
    - Uses the `$jobstarter` option in MOM config.

Note: Go because Docker, rkt, libcontainer and easy to call C code.

---

##How does it work
- Prepare container.
- Start and attach to container.
- Clean up.

---

##Prepare container
- Retrieve image if needed.
- Create a container based on the image.
- Mount files on the host into container:
    - `/etc/passwd` and `/etc/group`
    - `$HOME` (for the user running the job)
    - The batch job.

---

##Start the container
- Start the container
- Stdin, Stdout and Stderr is passed on from jobstarter to container.

---

##Clean up
- Detach from the container.
- Remove the container.

---

##Further work
- Create an image for running jobs.
- Test
    - Complex jobs
    - Lengthy jobs
- Evaluate security risks
- Support other container implementations
    - LXC
    - rkt
    - runC


---

##Help wanted

- How can the end user pass config to the jobstarter?
    - Image to use.
    - Where to get the image.
- Complex example jobs.
- To what extent can containers namespaces and cgroups interfere with existing infrastructure?
- Make it possible for containers to communicate (jobs with several processes).

---

#Questions
