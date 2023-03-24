package utils

import (
	"reflect"
	"testing"
)

func TestGetValueFormUnstructuredContent(t *testing.T) {
	type args struct {
		uc   map[string]interface{}
		path string
	}
	tests := []struct {
		name    string
		args    args
		want    interface{}
		wantErr bool
	}{
		{
			name: "test1",
			args: args{
				uc: map[string]interface{}{
					"spec": map[string]interface{}{
						"replicas": 1,
					},
				},
				path: "spec.replicas",
			},
			want:    1,
			wantErr: false,
		},
		{
			name: "test2",
			args: args{
				uc: map[string]interface{}{
					"spec": map[string]interface{}{
						"user": map[string]interface{}{
							"name": "yy",
						},
					},
				},
				path: "spec.user.name",
			},
			want:    "yy",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetValueFormUnstructuredContent(tt.args.uc, tt.args.path)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetValueFormUnstructuredContent() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetValueFormUnstructuredContent() got = %v, want %v", got, tt.want)
			}
		})
	}
}
