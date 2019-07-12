CREATE DATABASE bookbasket;

CREATE TABLE bookbasket.bookInfo(
    id INT AUTO_INCREMENT NOT NULL PRIMARY KEY,
    title VARCHAR(50),
    description TEXT,
    ISBN INT 
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
    userID INT,
    title VARCHAR(50),
    ISBN INT 
);

INSERT INTO bookbasket.threadMetaInfo (userID, title, ISBN) VALUES(
    '1',
    "I don't understand p.32 at all.",
    '100'
),
(
    '2',
    "there is an awful typo on p.55",
    '100'
);


CREATE TABLE bookbasket.threadMessage(
    id INT AUTO_INCREMENT NOT NULL PRIMARY KEY,
    userID INT,
    message TEXT,
    threadID INT 
);

INSERT INTO bookbasket.threadMessage (userID, message, threadID) VALUES(
    '11',
    'Me neither.',
    '1'
),
(
    '12',
    'I think the author tries to say ...',
    '1'
);


CREATE TABLE bookbasket.userInfo(
    id INT NOT NULL PRIMARY KEY,
    userName VARCHAR(50),
    password VARCHAR(50)
);

INSERT INTO bookbasket.userInfo VALUES(
    '1',
    'Alice',
    'pass'
),
(
    '2',
    'Bob',
    'word'
),
(
    '11',
    'Carol',
    'qwer'
),
(
    '12',
    'Charlie',
    'tyui'
);