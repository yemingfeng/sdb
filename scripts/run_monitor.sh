docker ps | grep grafana | awk '{print "docker kill " $1}'
docker ps | grep prometheus | awk '{print "docker kill " $1}'

echo "input http://host.docker.internal:9090 to grafana datasource"

docker run -d -p 3000:3000 -v ~/grafana:/var/lib/grafana grafana/grafana
docker run  -d -p 9090:9090 -v $(pwd)"/scripts/prometheus.yml":/etc/prometheus/prometheus.yml prom/prometheus
