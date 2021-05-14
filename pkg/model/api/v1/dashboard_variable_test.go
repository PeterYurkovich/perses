// Copyright 2021 Amadeus s.a.s
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

package v1

import (
	"encoding/json"
	"fmt"
	"regexp"
	"testing"

	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v2"
)

func TestUnmarshallJSONVariable(t *testing.T) {
	testSuite := []struct {
		title  string
		jason  string
		result *DashboardVariable
	}{
		{
			title: "simple ConstantVariable",
			jason: `
{
  "kind": "Constant",
  "parameter": {
    "values": [
      "myVariable"
    ]
  }
}
`,
			result: &DashboardVariable{
				Kind: KindConstantVariable,
				Parameter: &ConstantVariableParameter{
					Values: []string{"myVariable"},
				},
			},
		},
		{
			title: "query variable by label_names",
			jason: `
{
  "kind": "LabelNamesQuery",
  "parameter": {
    "capturing_regexp": ".*"
  }
}
`,
			result: &DashboardVariable{
				Kind: KindLabelNamesQueryVariable,
				Parameter: &LabelNamesQueryVariableParameter{
					CapturingRegexp: (*CapturingRegexp)(regexp.MustCompile(`.*`)),
				},
			},
		},
		{
			title: "query variable by label_names with matcher",
			jason: `
{
  "kind": "LabelNamesQuery",
  "parameter": { 
    "matchers": [
      "up"
    ],
    "capturing_regexp": ".*"
  }
}
`,
			result: &DashboardVariable{
				Kind: KindLabelNamesQueryVariable,
				Parameter: &LabelNamesQueryVariableParameter{
					Matchers:        []string{"up"},
					CapturingRegexp: (*CapturingRegexp)(regexp.MustCompile(`.*`)),
				},
			},
		},
		{
			title: "query variable with label_values and matcher",
			jason: `
{
  "kind": "LabelValuesQuery",
  "parameter": {
    "label_name": "instance",
    "matchers": [
    "up"
    ],
    "capturing_regexp": ".*"
  }
}
`,
			result: &DashboardVariable{
				Kind: KindLabelValuesQueryVariable,
				Parameter: &LabelValuesQueryVariableParameter{
					LabelName:       "instance",
					Matchers:        []string{"up"},
					CapturingRegexp: (*CapturingRegexp)(regexp.MustCompile(`.*`)),
				},
			},
		},
		{
			title: "query variable with expr",
			jason: `
{
  "kind": "PromQLQuery",
  "parameter": {
    "expr": "up{instance='localhost:8080'}",
    "filter": {
      "label_name": "instance",
      "label_value": ".*"
    }
  }
}
`,
			result: &DashboardVariable{
				Kind: KindPromQLQueryVariable,
				Parameter: &PromQLQueryVariableParameter{
					Expr: "up{instance='localhost:8080'}",
					Filter: PromQLQueryFilter{
						LabelName:  "instance",
						LabelValue: (*CapturingRegexp)(regexp.MustCompile(`.*`)),
					},
				},
			},
		},
	}
	for _, test := range testSuite {
		t.Run(test.title, func(t *testing.T) {
			result := &DashboardVariable{}
			assert.NoError(t, json.Unmarshal([]byte(test.jason), result))
			assert.Equal(t, test.result, result)
		})
	}
}

