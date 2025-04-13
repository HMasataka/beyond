#!/bin/bash

# curl -X POST -H "Authorization: Bearer $EMULATOR_ID_TOKEN" \
#     -d '{"name": "user_name"}' \
#     localhost:8080/user

curl -X GET -H "Authorization: Bearer $EMULATOR_ID_TOKEN" \
    http://localhost:8080/user/cvtt4pjq7buc3qu844v0
