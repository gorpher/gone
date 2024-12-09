package conv

import (
	"testing"
)

type PageInfo struct {
	Page     int `form:"page" json:"page" binding:"required,number"`
	PageSize int `form:"page_size" json:"page_size" binding:"required,number"`
}

func TestStructToMap(t *testing.T) {
	var in = struct {
		PageInfo `json:"page_info"`
		Name     string `json:"name"`
		Age      int    `json:"age"`
		Score    int    `json:"score,omitempty"`
		Height   int    `json:"height,string"`
		Weight   int
	}{
		Name:   "test",
		Age:    18,
		Score:  98,
		Height: 185,
		Weight: 75,
	}
	m := StructToMap(&in)
	Equal(t, m["name"].(string), in.Name)
	Equal(t, m["age"].(int), in.Age)
	Equal(t, m["score"].(int), in.Score)
	Equal(t, m["height"].(int), in.Height)
	Equal(t, m["Weight"].(int), in.Weight)

}
func Equal[T comparable](t *testing.T, a, b T) {
	if a != b {
		t.Fatalf("got != want, got: %v, want: %v", a, b)
	}
}
func TestMergeMap(t *testing.T) {
	var src = map[string]any{
		"name": "test",
		"age":  18,
	}
	var dst = map[string]any{
		"age":  19,
		"name": "test2",
	}
	MergeMap(src, dst)
	Equal(t, dst["name"], "test")
}
