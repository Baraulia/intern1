CREATE TABLE IF NOT EXISTS countries
(
    id serial PRIMARY KEY,
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
