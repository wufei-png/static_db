# Changelog

## 0.4.0 20190513

- require licensesdk >= 0.4.0
- remove Run() API
- remove SetRetry API
- caller should initialize all used const in NewAuthorizorFromFile
- GetConstValueBykey -> GetConstValueByKey, GetConstValueByKey can only fetch
  specified consts in NewAuthorizorFromFile

## 0.4.1 20190610

- add NewReloadableLimiter() to auto reload rate limiter on license change

