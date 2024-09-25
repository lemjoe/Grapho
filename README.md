# md-blog

## Overview

**md-blog** is a simple personal blog/wiki engine with markdown support.

Here you can store, organize and collaborate on information in a way that suits you best. Create, explore and share your knowledge with ease!

- Quickly access and edit your notes using the intuitive web-based interface
- Customize your pages using Markdown markup language
- Collaborate with friends and colleagues by inviting them to view or edit specific pages
- Write your thoughts and share them with everyone

## Config file

- Create a config file `.env` with the following content:

For CloverDB:
```
DB_TYPE=cloverdb
DB_PATH=./db
```

For MongoDB:
```
DB_TYPE=mongodb
DB_HOST=localhost
DB_PORT=5432
DB_NAME=md_blog
DB_USER=user
DB_PASSWD=password
```

## Usage

Just clone this repo and run the file `/cmd/main.go`.

```
git clone https://github.com/lemjoe/md-blog.git
cd md-blog
go mod tidy
go run cmd/main.go
```

Then type `localhost:4007` in your browser's address bar.

## Roadmap to beta:

| Milestone                   |       Ready?       |
| --------------------------- | :----------------: |
| Translation en-ru           |        :x:         |
| Integration with database   | :heavy_check_mark: |
| Authentication              | :heavy_check_mark: |
| Working editor with preview | :heavy_check_mark: |
| Code refactoring            | :heavy_check_mark: |
| Docker-ready                |        :x:         |
| Dark/light theme            | :heavy_check_mark: |