## und init

Initialize private validator, p2p, genesis, and application configuration files

### Synopsis

Initialize validators's and node's configuration files.

```
und init [moniker] [flags]
```

### Options

```
      --chain-id string   genesis file chain-id, if left blank will be randomly created
  -h, --help              help for init
      --home string       node's home directory (default "/home/hodge/.und_mainchain")
  -o, --overwrite         overwrite the genesis.json file
      --recover           provide seed phrase to recover existing key instead of creating
```

### SEE ALSO

* [und](und.md)	 - Unification Mainchain Daemon (server)

###### Auto generated by spf13/cobra on 28-Feb-2022
