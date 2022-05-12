build:
	sudo docker build -t country .

run:
	sudo docker run -d -p 8090:8090 country