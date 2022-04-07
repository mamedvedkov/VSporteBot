Get envs from system
---
---

Usage
---

```go
key := "TEST_ENV"
value := "42"

os.Setenv(key, expectedValue)


getedValue := Get(key).Int(100) // getedValue == 42
getedValue = Get(key).MustInt() // getedValue == 42
getedValue = Get("unseted_key").Int(100) // getedValue == 100
getedValue = Get("unseted_key").MustInt() // panic
```

---