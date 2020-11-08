# delivery-order-api

## Running test cases

### Unit test
```
docker-compose -f docker-compose.unit.test.yml up
```

### Integration test

```
docker-compose -f docker-compose.yml -f docker-compose.integration.test.yml up  --abort-on-container-exit
```

## Todo list
### Must
- [x] list order pagination 
- [x] unit test
- [ ] swagger.yaml
- [x] dockerize service
- [x] integration test
- [ ] readme instruction

### Nice to have
- [ ] flexible port
- [ ] remove db name hardcoding
- [x] seperete service and repo implementation file
