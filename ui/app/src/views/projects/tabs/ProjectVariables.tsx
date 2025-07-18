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

import { Card } from '@mui/material';
import { ReactElement, useCallback } from 'react';
import { getResourceExtendedDisplayName, Variable, VariableResource } from '@perses-dev/core';
import { useSnackbar } from '@perses-dev/components';
import { VariableList } from '../../../components/variable/VariableList';
import {
  useCreateVariableMutation,
  useDeleteVariableMutation,
  useUpdateVariableMutation,
  useVariableList,
} from '../../../model/variable-client';

interface ProjectVariablesProps {
  projectName: string;
  id?: string;
}

export function ProjectVariables(props: ProjectVariablesProps): ReactElement {
  const { projectName, id } = props;

  const { data, isLoading } = useVariableList(projectName);

  const { successSnackbar, exceptionSnackbar } = useSnackbar();
  const createVariableMutation = useCreateVariableMutation(projectName);
  const updateVariableMutation = useUpdateVariableMutation(projectName);
  const deleteVariableMutation = useDeleteVariableMutation(projectName);

  const handleVariableCreate = useCallback(
    (variable: VariableResource): Promise<void> =>
      new Promise((resolve, reject) => {
        createVariableMutation.mutate(variable, {
          onSuccess: (createdVariable: VariableResource) => {
            successSnackbar(
              `Variable ${getResourceExtendedDisplayName(createdVariable)} has been successfully created`
            );
            resolve();
          },
          onError: (err) => {
            exceptionSnackbar(err);
            reject();
            throw err;
          },
        });
      }),
    [exceptionSnackbar, successSnackbar, createVariableMutation]
  );

  const handleVariableUpdate = useCallback(
    (variable: VariableResource): Promise<void> =>
      new Promise((resolve, reject) => {
        updateVariableMutation.mutate(variable, {
          onSuccess: (updatedVariable: VariableResource) => {
            successSnackbar(
              `Variable ${getResourceExtendedDisplayName(updatedVariable)} has been successfully updated`
            );
            resolve();
          },
          onError: (err) => {
            exceptionSnackbar(err);
            reject();
            throw err;
          },
        });
      }),
    [exceptionSnackbar, successSnackbar, updateVariableMutation]
  );

  const handleVariableDelete = useCallback(
    (variable: VariableResource): Promise<void> =>
      new Promise((resolve, reject) => {
        deleteVariableMutation.mutate(variable, {
          onSuccess: (deletedVariable: Variable) => {
            successSnackbar(
              `Variable ${getResourceExtendedDisplayName(deletedVariable)} has been successfully deleted`
            );
            resolve();
          },
          onError: (err) => {
            exceptionSnackbar(err);
            reject();
            throw err;
          },
        });
      }),
    [exceptionSnackbar, successSnackbar, deleteVariableMutation]
  );

  return (
    <Card id={id}>
      <VariableList
        data={data ?? []}
        isLoading={isLoading}
        onCreate={handleVariableCreate}
        onUpdate={handleVariableUpdate}
        onDelete={handleVariableDelete}
        initialState={{
          columns: {
            columnVisibilityModel: {
              id: false,
              project: false,
              name: false,
              version: false,
              createdAt: false,
              updatedAt: false,
            },
          },
          sorting: {
            sortModel: [{ field: 'displayName', sort: 'asc' }],
          },
        }}
      />
    </Card>
  );
}
