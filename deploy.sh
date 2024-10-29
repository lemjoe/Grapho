#!/bin/bash

IMAGE_NAME="grapho:latest"
CONTAINER_NAME="grapho_container"

CONTAINER_TOOL=""

check_container_tool() {
    if command -v podman &> /dev/null; then
        echo "podman"
    elif command -v docker &> /dev/null; then
        echo "docker"
    else
        echo "none"
    fi
}

CONTAINER_TOOL=$(check_container_tool)

if [[ "$CONTAINER_TOOL" == "none" ]]; then
    echo "Error: no podman or docker installed."
    exit 1
fi

go build -o grapho ./cmd/main.go

echo "Building and deploying container with $CONTAINER_TOOL..."

$CONTAINER_TOOL build -t $IMAGE_NAME .

if [ $? -ne 0 ]; then
    echo "Error building image."
    exit 1
fi

$CONTAINER_TOOL stop $CONTAINER_NAME 2>/dev/null || true
$CONTAINER_TOOL rm $CONTAINER_NAME 2>/dev/null || true

mkdir -p ~/Grapho/db ~/Grapho/articles
touch ~/Grapho/main.log

$CONTAINER_TOOL run -d --name $CONTAINER_NAME -p 4007:4007 \
    --restart unless-stopped \
    -e JWT_SECRET="your_jwt_secret" \
    -e ADMIN_PASSWD="Admin111" \
    -e MAIN_LOG="main.log" \
    -e DB_PATH="./db" \
    -e DB_TYPE="cloverdb" \
    -v ~/Grapho/db:/opt/Grapho/db:z \
    -v ~/Grapho/articles:/opt/Grapho/articles:z \
    -v ~/Grapho/main.log:/opt/Grapho/main.log:z \
    $IMAGE_NAME

if [ $? -eq 0 ]; then
    echo "Container was successfully built and launched."
else
    echo "Something went wrong."
    exit 1
fi
