package cmd

import (
	"commandCenter/styles"
	"commandCenter/validators"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"

	"github.com/spf13/cobra"
)

var metaDataIngestCmd = &cobra.Command{
	Use:   "ingest",
	Short: "Ingest metadata (requires Connector and yaml config).",
	Long: `Allows the ingestion of a new service depending on the connector
    and config provided.`,
	Example: `
      # Ingest metadata using default pyproject.toml and config.yaml files
      ops ingest

      # Ingest metadata specifying a custom pyproject.toml file and config.yaml
      ops ingest -p path/to/pyproject.toml -c path/to/config.yaml

      # Show help for the ingest command
      ops ingest --help
    `,

	Run: ingest,
}

// init initializes the ingest command and its flags.
//
// Args:
//   - None
//
// Returns:
//   - None
func init() {
	metadataCmd.AddCommand(metaDataIngestCmd)

	metaDataIngestCmd.Flags().StringP("project_toml", "p", "pyproject.toml", "Path to your pyproject.toml file")
	metaDataIngestCmd.Flags().StringP("config", "c", "config.yaml", "Path to yaml config file")
}

type VirtualPythonEnv struct {
	executablePath           string
	pathToPyProjectToml      string
	pathToOpenMetaDataConfig string
	venvDirName              string
}

type PythonVirtualEnvSetup interface {
	runOpenMetaDataIngestion()
}

// checkDependency checks if a binary dependency is installed and in the system's PATH.
//
// Args:
//   - binaryDependencyName: The name of the binary to check.
//
// Returns:
//   - None
func checkDependency(binaryDependencyName string) {
	if _, err := exec.LookPath(binaryDependencyName); err != nil {
		message := fmt.Sprintf("%s, is not installed or not in $PATH", binaryDependencyName)
		log.Fatalln(styles.StyliseMessage(message, styles.FormatStyle.Error))
	} else {
		message := fmt.Sprintf("%s is installed", binaryDependencyName)
		log.Println(styles.StyliseMessage(message, styles.FormatStyle.Highlight))
	}
}

// getPwd gets the current working directory.
//
// Args:
//   - None
//
// Returns:
//   - string: The current working directory.
func getPwd() string {
	pwd, err := os.Getwd()
	if err != nil {
		message := fmt.Sprintf("Failed to get working directory: %s", err.Error())
		log.Fatalln(styles.StyliseMessage(message, styles.FormatStyle.Error))
	}
	return pwd
}

// executeCommand executes a shell command and logs the output.
//
// Args:
//   - command: The command to execute.
//   - args: The arguments to the command.
//
// Returns:
//   - None
func executeCommand(command string, args ...string) {
	cmd := exec.Command(command, args...)
	if strings.Contains(command, "metadata") {
		cmd.Env = append(os.Environ(), "PYTHONPATH="+getPwd())
	}
	output, err := cmd.CombinedOutput()
	if err != nil {
		log.Println(styles.StyliseMessage("Error encountered, cleaning up temporary resources", styles.FormatStyle.Highlight))
		cleanUp(".venv")
		message := fmt.Sprintf("Error executing command %s - %s\nCommand output (stderr/stdout):\n%s\n", command, err, output)
		log.Fatalf(styles.StyliseMessage(message, styles.FormatStyle.Error))
	}

	log.Printf(styles.StyliseMessage(fmt.Sprintf("Command executed successfully: %s%s", command, args), styles.FormatStyle.Highlight))
}

// runOpenMetaDataIngestion runs the Open-Metadata ingestion process.
//
// Args:
//   - None
//
// Returns:
//   - None
func (V VirtualPythonEnv) runOpenMetaDataIngestion() {
	executeCommand(V.executablePath, "ingest", "-c", V.pathToOpenMetaDataConfig)
}

// cleanUp removes a directory.
//
// Args:
//   - dirToRemove: The directory to remove.
//
// Returns:
//   - None
func cleanUp(dirToRemove string) {
	err := os.RemoveAll(dirToRemove)
	if err != nil {
		log.Fatal(err)
	}

	log.Println(styles.StyliseMessage("Temporary resources cleaned up", styles.FormatStyle.Highlight))
}

// ingestMetaData runs the Open-Metadata ingestion process.
//
// Args:
//   - P: The PythonVirtualEnvSetup interface.
//
// Returns:
//   - None
func ingestMetaData(P PythonVirtualEnvSetup) {
	log.Println(styles.StyliseMessage("Starting MetaData Ingestion...", styles.FormatStyle.Highlight))
	P.runOpenMetaDataIngestion()
}

// ingest is the main function for the ingest command.
//
// Args:
//   - cmd: The cobra command.
//   - args: The command arguments.
//
// Returns:
//   - None
func ingest(cmd *cobra.Command, args []string) {
	projectToml, err := validators.VerifyStringInputs(cmd, "project_toml")
	if err != nil {
		log.Fatalln(err)
	}

	yamlConfig, err := validators.VerifyStringInputs(cmd, "config")
	if err != nil {
		log.Fatalln(err)
	}

	pythonVirtualEnv := VirtualPythonEnv{
		executablePath:           "./.venv/bin/metadata",
		venvDirName:              ".venv",
		pathToPyProjectToml:      projectToml,
		pathToOpenMetaDataConfig: yamlConfig,
	}

	checkDependency("uv")
	executeCommand("uv", "venv")
	executeCommand("uv", "add", ".", "--dev")
	ingestMetaData(pythonVirtualEnv)
	cleanUp(".venv")
}
