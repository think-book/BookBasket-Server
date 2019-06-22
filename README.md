BookBasket-Server
===

[![CircleCI](https://circleci.com/gh/think-book/BookBasket-Server.svg?style=svg)](https://circleci.com/gh/think-book/BookBasket-Server)

# Overview

メモリ上にあらかじめ格納された本情報をGETRequestで取得できます。
POSTも実装しました。


# Description

メモリ上に以下のデータがあるので、これをGETRequestで取得できます。

### メタデータ
```
{"id": 1, "title": "cool book", "ISBN": 100},
{"id": 2, "title": "awesome book", "ISBN": 200}
```

### 本の詳細データ
```
{"ISBN": 100, "title": "cool book", "story": "A super hero beats monsters."},
{"ISBN": 200, "title": "awesome book", "story": "A text book of go langage."}
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
メタ情報が、   
`{"title": "~", "ISBN": xxx}`  

詳細情報が、  
`{"ISBN": xxx, "title": "~", "story": "~"}`  

で登録できます。  

先に対応するISBNを持つメタ情報が登録されていないと、詳細情報は登録できません。  

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
`{"ISBN": 100, "title": "cool book", "story": "A super hero beats monsters."}`
が取得できる。

`$ curl {ホストのIPアドレス}:8080/books/300`
は、対応するISBNの本を登録していなければ、    
`Not Found`  
が返ります。  

## POSTリクエスト  

POSTリクエストは、メタ情報が、  
`$ curl -X POST -H "Content-Type: application/json" -d '{"ISBN":, ...}' {ホストのIPアドレス}:8080/books`  
で行えます。  

詳細情報が、  
`$ curl -X POST -H "Content-Type: application/json" -d '{"ISBN":, ...}' {ホストのIPアドレス}:8080/books/:ISBN`  
で行えます。  

対応するISBNのメタ情報が未登録の場合、  
`Book Meta Data Not Found`  
が返ります。 

urlとデータのISBNが不一致の場合、  
`ISBN is inconsistent`  
が返ります。  

もしJSONがフォーマット通りでない場合、
`Invalid Post Format`  
が返ります。
