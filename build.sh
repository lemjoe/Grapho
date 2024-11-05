#!/bin/bash

platforms=("windows/amd64" "linux/amd64")
if [[ -z "$1" ]]; then
  version='development'
else
    version=$1
fi

for platform in "${platforms[@]}"
do
	platform_split=(${platform//\// })
	GOOS=${platform_split[0]}
	GOARCH=${platform_split[1]}
	archive_name='Grapho-'$version'-'$GOOS'-'$GOARCH
    echo 'Building Grapho release version '$version' for '$GOOS'-'$GOARCH'...'
	if [ $GOOS = "windows" ]; then
		env GOOS=$GOOS GOARCH=$GOARCH go build -o grapho.exe -ldflags "-X 'main.Version=$version' -X 'main.BuildTime=$(date)' -H=windowsgui" ./cmd/main.go
        mkdir $archive_name
        mv grapho.exe $archive_name
        cp -r deploy.sh Dockerfile .env.default grapho-compose.yml images lang lib LICENSE README.md $archive_name
        zip -r 'builds/'$archive_name'.zip' $archive_name
        rm -fr $archive_name
        if [ $? -ne 0 ]; then
            echo 'An error has occurred! Aborting the script execution...'
            exit 1
	    fi
	else
        env GOOS=$GOOS GOARCH=$GOARCH go build -o grapho -ldflags "-X 'main.Version=$version' -X 'main.BuildTime=$(date)'" ./cmd/main.go
        mkdir $archive_name
        mv grapho $archive_name
        cp -r deploy.sh Dockerfile .env.default grapho-compose.yml images lang lib LICENSE README.md $archive_name
        tar -czvf 'builds/'$archive_name'.tar.gz' $archive_name
        rm -fr $archive_name
        if [ $? -ne 0 ]; then
            echo 'An error has occurred! Aborting the script execution...'
            exit 1
	    fi
    fi
done