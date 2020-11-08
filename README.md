# example-service
a simple gRPC server with google cloud profiler integration

## Running Server
You can run the server with 
```bash
go run main.go -p $GOOGLE_CLOUD_PROJECT_ID -v $SERVICE_VERSION
```

## Running client
Once the server is running, you can then run the client inside `tools` folder,
since cloud profiler samples one profile every minute, you can call the client 
multiple times to ensure samples has functions captured, e.g., 
```bash
for i in $(seq 1 100); do go run tools/connect.go ; done
```

