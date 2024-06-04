# CLI-Query-Performance
Command line tool for benchmarking query performance across multiple workers against a TimescaleDB instance.

Important to note regarding specific implementations: This project was created in a way to be as easy as possible to run and test. As such, certain security decision were made. If a similar deployment is ever used, ensure to create a .env file to hold necessary secrets and parameterize necessary docker and compose files to utilize. The database is also initialized to contain test data so take care to edit the compose file if that is not desired.

# Setup

Once you have cloned the repository, first ensure you have the necessary applications installed (docker, docker compose). Then run the following

    docker compose up -d

Ensure that the database is finished initalizing before attempting to run anything.

# Execution

For ease, the command line tool container is set to stay up run, thus to test the program, the format of the command will be

    docker exec -it cli-tool CLI-query-performance --workers=$numberOfWorkers $sourceOfCSVData

where:
    - numberOfWorkers = number of workers you'd like the tool to use while running the queries
    - sourceOfCSVData = either a csv-formatted file or string in csv-format

Examples:

Using csv file

    docker exec -it cli-tool CLI-query-performance --workers=5 query_params.csv

Using string

    docker exec -it cli-tool CLI-query-performance --workers=3 hostname,start_time,end_time\nhost_000008,2017-01-01 08:59:22,2017-01-01 09:59:22\nhost_000001,2017-01-02 13:02:02,2017-01-02 14:02:02\nhost_000008,2017-01-02 18:50:28,2017-01-02 19:50:28
