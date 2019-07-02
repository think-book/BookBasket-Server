BookBasket-Server
===

[![CircleCI](https://circleci.com/gh/think-book/BookBasket-Server.svg?style=shield)](https://circleci.com/gh/think-book/BookBasket-Server)
[![codecov](https://codecov.io/gh/think-book/BookBasket-Server/branch/master/graph/badge.svg)](https://codecov.io/gh/think-book/BookBasket-Server)

# Overview

メモリ上にあらかじめ格納された本情報をGETRequestで取得できます。
POSTも実装しました。
フォーラム情報のGETを実装ました。
スレッドタイトルのPOSTもできるようになりました。


# Description

メモリ上に以下のデータがあるので、これをGETRequestで取得できます。

### 本のデータ
```
{"id": 1, "title": "cool book", "description": "A super hero beats monsters.", "ISBN": 100},
{"id": 2, "title": "awesome book", "description": "A text book of go langage.", "ISBN": 200}
```

### スレッドのデータ

#### メタ情報
ISBN:100の本に対するスレッドのタイトルリスト
```
{"id":1,"userID":1,"title":"I don't understand p.32 at all.","ISBN":100},
{"id":2,"userID":2,"title":"there is an awful typo on p.55","ISBN":100}
```

#### 発言情報
threadID:1のスレッドタイトル（上のメタ情報のid = 1のもの）に対するスレッドの発言リスト
```
{"id":1,"userID":11,"message":"Me neither.","threadID":1},
{"id":2,"userID":12,"message":"I think the author tries to say ...","threadID":1}
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

### 本情報
`{"title":"~","ISBN":xxx,"description":"~"}`

### スレッドタイトル
`{"userID":xxx,"title":"~","ISBN":xxx}`

で登録できます。

# Example

## GET リクエスト(本情報)
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

## GET リクエスト(フォーラム情報)
あるISBNの本のスレッドタイトルのリストを取得する場合、
`$ curl {ホストのIPアドレス}:8080/books/:ISBN/threads`
で取得できる。

あるスレッドタイトルに対する発言リストを取得する場合、
`$ curl {ホストのIPアドレス}:8080/threads/:threadID`
で取得できる。

いずれも、対応するISBNもしくはthreadIDが存在しなかった場合は、
`Not Found`
が返ります。

## POSTリクエスト（本情報）

POSTリクエストは、
`$ curl -X POST -H "Content-Type: application/json" -d '{"title":"~", ...}' {ホストのIPアドレス}:8080/books`
で行えます。

もしJSONがフォーマット通りでない場合、
`Invalid Post Format`
が返ります。

もし詳細情報がすでに存在している場合、
`Book info already exists`
が返ります。


## POSTリクエスト（スレッドタイトル）

POSTリクエストは、
`$ curl -X POST -H "Content-Type: application/json" -d '{"userID":xxx, ...}' {ホストのIPアドレス}:8080/books/:ISBN/threads`
で行えます。

もしJSONがフォーマット通りでない場合、
`Invalid Post Format`
が返ります。

もしurlとPOSTデータのISBNが一致しない場合、
`Inconsistent ISBN`
が返ります。

もしスレッドタイトルがすでに存在している場合（同じ本に同名のスレッドがある場合）、
`Thread title already exists`
が返ります。