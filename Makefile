
build: build-server build-simulation

build-server:
	docker build -f server.Dockerfile -t fleetwatch .

build-simulation:
	docker build -f simulator.Dockerfile -t fleetwatch-sim .

run-server:
	docker run -p 6733:6733 fleetwatch

run-simulation:
	docker run --name fleetwatch-sim fleetwatch-sim -host host.docker.internal
	docker cp fleetwatch-sim:/results.txt ./results.txt

clean:
	docker rmi -f fleetwatch fleetwatch-sim