func TestUnmarshallYAMLVariable(t *testing.T) {
	testSuite := []struct {
		title  string
		yamele string
		result *DashboardVariable
	}{
		{
			title: "simple ConstantVariable",
			yamele: `
kind: "Constant"
parameter:
  values:
  - "myVariable"
`,
			result: &DashboardVariable{
				Kind: KindConstantVariable,
				Parameter: &ConstantVariableParameter{
					Values: []string{"myVariable"},
				},
			},
		},
		{
			title: "query variable by label_names",
			yamele: `
kind: "LabelNamesQuery"
parameter:
  capturing_regexp: ".*"
`,
			result: &DashboardVariable{
				Kind: KindLabelNamesQueryVariable,
				Parameter: &LabelNamesQueryVariableParameter{
					CapturingRegexp: (*CapturingRegexp)(regexp.MustCompile(`.*`)),
				},
			},
		},
		{
			title: "query variable by label_names with matcher",
			yamele: `
kind: "LabelNamesQuery"
parameter:
  matchers:
  - "up"
  capturing_regexp: ".*"
`,
			result: &DashboardVariable{
				Kind: KindLabelNamesQueryVariable,
				Parameter: &LabelNamesQueryVariableParameter{
					Matchers:        []string{"up"},
					CapturingRegexp: (*CapturingRegexp)(regexp.MustCompile(`.*`)),
				},
			},
		},
		{
			title: "query variable with label_values and matcher",
			yamele: `
kind: "LabelValuesQuery"
parameter:
  label_name: "instance"
  matchers:
  - "up"
  capturing_regexp: ".*"
`,
			result: &DashboardVariable{
				Kind: KindLabelValuesQueryVariable,
				Parameter: &LabelValuesQueryVariableParameter{
					LabelName:       "instance",
					Matchers:        []string{"up"},
					CapturingRegexp: (*CapturingRegexp)(regexp.MustCompile(`.*`)),
				},
			},
		},
		{
			title: "query variable with expr",
			yamele: `
kind: "PromQLQuery"
parameter:
  expr: "up{instance='localhost:8080'}"
  filter:
    label_name: "instance"
    label_value: ".*"
`,
			result: &DashboardVariable{
				Kind: KindPromQLQueryVariable,
				Parameter: &PromQLQueryVariableParameter{
					Expr: "up{instance='localhost:8080'}",
					Filter: PromQLQueryFilter{
						LabelName:  "instance",
						LabelValue: (*CapturingRegexp)(regexp.MustCompile(`.*`)),
					},
				},
			},
		},
	}
	for _, test := range testSuite {
		t.Run(test.title, func(t *testing.T) {
			result := &DashboardVariable{}
			assert.NoError(t, yaml.Unmarshal([]byte(test.yamele), result))
			assert.Equal(t, test.result, result)
		})
	}
}

func TestUnmarshallVariableError(t *testing.T) {
	testSuite := []struct {
		title string
		jsone string
		err   error
	}{
		{
			title: "unsupported variable kind",
			jsone: `
{
  "kind": "Awkward",
  "parameter": "insane"
}
`,
			err: fmt.Errorf("unknown variable.kind 'Awkward' used"),
		},
		{
			title: "constant variable with no values",
			jsone: `
{
  "kind": "Constant",
  "parameter": {}
}
`,
			err: fmt.Errorf("parameter.values cannot be empty for a constant variable"),
		},
		{
			title: "label names query variable with no regexp",
			jsone: `
{
  "kind": "LabelNamesQuery",
  "parameter": {}
}
`,
			err: fmt.Errorf("'parameter.capturing_regexp' cannot be empty for a LabelNamesQuery"),
		},
		{
			title: "LabelValuesQuery variable with empty label_name",
			jsone: `
{
  "kind": "LabelValuesQuery",
  "parameter": {
    "capturing_regexp": ".*"
  }
}
`,
			err: fmt.Errorf("'parameter.label_name' cannot be empty for a LabelValuesQuery"),
		},
		{
			title: "LabelValuesQuery variable with empty regexp",
			jsone: `
{
  "kind": "LabelValuesQuery",
  "parameter": {
    "label_name": "test"
  }
}
`,
			err: fmt.Errorf("'parameter.capturing_regexp' cannot be empty for a LabelValuesQuery"),
		},
		{
			title: "PromQLQuery variable with empty expr",
			jsone: `
{
  "kind": "PromQLQuery",
  "parameter": {
  }
}
`,
			err: fmt.Errorf("parameter.expr cannot be empty for a PromQLQuery"),
		},
		{
			title: "PromQLQuery variable with empty label_name filter",
			jsone: `
{
  "kind": "PromQLQuery",
  "parameter": {
    "expr": "1"
  }
}
`,
			err: fmt.Errorf("parameter.filter.label_name cannot be empty for a PromQLQuery"),
		},
		{
			title: "PromQLQuery variable with empty label_value regexp",
			jsone: `
{
  "kind": "PromQLQuery",
  "parameter": {
    "expr": "1",
    "filter": {
      "label_name" :"test"
    }
  }
}
`,
			err: fmt.Errorf("parameter.filter.label_value cannot be empty for a PromQLQuery"),
		},
	}
	for _, test := range testSuite {
		t.Run(test.title, func(t *testing.T) {
			result := &DashboardVariable{}
			assert.Equal(t, test.err, json.Unmarshal([]byte(test.jsone), result))
		})
	}
}
