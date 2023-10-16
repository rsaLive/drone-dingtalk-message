go env -w GOOS=linux
cd ./build
del drone-ding
cd ..
go build  -ldflags "-s -w" -o ./build/drone-ding .
docker rmi drone-ding
docker rmi testhub.szjixun.cn:9043/public/drone-ding
docker build . -f .\Dockerfile -t drone-ding
docker tag drone-ding testhub.szjixun.cn:9043/public/drone-ding
docker push testhub.szjixun.cn:9043/public/drone-ding
pause