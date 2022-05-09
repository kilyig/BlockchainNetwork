miner_count=$1
for ((i=1; i<=$miner_count; i++))
do
    go run cmd/miner/miner.go --miner-name "miner $i" &
done