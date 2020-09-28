
# Pager service!

## Running tests
Requirements:
- Golang version: >=1.14.4

Run test (root project folder):  
$ `go test ./... -cover`


## Highlights: 
Note about naming:
> I try to be very descriptive with the names of the interfaces, entities, services, etc.

Note about tests:
> The project has 100% code coverage (excluding mocks folder).


## Entities
If want more details, each entity has a json schema (links below).
### Monitored service [(schema)](https://www.jsonschemavalidator.net/s/XfVpdHvc):
```
{
    "id": "UUID_V4",
    "name": "Monitored service test",
    "status": "healthy"
}
```
### Escalation Policy [(schema)](https://www.jsonschemavalidator.net/s/9X6WDAvW):
```
{
  "id": "UUID_V4",
  "monitored_service_id": "UUID_V4_M_SERVICE",
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
    "monitored_service_id": "UUID_V4_M_SERVICE",
    "message": "problem in production",
    "notified_level_id": 0,
    "status": "open"
}
```

## Scaffolding 
```
├─ README.md
├─ app
│  └─ src
│     ├─ constants
│     │  └─ constants.go
│     ├─ entities // Entities and their logic.
│     │  ├─ alert.go
│     │  ├─ escalation_policy.go
│     │  └─ monitored_service.go
│     ├─ interfaces 
│     │  └─ interfaces.go
│     ├─ mocks // repositories and services mocks
│     └─ services
│        ├─ base // base services
│        │  └─ notification_service.go
│        └─ use_cases // Most important package with use case logic
│           ├─ alert_acknowledgement_service.go
│           ├─ alert_acknowledgement_timeout_service.go
│           ├─ monitored_service_healthy_event.go
│           └─ pager_receives_alert_service.go
├─ go.mod
└─ go.sum
``` 
> files with logic has its own unit test file (*_test.go).

## Concurrency issues
To solve concurrency issues, I implement a lockService to lock entities when is necessary. Another option, if we use DynamoDB is activate the feature [Optimistic Locking](https://docs.aws.amazon.com/amazondynamodb/latest/developerguide/DynamoDBMapper.OptimisticLocking.html) with version number.


## Use Case Scenarios
(Copy from main repo)
```
- Given a Monitored Service in a Healthy State,
when the Pager receives an Alert related to this Monitored Service,
then the Monitored Service becomes Unhealthy,
the Pager notifies all targets of the first level of the escalation policy,
and sets a 15-minutes acknowledgement delay.

- Given a Monitored Service in an Unhealthy State,
the corresponding Alert is not Acknowledged
and the last level has not been notified,
when the Pager receives the Acknowledgement Timeout,
then the Pager notifies all targets of the next level of the escalation policy
and sets a 15-minutes acknowledgement delay.

- Given a Monitored Service in an Unhealthy State
when the Pager receives the Acknowledgement
and later receives the Acknowledgement Timeout,
then the Pager doesn't notify any Target
and doesn't set an acknowledgement delay.

- Given a Monitored Service in an Unhealthy State,
when the Pager receives an Alert related to this Monitored Service,
then the Pager doesn’t notify any Target
and doesn’t set an acknowledgement delay

- Given a Monitored Service in an Unhealthy State,
when the Pager receives a Healthy event related to this Monitored Service
and later receives the Acknowledgement Timeout,
then the Monitored Service becomes Healthy,
the Pager doesn’t notify any Target
and doesn’t set an acknowledgement delay
```

According to the life cycle of the pager, I have created 4 services to satisfy these use cases:
- Pager receives alert
- Alert acknowledgement timeout
- Alert acknowledgement
- Monitored Service healthy

## Observability and Monitoring
To keep code simplicity, I didn't any of this subject, but I recommend:  

Monitoring, add a APM such as Newrelic or Elastic APM.  
Metrics, add a metric service such as Datadog, Dynatrace or Prometheus.  
Add some logs in some places of the code.  

## Me
More information about me: [http://santi.fun](http://santi.fun)
