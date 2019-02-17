# Sync
Sync is a simple file sync utility that replaces the contents of a remote directory with those of a local one using SFTP.
# Build
```sh
go build -o sync cmd/sync/main.go
```

# Usage
Replaces contents of remote directory `destination` with contents of local directory `source`
```sh
sync user:password@host:port destination source
```

For example, to replace the contents of remote directory `/www` with those of local directory `dist`
```sh
sync user:password@example.com:22 /www dist
```
