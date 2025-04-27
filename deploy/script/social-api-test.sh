#!/bin/bash

repo_addr=registry.cn-shenzhen.aliyuncs.com/zeroim/test-social-api

tag="latest"

container_name="zeroim-test-social-api"

docker stop ${container_name}
docker rm ${container_name}

docker rmi ${repo_addr}:${tag}

docker pull ${repo_addr}:${tag}

docker run -d --name ${container_name} \
  --restart=always \
  -p 8881:8881 \
  -v $(pwd)/apps/social/api/etc/dev/social.yaml:/app/etc/social.yaml \
  -e TZ=Asia/Shanghai \
  --network zeroim-network \
  ${repo_addr}:${tag}
