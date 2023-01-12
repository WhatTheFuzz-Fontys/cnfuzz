package restler

import (
	"github.com/suecodelabs/cnfuzz/src/pkg/logger"
	"os"
	"path/filepath"
)

type RestlerOutput struct {
}

func ParseRestlerOutput(l logger.Logger, outputDir string) RestlerOutput {
	// TODO parse restler output
	files, err := filepath.Glob(filepath.Join(outputDir, "RestlerResults/experiment*/logs/testing_summary.json"))
	if err != nil {
		l.FatalError(err, "failed to find restler output files")
	}
	jsonString, err := os.ReadFile(files[0])
	var _, _ = jsonString, err

	// EngineStdErr.txt
	// EngineStdOut.txt
	// ResponseBuckets/
	// RestlerResults/
	// ResultsAnalyzerStdErr.txt
	// ResultsAnalyzerStdOut.txt
	// restler-yyyymmdd-xxxxxx.log

	return RestlerOutput{}
}
