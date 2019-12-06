CREATE DATABASE bookbasket;

DROP USER 'root'@'%';
DROP USER 'mysql.session'@'localhost'
DROP USER 'mysql.sys'@'localhost'
CREATE USER 'root'@'172.19.0.3' IDENTIFIED BY 'password';
GRANT SELECT ON bookbasket.* TO 'root'@'172.19.0.3' IDENTIFIED BY 'password';
GRANT INSERT ON bookbasket.* TO 'root'@'172.19.0.3' IDENTIFIED BY 'password';

CREATE TABLE bookbasket.bookInfo(
    ISBN BIGINT UNSIGNED NOT NULL PRIMARY KEY,
    title VARCHAR(50),
    description TEXT
);

INSERT INTO bookbasket.bookInfo (title, description, ISBN) VALUES(
    'cool book',
    'A super hero beats monsters.',
    '100'
),
(
    'awesome book',
    'A text book of go langage.',
    '200'
);

CREATE TABLE bookbasket.userBookRelation(
    userID INT NOT NULL,
    ISBN BIGINT UNSIGNED NOT NULL,
    PRIMARY KEY(userID, ISBN)
);

INSERT INTO bookbasket.userBookRelation (userID, ISBN) VALUES(
    '1',
    '100'
),
(
    '1',
    '200'
);

CREATE TABLE bookbasket.threadMetaInfo(
    id INT AUTO_INCREMENT NOT NULL PRIMARY KEY,
    userName VARCHAR(50),
    title VARCHAR(50),
    ISBN BIGINT UNSIGNED
);

INSERT INTO bookbasket.threadMetaInfo (userName, title, ISBN) VALUES(
    'Alice',
    "I don't understand p.32 at all.",
    '100'
),
(
    'Bob',
    "there is an awful typo on p.55",
    '100'
);


CREATE TABLE bookbasket.threadMessage(
    id INT AUTO_INCREMENT NOT NULL PRIMARY KEY,
    userName VARCHAR(50),
    message TEXT,
    threadID INT
);

INSERT INTO bookbasket.threadMessage (userName, message, threadID) VALUES(
    'Carol',
    'Me neither.',
    '1'
),
(
    'Charlie',
    'I think the author tries to say ...',
    '1'
);


CREATE TABLE bookbasket.userInfo(
    id INT AUTO_INCREMENT PRIMARY KEY,
    userName VARCHAR(50) NOT NULL,
    password VARCHAR(60)
);

INSERT INTO bookbasket.userInfo (userName, password) VALUES(
    'Alice',
    'pass'
),
(
    'Bob',
    'word'
),
(
    'Carol',
    'qwer'
),
(
    'Charlie',
    'tyui'
);
