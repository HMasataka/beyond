// Package pet provides primitives to interact with the openapi HTTP API.
//
// Code generated by github.com/deepmap/oapi-codegen version v1.9.0 DO NOT EDIT.
package pet

import (
	externalRef0 "github.com/caravan-inc/fankey-api-generator/domain/model"
)

// ListPetResponse defines model for ListPetResponse.
type ListPetResponse struct {
	Id *string `json:"id,omitempty"`
}

// PostPetRequest defines model for PostPetRequest.
type PostPetRequest struct {
	Id *string `json:"id,omitempty"`
}

// PostPetResponse defines model for PostPetResponse.
type PostPetResponse struct {
	Id *string `json:"id,omitempty"`
}

// ListPetsParams defines parameters for ListPets.
type ListPetsParams struct {
	SortBy *externalRef0.SortBy `json:"sortBy,omitempty"`
}

// FindPetsJSONBody defines parameters for FindPets.
type FindPetsJSONBody PostPetRequest

// FindPetsJSONRequestBody defines body for FindPets for application/json ContentType.
type FindPetsJSONRequestBody FindPetsJSONBody

