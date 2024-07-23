package models

import (
	"reflect"
	"testing"
)

func TestBasedFilter_DefaultQuery(t *testing.T) {
	tests := []struct {
		name string
		c    *BasedFilter
		want BasedFilter
	}{
		{
			name: "success",
			c:    &BasedFilter{},
			want: BasedFilter{
				Limit:  10,
				Offset: 0,
				Page:   1,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.c.DefaultQuery(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("BasedFilter.DefaultQuery() = %v, want %v", got, tt.want)
			}
		})
	}
}
