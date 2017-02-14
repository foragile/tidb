// Copyright 2017 PingCAP, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// See the License for the specific language governing permissions and
// limitations under the License.

package expression

import (
	"encoding/hex"
	"strings"

	. "github.com/pingcap/check"
	"github.com/pingcap/tidb/ast"
	"github.com/pingcap/tidb/util/hack"
	"github.com/pingcap/tidb/util/testleak"
	"github.com/pingcap/tidb/util/types"
)

func (s *testEvaluatorSuite) TestAESEncrypt(c *C) {
	defer testleak.AfterTest(c)()

	tests := []struct {
		str    string
		key    string
		expect string
	}{
		{"pingcap", "1234567890123456", "697BFE9B3F8C2F289DD82C88C7BC95C4"},
		{"pingcap123", "1234567890123456", "CEC348F4EF5F84D3AA6C4FA184C65766"},
		{"pingcap", "123456789012345678901234", "E435438AC6798B4718533096436EC342"},
	}
	for _, test := range tests {
		fc := funcs[ast.AesEncrypt]
		str := types.NewStringDatum(test.str)
		key := types.NewStringDatum(test.key)
		f, err := fc.getFunction(datumsToConstants([]types.Datum{str, key}), s.ctx)
		crypt, err := f.eval(nil)
		c.Assert(err, IsNil)
		c.Assert(toHex(crypt), Equals, test.expect)
	}
	// Test case for null input
	fc := funcs[ast.AesEncrypt]
	arg := types.NewStringDatum("str")
	var argNull types.Datum
	f, err := fc.getFunction(datumsToConstants([]types.Datum{arg, argNull}), s.ctx)
	crypt, err := f.eval(nil)
	c.Assert(err, IsNil)
	c.Assert(crypt.IsNull(), IsTrue)

	f, err = fc.getFunction(datumsToConstants([]types.Datum{argNull, arg}), s.ctx)
	crypt, err = f.eval(nil)
	c.Assert(err, IsNil)
	c.Assert(crypt.IsNull(), IsTrue)
}

func (s *testEvaluatorSuite) TestAESDecrypt(c *C) {
	defer testleak.AfterTest(c)()

	fc := funcs[ast.AesDecrypt]
	tests := []struct {
		expect      string
		key         string
		hexCryptStr string
	}{
		{"pingcap", "1234567890123456", "697BFE9B3F8C2F289DD82C88C7BC95C4"},
		{"pingcap123", "1234567890123456", "CEC348F4EF5F84D3AA6C4FA184C65766"},
		{"pingcap", "123456789012345678901234", "E435438AC6798B4718533096436EC342"},
	}
	for _, test := range tests {
		cryptStr := fromHex(test.hexCryptStr)
		key := types.NewStringDatum(test.key)
		f, err := fc.getFunction(datumsToConstants([]types.Datum{cryptStr, key}), s.ctx)
		str, err := f.eval(nil)
		c.Assert(err, IsNil)
		c.Assert(str.GetString(), Equals, test.expect)
	}
	// Test case for null input
	arg := types.NewStringDatum("str")
	var argNull types.Datum
	f, err := fc.getFunction(datumsToConstants([]types.Datum{arg, argNull}), s.ctx)
	crypt, err := f.eval(nil)
	c.Assert(err, IsNil)
	c.Assert(crypt.IsNull(), IsTrue)

	f, err = fc.getFunction(datumsToConstants([]types.Datum{argNull, arg}), s.ctx)
	crypt, err = f.eval(nil)
	c.Assert(err, IsNil)
	c.Assert(crypt.IsNull(), IsTrue)

	// For invalid key length, we also get null
	f, err = fc.getFunction(datumsToConstants([]types.Datum{arg, arg}), s.ctx)
	crypt, err = f.eval(nil)
	c.Assert(err, IsNil)
	c.Assert(crypt.IsNull(), IsTrue)
}

func toHex(d types.Datum) string {
	x, _ := d.ToString()
	return strings.ToUpper(hex.EncodeToString(hack.Slice(x)))
}

func fromHex(str string) (d types.Datum) {
	h, _ := hex.DecodeString(str)
	d.SetBytes(h)
	return d
}
