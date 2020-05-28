/*
 * Copyright 2020 Torben Schinke
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package v3

import (
	"encoding/json"
	"net/url"
	"strings"
)

type URL struct {
	*url.URL
}

func (u URL) MarshalJSON() ([]byte, error) {
	return []byte("\"" + u.URL.String() + "\""), nil
}

type Location string

const (
	QueryLocation  Location = "query"
	HeaderLocation Location = "header"
	PathLocation   Location = "path"
	CookieLocation Location = "cookie"
)

// A Document represents the root of an OpenAPI 3.x.x specification (OAS). The file name for the document should
// be openapi.json. See also https://github.com/OAI/OpenAPI-Specification/blob/master/versions/3.0.2.md#openapi-object.
type Document struct {
	OpenAPI    string              `json:"openapi"`           // OpenAPI version, e.g. 3.0.1 which is required
	Info       Info                `json:"info"`              // Info contains required metadata about the defined API
	Servers    []Server            `json:"servers,omitempty"` // Servers contains the target servers or / if empty
	Paths      map[string]PathItem `json:"paths"`             // Paths contains each endpoint specification
	Components *Components         `json:"components,omitempty"`
}

// ResolveRef tries to resolve the referenced schema.
// Currently only searching this document and definitions in components are resolvable.
func (d *Document) ResolveRef(ref string) (string, *Schema) {
	prefix := "#/components/schemas/"
	if strings.HasPrefix(ref, prefix) {
		name := ref[len(prefix):]
		if d.Components != nil && d.Components.Schemas != nil {
			schema, hasSchema := d.Components.Schemas[name]
			if hasSchema {
				return name, &schema
			}
			return "", nil
		}
	}
	return "", nil
}

// NewDocument returns a 3.0.n document
func NewDocument() *Document {
	return &Document{Paths: map[string]PathItem{}, OpenAPI: "3.0.1"}
}

func (d *Document) String() string {
	b, err := json.Marshal(d)
	if err != nil {
		panic(err)
	}
	return string(b)
}

// Info describes the API and may be required by some client. It is mainly presented for convenience.
type Info struct {
	Title          string  `json:"title"`                    // Title of the specified API and is required
	Description    string  `json:"description,omitempty"`    // Description is a short Markdown enriched text
	TermsOfService *URL    `json:"termsOfService,omitempty"` // TermsOfService is an URL or nil
	Contact        Contact `json:"license,omitempty"`        // Contact to the API maintainer
	License        License `json:"license,omitempty"`        // License information for the API
	Version        string  `json:"version"`                  // Version is for the specified API and is required
}

// Contact contains just some information about the maintainer of the API.
type Contact struct {
	Name  string `json:"name,omitempty"`  // The name of the contact person
	Url   *URL   `json:"url,omitempty"`   // URL to e.g. a vcard or a persons page or profile
	Email string `json:"email,omitempty"` // Email is the actual mail address

}

// License describes the license for the described API
type License struct {
	Name string `json:"name"`          // Name is the required identifier for the license
	Url  URL    `json:"url,omitempty"` // Url is an optional url to the license text
}

// Server represents a service endpoint behind a specific URL
type Server struct {
	Url         string                    `json:"url"`                   // Url is the required target host
	Description string                    `json:"description,omitempty"` // Description is the optional Markdown text
	Variables   map[string]ServerVariable `json:"variables,omitempty"`   // Variables define substitutions for url
}

// ServerVariable represents a Server url substitution rule. Variables in an URL must be declared in curly braces.
type ServerVariable struct {
	Enum        []string `json:"enum,omitempty"` // Enum of values for substitution
	Default     string   `json:"default"`        // Default is the required default value
	Description string   `json:"description"`    // Description is the optional Markdown text for this variable‚
}

// A PathItem describes the available operations for a specific path.
type PathItem struct {
	Get    *Operation `json:"get,omitempty"`    // Get defines‚ the get Verb
	Post   *Operation `json:"post,omitempty"`   // Get defines‚ the get Verb
	Delete *Operation `json:"delete,omitempty"` // Get defines‚ the get Verb
	Put    *Operation `json:"put,omitempty"`    // Get defines‚ the get Verb
	Patch  *Operation `json:"patch,omitempty"`  // Get defines‚ the get Verb
}

func (p *PathItem) Map() map[string]*Operation {
	r := map[string]*Operation{}
	if p.Get != nil {
		r["GET"] = p.Get
	}

	if p.Delete != nil {
		r["DELETE"] = p.Delete
	}

	if p.Patch != nil {
		r["PATCH"] = p.Patch
	}

	if p.Post != nil {
		r["POST"] = p.Post
	}

	if p.Put != nil {
		r["PUT"] = p.Put
	}

	return r
}

// An Operation is the http Verb specifier
type Operation struct {
	Tags        []string            `json:"tags,omitempty"`        // Tags are used for logical grouping
	Summary     string              `json:"summary,omitempty"`     // Summary is a short text for what this is
	Description string              `json:"description,omitempty"` // Description is like summary but Markdown and longer
	Parameters  []Parameter         `json:"parameters,omitempty"`  // Parameters for different locations
	Responses   map[string]Response `json:"responses"`             // Responses is required and defines the results
}

// Parameter is used for path, query, header and cookie parameters. It is only unique per name and location.
type Parameter struct {
	Name        string               `json:"name"`                 // Name is the required parameter identifier
	In          Location             `json:"in"`                   // In is the required location specifier
	Description string               `json:"description"`          // Description is the optional markdown text
	Required    bool                 `json:"required,omitempty"`   // Required is obligatory for *path* and must be true
	Deprecated  bool                 `json:"deprecated,omitempty"` // Deprecated declares that it should not be used
	Schema      Schema               `json:"schema,omitempty"`     // Schema should be used to describe the data type
	Content     map[string]MediaType `json:"content,omitempty"`    // Content should be used to describe the data type‚
	// allowEmptyValue is deprecated and should not be used

}

// Response specifies a single response from an API endpoint
type Response struct {
	Description string               `json:"description"`       // Description is required, for a change
	Headers     map[string]Header    `json:"headers,omitempty"` // Headers may contain additional information
	Content     map[string]MediaType `json:"content,omitempty"` // Content describes potential response types

}

// A Reference is a string referring to a component within this document (prefixed with #) or
// to an external file (e.g. MySchema.json). The key must be always $ref
type Reference string

// A Header is like a Parameter but without Name and In fields‚
type Header struct {
	Description string `json:"description"`          // Description is the optional markdown text
	Required    bool   `json:"required,omitempty"`   // Required is obligatory for *path* and must be true
	Deprecated  bool   `json:"deprecated,omitempty"` // Deprecated declares that a parameter should not be used
	Schema      Schema `json:"schema,omitempty"`     // Schema used to describe the content
}

// MediaType provides a schema and an example for it.
type MediaType struct {
	Schema Schema `json:"schema"` // Schema is required
	//	Encoding map[string]Encoding `json:"encoding,omitempty"` // Encoding maps between a property and its encoding.
}

// An Encoding is applied to a specific schema property.
type Encoding struct {
	ContentType   string            `json:"contentType,omitempty"`   // ContentType like application/json etc
	Headers       map[string]Header `json:"headers,omitempty"`       // Headers may contain additional information
	Style         string            `json:"style,omitempty"`         // Style is a hint with various typed meanings
	Explode       bool              `json:"explode,omitempty"`       // Explode generates parameters for each value
	AllowReserved bool              `json:"allowReserved,omitempty"` // AllowReserved allows RFC3986 characters
}

// Schema defines a data type or a union of data types.
type Schema struct {
	Type          Type              `json:"type,omitempty"`
	Format        string            `json:"format,omitempty"`        // Format may contain an arbitrary hint for the format
	Minimum       int64             `json:"minimum,omitempty"`       // Minimum is inclusive
	Maximum       int64             `json:"maximum,omitempty"`       // Maximum is inclusive
	MaxLength     int               `json:"maxLength,omitempty"`     // MaxLength in bytes
	MinLength     int               `json:"minLength,omitempty"`     // MinLength in bytes
	MaxItems      int               `json:"maxItems,omitempty"`      // MaxItems of an array
	MinItems      int               `json:"minItems,omitempty"`      // MinItems for an array
	Nullable      bool              `json:"nullable,omitempty"`      // Nullable allows a null value
	Pattern       string            `json:"pattern,omitempty"`       // Pattern should be a valid regex
	Discriminator *Discriminator    `json:"discriminator,omitempty"` // Discriminator allows union types
	ReadOnly      bool              `json:"readOnly,omitempty"`      // ReadOnly declares a read only property
	WriteOnly     bool              `json:"writeOnly,omitempty"`     // WriteOnly declares a write only property
	Deprecated    bool              `json:"deprecated,omitempty"`    // Deprecated, if true should not be used
	Properties    map[string]Schema `json:"properties,omitempty"`    // Properties is only valid for type Object
	Ref           *string           `json:"$ref,omitempty"`          // Ref is a reference to another schema, e.g. #/components/schemas/MySchema
	Items         *Items            `json:"items,omitempty"`
	Description   string            `json:"description,omitempty"`
	XType         *string           `json:"x-ee.type,omitempty"`
}

type Items struct {
	*Schema
}

// Components defines various central specifications
type Components struct {
	Schemas map[string]Schema `json:"schemas,omitempty"`
}

// Type of a schema, see https://swagger.io/docs/specification/data-models/data-types/
type Type string

const (
	String  Type = "string"
	Number  Type = "number"
	Integer Type = "integer"
	Boolean Type = "bool"
	Array   Type = "array"
	Object  Type = "object"
)

// Format hint, may be anything, e.g. Regex
type Format string

const (
	Int32    Format = "int32"
	Int64    Format = "int64"
	Float    Format = "float"
	Double   Format = "double"
	Binary   Format = "binary"    //sequence of octets
	Byte     Format = "byte"      // base64
	Date     Format = "date"      // full-date RFC3339
	DateTime Format = "date-time" // date-time RFC3339
	Password Format = "password"
)

// A Discriminator specifies a field which maps between values and (polymorphic) types.
type Discriminator struct {
	PropertyName string            `json:"propertyName"`      // The required field name
	Mapping      map[string]string `json:"mapping,omitempty"` // Mapping holds property values and schema or references
}

// FromJson tries to parse the document
func FromJson(str []byte) (*Document, error) {
	doc := &Document{}
	err := json.Unmarshal(str, doc)
	if err != nil {
		return doc, err
	}
	return doc, nil
}
