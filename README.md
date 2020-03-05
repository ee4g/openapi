# openapi ![wip](https://img.shields.io/badge/-work%20in%20progress-red) ![draft](https://img.shields.io/badge/-draft-red)
Contains type safe models for the openapi specification.
It is currently intended only for serialization (one way)
generation of a subset of OpenAPI documents.

See https://swagger.io/docs/specification/about/.

Example:
```go
spec := Document{
		OpenAPI: "3.0.1",
		Info: Info{
			Title:          "Demo API",
			Description:    "Short summary of the Demo API",
			TermsOfService: mustParse("https://raw.githubusercontent.com/ee4g/openapi/master/LICENSE"),
			Contact: Contact{
				Name:  "Torben Schinke",
				Url:   mustParse("https://github.com/torbenschinke"),
				Email: "tschinke@localhost",
			},
			License: License{
				Name: "Apache 2",
				Url:  mustParse("https://raw.githubusercontent.com/ee4g/openapi/master/LICENSE"),
			},
			Version: "0.0.1",
		},
		Servers: []Server{
			{
				Url:         mustParse("localhost:{port}"),
				Description: "For your local development experience",
				Variables: map[string]ServerVariable{
					"port": {
						Enum:        []string{"8080", "8181"},
						Default:     "8080",
						Description: "The port",
					},
				},
			},
		},
		Paths: map[string]PathItem{
			"/auth/session": {
				Summary:     "Authentication",
				Description: "The Session endpoint‚",
				Get: Operation{
					Tags:        []string{"Tag A", "Tag B"},
					Summary:     "A summary for the GET session",
					Description: "A more *lengthy* text for the description of GET",
					Parameters: []Parameter{
						{
							Name:        "limit",
							In:          QueryLocation,
							Description: "Limit query parameter",
							Required:    false,
							Deprecated:  false,
							Schema: Schema{
								Type: Integer,
							},
						},
					},
					Responses: map[string]Response{
						"default": {
							Description: "the default result",
							Headers: map[string]Header{
								"SOME_HEADER": {
									Description: "whatever",
									Schema: Schema{
										Type: String,
									},
								},
							},
						},
						"200": {
							Description: "ok",
						},
					},
				},
			},
		},
	}

	b, err := json.Marshal(spec)
```