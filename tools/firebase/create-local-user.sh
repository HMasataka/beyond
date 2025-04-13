#!/bin/bash

echo "--- Sign up Firebase Auth ---"
response=$(curl -s -X POST -H "Content-Type: application/json" -d '{"returnSecureToken": true}' "http://localhost:9099/identitytoolkit.googleapis.com/v1/accounts:signUp?key=emulator-key")

id_token=$(echo $response | jq -r '.idToken')
refresh_token=$(echo $response | jq -r '.refreshToken')

export EMULATOR_ID_TOKEN=$id_token
export EMULATOR_REFRESH_TOKEN=$refresh_token

echo "--- ID Token ---"
echo "export EMULATOR_ID_TOKEN=$id_token"

echo "--- Refresh Token ---"
echo "export EMULATOR_REFRESH_TOKEN=$refresh_token"

echo "--- Create Server Account ---"
curl -X POST http://localhost:8080/account -H "Authorization: Bearer $id_token" \
