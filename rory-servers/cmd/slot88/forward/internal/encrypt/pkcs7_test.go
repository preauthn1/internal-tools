package encrypt

import (
	"reflect"
	"testing"
)

func TestPKCS7Decode(t *testing.T) {
	type args struct {
		src []byte
	}
	tests := []struct {
		name    string
		args    args
		wantDst []byte
	}{
		{
			name: "Decode 16 size",
			args: args{
				src: []byte{123, 34, 112, 97, 114, 97, 109, 49, 34, 58, 34, 91, 48, 93, 34, 44, 34, 112, 97, 114, 97, 109, 50, 34, 58, 34, 34, 44, 34, 112, 97, 114, 97, 109, 51, 34, 58, 34, 34, 44, 34, 105, 100, 34, 58, 34, 50, 34, 44, 34, 97, 99, 116, 105, 111, 110, 73, 100, 34, 58, 34, 54, 34, 125, 16, 16, 16, 16, 16, 16, 16, 16, 16, 16, 16, 16, 16, 16, 16, 16},
			},
			wantDst: []byte{123, 34, 112, 97, 114, 97, 109, 49, 34, 58, 34, 91, 48, 93, 34, 44, 34, 112, 97, 114, 97, 109, 50, 34, 58, 34, 34, 44, 34, 112, 97, 114, 97, 109, 51, 34, 58, 34, 34, 44, 34, 105, 100, 34, 58, 34, 50, 34, 44, 34, 97, 99, 116, 105, 111, 110, 73, 100, 34, 58, 34, 54, 34, 125},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotDst := PKCS7Decode(tt.args.src); !reflect.DeepEqual(gotDst, tt.wantDst) {
				t.Errorf("PKCS7Decode() = %v, want %v", gotDst, tt.wantDst)
			}
		})
	}
}
