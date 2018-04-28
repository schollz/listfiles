# listfiles

List files recursively, requires cgo.

There are three methods here, it seems one of them is much better on Windows and one is much better on Linux (no tests on OSX). Neither of the better methods if the std `filepath.Walk` (StdLib).

## Benchmarks

On Windows (8-core) with SSD and tens of thousands files:

- ListFilesUsingC: 132 seconds
- ListStdLib: 24 seconds
- **ListInParallel: 1.8 seconds**

On Linux (8-core) with SSD and a million files: 

- **ListFilesUsingC: 5.1 seconds**
- ListStdLib: 5.5 seconds
- ListInParallel: 7.2 seconds

On Linux (8-core) with 7200rpm disk and hundred thousand files: 

- **ListFilesUsingC: 1.8 seconds**
- ListStdLib: 3.0 seconds
- ListInParallel: 3.9 seconds



