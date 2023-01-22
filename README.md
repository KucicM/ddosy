# ddosy

## start the server

```bash
./run.sh
```

## API

### Schedule new task

```bash
curl --request POST\
  --url 'localhost:4000/schedule'\
  --header 'Content-Type: application/json'
```

### Get status of the task or results if the task is done

```bash
curl --request GET --url 'localhost:4000/status?id=1' 
```

### Kill the running task

```bash
curl --request DEL --url 'localhost:4000/kill'
```
