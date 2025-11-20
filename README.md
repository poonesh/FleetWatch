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

# Write Up

The entire project took me 3 hours to get a basic example working. I then took a bit more time to write tests and set up a build system with Docker. The most difficult part of the problem was figuring out the overall design I wanted to use, including how devices would be stored and managed. Once I decided on a design of a central device manager, it came together nicely.

To modify the data model, I would extend the DeviceData struct in `devices.go`. This could be modified to include other metrics. Additionally, new methods could be added to the `DeviceManager` class to record or compute these metrics. If the metrics were significantly more complex, I would consider adding new classes or managers the `DeviceManager` could interact with.

The runtime complexity of my solution is O(1) for both adding and retrieving data.
Since first and last heartbeats and upload time are updated with every heartbeat, we do not need to iterate over the full heartbeat or upload time array each time we compute the relevant statistic, which would be O(n). I considered removing these arrays entirely, as they were unnecessary for computing the relevant values, but kept them in in case the requirements were to change in the future and the data became useful.
We could also consider just keeping the last N upload times/heartbeats to reduce the necessary storage while still keeping some data.
Additionally, I handle the case where heartbeats arrive out of order by always checking against the first and last heartbeat times. In a real world scenario, I would clarify if this should be considered an error.

