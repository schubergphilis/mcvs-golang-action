package application

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"log/slog"
	"os"
	"path/filepath"
	"strings"

	"github.com/schubergphilis/mcvs-golang-action/internal/app/mcvs_golang_action_bootstrap/data"
)

type Executor interface {
	DownloadFiles(jobs []DownloadJob) error
	GenerateMockeryYaml(moduleName, dir, outPath string, logger *slog.Logger) error
}

type Downloader struct {
	Storer data.Storer
	Logger *slog.Logger
}

func (d *Downloader) DownloadFiles(jobs []DownloadJob) error {
	for _, job := range jobs {
		d.Logger.Info("Processing download job", "url", job.URL, "destination", job.Destination)

		content, err := d.Storer.DownloadFileFromURL(job.URL)
		if err != nil {
			d.Logger.Error("Failed to download", "url", job.URL, "err", err)

			return fmt.Errorf("download failed for %s: %w", job.URL, err)
		}

		err = d.Storer.WriteToFile(job.Destination, content)
		if err != nil {
			d.Logger.Error("Failed to write file", "dest", job.Destination, "err", err)

			return fmt.Errorf("write failed for %s: %w", job.Destination, err)
		}
	}

	return nil
}

func FindInterfacesInPackages(moduleName, rootDir string, logger *slog.Logger) (map[string][]string, error) {
	pkgInterfaces := make(map[string][]string)
	err := filepath.WalkDir(rootDir, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if isSkippableDir(d) {
			return filepath.SkipDir
		}

		if !isGoSourceFile(path) {
			return nil
		}

		interfaces, parseErr := parseInterfacesInFile(path)
		if parseErr != nil {
			logger.Warn("Could not parse file", "file", path, "err", parseErr)

			return nil
		}

		if len(interfaces) == 0 {
			return nil
		}

		relDir, _ := filepath.Rel(rootDir, filepath.Dir(path))
		if relDir == "." {
			relDir = ""
		}

		importPath := moduleName
		if relDir != "" {
			importPath = moduleName + "/" + filepath.ToSlash(relDir)
		}

		for _, iface := range interfaces {
			if !contains(pkgInterfaces[importPath], iface) {
				pkgInterfaces[importPath] = append(pkgInterfaces[importPath], iface)
			}
		}

		return nil
	})

	return pkgInterfaces, err
}

func parseInterfacesInFile(filename string) ([]string, error) {
	interfaces := []string{}

	file, err := parser.ParseFile(token.NewFileSet(), filename, nil, parser.ParseComments)
	if err != nil {
		return interfaces, err
	}

	for _, decl := range file.Decls {
		gen, ok := decl.(*ast.GenDecl)
		if !ok || gen.Tok != token.TYPE {
			continue
		}

		for _, spec := range gen.Specs {
			tSpec, ok := spec.(*ast.TypeSpec)
			if !ok {
				continue
			}

			if _, ok := tSpec.Type.(*ast.InterfaceType); ok {
				interfaces = append(interfaces, tSpec.Name.Name)
			}
		}
	}

	return interfaces, nil
}

func isSkippableDir(d os.DirEntry) bool {
	return d.IsDir() && (d.Name() == "vendor" || d.Name() == "testdata" || strings.HasPrefix(d.Name(), "."))
}

func isGoSourceFile(path string) bool {
	return strings.HasSuffix(path, ".go") && !strings.HasSuffix(path, "_test.go")
}

func contains(list []string, val string) bool {
	for _, v := range list {
		if v == val {
			return true
		}
	}

	return false
}

func (d *Downloader) GenerateMockeryYaml(moduleName, dir, outPath string, logger *slog.Logger) error {
	logger.Info("Scanning for interfaces and generating .mockery.yaml", "dir", dir)

	pkgIfaces, err := FindInterfacesInPackages(moduleName, dir, logger)
	if err != nil {
		return err
	}

	var stringsBuilder strings.Builder
	stringsBuilder.WriteString(`---
dir: "{{.InterfaceDir}}/mocks"
filename: "{{.InterfaceName | snakecase}}.go"
structname: "{{.InterfaceName}}"
pkgname: mocks
template: testify
packages:
`)

	for pkg, ifaces := range pkgIfaces {
		stringsBuilder.WriteString(fmt.Sprintf("  %s:\n", pkg))
		stringsBuilder.WriteString("    interfaces:\n")

		for _, iface := range ifaces {
			stringsBuilder.WriteString(fmt.Sprintf("      %s: {}\n", iface))
		}
	}

	yaml := stringsBuilder.String()

	return d.Storer.WriteTextFile(outPath, yaml)
}
