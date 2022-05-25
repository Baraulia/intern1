CREATE TABLE IF NOT EXISTS countries
(
    id integer PRIMARY KEY AUTO_INCREMENT,
    name varchar(100) not null UNIQUE,
	full_name varchar(150) not null,
	english_name varchar(150) not null,
	alpha_2 varchar(2) not null,
	alpha_3 varchar(3) not null UNIQUE,
	iso int not null UNIQUE,
	location varchar(150) not null,
	location_precise varchar(150) not null,
	url varchar(255) not null
);

CREATE TABLE IF NOT EXISTS users
(
    id integer PRIMARY KEY AUTO_INCREMENT,
    name varchar(100) not null,
	email varchar(150) not null UNIQUE,
	description text not null,
	country_id integer not null,
	FOREIGN KEY (country_id) REFERENCES countries(id)
);

CREATE TABLE IF NOT EXISTS hobbies
(
    id integer PRIMARY KEY AUTO_INCREMENT,
    name varchar(100) not null UNIQUE
);

CREATE TABLE IF NOT EXISTS users_hobbies
(
    id serial,
    user_id integer not null,
    hobby_id integer not null,
    PRIMARY KEY (id),
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (hobby_id) REFERENCES hobbies(id) ON DELETE CASCADE
);