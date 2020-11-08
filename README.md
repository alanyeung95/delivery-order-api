# Delivery Order API

-------------------------
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

You should able to see the message before start using/testing the service
```
app_1    | Service is running on :8080
```

### Running unit tests
```
make unit-test
```

### Running integration tests

```
make integration-test
```

-------------------------
## Todo list
### Must
- [x] list order pagination 
- [x] unit test
- [x] swagger.yaml
- [x] dockerize service
- [x] integration test
- [x] readme instruction

### Improvement
- [ ] Config object for validation
