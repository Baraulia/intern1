# To run the app for the first time:
```
make build && make run
```

# To run the application for the second and subsequent times:
```
make run
```

## TESTING APPLICATION API USING CURL:

### Getting one country using curl:
```
curl http://127.0.0.1:8090/countries/AB
```
### Getting all country using curl:
```
curl http://127.0.0.1:8090/countries
```
### Getting all country using curl with pagination:
```
curl "http://127.0.0.1:8090/countries?pages=1&limit=10"
```
### Getting all country using curl with chunk:
```
curl http://127.0.0.1:8090/countries?chunk=true
```
