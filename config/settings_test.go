package config

import (
	"testing"
)

// TestIsMattermostConfigCorrect : Test si la configuration de Mattermost est correcte
func TestIsMattermostConfigCorrect(t *testing.T) {
	MattermostHookURL = ""
	MattermostHookPayload = ""
	isCorrect := IsMattermostConfigCorrect()

	if isCorrect {
		t.Errorf("IsMattermostConfigCorrect - got: %t, want: %t.", isCorrect, false)
	}

	MattermostHookURL = "aaa"
	MattermostHookPayload = ""
	isCorrect = IsMattermostConfigCorrect()

	if isCorrect {
		t.Errorf("IsMattermostConfigCorrect - got: %t, want: %t.", isCorrect, false)
	}

	MattermostHookURL = ""
	MattermostHookPayload = "aaa"
	isCorrect = IsMattermostConfigCorrect()

	if isCorrect {
		t.Errorf("IsMattermostConfigCorrect - got: %t, want: %t.", isCorrect, false)
	}

	MattermostHookURL = "aaa"
	MattermostHookPayload = "aaa"
	isCorrect = IsMattermostConfigCorrect()

	if !isCorrect {
		t.Errorf("IsMattermostConfigCorrect - got: %t, want: %t.", !isCorrect, true)
	}
}

// TestIsSlackConfigCorrect : Test si la configuration de Slack est correcte
func TestIsSlackConfigCorrect(t *testing.T) {
	SlackHookURL = ""
	SlackHookPayload = ""
	isCorrect := IsSlackConfigCorrect()

	if isCorrect {
		t.Errorf("IsSlackConfigCorrect - got: %t, want: %t.", isCorrect, false)
	}

	SlackHookURL = "aaa"
	SlackHookPayload = ""
	isCorrect = IsSlackConfigCorrect()

	if isCorrect {
		t.Errorf("IsSlackConfigCorrect - got: %t, want: %t.", isCorrect, false)
	}

	SlackHookURL = ""
	SlackHookPayload = "aaa"
	isCorrect = IsSlackConfigCorrect()

	if isCorrect {
		t.Errorf("IsSlackConfigCorrect - got: %t, want: %t.", isCorrect, false)
	}

	SlackHookURL = "aaa"
	SlackHookPayload = "aaa"
	isCorrect = IsSlackConfigCorrect()

	if !isCorrect {
		t.Errorf("IsSlackConfigCorrect - got: %t, want: %t.", !isCorrect, true)
	}
}

// TestIsDatebaseConfigCorrect : Test si la configuration de la base de données est correcte
func TestIsDatebaseConfigCorrect(t *testing.T) {
	DatabaseDriver = ""
	DatabaseName = ""
	DatabaseUser = ""
	isCorrect := IsDatabaseConfigCorrect()

	if isCorrect {
		t.Errorf("IsDatabaseConfigCorrect - got: %t, want: %t.", isCorrect, false)
	}

	DatabaseDriver = "aaa"
	DatabaseName = ""
	DatabaseUser = ""
	isCorrect = IsDatabaseConfigCorrect()

	if isCorrect {
		t.Errorf("IsDatabaseConfigCorrect - got: %t, want: %t.", isCorrect, false)
	}

	DatabaseDriver = "aaa"
	DatabaseName = "aaa"
	DatabaseUser = ""
	isCorrect = IsDatabaseConfigCorrect()

	if isCorrect {
		t.Errorf("IsDatabaseConfigCorrect - got: %t, want: %t.", isCorrect, false)
	}

	DatabaseDriver = "aaa"
	DatabaseName = ""
	DatabaseUser = "aaa"
	isCorrect = IsDatabaseConfigCorrect()

	if isCorrect {
		t.Errorf("IsDatabaseConfigCorrect - got: %t, want: %t.", isCorrect, false)
	}

	DatabaseDriver = ""
	DatabaseName = "aaa"
	DatabaseUser = ""
	isCorrect = IsDatabaseConfigCorrect()

	if isCorrect {
		t.Errorf("IsDatabaseConfigCorrect - got: %t, want: %t.", isCorrect, false)
	}

	DatabaseDriver = ""
	DatabaseName = "aaa"
	DatabaseUser = "aaa"
	isCorrect = IsDatabaseConfigCorrect()

	if isCorrect {
		t.Errorf("IsDatabaseConfigCorrect - got: %t, want: %t.", isCorrect, false)
	}

	DatabaseDriver = ""
	DatabaseName = ""
	DatabaseUser = "aaa"
	isCorrect = IsDatabaseConfigCorrect()

	if isCorrect {
		t.Errorf("IsDatabaseConfigCorrect - got: %t, want: %t.", isCorrect, false)
	}

	DatabaseDriver = "aaa"
	DatabaseName = "aaa"
	DatabaseUser = "aaa"
	isCorrect = IsDatabaseConfigCorrect()

	if !isCorrect {
		t.Errorf("IsDatabaseConfigCorrect - got: %t, want: %t.", !isCorrect, false)
	}
}

// TestIsSMTPServerConfigValid : Test si la configuration de la base de données est correcte
func TestIsSMTPServerConfigValid(t *testing.T) {
	SMTPPort = ""
	SMTPUsername = ""
	SMTPPassword = ""
	SMTPHost = ""
	isCorrect := IsSMTPServerConfigValid()

	if isCorrect {
		t.Errorf("IsSMTPServerConfigValid - got: %t, want: %t.", isCorrect, false)
	}

	SMTPHost = "localhost"
	isCorrect = IsSMTPServerConfigValid()

	if !isCorrect {
		t.Errorf("IsSMTPServerConfigValid - got: %t, want: %t.", isCorrect, true)
	}
}
