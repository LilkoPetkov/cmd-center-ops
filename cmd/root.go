/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"commandCenter/styles"
	"context"
	"github.com/charmbracelet/fang"
	"github.com/charmbracelet/lipgloss"
	"github.com/spf13/cobra"
	"os"
)

func myColorScheme() fang.ColorScheme {
	base := lipgloss.Color("#FFFFFF")
	return fang.ColorScheme{
		Base:         base,
		Title:        lipgloss.Color("#FF0000"),
		Description:  base,
		Comment:      lipgloss.Color("#AAAAAA"),
		Flag:         lipgloss.Color("#FF00FF"),
		FlagDefault:  lipgloss.Color("#FF77FF"),
		Command:      lipgloss.Color("#00FFFF"),
		QuotedString: lipgloss.Color("#00FF00"),
		Argument:     base,
		Help:         base,
		Dash:         base,
		ErrorDetails: lipgloss.Color("#FF0000"),
	}
}

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "ops",
	Short: styles.NewStyles().Title.Render("🔥 OPS🔥"),
	Long: `🌟 ops is a powerful, user-friendly Command Line Interface (CLI) application built for developers, system administrators, network engineers, and cybersecurity professionals. It provides precise and flexible tools for domain name resolution, DNS diagnostics, and DNS record retrieval directly from the terminal.🌟
`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) { },
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	opts := []fang.Option{
		fang.WithVersion("v0.0.1"),
	}

	if err := fang.Execute(context.Background(), rootCmd, opts...); err != nil {
		os.Exit(1)
	}

	// err := rootCmd.Execute()
	// if err != nil {
	// 	os.Exit(1)
	// }
}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.commandCenter.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
