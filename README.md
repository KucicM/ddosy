# ddosy

Service which can perform load tests on other services.

Overengineed wrapper around [Vegeta](https://github.com/tsenart/vegeta).

## start the server

```bash
./run.sh
```

## API

### Run new task

```bash
curl --request POST \
  --url http://localhost:4000/run \
  --header 'Content-Type: application/json' \
  --data '
{
 "endpoint": "localhost:4000/status?id=1",
 "load": [
  {
   "duration": "10s",
   "linear": {
    "startRate": 1,
    "endRate": 1
   }
  }
 ],
 "traffic": [
  {
   "weight": 1,
   "payload": ""
  }
 ]
}'
```

### Get status and results of a task

```bash
curl --request GET \
  --url 'http://localhost:4000/status?id=1' 
```

### Kill the running task

```bash
curl --request DELETE \
  --url http://localhost:4000/kill
```
