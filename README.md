BookBasket-Server
===

[![CircleCI](https://circleci.com/gh/think-book/BookBasket-Server.svg?style=shield)](https://circleci.com/gh/think-book/BookBasket-Server)
[![codecov](https://codecov.io/gh/think-book/BookBasket-Server/branch/master/graph/badge.svg)](https://codecov.io/gh/think-book/BookBasket-Server)

# Overview

データベース上にある本情報のGET, POSTができます。
フォーラム情報のGET, POSTができます。


# Description

データベース上に以下のデータがあるので、これをGETRequestで取得できます。(ユーザ情報はまだ)


### 本のデータ
```
{"ISBN": 100, "title": "cool book", "description": "A super hero beats monsters."},
{"ISBN": 200, "title": "awesome book", "description": "A text book of go langage."}
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

### ユーザのデータ

```
{"id":1, "userName":"Alice", "password": "pass"},
{"id":2, "userName":"Bob", "password": "word"},
{"id":11, "userName":"Carol", "password": "qwer"},
{"id":12, "userName":"Charlie", "password": "tyui"}
```


# Requirement

Docker上で動きます。

# Usage

## サーバ

docker-compose.ymlと同じ場所で、
```
$ docker-compose up --build
```
でdocker上にサーバとmysqlサーバが立ち上がります。

サーバはホスト側の8080番ポートでアクセスできます。

`$ docker-compose down -v`
でデータベース初期化してコンテナ終了
(-v しないとvolumeがどんどん溜まっていく。VPSは-vいらない)

データベースを初期化しない場合は、

`$ docker-compose stop`

## テスト

mysqlのコンテナ立ち上げ（buildは初回のみ）
```
$ docker build -t (イメージのタグ名) (database/ の場所)
$ docker run --name (コンテナ名) -p 3306:3306 (イメージのタグ名)
```

mysqlサーバを立ち上げたら、ローカルマシンのserver/で、
`$ go test -v ./...`
でテスト実行

データベースを初期化してmysqlサーバを終了するには(データベース初期化しないなら、stopだけでOK)、
```
$ docker stop -v (コンテナ名)
$ docker rm (コンテナ名)
```


## POSTフォーマット

### 本情報
`{"ISBN":xxx,"title":"~","description":"~"}`

### スレッドタイトル
`{"title":"~"}`

### スレッドメッセージ
`{"message":"~"}`

### ユーザ登録
`{"userName":"~", "password":"~"}`

で登録できます。

# Example

## GET リクエスト(本情報)
サーバ立ち上げ後、
`$ curl {ホストのIPアドレス}:8080/books`
で
```
[
    {"ISBN": 100, "title": "cool book"},
    {"ISBN": 200, "title": "awesome book"}
]
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

POSTリクエスト（ユーザの本棚への登録）は、
`$ curl -X POST -H "Content-Type: application/json" -d '{"ISBN":xxx, ...}' {ホストのIPアドレス}:8080/books`
で行えます。

登録が成功した場合、
`{"ISBN":xxx,"title":"~","description":"~"}\n`
が返ります。

もしJSONがフォーマット通りでない場合、
`Invalid Post Format`
が返ります。

もしユーザがその本を既に登録している場合、
`Book has already been registerd`
が返ります。


## POSTリクエスト（スレッドタイトル）

POSTリクエストは、
`$ curl -X POST -H "Content-Type: application/json" -d '{"userID":xxx, ...}' {ホストのIPアドレス}:8080/books/:ISBN/threads`
で行えます。

登録が成功した場合、
`{"id":x, "userID":x,"title":"~","ISBN":xxx}\n`
が返ります。

もしJSONがフォーマット通りでない場合、
`Invalid Post Format`
が返ります。

もし指定したISBNの本がデータベースに存在しない場合、
`Book doesn't exist`
が返ります。

もしログインしたuserIDのユーザがデータベースに存在しない場合、
`User doesn't exist`
が返ります。


## POSTリクエスト（スレッドメッセージ）

POSTリクエストは、
`$ curl -X POST -H "Content-Type: application/json" -d '{"userID":xxx, ...}' {ホストのIPアドレス}:8080/threads/:threadID`
で行えます。

登録が成功した場合、
`{"userID":x,"message":"~","threadID":xxx}\n`
が返ります。

もしJSONがフォーマット通りでない場合、
`Invalid Post Format`
が返ります。

もし指定したthreadIDのスレッドがデータベースに存在しない場合、
`Thread doesn't exist`
が返ります。

もしログインしたuserIDのユーザがデータベースに存在しない場合、
`User doesn't exist`
が返ります。

## ユーザ登録

POSTリクエストは、
`$ curl -X POST -H "Content-Type: application/json" -d '{"userName":"~", ...}' {ホストのIPアドレス}/users/registration`
で行えます。

ユーザ名は大文字小文字の区別なく、他人と重複してはいけません。
重複すると、
`User already exists`
が返ります。

成功すると、
`{"id":x,"userName":"~"}`
が返ります。

パスワードはbcryptで暗号化されます。

## ユーザ認証

POSTリクエストは、
`$ curl -X POST -H "Content-Type: application/json" -d '{"userName":"~", ...}' {ホストのIPアドレス}/users/login`
で行えます。

失敗すると、
`Login Failed`
が返ります。

成功すると、
`{"id":x,"userName":"~"}`
が返ります。
