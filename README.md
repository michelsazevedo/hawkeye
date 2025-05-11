# Hawkeye

Hawkeye is a scalable and high-performance micro-service designed for efficient and flexible search operations built with [Go](https://golang.org/). It enables efficient [indexing and searching](https://www.elastic.co/enterprise-search) of data, supporting real-time updates and high-performance queries.

The service is designed with a clean [hexagonal architecture](https://netflixtechblog.com/ready-for-changes-with-hexagonal-architecture-b315ec967749), promoting separation of concerns, ease of testing, and adaptability for future extensions. It seamlessly integrates with a broader event-driven ecosystem via [NATS](https://nats.io/) to react to new data creation events and update its search index in near real-time.

## Instructions

### Running with Docker
[Docker](www.docker.com) is an open platform for developers and sysadmins to build, ship, and run distributed applications, whether on laptops, data center VMs, or the cloud.

If you haven't used Docker before, it would be good idea to read this article first: Install [Docker Engine](https://docs.docker.com/engine/installation/)

1. Install [Docker](https://www.docker.com/what-docker) and then [Docker Compose](https://docs.docker.com/compose/):

2. Run `docker compose build --no-cache` to build the images for the project.

3. Finally, run the local app with `docker compose up web` and hawkeye will perform requests.

4. Aaaaand, you can run the automated tests suite running a `docker compose run --rm test` with no other parameters!

## License
Copyright © 2025
