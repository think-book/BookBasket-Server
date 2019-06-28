BookBasket-Server
===

[![CircleCI](https://circleci.com/gh/think-book/BookBasket-Server.svg?style=shield)](https://circleci.com/gh/think-book/BookBasket-Server)
[![codecov](https://codecov.io/gh/think-book/BookBasket-Server/branch/master/graph/badge.svg)](https://codecov.io/gh/think-book/BookBasket-Server)

# Overview

メモリ上にあらかじめ格納された本情報をGETRequestで取得できます。
POSTも実装しました。
フォーラム情報のGETもできるようになりました。


# Description

メモリ上に以下のデータがあるので、これをGETRequestで取得できます。

### 本のデータ
```
{"id": 1, "title": "cool book", "description": "A super hero beats monsters.", "ISBN": 100},
{"id": 2, "title": "awesome book", "description": "A text book of go langage.", "ISBN": 200}
```

### フォーラムのデータ

#### メタ情報
ISBN:100の本に対するフォーラムのタイトルリスト
```
{"id":1,"user":"user_X","title":"I don't understand p.32 at all.","ISBN":100},
{"id":2,"user":"user_Y","title":"there is an awful typo on p.55","ISBN":100}
```

#### 発言情報
forumID:1のフォーラムタイトル（上のメタ情報のid = 1のもの）に対するフォーラムの発言リスト
```
{"id":1,"user":"user_A","message":"Me neither.","forumID":1},
{"id":2,"user":"user_B","message":"I think the author tries to say ...","forumID":1}
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
あるISBNの本のフォーラムタイトルのリストを取得する場合、
`$ curl {ホストのIPアドレス}:8080/books/:ISBN/forum`
で取得できる。

あるフォーラムタイトルに対する発言リストを取得する場合、
`$ curl {ホストのIPアドレス}:8080/forum/:forumID`
で取得できる。

いずれも、対応するISBNもしくはforumIDが存在しなかった場合は、
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
