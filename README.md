# CLI-Query-Performance
Command line tool for benchmarking query performance across multiple workers against a TimescaleDB instance.

Important to note regarding specific implementations: This project was created in a way to be as easy as possible to run and test. As such, certain security decision were made. If a similar deployment is ever used, ensure to create a .env file to hold necessary secrets and parameterize necessary docker and compose files to utilize. The database is also initialized to contain test data so take care to edit the compose file if that is not desired.

# Setup

Once you have cloned the repository, first ensure you have the necessary applications installed (docker, docker compose). Then run the following

    docker compose up -d

Ensure that the database is finished initalizing before attempting to run anything.

# Execution

For ease, the command line tool container is set to stay up running, thus to test the program, the format of the command will be

    docker exec -it cli-tool CLI-query-performance --workers=$numberOfWorkers $sourceOfCSVData

where:

numberOfWorkers = number of workers you'd like the tool to use while running the queries
    
sourceOfCSVData = either a csv-formatted file or string in csv-format

Examples:

Using csv file

    docker exec -it cli-tool CLI-query-performance --workers=5 query_params.csv

Using string

    docker exec -it cli-tool CLI-query-performance --workers=3 "hostname,start_time,end_time
    host_000008,2017-01-01 08:59:22,2017-01-01 09:59:22
    host_000001,2017-01-02 13:02:02,2017-01-02 14:02:02
    host_000008,2017-01-02 18:50:28,2017-01-02 19:50:28
    host_000002,2017-01-02 15:16:29,2017-01-02 16:16:29
    host_000003,2017-01-01 08:52:14,2017-01-01 09:52:14
    host_000002,2017-01-02 00:25:56,2017-01-02 01:25:56
    host_000008,2017-01-01 07:36:28,2017-01-01 08:36:28
    host_000000,2017-01-02 12:54:10,2017-01-02 13:54:10
    host_000005,2017-01-02 11:29:42,2017-01-02 12:29:42
    host_000006,2017-01-02 01:18:53,2017-01-02 02:18:53
    host_000000,2017-01-02 15:44:45,2017-01-02 16:44:45"

Note: csv-format must follow (hostname,start_time,end_time). To see more details about the database format, check db_setup\100_cpu_usage.sql.
