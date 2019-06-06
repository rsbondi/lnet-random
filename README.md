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

# the above will set config and launch nodesdebug GUI
# optionally you can use the CLI UI (WIP) instead of above script, see link below

# optionally, generate random activity
# this will create random invoices and pay them to/from random nodes
# it will run continuously until you manually stop
node activity.js

# when finished, close the UI and run
lnet-cli shutdown

```

[GUI video](https://youtu.be/Z6EAhRpU2Nw)

[Video for CLI UI](https://youtu.be/Hb2-DwtqYYk)

[download CLI UI linux-x64](https://moonbreeze.richardbondi.net/lnui)
