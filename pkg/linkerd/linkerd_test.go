package linkerd

import (
	"testing"

	"github.com/Masterminds/semver"
)

func Test_getLinkerdSemver(t *testing.T) {
	type args struct {
		globalVersionString string
		versionString       string
	}
	tests := []struct {
		name    string
		args    args
		want    *semver.Version
		wantErr bool
	}{
		{
			name: "stable-2.10.2",
			args: args{
				globalVersionString: "",
				versionString:       "stable-2.10.2",
			},
			want:    semver.MustParse("2.10.2"),
			wantErr: false,
		},
		{
			name: "stable-2.9.5",
			args: args{
				globalVersionString: "stable-2.9.5",
				versionString:       "",
			},
			want:    semver.MustParse("2.9.5"),
			wantErr: false,
		},
		{
			name: "edge-21.6.1",
			args: args{
				globalVersionString: "",
				versionString:       "edge-21.6.1",
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := getLinkerdSemver(tt.args.globalVersionString, tt.args.versionString)
			if (err != nil) != tt.wantErr {
				t.Errorf("getLinkerdSemver() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && !got.Equal(tt.want) {
				t.Errorf("getLinkerdSemver() = %v, want %v", got, tt.want)
			}
		})
	}
}
