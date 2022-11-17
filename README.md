# ddosy

## start the server
```
./ddosy
```

## API

### Schedule new task

```
curl --request POST\
  --url 'localhost:4000/schedule'\
  --header 'Content-Type: application/json'
```


### Get status of the task or results if the task is done

```
curl --request GET --url 'localhost:4000/status?id=1' 
```

### Kill the running task

```
curl --request DEL --url 'localhost:4000/kill'
```
