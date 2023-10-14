echo "partition1" | kafkacat -P -p 0 -b localhost:9093 -t test
echo "partition2" | kafkacat -P -p 1 -b localhost:9093 -t test
echo "partition3" | kafkacat -P -p 2 -b localhost:9093 -t test

count=0
while [ $count -lt 100000 ]; do
  message="message$count"
  echo $message | kafkacat -P -b localhost:9093 -t test
  count=$((count+1))
  echo $message
  sleep 0.1
done