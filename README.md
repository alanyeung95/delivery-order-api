# delivery-order-api

## Setup / Running Instructions
### Google Map API Config
Place your own custom key in `.env` file like

```
GOOGLE_MAP_API_KEY=AIzaSyDXJyfA6jxxxxxxxxxxxxxxxxxxxx
```

### Run the service
```
make run
```

## Running test cases

### Unit test
```
make unit-test
```

### Integration test

```
make integration-test
```

## Todo list
### Must
- [x] list order pagination 
- [x] unit test
- [ ] swagger.yaml
- [x] dockerize service
- [x] integration test
- [x] readme instruction

### Improvement
- [ ] Config object for validation