BookBasket-Server
===

[![CircleCI](https://circleci.com/gh/think-book/BookBasket-Server.svg?style=svg)](https://circleci.com/gh/think-book/BookBasket-Server)

# Overview

メモリ上にあらかじめ格納された本情報をGETRequestで取得できます。
POSTは実装できませんでした...
エラーハンドリングもできてません。

# Description

メモリ上に以下のデータがあるので、これをGETRequestで取得できます。
```
"{"id": 1, "title": "cool book", "memo": "foo"},
"{"id": 2, "title": "awesome book", "memo": "bar"}
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

# Example

サーバ立ち上げ後、
`$ curl {ホストのIPアドレス}:8080/api/v1/event/1`
で
`{"id": 1, "title": "cool book", "memo": "foo"}`
が取得できる。

`$ curl {ホストのIPアドレス}:8080/api/v1/event`
で
```
{"id": 1, "title": "cool book", "memo": "foo"},
{"id": 2, "title": "awesome book", "memo": "bar"}
```
が取得できる。

`$ curl {ホストのIPアドレス}:8080/api/v1/event/3`
は、インデックスが大きすぎるので、
`Not Found`
が返ります。
