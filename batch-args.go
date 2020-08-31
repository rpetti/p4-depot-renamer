/**
 *
 * MIT License
 *
 * Copyright (c) 2020 Google LLC
 *
 * Permission is hereby granted, free of charge, to any person obtaining a copy of
 * this software and associated documentation files (the "Software"), to deal in
 * the Software without restriction, including without limitation the rights to
 * use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of
 * the Software, and to permit persons to whom the Software is furnished to do so,
 * subject to the following conditions:
 *
 * The above copyright notice and this permission notice shall be included in all
 * copies or substantial portions of the Software.
 *
 * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
 * IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS
 * FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR
 * COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER
 * IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN
 * CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
 */

package main

import (
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

type BatchArgument struct {
	PathFrom              string   `yaml:"path_from"`
	PathTo                string   `yaml:"path_to"`
	IncludedTransforms    []string `yaml:"included_transforms"`
	ExcludedTransforms    []string `yaml:"excluded_transforms"`
	IncludedTransformsMap map[string]bool
	ExcludedTransformsMap map[string]bool
}

func BatchArguments(yamlPath string) ([]BatchArgument, error) {
	content, err := ioutil.ReadFile(yamlPath)
	if err != nil {
		return nil, err
	}

	var batchArguments []BatchArgument
	err = yaml.Unmarshal(content, &batchArguments)
	// Convert exclusion/inclusion lists to maps for quick lookup
	for index := range batchArguments {
		batchArguments[index].IncludedTransformsMap = map[string]bool{}
		batchArguments[index].ExcludedTransformsMap = map[string]bool{}
		for _, includedTransform := range batchArguments[index].IncludedTransforms {
			batchArguments[index].IncludedTransformsMap[includedTransform] = true
		}
		for _, excludedTransform := range batchArguments[index].ExcludedTransforms {
			batchArguments[index].ExcludedTransformsMap[excludedTransform] = true
		}
	}

	return batchArguments, err
}
