# FleetWatch

FleetWatch is a demo server handling statistics for a fleet of devices.

# Instructions

Run `make build` to build the server and simulator docker containers.

Then, run
`make run-server` to run the server

and

`make run-simulation` to run the simulation, which will output the results into `results.txt`.

Note that running the simulation multiple times without restarting the server
will give incorrect results as devices are not cleared between runs.

If you do not have docker, you can run the server with `go run .`


