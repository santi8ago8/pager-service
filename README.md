
# Pager service!

## Runing tests
Requirements:
- Golang version: >=1.14.4

Run test (root project folder):  
$ `go test ./... -cover`

Note about tests:
> The project has 100% code coverage.


## Entities
If want more details, each entity has a json schema (links below).
### Monitored service [(schema)](https://www.jsonschemavalidator.net/s/XfVpdHvc):
```
{
    "id": "test",
    "name": "Monitored service test",
    "status": "healthy"
}
```
### Escalation Policy [(schema)](https://www.jsonschemavalidator.net/s/9X6WDAvW):
```
{
  "id": "UUID_V4",
  "monitored_service_id": "UUID_V4_MSERV",
  "levels": [
    {
      "id": 0,
      "targets": [
        {
          "id": "UUID_V4_TARGET",
          "type": "email",
          "email": "example@aaa.com"
        },
        {
          "id": "UUID_V4_TARGET",
          "type": "SMS",
          "phone_number": "+5492664295169"
        }
      ]
    }
  ]
}
```
### Alert [(schema)](https://www.jsonschemavalidator.net/s/J0cyzg1F):
```
{
    "id": "UUID_V4",
    "monitored_service_id": "UUID_V4_MSERV",
    "message": "problem in production",
    "notified_level_id": 0,
    "status": "open",
}
```