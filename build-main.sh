echo "compiling go project ..."
CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .
echo "compiled"
echo "building docker container ..."
docker build -t mz47/lunchomat .
echo "built"
echo "pushing to docker hub ...."
docker push mz47/lunchomat
echo "pushed"
echo "finished!"