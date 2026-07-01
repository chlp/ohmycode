package config

import (
	"os"
	"testing"
)

func TestApplyEnvOverrides(t *testing.T) {
	t.Setenv("OHMYCODE_API_URL", "ws://override:1234")
	t.Setenv("OHMYCODE_RUNNER_ID", "override-id")
	t.Setenv("OHMYCODE_RUNNER_TOKEN", "override-token")

	conf := RunnerConf{ApiUrl: "ws://original", RunnerId: "original-id"}
	applyEnvOverrides(&conf)

	if conf.ApiUrl != "ws://override:1234" {
		t.Errorf("ApiUrl = %q, want override", conf.ApiUrl)
	}
	if conf.RunnerId != "override-id" {
		t.Errorf("RunnerId = %q, want override", conf.RunnerId)
	}
	if conf.RunnerToken != "override-token" {
		t.Errorf("RunnerToken = %q, want override-token", conf.RunnerToken)
	}
}

func TestApplyEnvOverridesLeavesUnsetFieldsAlone(t *testing.T) {
	for _, key := range []string{"OHMYCODE_API_URL", "OHMYCODE_RUNNER_ID", "OHMYCODE_RUNNER_TOKEN"} {
		if err := os.Unsetenv(key); err != nil {
			t.Fatal(err)
		}
	}

	conf := RunnerConf{ApiUrl: "ws://original", RunnerId: "original-id", RunnerToken: "original-token"}
	applyEnvOverrides(&conf)

	if conf.ApiUrl != "ws://original" || conf.RunnerId != "original-id" || conf.RunnerToken != "original-token" {
		t.Errorf("applyEnvOverrides changed unset fields: %+v", conf)
	}
}
