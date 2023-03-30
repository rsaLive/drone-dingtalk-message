go env -w GOOS=linux
cd ./build
del drone-ding
cd ..
go build  -ldflags "-s -w" -o ./build/drone-ding .
docker rmi drone-ding
docker rmi 172.16.100.99:9006/drone-ding
docker build . -f .\Dockerfile -t drone-ding
docker tag drone-ding 172.16.100.99:9006/drone-ding
docker push 172.16.100.99:9006/drone-ding
pause