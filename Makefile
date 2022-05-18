mysql:
	sudo docker run --name my_mysql -e MYSQL_DATABASE=hobby -e MYSQL_USER=intern_1 -e MYSQL_PASSWORD=qwerty -e MYSQL_ROOT_PASSWORD=qwerty -p 3310:3306 -d mysql

build:
	sudo docker build -t country .

run:
	sudo docker run -d -p 8090:8090 country