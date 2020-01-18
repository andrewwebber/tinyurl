untilsuccessful() {
  "$@"
  while [ $? -ne 0 ]
  do
    echo Retrying...
    sleep 1
    "$@"
  done
}

docker stop couchbase
docker rm couchbase
docker run -d --name couchbase --net=host couchbase/server:enterprise-6.0.3
sleep 5
untilsuccessful docker exec -it couchbase /opt/couchbase/bin/couchbase-cli cluster-init -c localhost:8091 --cluster-username=Administrator --cluster-password=password --cluster-ramsize=4000 --service="data,index,query,fts"
untilsuccessful docker exec -it couchbase /opt/couchbase/bin/couchbase-cli bucket-create -c localhost:8091 -u Administrator -p password  --bucket-type=couchbase --bucket=tinyurl -c localhost:8091 --bucket-ramsize=1000 --wait

