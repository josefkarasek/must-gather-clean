package reporting

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/openshift/must-gather-clean/pkg/obfuscator"
	"gopkg.in/yaml.v3"
	"k8s.io/klog/v2"
)

type Report struct {
	Replacements []obfuscator.Replacement `yaml:"replacements,omitempty"`
	Omissions    []string                 `yaml:"omissions,omitempty"`
}

type Reporter interface {
	// WriteReport writes the final report into the given path, will create folders if necessary.
	WriteReport(path string) error

	// CollectOmitterReport collects the omitter's omission results.
	CollectOmitterReport(omitter []string)

	// CollectObfuscatorReport will call the Report method on the obfuscator and collect the individual obfuscation results.
	CollectObfuscatorReport(obfuscatorReport []obfuscator.ReplacementReport)
}

type SimpleReporter struct {
	replacements []obfuscator.Replacement
	omissions    []string
}

var _ Reporter = (*SimpleReporter)(nil)

func (s *SimpleReporter) WriteReport(path string) error {

	reportingFolder := filepath.Dir(path)
	err := os.MkdirAll(reportingFolder, 0700)
	if err != nil {
		return fmt.Errorf("failed to create reporting output folder: %w", err)
	}

	reportFile, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("failed to open report file %s: %w", path, err)
	}
	rEncoder := yaml.NewEncoder(reportFile)
	err = rEncoder.Encode(Report{
		Replacements: s.replacements,
		Omissions:    s.omissions,
	})
	if err != nil {
		return fmt.Errorf("failed to write report at %s: %w", path, err)
	}

	klog.V(2).Infof("successfully saved obfuscation report in %s", path)

	return nil
}

func (s *SimpleReporter) CollectOmitterReport(report []string) {
	s.omissions = append(s.omissions, report...)
}

func (s *SimpleReporter) CollectObfuscatorReport(obfuscatorReport []obfuscator.ReplacementReport) {
	for _, report := range obfuscatorReport {
		s.replacements = append(s.replacements, report.Replacements...)
	}
}

func NewSimpleReporter() *SimpleReporter {
	return &SimpleReporter{
		replacements: []obfuscator.Replacement{},
		omissions:    []string{},
	}
}
