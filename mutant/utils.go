package mutant

import (
	"encoding/json"
	"math/rand"
	"time"
)

type WeightedItem struct {
	Value  string
	Weight int
}

type JSONPatchOpt struct {
	// The Operation to be applied.
	// Can be one of the following:
	// add, remove, replace, copy, move, test.
	Op string `json:"op"`
	// The JSON Pointer to the value on which to operate.
	Path  string      `json:"path"`
	Value interface{} `json:"value,omitempty"`
}

// A JSONPatch is a collection of JSONPatchOperations
type JSONPatch []JSONPatchOpt

func PrettyJson(v any) string {
	prettyJson, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return err.Error()
	}
	return string(prettyJson)
}

func WeightedRandomSelect(items []WeightedItem) WeightedItem {
	totalWeight := 0
	for _, item := range items {
		totalWeight += item.Weight
	}

	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	randomWeight := r.Intn(totalWeight)

	for _, item := range items {
		if randomWeight < item.Weight {
			return item
		}
		randomWeight -= item.Weight
	}

	return items[0]
}

// func chooseStorageClass()
