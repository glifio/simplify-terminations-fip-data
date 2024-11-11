simplify-terminations-fip-data
==============================

Tool to collect data for analysis.

* https://github.com/filecoin-project/FIPs/discussions/1036

# Example usage

First, copy mainnet.env.sample to mainnet.env and configure node address and token.

```
time go run . collect 4435141 > /tmp/4435141.csv
```

# License

MIT/Apache 2

