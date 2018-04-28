# listfiles

List files recursively

On Windows (8-core) with SSD and tens of thousands files:

- ListFromFiles (using C): 132 seconds
- ListStdLib: 24 seconds
- ListInParallel: 1.8 seconds

On Linux (8-core) with a million files: 

- ListFromFiles (using C): 132 seconds
- ListStdLib: 
- ListInParallel

