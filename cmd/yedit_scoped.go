package cmd

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"slices"
	"strings"
	"sync"

	"commandCenter/styles"
	"commandCenter/validators"

	"github.com/goccy/go-yaml/ast"
	"github.com/goccy/go-yaml/parser"
	"github.com/spf13/cobra"
)

type YamlConfigDetails struct {
	Path       string
	Substr     *regexp.Regexp
	Key        string
	Value      string
	Exceptions []string
}

var yamlUpdateScopedCmd = &cobra.Command{
	Use:        "yedit_scoped",
	Short:      "Update yaml's scoped value for more fine-grained control.",
	Aliases:    []string{"yedit_scope", "yaml_edit_scoped", "edit_yaml_scoped"},
	SuggestFor: []string{"edit", "yaml_edi", "editor"},
	Example: `
      # Update the value of a specific key inside a matched YAML block
      ops yaml yedit_scoped -p ./config.yaml -s "database" -k "host" -v "localhost"

      # Apply scoped update across all YAML files in a directory
      ops yaml yedit_scoped -p ./configs/ -s "authService" -k "port" -v "8081"

      # Update a nested key where block contains the word "cache"
      ops yaml yedit_scoped -p settings.yml -s "cache" -k "enabled" -v "true"

      # Match a block with a regex (e.g., block name contains "prod")
      ops yaml yedit_scoped -p app.yaml -s "prod.*" -k "url" -v "https://prod.example.com"

      # Preview available flags and help
      ops yaml yedit_scoped --help
    `,

	Run: processYaml,
}

// init initializes the yamlUpdateScopedCmd and its flags.
//
// Args:
//   - None
//
// Returns:
//   - None
func init() {
	yamlCmd.AddCommand(yamlUpdateScopedCmd)

	yamlUpdateScopedCmd.Flags().StringP("path", "p", "examplePath", "file or directory path to the yaml files that should be updated")
	yamlUpdateScopedCmd.Flags().StringP("substr", "s", "exampleSubstr", "string or substring to match in the yaml and update only its block")
	yamlUpdateScopedCmd.Flags().StringP("key", "k", "exampleKey", "key that should be matched and updated")
	yamlUpdateScopedCmd.Flags().StringP("value", "v", "exampleValue", "value that will be used for the matched key")

	yamlUpdateScopedCmd.Flags().StringSliceP("exceptions", "e", []string{}, "exception regexes of a node that should be skipped completely")
}

// walkAndEditFromMatchedKey traverses the AST and edits values based on a matched key.
//
// Args:
//   - node: The current AST node.
//
// Returns:
//   - None
func (Y YamlConfigDetails) walkAndEditFromMatchedKey(
	node ast.Node,
) {
	switch n := node.(type) {
	case *ast.MappingNode:
		for _, pair := range n.Values {
			keyNode, ok := pair.Key.(*ast.StringNode)
			if !ok {
				continue
			}
			key := keyNode.Value

			if slices.Contains(Y.Exceptions, key) {
				continue
			}

			if Y.Substr.MatchString(key) {
				Y.drillAndEdit(pair.Value)
			}

			Y.walkAndEditFromMatchedKey(pair.Value)
		}

	case *ast.SequenceNode:
		for _, elem := range n.Values {
			Y.walkAndEditFromMatchedKey(elem)
		}

	case *ast.AnchorNode:
		Y.walkAndEditFromMatchedKey(n.Value)
	}
}

// drillAndEdit drills down into the AST and edits the value of a specific key.
//
// Args:
//   - node: The current AST node.
//
// Returns:
//   - None
func (Y YamlConfigDetails) drillAndEdit(
	node ast.Node,
) {
	switch n := node.(type) {
	case *ast.MappingNode:
		for _, pair := range n.Values {
			if keyNode, ok := pair.Key.(*ast.StringNode); ok && keyNode.Value == Y.Key {
				if valNode, ok := pair.Value.(*ast.StringNode); ok {
					valNode.Value = Y.Value
				}
			}
			Y.drillAndEdit(pair.Value)
		}

	case *ast.SequenceNode:
		for _, elem := range n.Values {
			Y.drillAndEdit(elem)
		}

	case *ast.AnchorNode:
		Y.drillAndEdit(n.Value)
	}
}

// parseYamlFilter parses a YAML file and applies the scoped edit.
//
// Args:
//   - fullFilePath: The full path to the YAML file.
//
// Returns:
//   - None
func (Y YamlConfigDetails) parseYamlFilter(fullFilePath string) {
	data, err := os.ReadFile(fullFilePath)
	if err != nil {
		log.Fatalf("failed to read file: %s", err)
	}

	hasStartingBlock := strings.HasPrefix(strings.TrimSpace(string(data)), "---")

	parsedFile, err := parser.ParseBytes(data, 0)
	if err != nil {
		log.Fatalf("failed to parse YAML: %s", err)
	}
	rootNode := parsedFile.Docs[0].Body

	if len(parsedFile.Docs) > 0 && rootNode != nil {
		Y.walkAndEditFromMatchedKey(rootNode)
		writeToYaml(rootNode, fullFilePath, hasStartingBlock)
	}
}

// processYamlsFilter processes YAML files based on the configured details.
//
// Args:
//   - None
//
// Returns:
//   - None
func (Y YamlConfigDetails) processYamlsFilter() {
	entry, err := os.Stat(Y.Path)
	if err != nil {
		log.Fatalf("error fetching system path: %s", err)
	}

	var wg sync.WaitGroup

	switch mode := entry.Mode(); {
	case mode.IsDir():
		files, err := os.ReadDir(Y.Path)
		if err != nil {
			log.Fatalf("could not read directory: %s", err)
		}

		for _, file := range files {
			currentFile := file

			wg.Add(1)

			go func() {
				defer wg.Done()

				fileName := currentFile.Name()

				if strings.Contains(fileName, ".yaml") || strings.Contains(fileName, ".yml") {
					fullFilePath := filepath.Join(Y.Path, file.Name())
					Y.parseYamlFilter(fullFilePath)
				}
			}()

		}

		wg.Wait()
	case mode.IsRegular():
		Y.parseYamlFilter(Y.Path)
	}
}

func processYaml(cmd *cobra.Command, args []string) {
	path, err := validators.VerifyStringInputs(cmd, "path")
	if err != nil {
		log.Fatalln(err)
	}

	substr, err := validators.VerifyStringInputs(cmd, "substr")
	if err != nil {
		log.Fatalln(err)
	}

	key, err := validators.VerifyStringInputs(cmd, "key")
	if err != nil {
		log.Fatalln(err)
	}

	value, err := validators.VerifyStringInputs(cmd, "value")
	if err != nil {
		log.Fatalln(err)
	}

	compiledPattern, err := regexp.Compile(substr)
	if err != nil {
		fmt.Printf(styles.NewStyles().Error.Render("error compiling regex pattern: %s"), err)
	}

	exceptions, err := cmd.Flags().GetStringSlice("exceptions")
	if err != nil {
		fmt.Printf(styles.NewStyles().Error.Render("error parsing 'exceptions': %s"), err)
	}

	yamlConfig := YamlConfigDetails{
		Path:       path,
		Substr:     compiledPattern,
		Key:        key,
		Value:      value,
		Exceptions: exceptions,
	}

	yamlConfig.processYamlsFilter()
}
