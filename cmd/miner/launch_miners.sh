for ((i=1; i <=7; i++))
do
    go run cmd/miner/miner.go --miner-name "miner $i" &
done