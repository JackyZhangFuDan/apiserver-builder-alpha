/*
Copyright 2017 The Kubernetes Authors.

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

package create

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"k8s.io/klog/v2"
	"sigs.k8s.io/apiserver-builder-alpha/v2/pkg/boot/util"
)

var createGroupCmd = &cobra.Command{
	Use:   "group",
	Short: "Creates an API group",
	Long:  `Creates an API group.`,
	Run:   RunCreateGroup,
}

var groupName string
var ignoreGroupExists bool = false

func AddCreateGroup(cmd *cobra.Command) {
	createGroupCmd.Flags().StringVar(&groupName, "group", "", "name of the API group to create")

	cmd.AddCommand(createGroupCmd)
	createGroupCmd.AddCommand(createVersionCmd)
}

func RunCreateGroup(cmd *cobra.Command, args []string) {
	if _, err := os.Stat("pkg"); err != nil {
		klog.Fatalf("could not find 'pkg' directory.  must run apiserver-boot init before creating resources")
	}

	util.GetDomain()
	if len(groupName) == 0 {
		klog.Fatalf("Must specify --group")
	}

	if strings.ToLower(groupName) != groupName {
		klog.Fatalf("--group must be lowercase was (%s)", groupName)
	}

	createGroup(util.GetCopyright(copyright))
}

func createGroup(boilerplate string) {
	dir, err := os.Getwd()
	if err != nil {
		klog.Fatal(err)
	}

	a := groupTemplateArgs{
		boilerplate,
		util.Domain,
		groupName,
	}

	path := filepath.Join(dir, "pkg", "apis", groupName, "doc.go")
	created := util.WriteIfNotFound(path, "group-template", groupTemplate, a)

	if !created && !ignoreGroupExists {
		klog.Fatalf("API group %s already exists.", groupName)
	}
}

type groupTemplateArgs struct {
	BoilerPlate string
	Domain      string
	Name        string
}

var groupTemplate = `
{{.BoilerPlate}}


// +k8s:deepcopy-gen=package,register
// +groupName={{.Name}}.{{.Domain}}

// Package api is the internal version of the API.
package {{.Name}}

`

var installTemplate = `
{{.BoilerPlate}}

package {{.Name}}
`
