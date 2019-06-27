BookBasket-Server
===

[![CircleCI](https://circleci.com/gh/think-book/BookBasket-Server.svg?style=shield)](https://circleci.com/gh/think-book/BookBasket-Server)
[![codecov](https://codecov.io/gh/think-book/BookBasket-Server/branch/master/graph/badge.svg)](https://codecov.io/gh/think-book/BookBasket-Server)

# Overview

メモリ上にあらかじめ格納された本情報をGETRequestで取得できます。
POSTも実装しました。


# Description

メモリ上に以下のデータがあるので、これをGETRequestで取得できます。

### 本のデータ
```
{"id": 1, "title": "cool book", "description": "A super hero beats monsters.", "ISBN": 100},
{"id": 2, "title": "awesome book", "description": "A text book of go langage.", "ISBN": 200}
```

# Requirement

Docker上で動きます。

# Usage

docker-compose.ymlと同じ場所で、
```
$ docker-compose up --build
```
でdocker上にサーバが立ち上がります。

ホスト側の8080番ポートでアクセスできます。

`$ docker-compose down`
でコンテナ終了


## POSTフォーマット
`{"title":"~","ISBN":xxx,"description":"~"}`

で登録できます。

# Example

## GET リクエスト
サーバ立ち上げ後、
`$ curl {ホストのIPアドレス}:8080/books`
で
```
{"id": 1, "title": "cool book", "ISBN": 100},
{"id": 2, "title": "awesome book", "ISBN": 200}
```
が取得できる。

ISBNでの取得は、
`$ curl {ホストのIPアドレス}:8080/books/:ISBN`
の書式

例えば、
`$ curl {ホストのIPアドレス}:8080/books/100`
でISBN = 100の本の詳細、
`{"ISBN": 100, "title": "cool book", "description": "A super hero beats monsters."}`
が取得できる。

`$ curl {ホストのIPアドレス}:8080/books/300`
は、対応するISBNの本を登録していなければ、
`Not Found`
が返ります。

## POSTリクエスト

POSTリクエストは、
`$ curl -X POST -H "Content-Type: application/json" -d '{"title":"~", ...}' {ホストのIPアドレス}:8080/books`
で行えます。

もしJSONがフォーマット通りでない場合、
`Invalid Post Format`
が返ります。

もし詳細情報がすでに存在している場合、
`Book info already exists`
が返ります。
