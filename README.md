
![Grapho.](/images/dark/logo.png)

[![Go Report Card](https://goreportcard.com/badge/github.com/lemjoe/Grapho)](https://goreportcard.com/report/github.com/lemjoe/Grapho) [![License](https://img.shields.io/badge/License-BSD%203--Clause-blue.svg)](https://github.com/lemjoe/Grapho/blob/main/LICENSE)

## Overview

**Grapho** is a simple personal blog/wiki engine with markdown support.

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
DB_NAME=grapho
DB_USER=user
DB_PASSWD=password
```

## Usage

Just clone this repo and run the file `/cmd/main.go`.

```
git clone https://github.com/lemjoe/Grapho.git
cd Grapho
go mod tidy
go run cmd/main.go
```

Then type `localhost:4007` in your browser's address bar.

## Roadmap to release:

| Milestone                   |       Ready?       |
| --------------------------- | :----------------: |
| Integration with database   | :heavy_check_mark: |
| Authentication              | :heavy_check_mark: |
| Working editor with preview | :heavy_check_mark: |
| Code refactoring            | :heavy_check_mark: |
| Dark/light theme            | :heavy_check_mark: |
| Docker-ready                | :heavy_check_mark: |
| Translation en-ru           | :heavy_check_mark: |
| Admin panel                 | :heavy_check_mark: |
