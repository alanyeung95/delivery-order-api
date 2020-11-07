# delivery-order-api

## Running test cases

### Unit test
```
go test -v ./...
```

### Integration test

```
go test --tags=integrationa -v ./...
```

## Todo list
### Must
- [x] list order pagination 
- [x] unit test
- [ ] swagger.yaml
- [ ] dockerize service
- [x] integration test
- [ ] readme instruction

### Nice to have
- [ ] flexible port
- [ ] remove db name hardcoding
- [x] seperete service and repo implementation file
