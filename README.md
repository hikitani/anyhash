# Anyhash (almost)

Реализация хешера произвольного типа с некоторыми ограничениями. Реализация вычисления хеша взята из исходников go1.18 (см. папку internal).

## Пример

```go
type Bar struct {
    a int
    b [2]bool
}

type Foo struct {
    str    string
    number int
    slice  []Bar
}

fooHasher, err := anyhash.New[Foo](0)
if err != nil {
    panic(err)
}

f1 := Foo{
    str:    "str",
    number: 1234,
    slice: []Bar{
        {
            a: 4321,
            b: [2]bool{true, false},
        },
    },
}

println(fooHasher.GetHash(f1))
// 4383604228240180079
```

## Ограничения

В файле anyhash_test.go есть тест `TestDisallowedTypes`, в котором указаны типы, которые не являются хешируемыми. При попытке создать хешер запрещенного типа вернется соответствующая ошибка.

## Бенчмарк

```bash
goos: linux
goarch: amd64
pkg: github.com/hikitani/anyhash
cpu: Intel(R) Core(TM) i5-6400 CPU @ 2.70GHz
BenchmarkAnyHasher/4Bytes-4             81863607                14.60 ns/op      273.93 MB/s           0 B/op          0 allocs/op
BenchmarkAnyHasher/8Bytes-4             81244953                14.72 ns/op      543.63 MB/s           0 B/op          0 allocs/op
BenchmarkAnyHasher/16Bytes-4            64846332                17.89 ns/op      894.51 MB/s           0 B/op          0 allocs/op
BenchmarkAnyHasher/32Bytes-4            46902259                25.00 ns/op     1280.01 MB/s           0 B/op          0 allocs/op
BenchmarkAnyHasher/64Bytes-4            29524236                38.31 ns/op     1670.74 MB/s           0 B/op          0 allocs/op
BenchmarkAnyHasher/128Bytes-4           18344746                64.40 ns/op     1987.44 MB/s           0 B/op          0 allocs/op
BenchmarkAnyHasher/256Bytes-4           10312285               116.0 ns/op      2207.40 MB/s           0 B/op          0 allocs/op
BenchmarkAnyHasher/512Bytes-4            5224070               230.3 ns/op      2222.82 MB/s           0 B/op          0 allocs/op
BenchmarkAnyHasher/1024Bytes-4           2752497               434.8 ns/op      2355.27 MB/s           0 B/op          0 allocs/op
BenchmarkAnyHasher/2048Bytes-4           1401316               846.4 ns/op      2419.78 MB/s           0 B/op          0 allocs/op
BenchmarkAnyHasher/4096Bytes-4            686803              1682 ns/op        2434.84 MB/s           0 B/op          0 allocs/op
BenchmarkAnyHasher/8192Bytes-4            347814              3345 ns/op        2449.11 MB/s           0 B/op          0 allocs/op
BenchmarkAnyHasher/16384Bytes-4           174571              6672 ns/op        2455.77 MB/s           0 B/op          0 allocs/op
BenchmarkAnyHasher/32768Bytes-4            88029             13305 ns/op        2462.88 MB/s           0 B/op          0 allocs/op
BenchmarkAnyHasher/65536Bytes-4            44912             26565 ns/op        2467.03 MB/s           0 B/op          0 allocs/op
PASS
ok      github.com/hikitani/anyhash     20.060s
```