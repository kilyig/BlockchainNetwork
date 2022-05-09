node_count=$1
for ((i=8080; i<=$miner_count + 8079; i++))
do
    go run cmd/node/node.go --node-addr "[::1]:$i" &
done