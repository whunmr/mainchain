## und query beacon search

Query all BEACONs with optional filters

### Synopsis

Query for all paginated BEACONs that match optional filters:

Example:
$ und query beacon search --moniker beacon1
$ und query beacon search --owner und1chknpc8nf2tmj5582vhlvphnjyekc9ypspx5ay
$ und query beacon search --page=2 --limit=100

```
und query beacon search [flags]
```

### Options

```
      --height int       Use a specific height to query state at (this can error if the node is pruning state)
  -h, --help             help for search
      --moniker string   (optional) filter beacons by name
      --node string      <host>:<port> to Tendermint RPC interface for this chain (default "tcp://localhost:26657")
  -o, --output string    Output format (text|json) (default "text")
      --owner string     (optional) filter beacons by owner address
```

### Options inherited from parent commands

```
      --chain-id string   The network chain ID
```

### SEE ALSO

* [und query beacon](und_query_beacon.md)	 - Querying commands for the beacon module

###### Auto generated by spf13/cobra on 28-Feb-2022
