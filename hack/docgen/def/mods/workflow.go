/*
 Copyright 2022 The KubeVela Authors.

 Licensed under the Apache License, Version 2.0 (the "License");
 you may not use this file except in compliance with the License.
 You may obtain a copy of the License at

 	http://www.apache.org/licenses/LICENSE-2.0

 Unless required by applicable law or agreed to in writing, software
 distributed under the License is distributed on an "AS IS" BASIS,
 WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 See the License for the specific language governing permissions and
 limitations under the License.
*/

package mods

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"time"

	"github.com/oam-dev/kubevela/apis/types"
	"github.com/oam-dev/kubevela/pkg/utils/common"
	"github.com/oam-dev/kubevela/references/plugins"
)

const (
	// WorkflowDefRefPath is the target path for kubevela.io workflow ref docs
	WorkflowDefRefPath = "../kubevela.io/docs/end-user/workflow/built-in-workflow-defs.md"
	// WorkflowDefRefPathZh is the target path for kubevela.io workflow ref docs in Chinese
	WorkflowDefRefPathZh = "../kubevela.io/i18n/zh/docusaurus-plugin-content-docs/current/end-user/workflow/built-in-workflow-defs.md"

	// WorkflowDefDir store inner CUE definition
	WorkflowDefDir = "./vela-templates/definitions/internal/workflowstep/"
)

// CustomWorkflowHeaderEN .
var CustomWorkflowHeaderEN = `---
title: Built-in WorkflowStep Type
---

This documentation will walk through all the built-in workflow step types sorted alphabetically.

` + fmt.Sprintf("> It was generated automatically by [scripts](../../contributor/cli-ref-doc), please don't update manually, last updated at %s.\n\n", time.Now().Format(time.RFC3339))

// CustomWorkflowHeaderZH .
var CustomWorkflowHeaderZH = `---
title: 内置工作流步骤列表
---

本文档将**按字典序**展示所有内置工作流步骤的参数列表。

` + fmt.Sprintf("> 本文档由[脚本](../../contributor/cli-ref-doc)自动生成，请勿手动修改，上次更新于 %s。\n\n", time.Now().Format(time.RFC3339))

// WorkflowDef generate workflow def reference doc
func WorkflowDef(ctx context.Context, c common.Args, path, location *string, defdir string) {
	if defdir == "" {
		defdir = WorkflowDefDir
	}
	ref := &plugins.MarkdownReference{
		AllInOne: true,
		Filter: func(capability types.Capability) bool {

			if capability.Type != types.TypeWorkflowStep || capability.Category != types.CUECategory {
				return false
			}

			if capability.Labels != nil && (capability.Labels[types.LabelDefinitionHidden] == "true" || capability.Labels[types.LabelDefinitionDeprecated] == "true") {
				return false
			}
			// only print capability which contained in cue def
			files, err := ioutil.ReadDir(defdir)
			if err != nil {
				fmt.Println("read dir err", defdir, err)
				return false
			}
			for _, f := range files {
				if strings.Contains(f.Name(), capability.Name) {
					return true
				}
			}
			return false
		},
		CustomDocHeader: CustomWorkflowHeaderEN,
	}
	ref.Remote = &plugins.FromCluster{Namespace: types.DefaultKubeVelaNS}

	if *path != "" {
		ref.I18N = &plugins.En
		if strings.Contains(*location, "zh") || strings.Contains(*location, "chinese") {
			ref.I18N = &plugins.Zh
			ref.CustomDocHeader = CustomWorkflowHeaderZH
		}
		if err := ref.GenerateReferenceDocs(ctx, c, *path); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		fmt.Printf("workflow reference docs (%s) successfully generated in %s \n", ref.I18N.Language(), *path)
	}
	if *location == "" || *location == "en" {
		ref.I18N = &plugins.En
		if err := ref.GenerateReferenceDocs(ctx, c, WorkflowDefRefPath); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		fmt.Printf("workflow reference docs (%s) successfully generated in %s \n", ref.I18N.Language(), WorkflowDefRefPath)
	}
	if *location == "" || *location == "zh" {
		ref.I18N = &plugins.Zh
		ref.CustomDocHeader = CustomWorkflowHeaderZH
		if err := ref.GenerateReferenceDocs(ctx, c, WorkflowDefRefPathZh); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		fmt.Printf("workflow reference docs (%s) successfully generated in %s \n", ref.I18N.Language(), WorkflowDefRefPathZh)
	}
}
