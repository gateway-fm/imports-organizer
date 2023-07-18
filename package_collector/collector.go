package package_collector

import (
	"golang.org/x/tools/go/packages"
	"strings"
)

func init() {
	LoadModulePackages()
	removeDuplicateStr()
}

var Packages = make(map[string]int)

func Populate() {
	for i := range PackagesSlice {
		Packages[PackagesSlice[i]] = 3
	}
}

func ExternalPackages() {
	for i := range PackagesSlice {
		path := strings.Split(PackagesSlice[i], "/")
		if len(path) >= 2 {
			if path[1] == "gateway-fm" {
				str := strings.Join(path, "/")
				Packages[str] = 2
			} else if path[0] == "github.com" || path[0] == "google.golang.org" || path[0] == "golang.org" || path[0] == "go.uber.org" || path[0] == "gopkg.in" || path[0] == "go.opentelemetry.io" && path[1] != "your-project" {
				str := strings.Join(path, "/")
				Packages[str] = 1
			}
		}
	}
}
func Std() {
	cfg := &packages.Config{Mode: packages.NeedFiles | packages.NeedSyntax}
	pkgs, err := packages.Load(cfg, "std")
	if err == nil {
		for _, p := range pkgs {
			Packages[p.ID] = 0
		}
	}
}
func IsProto(name string) bool {
	path := strings.Split(name, "/")

	for i := range path {
		if strings.Contains(path[i], ".") {
			bbb := strings.Split(path[i], ".")
			for ii := range bbb {
				if strings.Contains(bbb[ii], "pb") {
					return true
				}
			}

		}
	}
	return false
}

var Goes []string

func GoFiles() []string {
	cfg := &packages.Config{Mode: packages.NeedFiles}

	pkgs, err := packages.Load(cfg, "./...")
	if err == nil {
		for _, pkg := range pkgs {
			for i := range pkg.GoFiles {
				if !IsProto(pkg.GoFiles[i]) {
					Goes = append(Goes, pkg.GoFiles[i])
				}
			}

		}
	}
	return Goes
}

var PackagesSlice []string

func LoadModulePackages() {
	cfg := &packages.Config{Mode: packages.NeedImports}

	pkgs, err := packages.Load(cfg, "./...")
	if err == nil {
		for _, pkg := range pkgs {
			for i := range pkg.Imports {
				PackagesSlice = append(PackagesSlice, pkg.Imports[i].ID)
			}
		}
	}

}
func removeDuplicateStr() {
	allKeys := make(map[string]bool)
	for _, item := range PackagesSlice {
		if _, value := allKeys[item]; !value {
			allKeys[item] = true
			PackagesSlice = append(PackagesSlice, item)
		}
	}
}
