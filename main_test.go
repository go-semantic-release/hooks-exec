package main

import (
	"bytes"
	"log"
	"testing"

	"github.com/go-semantic-release/semantic-release/v2/pkg/hooks"
	"github.com/go-semantic-release/semantic-release/v2/pkg/provider"
	"github.com/go-semantic-release/semantic-release/v2/pkg/semrel"
	"github.com/stretchr/testify/require"
)

func TestInit(t *testing.T) {
	h := &Exec{}

	require.NoError(t, h.Init(map[string]string{}))
	require.NoError(t, h.Init(map[string]string{"exec_on_success": "echo hello"}))
	require.NoError(t, h.Init(map[string]string{"exec_on_no_release": "echo hello"}))
	require.ErrorContains(t, h.Init(map[string]string{"exec_on_success": "{{invalid(1)}}"}), "not defined")
}

func TestSuccess(t *testing.T) {
	ioBuf := &bytes.Buffer{}
	h := &Exec{
		log: log.New(ioBuf, "", 0),
	}
	require.NoError(t, h.Init(map[string]string{
		"exec_on_success": "echo {{.PrevRelease.Version}}-{{.NewRelease.Version}}-{{len .Commits}}-{{.Changelog}}-{{.RepoInfo.Owner}}-{{.RepoInfo.Repo}}-{{.RepoInfo.DefaultBranch}}-{{.RepoInfo.Private}}",
	}))
	require.NoError(t, h.Success(&hooks.SuccessHookConfig{
		Commits:     []*semrel.Commit{{}, {}, {}},
		PrevRelease: &semrel.Release{Version: "1.0.0"},
		NewRelease:  &semrel.Release{Version: "2.0.0"},
		Changelog:   "test:test:test",
		RepoInfo: &provider.RepositoryInfo{
			Owner:         "owner",
			Repo:          "repo",
			DefaultBranch: "main",
			Private:       false,
		},
	}))
	require.Contains(t, ioBuf.String(), "1.0.0-2.0.0-3-test:test:test-owner-repo-main-false")
}

func TestNoRelease(t *testing.T) {
	ioBuf := &bytes.Buffer{}
	h := &Exec{
		log: log.New(ioBuf, "", 0),
	}
	require.NoError(t, h.Init(map[string]string{
		"exec_on_no_release": "echo {{.Reason}}-{{.Message}}",
	}))
	require.NoError(t, h.NoRelease(&hooks.NoReleaseConfig{
		Reason:  hooks.NoReleaseReason_NO_CHANGE,
		Message: "message",
	}))
	require.Contains(t, ioBuf.String(), "NO_CHANGE-message")
}
