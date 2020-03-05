package v3

import (
	"encoding/json"
	"fmt"
	"net/url"
	"testing"
)

func mustParse(s string) URL {
	r, err := url.Parse(s)
	if err != nil {
		panic(err)
	}
	return URL{r}
}

func Test_model(t *testing.T) {
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

	b, err := json.MarshalIndent(spec, " ", " ")
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(string(b))
}
