/*
Copyright 2017 The Nuclio Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

// @nuclio.configure
//
// function.yaml:
//   apiVersion: "nuclio.io/v1"
//   kind: "NuclioFunction"
//   spec:
//     runtime: "golang"

package main

import (
	"math"
	"strconv"

	"github.com/nuclio/errors"
	"github.com/nuclio/nuclio-sdk-go"
)

func Handler(context *nuclio.Context, event nuclio.Event) (interface{}, error) {
	n, err := strconv.ParseUint(string(event.GetBody()), 10, 64)
	if err != nil {
		return nil, err
	}

	context.Logger.InfoWith("Calculating Fibonacci number", "n", n)

	result, err := fib(n)
	if err != nil {
		return nil, err
	}

	return nuclio.Response{
		StatusCode:  200,
		ContentType: "application/text",
		Body:        []byte(strconv.FormatUint(result, 10)),
	}, nil
}

func fib(n uint64) (uint64, error) {
	if n == 0 {
		return 0, nil
	}

	var a, b uint64 = 0, 1
	for i := uint64(0); i < n-1; i++ {
		if a < math.MaxUint64-b {
			a, b = b, a+b
		} else {
			return 0, errors.New("Overflow. The request size exceeds the maximum")
		}
	}

	return b, nil
}
