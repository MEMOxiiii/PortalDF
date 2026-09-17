package portaldf

import "testing"

func TestConfigValidate(t *testing.T) {
	tests := []struct {
		name    string
		config  Config
		wantErr bool
	}{
		{name: "raknet host:port", config: Config{ServerAddress: "127.0.0.1:19132", Transport: TransportRakNet}},
		{name: "raknet default transport", config: Config{ServerAddress: "127.0.0.1:19132"}},
		{name: "raknet empty address", config: Config{Transport: TransportRakNet}, wantErr: true},
		{
			name:   "nethernet http url",
			config: Config{ServerAddress: "http://127.0.0.1:19133", Transport: TransportNetherNet},
		},
		{
			name:    "nethernet host:port left over from raknet",
			config:  Config{ServerAddress: "127.0.0.1:19135", Transport: TransportNetherNet},
			wantErr: true,
		},
		{
			name:    "nethernet empty address",
			config:  Config{Transport: TransportNetherNet},
			wantErr: true,
		},
		{
			name:    "nethernet missing port",
			config:  Config{ServerAddress: "http://127.0.0.1", Transport: TransportNetherNet},
			wantErr: true,
		},
		{
			name:    "nethernet has a path",
			config:  Config{ServerAddress: "http://127.0.0.1:19133/nethernet", Transport: TransportNetherNet},
			wantErr: true,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := test.config.validate()
			if test.wantErr && err == nil {
				t.Fatalf("validate() error = nil, want error")
			}
			if !test.wantErr && err != nil {
				t.Fatalf("validate() error = %v, want nil", err)
			}
		})
	}
}
