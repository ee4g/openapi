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

// Package v3 of openapi contains a one-way-model of the OpenAPI, formerly known as Swagger.
// It is used to create an instance of the specification programmatically. It is not expected that a YAML format
// is worth the effort, because human consumers will use the Swagger UI anyway, so only the JSON format
// is currently supported.
//
// Note that each field which is annotated with *omitempty* is optional.
package v3
