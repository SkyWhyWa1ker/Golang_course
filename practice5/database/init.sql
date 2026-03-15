DROP TABLE IF EXISTS user_friends;
DROP TABLE IF EXISTS users;

CREATE TABLE users (
                       id SERIAL PRIMARY KEY,
                       name VARCHAR(255) NOT NULL,
                       email VARCHAR(255) NOT NULL UNIQUE,
                       gender VARCHAR(50) NOT NULL,
                       birth_date DATE NOT NULL
);

CREATE TABLE user_friends (
                              user_id INTEGER REFERENCES users(id) ON DELETE CASCADE,
                              friend_id INTEGER REFERENCES users(id) ON DELETE CASCADE,
                              PRIMARY KEY (user_id, friend_id),
                              CHECK (user_id <> friend_id)
);

INSERT INTO users (name, email, gender, birth_date) VALUES
                                                        ('Alex', 'alex@example.com', 'male', '2000-01-10'),
                                                        ('John', 'john@example.com', 'male', '1999-02-11'),
                                                        ('Anna', 'anna@example.com', 'female', '2001-03-12'),
                                                        ('Kate', 'kate@example.com', 'female', '1998-04-13'),
                                                        ('Mike', 'mike@example.com', 'male', '2002-05-14'),
                                                        ('Sara', 'sara@example.com', 'female', '2000-06-15'),
                                                        ('David', 'david@example.com', 'male', '1997-07-16'),
                                                        ('Emma', 'emma@example.com', 'female', '2001-08-17'),
                                                        ('Chris', 'chris@example.com', 'male', '1999-09-18'),
                                                        ('Olivia', 'olivia@example.com', 'female', '2003-10-19'),
                                                        ('Daniel', 'daniel@example.com', 'male', '2000-11-20'),
                                                        ('Sophia', 'sophia@example.com', 'female', '1998-12-21'),
                                                        ('James', 'james@example.com', 'male', '2001-01-22'),
                                                        ('Mia', 'mia@example.com', 'female', '2002-02-23'),
                                                        ('Liam', 'liam@example.com', 'male', '1997-03-24'),
                                                        ('Chloe', 'chloe@example.com', 'female', '1999-04-25'),
                                                        ('Noah', 'noah@example.com', 'male', '2000-05-26'),
                                                        ('Ava', 'ava@example.com', 'female', '2001-06-27'),
                                                        ('Ethan', 'ethan@example.com', 'male', '1998-07-28'),
                                                        ('Lily', 'lily@example.com', 'female', '2002-08-29');

-- user 1 friends
INSERT INTO user_friends (user_id, friend_id) VALUES
                                                  (1, 3), (1, 4), (1, 5), (1, 6), (1, 7);

-- user 2 friends
INSERT INTO user_friends (user_id, friend_id) VALUES
                                                  (2, 3), (2, 4), (2, 5), (2, 8), (2, 9);

-- other relations
INSERT INTO user_friends (user_id, friend_id) VALUES
                                                  (3, 10), (3, 11),
                                                  (4, 12), (4, 13),
                                                  (5, 14), (5, 15),
                                                  (6, 16), (6, 17),
                                                  (7, 18), (7, 19),
                                                  (8, 20), (9, 10),
                                                  (11, 12), (13, 14),
                                                  (15, 16), (17, 18),
                                                  (19, 20);

-- optional symmetric rows
INSERT INTO user_friends (user_id, friend_id) VALUES
                                                  (3, 1), (4, 1), (5, 1), (6, 1), (7, 1),
                                                  (3, 2), (4, 2), (5, 2), (8, 2), (9, 2);