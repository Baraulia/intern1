CREATE TABLE IF NOT EXISTS countries
(
    id integer PRIMARY KEY AUTO_INCREMENT,
    name varchar(100) NOT NULL UNIQUE,
	full_name varchar(150) NOT NULL,
	english_name varchar(150) NOT NULL,
	alpha_2 varchar(2) NOT NULL,
	alpha_3 varchar(3) NOT NULL UNIQUE,
	iso int NOT NULL UNIQUE,
	location varchar(150) NOT NULL,
	location_precise varchar(150) NOT NULL,
	url varchar(255) DEFAULT ''
);

CREATE TABLE IF NOT EXISTS users
(
    id integer PRIMARY KEY AUTO_INCREMENT,
    name varchar(100) NOT NULL,
	email varchar(150) NOT NULL UNIQUE,
	description text NOT NULL,
	country_id integer NOT NULL,
	FOREIGN KEY (country_id) REFERENCES countries(id)
);

CREATE TABLE IF NOT EXISTS hobbies
(
    id integer PRIMARY KEY AUTO_INCREMENT,
    name varchar(100) NOT NULL UNIQUE
);

CREATE TABLE IF NOT EXISTS users_hobbies
(
    id serial,
    user_id integer NOT NULL,
    hobby_id integer NOT NULL,
    PRIMARY KEY (id),
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (hobby_id) REFERENCES hobbies(id) ON DELETE CASCADE
);