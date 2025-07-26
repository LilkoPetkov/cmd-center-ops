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

func checkDependency(binaryDependencyName string) {
	if _, err := exec.LookPath(binaryDependencyName); err != nil {
		message := fmt.Sprintf("%s, is not installed or not in $PATH", binaryDependencyName)
		log.Fatalln(styles.StyliseMessage(message, styles.FormatStyle.Error))
	} else {
		message := fmt.Sprintf("%s is installed", binaryDependencyName)
		log.Println(styles.StyliseMessage(message, styles.FormatStyle.Highlight))
	}
}

func getPwd() string {
	pwd, err := os.Getwd()
	if err != nil {
		message := fmt.Sprintf("Failed to get working directory: %s", err.Error())
		log.Fatalln(styles.StyliseMessage(message, styles.FormatStyle.Error))
	}
	return pwd
}

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

func (V VirtualPythonEnv) runOpenMetaDataIngestion() {
	executeCommand(V.executablePath, "ingest", "-c", V.pathToOpenMetaDataConfig)
}

func cleanUp(dirToRemove string) {
	err := os.RemoveAll(dirToRemove)
	if err != nil {
		log.Fatal(err)
	}

	log.Println(styles.StyliseMessage("Temporary resources cleaned up", styles.FormatStyle.Highlight))
}

func ingestMetaData(P PythonVirtualEnvSetup) {
	log.Println(styles.StyliseMessage("Starting MetaData Ingestion...", styles.FormatStyle.Highlight))
	P.runOpenMetaDataIngestion()
}

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
