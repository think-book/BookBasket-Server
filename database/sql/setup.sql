CREATE DATABASE bookbasket;

CREATE TABLE bookbasket.bookInfo(
    ISBN INT NOT NULL PRIMARY KEY,
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


CREATE TABLE bookbasket.threadMetaInfo(
    id INT AUTO_INCREMENT NOT NULL PRIMARY KEY,
    userID VARCHAR(50),
    title VARCHAR(50),
    ISBN INT
);

INSERT INTO bookbasket.threadMetaInfo (userID, title, ISBN) VALUES(
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
    userID VARCHAR(50),
    message TEXT,
    threadID INT
);

INSERT INTO bookbasket.threadMessage (userID, message, threadID) VALUES(
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
    id VARCHAR(50)PRIMARY KEY,
    password VARCHAR(50)
);

INSERT INTO bookbasket.userInfo VALUES(
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
