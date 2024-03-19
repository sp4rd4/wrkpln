# 📁 Workers:
## End-point: Create Worker
### Request:
```shell
curl --location 'localhost:8080/worker' \
--header 'Content-Type: application/json' \
--data '{
    "name": "John Doe"
}'
```

### Response: 201
```json
{
    "id": "8e6599ba-3c94-4e1f-9f78-c5568ef74b65",
    "name": "John Doe"
}
```
⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃

## End-point: List Workers
### Request:
```shell
curl --location 'localhost:8080/workers?name=john'
```

### Response: 200
```json
[
    {
        "id": "a291a3b1-d14e-4812-a590-79fe2c88edd1",
        "name": "John Doe"
    },
    {
        "id": "903d317f-7f11-41bc-8d34-9c4e18294e65",
        "name": "John Smith"
    }
]
```

⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃
# 📁 Shifts:
## End-point: Create Shift
### Request:
```shell
curl --location 'localhost:8080/shift' \
--header 'Content-Type: application/json' \
--data '{
    "worker_id": "a291a3b1-d14e-4812-a590-79fe2c88edd1",
    "date": "2024-03-19T23:14:10+00:00",
    "start_hour": 16,
    "end_hour": 24
}'
```
### Response: 201
```json
{
    "id": "5b44593b-6296-4f91-9931-c2afa79b5bd3",
    "worker_id": "a291a3b1-d14e-4812-a590-79fe2c88edd1",
    "date": "2024-03-19T00:00:00Z",
    "start_hour": 16,
    "end_hour": 24
}
```
⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃ ⁃

## End-point: List Shifts
### Request:
```shell
curl --location 'localhost:8080/shifts?worker_id=a291a3b1-d14e-4812-a590-79fe2c88edd1&date=2024-03-19T00%3A00%3A00Z'
```
### Response: 200
```json
[
    {
        "id": "5b44593b-6296-4f91-9931-c2afa79b5bd3",
        "worker_id": "a291a3b1-d14e-4812-a590-79fe2c88edd1",
        "date": "2024-03-19T00:00:00Z",
        "start_hour": 16,
        "end_hour": 24
    }
]
```