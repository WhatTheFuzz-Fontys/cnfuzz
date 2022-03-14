package cmd

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	IsDebug = "debug"

	// TargetPodName start job with Kubernetes
	TargetPodName      = "pod"
	TargetPodNamespace = "namespace" // Namespace that target lives in

	// Misc Kubernetes flags
	InsideClusterFlag  = "inside-cluster"
	OnlyFuzzMarkedFlag = "only-marked"
	JobImageFlag       = "job-img"
	HomeNamespaceFlag  = "home-ns" // Namespace to start containers in (jobs etc.)

	// Fuzzing related flags
	RestlerInitImageFlag = "restler-init-img"
	RestlerImageFlag     = "restler-img"
	RestlerTimeBudget    = "restler-time-budget"

	// Auth related flags
	AuthUsername   = "username"
	AuthSecretFlag = "secret"
)

func SetupFlags(rootCmd *cobra.Command) {
	// Debug flag
	rootCmd.Flags().BoolP(IsDebug, "d", false, "Debug mode")
	_ = viper.BindPFlag(IsDebug, rootCmd.Flags().Lookup(IsDebug))

	registerDirectFuzzingFlags(rootCmd)

	registerKubernetesFlags(rootCmd)

	rootCmd.Flags().StringP(RestlerInitImageFlag, "", "curlimages/curl:7.81.0", "Init Image for preparing RESTler runtime")
	_ = viper.BindPFlag(RestlerInitImageFlag, rootCmd.Flags().Lookup(RestlerInitImageFlag))

	rootCmd.Flags().StringP(RestlerImageFlag, "", "mcr.microsoft.com/restlerfuzzer/restler:v7.4.0", "RESTler image to use (https://hub.docker.com/_/microsoft-restlerfuzzer-restler)")
	_ = viper.BindPFlag(RestlerImageFlag, rootCmd.Flags().Lookup(RestlerImageFlag))

	rootCmd.Flags().StringP(RestlerTimeBudget, "", "1", "Maximum hours a Fuzzing Job may take.")
	_ = viper.BindPFlag(RestlerTimeBudget, rootCmd.Flags().Lookup(RestlerTimeBudget))

	registerAuthFlags(rootCmd)
}

// registerDirectFuzzingFlags Register flags used when directly fuzzing a target
func registerDirectFuzzingFlags(rootCmd *cobra.Command) {
	rootCmd.Flags().StringP(TargetPodName, "", "", "Kubernetes pod to target for fuzzing")
	_ = viper.BindPFlag(TargetPodName, rootCmd.Flags().Lookup(TargetPodName))

	rootCmd.Flags().StringP(TargetPodNamespace, "n", "default", "Namespace for the target pod")
	_ = viper.BindPFlag(TargetPodNamespace, rootCmd.Flags().Lookup(TargetPodNamespace))
}

func registerKubernetesFlags(rootCmd *cobra.Command) {
	rootCmd.Flags().BoolP(InsideClusterFlag, "k", true, "Tells the fuzzer that it is running inside Kubernetes")
	_ = viper.BindPFlag(InsideClusterFlag, rootCmd.Flags().Lookup(InsideClusterFlag))

	rootCmd.Flags().BoolP(OnlyFuzzMarkedFlag, "m", false, "Only fuzz pods that have a '\"cnfuzz/enable\": \"true\"' annotation")
	_ = viper.BindPFlag(OnlyFuzzMarkedFlag, rootCmd.Flags().Lookup(OnlyFuzzMarkedFlag))

	// TODO change current temp default image to actual image ones it exists
	defaultJImg := ""
	rootCmd.PersistentFlags().StringP(JobImageFlag, "i", defaultJImg, "Image used for the Kubernetes job, you can use this to change image version or replace the entire image")
	_ = viper.BindPFlag(JobImageFlag, rootCmd.PersistentFlags().Lookup(JobImageFlag))

	rootCmd.Flags().StringP(HomeNamespaceFlag, "", "default", "Namespace to start fuzzing containers in")
	_ = viper.BindPFlag(HomeNamespaceFlag, rootCmd.Flags().Lookup(HomeNamespaceFlag))
}

func registerAuthFlags(rootCmd *cobra.Command) {
	rootCmd.PersistentFlags().StringP(AuthUsername, "", "fuzz-client", "Username to be used in auth (only works for basic auth)")
	rootCmd.PersistentFlags().StringP(AuthSecretFlag, "", "", "Secret to be used for authentication")
	_ = viper.BindPFlag(AuthUsername, rootCmd.PersistentFlags().Lookup(AuthUsername))
	_ = viper.BindPFlag(AuthSecretFlag, rootCmd.PersistentFlags().Lookup(AuthSecretFlag))
	// TODO add a scopes flag
}
