# Beyond

## Create Local User

```bash
task create-local-user
```

```bash
curl -X POST -H "Authorization: Bearer $EMULATOR_ID_TOKEN" \
    -d '{"name": "user_name"}' \
    localhost:8080/user
```

```bash
curl -X GET -H "Authorization: Bearer $EMULATOR_ID_TOKEN" \
    http://localhost:8080/user/$USER_ID
```
