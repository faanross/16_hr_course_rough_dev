package control

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
)

// Registry of valid commands with their validators and processors
var validCommands = map[string]struct {
	Validator CommandValidator
}{
	"shellcode": {
		Validator: validateShellcodeCommand,
	},
	"whoami": {},
}

// CommandValidator validates command-specific arguments
type CommandValidator func(json.RawMessage) error

// validateShellcodeCommand validates "shellcode" command arguments from client
func validateShellcodeCommand(rawArgs json.RawMessage) error {
	if len(rawArgs) == 0 {
		return fmt.Errorf("shellcode command requires arguments")
	}

	var args ShellcodeArgsClient

	if err := json.Unmarshal(rawArgs, &args); err != nil {
		return fmt.Errorf("invalid argument format: %w", err)
	}

	if args.FilePath == "" {
		return fmt.Errorf("file_path is required")
	}

	if args.ExportName == "" {
		return fmt.Errorf("export_name is required")
	}

	// Check if file exists
	if _, err := os.Stat(args.FilePath); os.IsNotExist(err) {
		return fmt.Errorf("file does not exist: %s", args.FilePath)
	}

	log.Printf("Validation passed: file_path=%s, export_name=%s", args.FilePath, args.ExportName)

	return nil
}
