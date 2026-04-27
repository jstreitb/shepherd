// Package executor provides secure command execution helpers.
package executor

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
	"sync"
	"time"

	"github.com/jstreitb/baa/internal/pkgmanager"
	"github.com/jstreitb/baa/internal/sanitize"
)

// ZeroBytes overwrites every element of b with zero.
// Call this immediately after the password is no longer needed.
func ZeroBytes(b []byte) {
	for i := range b {
		b[i] = 0
	}
}

// RunCommand executes a single command, optionally under sudo -S.
// Live output lines are sent to outputCh. Both stdout and stderr are
// merged and streamed. The function blocks until the command exits.
func RunCommand(password []byte, needsSudo bool, args []string, env []string, outputCh chan<- string) (string, error) {
	var cmd *exec.Cmd
	if needsSudo {
		// Build: sudo -S -- env VAR=val ... cmd args...
		// We use the `env` command to inject environment variables because
		// sudo resets the environment for the child process.
		sudoArgs := make([]string, 0, 4+len(env)+len(args))
		sudoArgs = append(sudoArgs, "-S", "--")
		if len(env) > 0 {
			sudoArgs = append(sudoArgs, "env")
			sudoArgs = append(sudoArgs, env...)
		}
		sudoArgs = append(sudoArgs, args...)
		cmd = exec.Command("sudo", sudoArgs...)
	} else {
		cmd = exec.Command(args[0], args[1:]...)
		// For non-sudo commands, overlay env directly.
		if len(env) > 0 {
			cmd.Env = append(os.Environ(), env...)
		}
	}

	// Set up stdin pipe for password delivery.
	stdinPipe, err := cmd.StdinPipe()
	if err != nil {
		return "", fmt.Errorf("stdin pipe: %w", err)
	}

	// Set up stdout and stderr pipes for live streaming.
	stdoutPipe, err := cmd.StdoutPipe()
	if err != nil {
		return "", fmt.Errorf("stdout pipe: %w", err)
	}
	stderrPipe, err := cmd.StderrPipe()
	if err != nil {
		return "", fmt.Errorf("stderr pipe: %w", err)
	}

	if err := cmd.Start(); err != nil {
		return "", fmt.Errorf("start: %w", err)
	}

	// Deliver the password to sudo -S via stdin, then close.
	if needsSudo && len(password) > 0 {
		_, _ = stdinPipe.Write(password)
		_, _ = stdinPipe.Write([]byte("\n"))
	}
	_ = stdinPipe.Close()

	// Merge stdout and stderr into a single ordered stream.
	lineCh := make(chan string, 64)
	var wg sync.WaitGroup
	wg.Add(2)

	readPipe := func(r io.Reader) {
		defer wg.Done()
		scanner := bufio.NewScanner(r)
		scanner.Buffer(make([]byte, 0, 64*1024), 1024*1024)
		for scanner.Scan() {
			lineCh <- scanner.Text()
		}
	}
	go readPipe(stdoutPipe)
	go readPipe(stderrPipe)

	go func() {
		wg.Wait()
		close(lineCh)
	}()

	var output strings.Builder
	for line := range lineCh {
		output.WriteString(line)
		output.WriteByte('\n')
		// Non-blocking send; drop lines if the TUI is slow.
		select {
		case outputCh <- line:
		default:
		}
	}

	if err := cmd.Wait(); err != nil {
		return output.String(), err
	}
	return output.String(), nil
}

// RunManagerUpdate sequentially executes every command for the given
// PackageManager, streaming output to outputCh. It closes outputCh
// when finished. The returned UpdateResult summarises the outcome.
func RunManagerUpdate(password []byte, mgr pkgmanager.PackageManager, outputCh chan<- string) pkgmanager.UpdateResult {
	defer close(outputCh)
	start := time.Now()
	var allOutput strings.Builder

	for _, cmdArgs := range mgr.Commands() {
		outputCh <- fmt.Sprintf("▶ Running: %s", strings.Join(cmdArgs, " "))
		out, err := RunCommand(password, mgr.NeedsSudo(), cmdArgs, mgr.Env(), outputCh)
		allOutput.WriteString(out)
		if err != nil {
			errStr := strings.TrimSpace(out)
			if errStr == "" {
				errStr = err.Error()
			}
			return pkgmanager.UpdateResult{
				Manager:  mgr.Name(),
				Success:  false,
				Output:   allOutput.String(),
				Error:    sanitize.SanitizeError(errStr, 500),
				Duration: time.Since(start),
			}
		}
	}

	return pkgmanager.UpdateResult{
		Manager:  mgr.Name(),
		Success:  true,
		Output:   allOutput.String(),
		Duration: time.Since(start),
	}
}
