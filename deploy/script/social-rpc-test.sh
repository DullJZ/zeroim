#!/bin/bash

repo_addr=registry.cn-shenzhen.aliyuncs.com/zeroim/test-social-rpc

tag="latest"

container_name="zeroim-test-social-rpc"

docker stop ${container_name}
docker rm ${container_name}

docker rmi ${repo_addr}:${tag}

docker pull ${repo_addr}:${tag}

