Script for generating configuration file `clientBinaries.json` for Mist and Ethereum Wallet. To use with go-ethereum forks and customizations of Mist for them. Swarm section is currently being skipped.

#### HOWTO

In your `go-ethereum` directory build archives with geth:

```
make geth-linux-386 geth-windows-386 geth-linux-amd64 geth-windows-amd64 geth-darwin-amd64
./build/env.sh go run build/ci.go archive -arch windows-4.0-amd64
./build/env.sh go run build/ci.go archive -arch windows-4.0-386
./build/env.sh go run build/ci.go archive -arch linux-amd64 -type tar
./build/env.sh go run build/ci.go archive -arch linux-386 -type tar
./build/env.sh go run build/ci.go archive -arch darwin-10.6-amd64 -type tar
```

Then run the script:

```
go run *.go -url https://example.com -path ~/path_to_go-ethereum/ -version 1.0 -binary geth
```

`clientBinaries.json` will be saved in current directory. Upload that json and geth archives to your website and configure Mist.
