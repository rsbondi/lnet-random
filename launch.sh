echo "creating random coniguration and launching lnet"
echo "once lnet lauch complete, run start.sh"

node random.js count=$1 maxchannels=$2 out=random.dot
lnet-cli start random.dot

