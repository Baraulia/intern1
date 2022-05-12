






# **_TESTING APPLICATION API USING CURL:_**

## Getting one country using curl:

curl http://127.0.0.1:8090/countries/AB

## Getting all country using curl:

curl http://127.0.0.1:8090/countries

## Getting all country using curl with pagination:

curl "http://127.0.0.1:8090/countries?pages=1&limit=10"

## Getting all country using curl with chunk:

curl http://127.0.0.1:8090/countries?chunk=true

