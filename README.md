Go-micro Apollo Plugin

## Example

```go
// or json encoder.
e := yaml.NewEncoder()
if err := config.Load(apollo.NewSource(source.WithEncoder(e))); err != nil {
    log.Error(err)
}
if err := config.Scan(&StatsConfig); err != nil {
    log.Error(err)
}
```

## Auto Update

```go
go func() {
    for {
        w, err := config.Watch(path...)
        if err != nil {
            log.Error(err)
        }
        // wait for next value
        v, err := w.Next()
        if err != nil {
            log.Error(err)
        }
        if err := v.Scan(&StatsConfig); err != nil {
            log.Error(err)
        }
        // TODO
        log.Info(StatsConfig.AppName)
    }
}()
```