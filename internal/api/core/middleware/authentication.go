// Copyright 2024 The Perses Authors
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
	"github.com/labstack/echo/v4"
	crypto "github.com/rhobs/perses/internal/api/crypto"
	"github.com/rhobs/perses/internal/api/rbac"
	"github.com/rhobs/perses/pkg/model/api/config"
)

func GetAuthenticationMiddleware(conf config.Config, rbac rbac.RBAC, security crypto.Security) echo.MiddlewareFunc {
	if conf.Security.Authentication.Providers.KubernetesProvider.Enabled {
		return K8sMiddleware(func(_ echo.Context) bool {
			return !conf.Security.EnableAuth
		}, rbac, security)
	}
	return JwtMiddleware(security.GetJWT(), func(_ echo.Context) bool {
		return !conf.Security.EnableAuth
	})
}
