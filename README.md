# Generate Proto File

```
# cd to your proto dir
cd ./proto

# generate proto file
protoc *.proto --go_out=plugins=grpc:. --go_opt=paths=source_relative
```