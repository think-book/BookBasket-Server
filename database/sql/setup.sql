CREATE DATABASE bookbasket;

DROP USER 'root'@'%';
CREATE USER 'root'@'172.19.0.3' IDENTIFIED BY 'password';
GRANT SELECT ON bookbasket.* TO 'root'@'172.19.0.3' IDENTIFIED BY 'password';
GRANT INSERT ON bookbasket.* TO 'root'@'172.19.0.3' IDENTIFIED BY 'password';

CREATE TABLE bookbasket.bookInfo(
    ISBN BIGINT UNSIGNED NOT NULL PRIMARY KEY,
    title VARCHAR(50),
    description TEXT
);

INSERT INTO bookbasket.bookInfo (title, description, ISBN) VALUES(
    '機械学習入門',
    'ボルツマン機械学習から深層学習まで',
    '9784274219986'
),
(
    'ブロックチェーン 仕組みと理論',
    'サンプルで学ぶFinTechのコア技術',
    '9784865940404'
),
(
    '入門 Python 3',
    'プログラミングが初めてという人を対象に書かれた本です。',
    '9784873117386'
),
(
    'あたらしい人工知能の教科書',
    '人工知能を利用した開発に必要な基礎知識がわかる！',
    '9784798145600'
),
(
    'SQL実践入門',
    '高速でわかりやすいクエリの書き方',
    '9784774173016'
);

CREATE TABLE bookbasket.userBookRelation(
    userID INT NOT NULL,
    ISBN BIGINT UNSIGNED NOT NULL,
    PRIMARY KEY(userID, ISBN)
);

CREATE TABLE bookbasket.threadMetaInfo(
    id INT AUTO_INCREMENT NOT NULL PRIMARY KEY,
    userName VARCHAR(50),
    title VARCHAR(50),
    ISBN BIGINT UNSIGNED
);

CREATE TABLE bookbasket.threadMessage(
    id INT AUTO_INCREMENT NOT NULL PRIMARY KEY,
    userName VARCHAR(50),
    message TEXT,
    threadID INT
);

CREATE TABLE bookbasket.userInfo(
    id INT AUTO_INCREMENT PRIMARY KEY,
    userName VARCHAR(50) NOT NULL,
    password VARCHAR(60)
);
