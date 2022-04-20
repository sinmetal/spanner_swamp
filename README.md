# spanner_swamp

## Run UnitTest

```
gcloud emulators spanner start --host-port localhost:9030
export SPANNER_EMULATOR_HOST=localhost:9030
go test ./...
```