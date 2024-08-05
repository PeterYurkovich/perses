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

import React, { ReactNode } from 'react';

import { createEnumParam, JsonParam, NumberParam, useQueryParams, withDefault } from 'use-query-params';
import { ExplorerManagerProvider } from './ExplorerManagerProvider';

const exploreQueryConfig = {
  explorer: withDefault(createEnumParam(['metrics', 'traces']), 'metrics'),
  tab: withDefault(NumberParam, 0),
  queries: withDefault(JsonParam, []),
};

interface ExplorerManagerProviderWithQueryParamsProps {
  children: ReactNode;
}

export function ExplorerManagerProviderWithQueryParams({ children }: ExplorerManagerProviderWithQueryParamsProps) {
  const [queryParams, setQueryParams] = useQueryParams(exploreQueryConfig);

  return <ExplorerManagerProvider store={[queryParams, setQueryParams]}>{children}</ExplorerManagerProvider>;
}
