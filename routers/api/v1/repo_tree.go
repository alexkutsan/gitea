// Copyright 2022 The Gitea Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package v1

import (
	"net/http"
	
	"code.gitea.io/gitea/modules/context"
	"code.gitea.io/gitea/modules/git"
)

// GetTreeByRef handles GET /repos/{owner}/{repo}/git/trees/{ref}
// Returns tree contents for a given reference
func GetTreeByRef(ctx *context.APIContext) {
	// swagger:operation GET /repos/{owner}/{repo}/git/trees/{ref} repository getTreeByRef
	// ---
	// summary: Get tree by reference
	// description: Returns the tree contents for a git reference (branch, tag, or commit)
	// parameters:
	// - name: owner
	//   in: path
	//   description: owner of the repo
	//   type: string
	//   required: true
	// - name: repo
	//   in: path
	//   description: name of the repo
	//   type: string
	//   required: true
	// - name: ref
	//   in: path
	//   description: git reference (branch, tag, or commit)
	//   type: string
	//   required: true
	// responses:
	//   "200":
	//     description: Tree
	//     schema:
	//       type: array
	//       items:
	//         type: string
	
	ref := ctx.Params(":ref")
	
	repo := ctx.Repo.Repository
	if repo == nil {
		ctx.Error(http.StatusInternalServerError, "GetTreeByRef", "repository not found")
		return
	}
	
	// Get git repository
	gitRepo, err := git.OpenRepositoryCtx(ctx.Req.Context(), repo.RepoPath())
	if err != nil {
		ctx.Error(http.StatusInternalServerError, "GetTreeByRef", err)
		return
	}
	defer gitRepo.Close()
	
	// Use the new flexible ref function
	files, err := gitRepo.GetTreeFilesByRef(ctx.Req.Context(), ref)
	if err != nil {
		ctx.Error(http.StatusBadRequest, "GetTreeByRef", err)
		return
	}
	
	ctx.JSON(http.StatusOK, map[string]interface{}{
		"ref":   ref,
		"files": files,
	})
}

// GetCommitByRef handles GET /repos/{owner}/{repo}/git/commits/{ref}
// Returns commit information for a given reference
func GetCommitByRef(ctx *context.APIContext) {
	// swagger:operation GET /repos/{owner}/{repo}/git/commits/{ref} repository getCommitByRef
	// ---
	// summary: Get commit by reference
	// description: Returns commit information for a git reference
	// responses:
	//   "200":
	//     description: Commit info
	
	ref := ctx.Params(":ref")
	
	repo := ctx.Repo.Repository
	if repo == nil {
		ctx.Error(http.StatusInternalServerError, "GetCommitByRef", "repository not found")
		return
	}
	
	gitRepo, err := git.OpenRepositoryCtx(ctx.Req.Context(), repo.RepoPath())
	if err != nil {
		ctx.Error(http.StatusInternalServerError, "GetCommitByRef", err)
		return
	}
	defer gitRepo.Close()
	
	// Use the new commit info function with flexible ref
	info, err := gitRepo.GetCommitInfoByRef(ctx.Req.Context(), ref)
	if err != nil {
		ctx.Error(http.StatusBadRequest, "GetCommitByRef", err)
		return
	}
	
	ctx.JSON(http.StatusOK, map[string]string{
		"ref":  ref,
		"info": info,
	})
}

// DiffRefs handles GET /repos/{owner}/{repo}/git/diff/{from}...{to}
// Returns diff between two references
func DiffRefs(ctx *context.APIContext) {
	// swagger:operation GET /repos/{owner}/{repo}/git/diff/{from}...{to} repository diffRefs
	// ---
	// summary: Diff between refs
	// description: Returns the diff between two git references
	// responses:
	//   "200":
	//     description: Diff output
	
	fromRef := ctx.Params(":from")
	toRef := ctx.Params(":to")
	
	repo := ctx.Repo.Repository
	if repo == nil {
		ctx.Error(http.StatusInternalServerError, "DiffRefs", "repository not found")
		return
	}
	
	gitRepo, err := git.OpenRepositoryCtx(ctx.Req.Context(), repo.RepoPath())
	if err != nil {
		ctx.Error(http.StatusInternalServerError, "DiffRefs", err)
		return
	}
	defer gitRepo.Close()
	
	// Use the new diff function with flexible refs
	diff, err := gitRepo.DiffRefs(ctx.Req.Context(), fromRef, toRef)
	if err != nil {
		ctx.Error(http.StatusBadRequest, "DiffRefs", err)
		return
	}
	
	ctx.PlainText(http.StatusOK, diff)
}
