//                           _       _
// __      _____  __ ___   ___  __ _| |_ ___
// \ \ /\ / / _ \/ _` \ \ / / |/ _` | __/ _ \
//  \ V  V /  __/ (_| |\ V /| | (_| | ||  __/
//   \_/\_/ \___|\__,_| \_/ |_|\__,_|\__\___|
//
//  Copyright © 2016 - 2023 Weaviate B.V. All rights reserved.
//
//  CONTACT: hello@weaviate.io
//

package test

import (
	"fmt"
	"testing"

	"github.com/go-openapi/strfmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/weaviate/weaviate/test/helper"
	graphqlhelper "github.com/weaviate/weaviate/test/helper/graphql"
)

func getWithScrollSearch(t *testing.T) {
	t.Run("listing objects using scroll api", func(t *testing.T) {
		tests := []struct {
			name             string
			after            string
			limit            int
			filter           string
			expectedIDs      []strfmt.UUID
			expectedErrorMsg string
		}{
			{
				name:  `scroll with after: "" limit: 2`,
				after: "",
				limit: 2,
				expectedIDs: []strfmt.UUID{
					scrollClassID1,
					scrollClassID2,
					scrollClassID3,
					scrollClassID4,
					scrollClassID5,
					scrollClassID6,
					scrollClassID7,
				},
			},
			{
				name:  fmt.Sprintf("scroll with after: \"%s\" limit: 1", scrollClassID4),
				after: scrollClassID4.String(),
				limit: 1,
				expectedIDs: []strfmt.UUID{
					scrollClassID5,
					scrollClassID6,
					scrollClassID7,
				},
			},
			{
				name:             "error with offset",
				filter:           "offset: 1",
				expectedErrorMsg: "invalid 'after' filter: offset cannot be set with after and limit parameters",
			},
			{
				name:             "error with nearObject",
				filter:           fmt.Sprintf("nearObject:{id:\"%s\"}", scrollClassID1),
				expectedErrorMsg: "invalid 'after' filter: other params cannot be set with after and limit parameters",
			},
			{
				name:             "error with nearVector",
				filter:           `nearVector:{vector:[0.1, 0.2]}`,
				expectedErrorMsg: "invalid 'after' filter: other params cannot be set with after and limit parameters",
			},
			{
				name:             "error with hybrid",
				filter:           `hybrid:{query:"scroll api"}`,
				expectedErrorMsg: "invalid 'after' filter: other params cannot be set with after and limit parameters",
			},
			{
				name:             "error with bm25",
				filter:           `bm25:{query:"scroll api"}`,
				expectedErrorMsg: "invalid 'after' filter: other params cannot be set with after and limit parameters",
			},
			{
				name:             "error with sort",
				filter:           `sort:{path:"name"}`,
				expectedErrorMsg: "invalid 'after' filter: sort cannot be set with after and limit parameters",
			},
			{
				name:             "error with where",
				filter:           `where:{path:"id" operator:Like valueString:"*"}`,
				expectedErrorMsg: "invalid 'after' filter: where cannot be set with after and limit parameters",
			},
			{
				name:             "error with bm25, hybrid and offset",
				filter:           `bm25:{query:"scroll api"} hybrid:{query:"scroll api"} offset:1`,
				expectedErrorMsg: "invalid 'after' filter: other params cannot be set with after and limit parameters",
			},
		}
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				query := "{ Get { ScrollClass %s { _additional { id } } } }"
				if len(tt.expectedErrorMsg) > 0 {
					errQuery := fmt.Sprintf(query, fmt.Sprintf("(limit: 1 after: \"\" %s)", tt.filter))
					result := graphqlhelper.ErrorGraphQL(t, helper.RootAuth, errQuery)
					assert.Len(t, result, 1)

					errMsg := result[0].Message
					assert.Equal(t, tt.expectedErrorMsg, errMsg)
				} else {
					parseResults := func(t *testing.T, cities []interface{}) []strfmt.UUID {
						var ids []strfmt.UUID
						for _, city := range cities {
							id, ok := city.(map[string]interface{})["_additional"].(map[string]interface{})["id"]
							require.True(t, ok)

							idString, ok := id.(string)
							require.True(t, ok)

							ids = append(ids, strfmt.UUID(idString))
						}
						return ids
					}
					// use scroll api
					scrollSearch := func(t *testing.T, after string, limit int) []strfmt.UUID {
						scroll := fmt.Sprintf(`(limit: %v after: "%s")`, limit, after)
						result := graphqlhelper.AssertGraphQL(t, helper.RootAuth, fmt.Sprintf(query, scroll))
						cities := result.Get("Get", "ScrollClass").AsSlice()
						return parseResults(t, cities)
					}

					var scrollIDs []strfmt.UUID
					after, limit := tt.after, tt.limit
					for {
						result := scrollSearch(t, after, limit)
						scrollIDs = append(scrollIDs, result...)
						if len(result) == 0 {
							break
						}
						after = result[len(result)-1].String()
					}

					require.Equal(t, len(tt.expectedIDs), len(scrollIDs))
					for i := range tt.expectedIDs {
						assert.Equal(t, tt.expectedIDs[i], scrollIDs[i])
					}
				}
			})
		}

	})
}
