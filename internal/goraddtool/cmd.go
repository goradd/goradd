package goraddtool

import (
	"github.com/spf13/cobra"
)

func MakeRootCommand() *cobra.Command {
	var overwrite bool
	var step int
	var browser bool
	var headless bool

	var rootCmd = &cobra.Command{
		Use:   "goradd",
		Short: "goradd is a tool for installing goradd, generating code, and building a web application",
		Long:  `goradd is a web application framework. This command-line tool is a helper for installing the framework, generating code, seeing intermediate results, and ultimately building your web application.`,
	}

	var cmdInstall = &cobra.Command{
		Use:   "install",
		Short: "Install the goradd-project directory in the current working directory",
		Long:  `Install the goradd-project directory in the current working directory. Use the -r flag to force replacement of the directory.`,
		Run: func(cmd *cobra.Command, args []string) {
			install(step, overwrite)
		},
	}

	cmdInstall.Flags().BoolVarP(&overwrite, "replace", "r", false, "Previous directories will be deleted first without asking. Use with extreme caution.")
	cmdInstall.Flags().IntVarP(&step, "step", "s", 0, "For debugging the tests, this will execute only the given step. 0 means execute the entire test process.")

	var cmdTest = &cobra.Command{
		Use:   "test",
		Short: "Installs the goradd-test directory in the current working directory, and then runs the tests, mimicking the continuous integration test",
		Long:  `Installs the goradd-test directory in the current working directory, and then runs the tests, mimicking the continuous integration test`,
		Run: func(cmd *cobra.Command, args []string) {
			test(step, browser, headless)
		},
	}

	cmdTest.Flags().IntVarP(&step, "step", "s", 0, "For debugging the tests, this will execute only the given step. 0 means execute the entire test process.")
	cmdTest.Flags().BoolVarP(&browser, "browser", "b", false, "Whether to launch the chrome browser for the browser based tests. If you don't specify this, you should already have a browser running at localhost:8000?all=1")
	//cmdTest.Flags().BoolVarP(&headless, "headless", "l", false, "Whether to launch the browser as a headless browser.")

	rootCmd.AddCommand(cmdInstall, cmdTest)

	return rootCmd
}
