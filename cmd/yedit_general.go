package cmd

import (
	"commandCenter/styles"
	"commandCenter/validators"
	"fmt"
	"log"
	"sync"

	"os"
	"path/filepath"
	"strings"

	"github.com/goccy/go-yaml/ast"
	"github.com/goccy/go-yaml/parser"
	"github.com/goccy/go-yaml/printer"
	"github.com/spf13/cobra"
)

var yamlUpdateGeneralCmd = &cobra.Command{
	Use:        "yedit",
	Short:      "Update yaml's general structure without fine-grained control over the edited blocks.",
	Aliases:    []string{"yedit", "general_yedit", "edit"},
	SuggestFor: []string{"yedit", "yaml_edit"},
	Example: `
      # Update the value of a top-level key in a YAML file
      ops yaml yedit -p ./config.yaml -k port -v 8080

      # Set the environment to production in a single YAML file
      ops yaml yedit -p ./app.yml -k environment -v production

      # Update the service name across all YAML files in a directory
      ops yaml yedit -p ./configs/ -k name -v my-updated-service

      # Enable a feature toggle in a YAML configuration
      ops yaml yedit -p ./settings.yaml -k featureX_enabled -v true

      # Get help for this command
      ops yaml yedit --help
    `,

	Run: processYamlGeneral,
}

// init initializes the yamlUpdateGeneralCmd and its flags.
//
// Args:
//   - None
//
// Returns:
//   - None
func init() {
	yamlCmd.AddCommand(yamlUpdateGeneralCmd)

	yamlUpdateGeneralCmd.Flags().StringP("path", "p", "examplePath", "file or directory path to the yaml files that should be updated")
	yamlUpdateGeneralCmd.Flags().StringP("key", "k", "exampleKey", "key that should be matched and updated")
	yamlUpdateGeneralCmd.Flags().StringP("value", "v", "exampleValue", "value that will be used for the matched key")

}

// walkAST traverses the AST of a YAML document.
//
// Args:
//   - node: The current AST node.
//   - visit: A function to apply to each mapping value node.
//
// Returns:
//   - None
func walkAST(node ast.Node, visit func(pair *ast.MappingValueNode, key string, val ast.Node)) {
	switch n := node.(type) {
	case *ast.MappingNode:
		for _, pair := range n.Values {
			keyNode, ok := pair.Key.(*ast.StringNode)
			if !ok {
				continue
			}
			visit(pair, keyNode.Value, pair.Value)

			walkAST(pair.Value, visit)
		}

	case *ast.SequenceNode:
		for _, elem := range n.Values {
			walkAST(elem, visit)
		}

	case *ast.AnchorNode:
		walkAST(n.Value, visit)
	}
}

// writeToYaml writes the given AST node to a YAML file.
//
// Args:
//   - rootNode: The root AST node to write.
//   - fileName: The name of the file to write to.
//   - hasStartingBlock: A boolean indicating if the original file had a "---" starting block.
//
// Returns:
//   - None
func writeToYaml(rootNode ast.Node, fileName string, hasStartingBlock bool) {
	pr := &printer.Printer{}
	src := pr.PrintNode(rootNode)

	var output []byte
	if hasStartingBlock && !strings.HasPrefix(string(src), "---") {
		output = append([]byte("---\n"), src...)
	} else {
		output = src
	}

	if err := os.WriteFile(fileName, []byte(output), 0644); err != nil {
		log.Fatalf("failed to write output file: %s", err)
	}

	fmt.Println(styles.NewStyles().Highlight.Render("Modified YAML written to" + fileName))
}

// parseYaml parses a YAML file and updates the specified key with the given value.
//
// Args:
//   - fullFilePath: The full path to the YAML file.
//   - k: The key to update.
//   - value: The new value for the key.
//
// Returns:
//   - None
func parseYaml(fullFilePath, k, value string) {
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
		walkAST(parsedFile.Docs[0].Body, func(pair *ast.MappingValueNode, key string, val ast.Node) {
			switch key {
			case k:
				if strVal, ok := val.(*ast.StringNode); ok {
					strVal.Value = value
				}
			}
		})
	}

	writeToYaml(rootNode, fullFilePath, hasStartingBlock)
}

// processAllYamlsNoFilter processes all YAML files in a given path (file or directory) without filtering.
//
// Args:
//   - path: The path to the YAML file or directory.
//   - k: The key to update.
//   - value: The new value for the key.
//
// Returns:
//   - None
func processAllYamlsNoFilter(path, k, value string) {
	entry, err := os.Stat(path)
	if err != nil {
		log.Fatalf("error fetching system path: %s", err)
	}

	var wg sync.WaitGroup

	switch mode := entry.Mode(); {
	case mode.IsDir():
		files, err := os.ReadDir(path)
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
					fullFilePath := filepath.Join(path, file.Name())
					parseYaml(fullFilePath, k, value)
				}
			}()
		}

		wg.Wait()
	case mode.IsRegular():
		parseYaml(path, k, value)
	}
}

func processYamlGeneral(cmd *cobra.Command, args []string) {
	path, err := validators.VerifyStringInputs(cmd, "path")
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

	processAllYamlsNoFilter(path, key, value)
}
