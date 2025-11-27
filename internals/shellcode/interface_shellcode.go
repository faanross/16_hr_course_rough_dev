package shellcode

import (
	"github.com/faanross/16_hr_course_rough_dev/internals/models"
)

// CommandShellcode is the interface for shellcode execution
type CommandShellcode interface {
	DoShellcode(dllBytes []byte, exportName string) (models.ShellcodeResult, error)
}
