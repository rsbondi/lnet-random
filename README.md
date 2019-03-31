### lnet-random

#### Description

Scripts for use for generating random lightning configurations with [lnet](https://github.com/cdecker/lnet)
and managing via [nodes-debug](https://github.com/rsbondi/nodes-debug)

#### Getting Started

install `bitcoind` and `clightning`

install `lnet` and `nodes-debug` per above links

```bash
git clone https://github.com/rsbondi/lnet-random.git

cd lnet-random

./launch.sh n m # where n is the number of random nodes to create and m is the maximum channels to a node, ex ./launch.sh 10 2 for 10 nodes with max 2 in

# open another terminal and run
./start.sh

# when finished, close the UI and run
lnet-cli shutdown

```

