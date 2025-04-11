#!/bin/bash

repo_addr=registry.cn-shenzhen.aliyuncs.com/zeroim/user-rpc-test

tag="latest"

container_name="zeroim-user-rpc-test"

docker stop ${container_name}
docker rm ${container_name}

docker rmi ${repo_addr}:${tag}

docker pull ${repo_addr}:${tag}
