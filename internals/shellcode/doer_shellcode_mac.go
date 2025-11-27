//go:build darwin

package shellcode

import (
	"fmt"
	"github.com/faanross/16_hr_course_rough_dev/internals/agent"
)

// macShellcode implements the CommandShellcode interface for Darwin/MacOS
type macShellcode struct{}

// New is the constructor for our Mac-specific Shellcode command
func New() CommandShellcode {
	return &macShellcode{}
}

// DoShellcode is the stub implementation for macOS
func (ms *macShellcode) DoShellcode(dllBytes []byte, exportName string) (agent.ShellcodeResult, error) {
	fmt.Println("|❗ SHELLCODE DOER MACOS| This feature has not yet been implemented for MacOS.")

	result := agent.ShellcodeResult{
		Message: "FAILURE: Not implemented on macOS",
	}
	return result, nil
}
