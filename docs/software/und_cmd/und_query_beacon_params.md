## und query beacon params

Query the current Beacon parameters

```
und query beacon params [flags]
```

### Options

```
      --height int      Use a specific height to query state at (this can error if the node is pruning state)
  -h, --help            help for params
      --node string     <host>:<port> to Tendermint RPC interface for this chain (default "tcp://localhost:26657")
  -o, --output string   Output format (text|json) (default "text")
```

### Options inherited from parent commands

```
      --chain-id string   The network chain ID
```

### SEE ALSO

* [und query beacon](und_query_beacon.md)	 - Querying commands for the beacon module

###### Auto generated by spf13/cobra on 28-Feb-2022
