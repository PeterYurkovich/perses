// Copyright 2023 The Perses Authors
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package middleware

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	apiInterface "github.com/perses/perses/internal/api/interface"

	"github.com/labstack/echo/v4"
	databaseModel "github.com/perses/perses/internal/api/database/model"
	"github.com/perses/perses/internal/api/interface/v1/project"
	"github.com/perses/perses/internal/api/utils"
)

type partialMetadata struct {
	Project string `json:"project"`
}

type partialObject struct {
	Metadata partialMetadata `json:"metadata"`
}

// CheckProject is a middleware that will verify if the project used for the request exists.
func CheckProject(svc project.Service) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// We don't need to verify if a project exists in case we are in a PUT / DELETE request since if the project doesn't exist, then the dashboard won't exist either.
			// Also, we avoid an additional query to the DB like that.
			// In case the body is nil, then there is nothing to do with it as well
			method := c.Request().Method
			if (method != http.MethodPost && method != http.MethodGet) || c.Request().Body == nil {
				return next(c)
			}
			projectName := utils.GetProjectParameter(c)
			if len(projectName) == 0 && method == http.MethodPost {
				// It's possible the HTTP Path doesn't contain the project because the user is calling the root endpoint to create a new resource.
				// So we need to ensure the project name exists in the resource, which is why we will partially decode the body to get the project name.
				// And just to avoid a non-necessary deserialization, we will ensure we are managing a resource that is part of a project by checking the HTTP Path.
				for _, path := range utils.ProjectResourcePathList {
					if strings.HasPrefix(c.Path(), fmt.Sprintf("%s/%s", utils.APIV1Prefix, path)) {
						// Parsing the body in Echo middleware may cause the error code=400, message=EOF.
						//
						// Context.Bind only can be called only once in the life of the request as it read the body which can only be read once.
						// The request data reader is running out, Context.Bind() function read request body data from the socket buffer, once you took it out, it is just gone
						// That’s why it returns EOF error.
						//
						// In this middleware we need to partially decode the body to see if the project is set.
						// So we read the body, and then we re-inject it in the request.
						bodyBytes, err := io.ReadAll(c.Request().Body)
						if err != nil {
							return apiInterface.HandleBadRequestError(err.Error())
						}
						// write back to request body
						c.Request().Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
						// now we can safely partially decode the body
						o := &partialObject{}
						if unmarshalErr := json.Unmarshal(bodyBytes, o); unmarshalErr != nil {
							return apiInterface.HandleBadRequestError(unmarshalErr.Error())
						}
						if len(o.Metadata.Project) == 0 {
							return apiInterface.HandleBadRequestError("metadata.project cannot be empty")
						}
						projectName = o.Metadata.Project
						break
					}
				}
			}
			if len(projectName) > 0 {
				if _, err := svc.Get(apiInterface.Parameters{Name: projectName}); err != nil {
					if databaseModel.IsKeyNotFound(err) {
						return apiInterface.HandleBadRequestError(fmt.Sprintf("metadata.project %q doesn't exist", projectName))
					}
					return err
				}
			}
			return next(c)
		}
	}
}
