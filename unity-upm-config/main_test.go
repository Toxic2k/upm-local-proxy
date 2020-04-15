package unity_upm_config

import "testing"

func TestLoad(t *testing.T) {
	txt := "[npmAuth.\"https://example.com\"]\ntoken = \"sampleToken\"\nemail = \"text@example.com\"\nalwaysAuth = true\n"
	cfg, err := loadBA([]byte(txt))
	if err != nil {
		t.Error(err)
		return
	}
	for k, v := range cfg.NpmAuth {
		t.Logf("Name: %s, Email: %s, Token: %s, AlwaysAuth: %v", k, v.Email, v.Token, v.AlwaysAuth)
	}
}

func TestLoad2(t *testing.T) {
	txt := "[npmAuth]\n[npmAuth.\"https://example.com\"]\ntoken = \"sampleToken\"\nemail = \"text@example.com\"\nalwaysAuth = true\n"
	cfg, err := loadBA([]byte(txt))
	if err != nil {
		t.Error(err)
		return
	}
	for k, v := range cfg.NpmAuth {
		t.Logf("Name: %s, Email: %s, Token: %s, AlwaysAuth: %v", k, v.Email, v.Token, v.AlwaysAuth)
	}
}

func TestSave(t *testing.T) {
	cfg := NewConfig()

	cfg.NpmAuth["https://example.com"] = ConfigElement{
		Token:      "MyToken",
		Email:      "sample@example.com",
		AlwaysAuth: false,
	}

	cfg.NpmAuth["https://zzz.com"] = ConfigElement{
		Token:      "MyTokenZZZ",
		Email:      "zzz@xxx.com",
		AlwaysAuth: true,
	}

	ba, err := cfg.saveBA()
	if err != nil {
		t.Error(err)
		return
	}

	t.Logf("Result:\n%s", string(ba))
}
