/*
 * Copyright 2022 Sue B.V.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *      http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package cmd

import (
	"github.com/go-logr/logr"
	"github.com/spf13/cobra"
	"github.com/suecodelabs/cnfuzz/src/config"
	"github.com/suecodelabs/cnfuzz/src/k8s"
	"github.com/suecodelabs/cnfuzz/src/logger"
	"github.com/suecodelabs/cnfuzz/src/persistence"
	"log"
	"os"
)

type CnFuzzCommand struct {
	command *cobra.Command
	*Args
}

type Args struct {
	isDebug     bool
	localConfig bool
	configFile  string
	printConfig bool
	dDocIp      string
	dDocPort    int32
}

// Execute starts the base command
func Execute() {
	cmd := CnFuzzCommand{
		command: &cobra.Command{
			Use:   "cnfuzz <flags>",
			Short: "Native Cloud Fuzzer is a fuzzer build for cloud native environments",
			Long: `Native Cloud Fuzzer is a fuzzer build for cloud native environments.
More info here:
https://github.com/suecodelabs/cnfuzz`,
		},
		Args: &Args{
			isDebug:     false,
			localConfig: false,
			configFile:  "/config/config.yaml",
			printConfig: false,
			dDocIp:      "",
			dDocPort:    0,
		},
	}

	cmd.command.PersistentFlags().BoolVarP(&cmd.Args.isDebug, "debug", "d", cmd.isDebug, "Enable debug mode")
	cmd.command.PersistentFlags().BoolVar(&cmd.Args.localConfig, "local-config", cmd.Args.localConfig, "Use the local kubeconfig instead of getting it from the cluster")
	cmd.command.PersistentFlags().BoolVar(&cmd.Args.printConfig, "print-config", cmd.Args.printConfig, "Print the config file")
	cmd.command.PersistentFlags().StringVar(&cmd.Args.configFile, "config", cmd.Args.configFile, "Location of the config file to use")
	cmd.command.PersistentFlags().StringVar(&cmd.Args.dDocIp, "ddoc-ip", cmd.Args.dDocIp, "Overwrite the IP address cnfuzz uses to get the discovery doc (useful for developers)")
	cmd.command.PersistentFlags().Int32Var(&cmd.Args.dDocPort, "ddoc-port", cmd.Args.dDocPort, "Overwrites the port cnfuzz uses to get the discovery doc (useful for developers)")

	cmd.command.Run = func(_ *cobra.Command, _ []string) {
		l := logger.CreateLogger(cmd.Args.isDebug)
		run(l, *cmd.Args)
	}

	if err := cmd.command.Execute(); err != nil {
		log.Fatalln(err)
	}
}

func run(l logr.Logger, args Args) {
	l.Info("starting cnfuzz")

	cnf := config.LoadConfigFile(l, args.configFile, args.printConfig)
	overwrites := config.Overwrites{
		DiscoveryDocIP:   args.dDocIp,
		DiscoveryDocPort: args.dDocPort,
	}

	var strg *persistence.Storage
	if cnf.CacheSolution == persistence.Redis.String() || true {
		l.V(logger.InfoLevel).Info("using redis for storage", "configStorageValue", cnf.CacheSolution)
		strg = persistence.InitRedisCache(l, cnf.RedisConfig.HostName, cnf.RedisConfig.Port)
	} else {
		l.V(logger.DebugLevel).Info("using in_memory for storage", "configStorageValue", cnf.CacheSolution)
		strg = persistence.InitMemoryCache(l)
	}

	client := k8s.CreateClientset(l, !args.localConfig)
	// Start fuzzing!
	err := k8s.StartController(l, strg, cnf, overwrites, client)
	if err != nil {
		l.Error(err, "error while starting cnfuzz controller")
		os.Exit(1)
	}
}
