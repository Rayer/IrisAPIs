package IrisAPIs

import "testing"

func TestConfiguration_LoadConfiguration(t *testing.T) {
	type fields struct {
		FixioApiKey      string
		ConnectionString string
		DatabaseType     string
		LogLevel         string
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name:    "BasicTest",
			fields:  fields{},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Configuration{
				FixerIoApiKey:    tt.fields.FixioApiKey,
				ConnectionString: tt.fields.ConnectionString,
				DatabaseType:     tt.fields.DatabaseType,
				LogLevel:         tt.fields.LogLevel,
			}
			if err := c.LoadConfiguration(); (err != nil) != tt.wantErr {
				t.Errorf("LoadConfiguration() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
