# ddosy

## start the server
```
./ddosy
```

## API

### Schedule new load test
```
curl --request POST\
  --url 'localhost:4000/schedule'\
  --header 'Content-Type: application/json'
```


### Get status
```
curl --request GET --url 'localhost:4000/status?id=1' 
```

### Kill running load test
```
curl --request DEL --url 'localhost:4000/kill'
```
