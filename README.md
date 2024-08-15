Directory watcher for s3, hdfs and local fs with websockets.
----
Running the project for local fs:
```bash
# master application
go run main.go --type local --path ./ --server 0.0.0.0:8097
```
```bash
# client application
go run client.go --server server.ip:8097
```