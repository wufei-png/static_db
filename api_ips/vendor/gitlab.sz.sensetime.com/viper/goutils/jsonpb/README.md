Faster JSONPB
==========

Faster JSONPB a fast alternative of proto/jsonpb.

And a drop-in replacement for grpc-gateway/JSONPb marshaler

## Compatibility

Passes all JSONPB unittests.

## Config:

For both marshaler & unmarshaller, set `SkipRequireCheck` to proto3 only
system.

## Performance

- Unmarshal: 2.3x faster than proto/jsonpb
- Marshal: 3x faster than proto/jsonpb

```
goos: darwin
goarch: amd64
pkg: gitlab.sz.sensetime.com/viper/goutils/jsonpb
BenchmarkUnmarshalFast-4            5000            296088 ns/op           96111 B/op       2178 allocs/op
BenchmarkUnmarshalRaw-4            50000             27708 ns/op            9947 B/op        193 allocs/op
BenchmarkUnmarshal-4                2000            660603 ns/op          150987 B/op       2182 allocs/op
BenchmarkMarshalFast-4             10000            126528 ns/op           32709 B/op        670 allocs/op
BenchmarkMarshalRaw-4              50000             30653 ns/op           22407 B/op         76 allocs/op
BenchmarkMarshal-4                  3000            362391 ns/op          124464 B/op       2340 allocs/op
PASS
ok      gitlab.sz.sensetime.com/viper/goutils/jsonpb    8.872s
```
