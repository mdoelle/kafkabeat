// Licensed to Elasticsearch B.V. under one or more contributor
// license agreements. See the NOTICE file distributed with
// this work for additional information regarding copyright
// ownership. Elasticsearch B.V. licenses this file to you under
// the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing,
// software distributed under the License is distributed on an
// "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
// KIND, either express or implied.  See the License for the
// specific language governing permissions and limitations
// under the License.

package mint

import (
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"testing/quick"
	"time"

	"github.com/urso/qcgen"
)

// SetDefaultGenerators sets default generators for use with QuickCheck.
func (t *T) SetDefaultGenerators(fns ...interface{}) {
	t.defaultGenerators = fns
}

// QuickCheck runs a quick/check test. The last function passed specifies the
// test to be run. The test function can accept any parameters, but must return
// a bool value.  The arguments are generated by random. Custom typed generator
// function of type `func(*rand.Rand) T` can be passed before the function
// under test. A generator with return value T will be used for every argument
// of type T.
// By default the random number generators seed is based on the current
// timestamp. Use the TEST_SEED environment value to configure a static seed value to be used by every test.
// The random number generator is not shared between tests.
func (t *T) QuickCheck(fns ...interface{}) {
	L := len(fns)

	check, generators := fns[L-1], fns[:L-1]
	if len(t.defaultGenerators) > 0 {
		if len(generators) > 0 {
			tmp := make([]interface{}, len(t.defaultGenerators)+len(generators))
			n := copy(tmp, t.defaultGenerators)
			copy(tmp[n:], generators)
			generators = tmp
		} else {
			generators = t.defaultGenerators
		}
	}

	seed := qcSeed()
	rng := NewRng(seed)
	t.Log("quick check rng seed: ", seed)

	t.NoError(quick.Check(check, &quick.Config{
		Rand:   rng,
		Values: qcgen.NewGenerator(check, generators...).Gen,
	}))
}

func NewRng(seed int64) *rand.Rand {
	if seed <= 0 {
		seed = qcSeed()
	}
	return rand.New(rand.NewSource(seed))
}

func RngSeed() int64 {
	return qcSeed()
}

func qcSeed() int64 {
	v := os.Getenv("TEST_SEED")
	if v == "" {
		return time.Now().UnixNano()
	}

	i, err := strconv.ParseInt(v, 10, 64)
	if err != nil {
		panic(fmt.Errorf("invalid seed '%v': %v", v, err))
	}
	return i
}