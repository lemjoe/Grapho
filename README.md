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

```
DB_TYPE=cloverdb # todo add mongodb
DB_PATH=./db # for cloverdb
DB_HOST=localhost # for mongodb or postgresql
DB_PORT=5432 # for mongodb or postgresql
DB_NAME=md_blog  # for mongodb or postgresql
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
| Authentication              |        :x:         |
| Working editor with preview | :heavy_check_mark: |
| Code refactoring            |        :x:         |
| Docker-ready                |        :x:         |

## Roadmap to refactor:

| Tasks                        |       Ready?       |
| ---------------------------- | :----------------: |
| init project structure       | :heavy_check_mark: |
| add interfaces for Tables db | :heavy_check_mark: |
| implement db methods         | :heavy_check_mark: |
| add interfaces for services  | :heavy_check_mark: |
| implement services methods   | :heavy_check_mark: |
| porting handlers             |        :x:         |
| add fs services              | :heavy_check_mark: |
