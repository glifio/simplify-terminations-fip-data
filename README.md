simplify-terminations-fip-data
==============================

Tool to collect data for analysis.

* https://github.com/filecoin-project/FIPs/discussions/1036

# Example usage

First, copy mainnet.env.sample to mainnet.env and configure node address and token.

```
time go run . collect 4435141 > /tmp/4435141.csv
```

# Penalty on Pledge usage

```
$ go run . penalty-on-pledge 4443985
32GiB, 180 days, 0.0% verified: 1.80% penalty (18.0 FIL / 999.9 FIL, 7142 sectors, pledge: 0.14 FIL/sector)
32GiB, 180 days, 100.0% verified: 1.80% penalty (18.0 FIL / 999.7 FIL, 714 sectors, pledge: 1.40 FIL/sector)
64GiB, 180 days, 0.0% verified: 1.80% penalty (18.0 FIL / 999.9 FIL, 3571 sectors, pledge: 0.28 FIL/sector)
64GiB, 180 days, 100.0% verified: 1.80% penalty (18.0 FIL / 999.7 FIL, 357 sectors, pledge: 2.80 FIL/sector)
32GiB, 360 days, 0.0% verified: 1.80% penalty (18.0 FIL / 999.9 FIL, 7142 sectors, pledge: 0.14 FIL/sector)
32GiB, 360 days, 100.0% verified: 1.80% penalty (18.0 FIL / 999.7 FIL, 714 sectors, pledge: 1.40 FIL/sector)
64GiB, 360 days, 0.0% verified: 1.80% penalty (18.0 FIL / 999.9 FIL, 3571 sectors, pledge: 0.28 FIL/sector)
64GiB, 360 days, 100.0% verified: 1.80% penalty (18.0 FIL / 999.7 FIL, 357 sectors, pledge: 2.80 FIL/sector)
32GiB, 540 days, 0.0% verified: 1.80% penalty (18.0 FIL / 999.9 FIL, 7142 sectors, pledge: 0.14 FIL/sector)
32GiB, 540 days, 100.0% verified: 1.80% penalty (18.0 FIL / 999.7 FIL, 714 sectors, pledge: 1.40 FIL/sector)
64GiB, 540 days, 0.0% verified: 1.80% penalty (18.0 FIL / 999.9 FIL, 3571 sectors, pledge: 0.28 FIL/sector)
64GiB, 540 days, 100.0% verified: 1.80% penalty (18.0 FIL / 999.7 FIL, 357 sectors, pledge: 2.80 FIL/sector)
32GiB, 720 days, 0.0% verified: 1.80% penalty (18.0 FIL / 999.9 FIL, 7142 sectors, pledge: 0.14 FIL/sector)
32GiB, 720 days, 100.0% verified: 1.80% penalty (18.0 FIL / 999.7 FIL, 714 sectors, pledge: 1.40 FIL/sector)
64GiB, 720 days, 0.0% verified: 1.80% penalty (18.0 FIL / 999.9 FIL, 3571 sectors, pledge: 0.28 FIL/sector)
64GiB, 720 days, 100.0% verified: 1.80% penalty (18.0 FIL / 999.7 FIL, 357 sectors, pledge: 2.80 FIL/sector)
32GiB, 900 days, 0.0% verified: 1.80% penalty (18.0 FIL / 999.9 FIL, 7142 sectors, pledge: 0.14 FIL/sector)
32GiB, 900 days, 100.0% verified: 1.80% penalty (18.0 FIL / 999.7 FIL, 714 sectors, pledge: 1.40 FIL/sector)
64GiB, 900 days, 0.0% verified: 1.80% penalty (18.0 FIL / 999.9 FIL, 3571 sectors, pledge: 0.28 FIL/sector)
64GiB, 900 days, 100.0% verified: 1.80% penalty (18.0 FIL / 999.7 FIL, 357 sectors, pledge: 2.80 FIL/sector)
32GiB, 1080 days, 0.0% verified: 1.80% penalty (18.0 FIL / 999.9 FIL, 7142 sectors, pledge: 0.14 FIL/sector)
32GiB, 1080 days, 100.0% verified: 1.80% penalty (18.0 FIL / 999.7 FIL, 714 sectors, pledge: 1.40 FIL/sector)
64GiB, 1080 days, 0.0% verified: 1.80% penalty (18.0 FIL / 999.9 FIL, 3571 sectors, pledge: 0.28 FIL/sector)
64GiB, 1080 days, 100.0% verified: 1.80% penalty (18.0 FIL / 999.7 FIL, 357 sectors, pledge: 2.80 FIL/sector)
32GiB, 1260 days, 0.0% verified: 1.80% penalty (18.0 FIL / 999.9 FIL, 7142 sectors, pledge: 0.14 FIL/sector)
32GiB, 1260 days, 100.0% verified: 1.80% penalty (18.0 FIL / 999.7 FIL, 714 sectors, pledge: 1.40 FIL/sector)
64GiB, 1260 days, 0.0% verified: 1.80% penalty (18.0 FIL / 999.9 FIL, 3571 sectors, pledge: 0.28 FIL/sector)
64GiB, 1260 days, 100.0% verified: 1.80% penalty (18.0 FIL / 999.7 FIL, 357 sectors, pledge: 2.80 FIL/sector)
```

# License

MIT/Apache 2

