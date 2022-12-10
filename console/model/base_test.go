package model

import (
	"reflect"
	"testing"
)

func TestBuildFuzzyQueryCondition(t *testing.T) {
	type condition map[string]interface{}
	type test struct {
		query condition
		want  []string
	}
	tests := []test{
		{
			query: condition{
				"name": "hello",
				"age":  18,
			},
			want: []string{"name like ? and age like ?", "%hello%", "%18%"},
		},
		{
			query: condition{
				"name": "hello",
				"age":  "18",
			},
			want: []string{"name like ? and age like ?", "%hello%", "%18%"},
		},
	}
	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			got := buildFuzzyQueryCondition(tt.query)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("BuildFuzzyQuery() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBuildFuzzyQueryConditionStruct(t *testing.T) {
	type A struct {
		name string
	}
	type B struct {
		age int
	}
	type test struct {
		query interface{}
		want  []string
	}
	tests := []test{
		{
			query: A{
				"hello",
			},
			want: []string{"name like ?", "%hello%"},
		},
		{
			query: B{
				18,
			},
			want: []string{"age like ?", "%18%"},
		},
	}
	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			got := buildFuzzyQueryCondition(tt.query)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("BuildFuzzyQuery() = %v, want %v", got, tt.want)
			}
		})
	}
}
