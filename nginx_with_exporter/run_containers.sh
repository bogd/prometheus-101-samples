#!/bin/bash

docker run -p 8080:80 -v $(pwd)/default.conf:/etc/nginx/conf.d/default.conf -d nginx

docker run -p 9113:9113 -d nginx/nginx-prometheus-exporter -nginx.scrape-uri=http://172.17.0.2/status

curl 172.17.0.3:9113/metrics


