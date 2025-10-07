# Build and deploy to docker hub
set -e

VERSION=1.1.2

docker build -t zebrash-api:$VERSION .
docker tag zebrash-api:$VERSION deinticketshop/zebrash-api:$VERSION
docker tag zebrash-api:$VERSION deinticketshop/zebrash-api:latest
docker push deinticketshop/zebrash-api:$VERSION
docker push deinticketshop/zebrash-api:latest