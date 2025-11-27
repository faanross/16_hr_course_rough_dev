package shellcode

import "github.com/faanross/16_hr_course_rough_dev/internals/agent"

// CommandShellcode is the interface for shellcode execution
type CommandShellcode interface {
	DoShellcode(dllBytes []byte, exportName string) (agent.ShellcodeResult, error)
}
