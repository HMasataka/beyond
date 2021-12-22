// Package user provides primitives to interact with the openapi HTTP API.
//
// Code generated by github.com/deepmap/oapi-codegen version v1.9.0 DO NOT EDIT.
package user

// PostUserRequest defines model for PostUserRequest.
type PostUserRequest struct {
	Id *string `json:"id,omitempty"`
}

// PostUserResponse defines model for PostUserResponse.
type PostUserResponse struct {
	Id *string `json:"id,omitempty"`
}

// PostUserJSONBody defines parameters for PostUser.
type PostUserJSONBody PostUserRequest

// PostUserJSONRequestBody defines body for PostUser for application/json ContentType.
type PostUserJSONRequestBody PostUserJSONBody
