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

package schema

type InvertedIndexConfig struct {
	BM25                BM25Config
	Stopwords           StopwordConfig
	IndexTimestamps     bool
	IndexNullState      bool
	IndexPropertyLength bool
}

type BM25Config struct {
	K1 float64
	B  float64
}

type StopwordConfig struct {
	Preset    string
	Additions []string
	Removals  []string
}
