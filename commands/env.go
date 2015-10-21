package commands

import (
	"github.com/blablacar/cnt/log"
	"github.com/blablacar/cnt/utils"
	"github.com/blablacar/green-garden/config"
	"github.com/blablacar/green-garden/work"
	"github.com/spf13/cobra"
	"strings"
)

func loadEnvCommands(rootCmd *cobra.Command) {
	log.Info("Work path is :", config.GetConfig().WorkPath)
	work := work.NewWork(config.GetConfig().WorkPath)

	for _, f := range work.ListEnvs() {
		var env = f
		var envCmd = &cobra.Command{
			Use:   env,
			Short: "Run command for " + env,
			Run: func(cmd *cobra.Command, args []string) {
				runner(cmd, args, work)
			},
		}

		var generateCmd = &cobra.Command{
			Use:   "generate",
			Short: "Generate units for " + env,
			Run: func(cmd *cobra.Command, args []string) {
				generate(cmd, args, work, env)
			},
		}
		envCmd.AddCommand(generateCmd)

		rootCmd.AddCommand(envCmd)
	}
}

func generate(cmd *cobra.Command, args []string, work *work.Work, env string) {
	log.Debug("Generate units for env " + env)
	work.LoadEnv(env).Generate()
}

func runner(cmd *cobra.Command, args []string, work *work.Work) {
	log.Info("Running command on " + cmd.Use)

	env := work.LoadEnv(cmd.Use)

	units, err := utils.ExecCmdGetOutput("fleetctl", "-strict-host-key-checking=false", "list-unit-files", "-no-legend", "-fields", "unit")
	if err != nil {
		panic("Cannot list unit files" + err.Error())
	}

	for _, unit := range strings.Split(units, "\n") {
		content, err := utils.ExecCmdGetOutput("fleetctl", "-strict-host-key-checking=false", "cat", unit)
		if err != nil {
			panic("Fleetctl failed to cat service content : " + unit)
		}
		unitInfo := strings.Split(unit, "_")
		if unitInfo[0] != cmd.Use {
			log.Warn("Unknown unit" + unit)
			continue
		}

		res, err := env.LoadService(unitInfo[1]).LoadUnit(unit).GetUnitContentAsFleeted()
		if err != nil {
			log.Warn("Cannot read unit file : " + unit)
			continue
		}
		if res != content {
			log.Info("Unit '" + unit + "' is not up to date")
			log.Debug(content)
			//			log.Trace(fmt.Sprintf("%x", content))
			log.Debug(res)
			//			log.Trace(fmt.Sprintf("%x", res))
		}
	}
}