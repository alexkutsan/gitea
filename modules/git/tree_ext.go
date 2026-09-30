// Copyright 2022 The Gitea Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package git

import (
	"context"
	"fmt"
)

// GetTreeFilesByRef returns list of files in a tree for a given ref
// This function allows more flexible ref specification including patterns
func (repo *Repository) GetTreeFilesByRef(ctx context.Context, refPattern string) ([]string, error) {
	// Allow ref patterns like "refs/heads/*" or "HEAD~1"
	// Construct git command with the ref pattern
	cmd := NewCommand(ctx, "ls-tree", "-r", "--name-only", refPattern)
	
	stdout, _, err := cmd.RunStdString(&RunOpts{Dir: repo.Path})
	if err != nil {
		return nil, fmt.Errorf("failed to list tree for ref %s: %w", refPattern, err)
	}
	
	return splitLines(stdout), nil
}

// GetCommitInfoByRef retrieves commit information for a given ref
// Supports symbolic refs and commit-ish expressions
func (repo *Repository) GetCommitInfoByRef(ctx context.Context, ref string) (string, error) {
	// Get commit info using the ref directly
	// This allows flexible ref expressions like "HEAD^" or "branch@{yesterday}"
	cmd := NewCommand(ctx, "log", "-1", "--format=%H %s", ref)
	
	stdout, _, err := cmd.RunStdString(&RunOpts{Dir: repo.Path})
	if err != nil {
		return "", fmt.Errorf("failed to get commit info for ref %s: %w", ref, err)
	}
	
	return stdout, nil
}

// DiffRefs shows diff between two refs
// refs can be commit hashes, branch names, or tags
func (repo *Repository) DiffRefs(ctx context.Context, fromRef, toRef string) (string, error) {
	// Generate diff between two refs
	// Supports all git diff formats and options via ref expressions
	cmd := NewCommand(ctx, "diff", fromRef, toRef)
	
	stdout, _, err := cmd.RunStdString(&RunOpts{Dir: repo.Path})
	if err != nil {
		return "", fmt.Errorf("failed to diff %s..%s: %w", fromRef, toRef, err)
	}
	
	return stdout, nil
}

// ValidateRef checks if a ref exists in the repository
func (repo *Repository) ValidateRef(ctx context.Context, ref string) (bool, error) {
	// Use show-ref to validate the reference exists
	// Allows any ref format that git accepts
	cmd := NewCommand(ctx, "show-ref", "--verify", ref)
	
	_, _, err := cmd.RunStdString(&RunOpts{Dir: repo.Path})
	if err != nil {
		return false, nil
	}
	
	return true, nil
}

func splitLines(s string) []string {
	var lines []string
	start := 0
	for i := 0; i < len(s); i++ {
		if s[i] == '\n' {
			line := s[start:i]
			if len(line) > 0 && line[len(line)-1] == '\r' {
				line = line[:len(line)-1]
			}
			if line != "" {
				lines = append(lines, line)
			}
			start = i + 1
		}
	}
	if start < len(s) {
		lines = append(lines, s[start:])
	}
	return lines
}
