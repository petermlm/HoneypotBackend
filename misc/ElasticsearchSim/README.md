This script does the following:

 * Starts an Elasticsearch docker container
 * Creates an index for users along with a it's mappings
 * Creates random data that looks like login data
 * Uses curl to make requests to Elasticsearch
 * Stores the results of those curl requests in a file.

The output files contain the HTTP headers of the answer. They are to be used by
the honeypot to trick actors that are trying to fetch data from unsafe
Elasticsearch instances.

The logic is that one thing that such actors might try to do is list the
indices of an Elasticsearch instance, and then see get mappings and even data.
By providing fake answer as they are generated via this script, one can observe
how such an attack evolves.

Run `make` to generate the output files and `make clean` to clean all. The
docker instance is called `temp-es` and is removed by the script.
