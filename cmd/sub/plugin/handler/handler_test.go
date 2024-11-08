package handler

import (
	"testing"
)

func TestGetHandler(t *testing.T) {
	tests := []struct {
		name    string
		url     string
		source  string
		wantErr string
	}{
		{
			"Hangar SayanVanish (HTTPS)",
			"https://hangar.papermc.io/Syrent/SayanVanish",
			"hangar",
			"",
		},
		{
			"Hangar SayanVanish (hangar)",
			"hangar:Syrent/SayanVanish",
			"hangar",
			"",
		},

		{
			"Modrinth SayanVanish (HTTPS)",
			"https://modrinth.com/plugin/sayanvanish",
			"modrinth",
			"",
		},
		{
			"Modrinth SayanVanish (modrinth)",
			"modrinth:sayanvanish",
			"modrinth",
			"",
		},
		{
			"Source not supported (spigot)",
			"spigot:name",
			"",
			`"spigot" is not a supported source`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetHandler(tt.url)
			if err != nil {
				if tt.wantErr == "" {
					t.Errorf("did not expected error but got=%v", err)
					return
				}

				if err.Error() != tt.wantErr {
					t.Errorf("errors do not match, want=%v got%v", err, tt.wantErr)
				}

				return
			}

			if got.Name() != tt.source {
				t.Errorf("Wrong handler, want=%q got=%q", tt.source, got.Name())
			}
		})
	}
}
