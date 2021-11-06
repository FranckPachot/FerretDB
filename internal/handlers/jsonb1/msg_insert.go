// Copyright 2021 Baltoro OÜ.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package jsonb1

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v4"

	"github.com/MangoDB-io/MangoDB/internal/bson"
	"github.com/MangoDB-io/MangoDB/internal/types"
	"github.com/MangoDB-io/MangoDB/internal/wire"
)

func (h *storage) MsgInsert(ctx context.Context, msg *wire.OpMsg) (*wire.OpMsg, error) {
	// TODO rework when sections are added

	document := msg.Documents[0]

	m := document.Map()
	collection := m[document.Command()].(string)
	db := m["$db"].(string)
	docs, _ := m["documents"].(types.Array)

	for _, d := range msg.Documents[1:] {
		docs = append(docs, d)
	}

	var inserted int32
	for _, doc := range docs {
		d := doc.(types.Document)
		sql := fmt.Sprintf("INSERT INTO %s (_jsonb) VALUES ($1)", pgx.Identifier{db, collection}.Sanitize())
		b, err := bson.MustNewDocument(d).MarshalJSON()
		if err != nil {
			return nil, err
		}

		if _, err = h.pgPool.Exec(ctx, sql, b); err != nil {
			return nil, err
		}

		inserted++
	}

	res := &wire.OpMsg{
		Documents: []types.Document{types.MustMakeDocument(
			"n", inserted,
			"ok", float64(1),
		)},
	}

	return res, nil
}