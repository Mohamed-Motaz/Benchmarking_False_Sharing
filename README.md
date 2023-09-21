# Benchmarking_False_Sharing
This is a simple script to test the effects of false sharing and cache coherency in golang.

# Steps to run:
- Make sure golang is installed.
- Run
```
go test -bench=. -count 5 -run=^#
```
