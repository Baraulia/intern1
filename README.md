# To run the app for the first time:
```
docker-compose run
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
### Create new country using curl:
```
curl -X POST -H "Content-Type: application/json" 
    -d '{"name": "ТестоваяСтрана","full_name": "Республика ТестоваяСтрана","english_name": "SdDDcEGDdaFREGfsvfDSF","alpha_2": "TT", "alpha_3": "TTT","iso": 1700,"location": "Азия","location_precise": "Закавказье"}' http://127.0.0.1:8090/countries
```
### Delete country by id using curl:
```
curl -X DELETE http://127.0.0.1:8090/countries/AH
```
### Update country using curl:
```
curl -X PUT -H "Content-Type: application/json" 
    -d '{"name": "ТестоваяСтрана","full_name": "Республика ТестоваяСтрана","english_name": "SdDDcEGDdaFREGfsvfDSF","alpha_2": "TT", "alpha_3": "TTT","iso": 1700,"location": "Азия","location_precise": "Закавказье"}' http://127.0.0.1:8090/countries/AH
```
