package osutil

import (
	"fmt"
	"os"
	"regexp"
)

// MakeFileCopyTargets copyies Make targets of a Makefile to another Makefile that can have a manually-writen parts until a `# REMARK Do not edit` line.
func MakeFileCopyTargets(sourceMakefilePath string, destinationMakefilePath string, makeTargets []string) (err error) {
	copiedMakefile, err := os.ReadFile(destinationMakefilePath)
	if err != nil {
		return err
	}

	// Only keep the manual written part of the Makefile.
	copiedMakefile = regexp.MustCompile(`(?ms)(\A.+# REMARK Do not edit.+?\n).+\z`).ReplaceAll(copiedMakefile, []byte("$1"))

	// Copy the Make targets from the original Makefile that should be included in the copied Makefile.
	originalMakefile, err := os.ReadFile(sourceMakefilePath)
	if err != nil {
		return err
	}
	for _, targetName := range makeTargets {
		matches := regexp.MustCompile(`(?ms)(` + targetName + `: .+?\.PHONY:.+?\n)`).FindSubmatch(originalMakefile)
		if matches == nil {
			return fmt.Errorf("could not find Make target %q", targetName)
		}

		copiedMakefile = append(copiedMakefile, []byte("\n")...)
		copiedMakefile = append(copiedMakefile, matches[1]...)
	}

	if err := os.WriteFile(destinationMakefilePath, copiedMakefile, 0644); err != nil {
		return err
	}

	return nil
}
