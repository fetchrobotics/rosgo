package ros

import (
	"bufio"
	"bytes"
	"fmt"
	"log"
	"strings"
	"testing"
)

func TestDefaultLogger(t *testing.T) {
	logger := NewDefaultLogger()
	textOne := "Lorem ipsum dolor sit amet, consectetur adipiscing elit"
	textTwo := "Neque porro quisquam est qui dolorem ipsum quia dolor sit amet."

	// out is the buffer used to store logs
	var out bytes.Buffer

	r := bufio.NewReader(&out) // Reader to read output logs
	log.SetOutput(&out)        // Set Output to store logs in out buffer
	log.SetFlags(0)            // Disable timestamps

	t.Run("Severity", func(t *testing.T) {
		t.Run("LogLevelFatal", func(*testing.T) {
			logger.SetSeverity(LogLevelFatal)
			if logger.Severity() != LogLevelFatal {
				t.Errorf("Log Level Wanted: %v Got: %v", LogLevelFatal, logger.Severity())
			}
		})

		t.Run("LogLevelError", func(*testing.T) {
			logger.SetSeverity(LogLevelError)
			if logger.Severity() != LogLevelError {
				t.Errorf("Log Level Wanted: %v Got: %v", LogLevelError, logger.Severity())
			}

			logger.Error(textOne)
			textWanted := fmt.Sprintf("[ERROR] %v\n", textOne)
			if err := checkLogLine(r, textWanted); err != nil {
				t.Error(err)
			}

			logger.Errorf("%v: %v", textOne, textTwo)
			textWanted = fmt.Sprintf("[ERROR] %v: %v\n", textOne, textTwo)
			if err := checkLogLine(r, textWanted); err != nil {
				t.Error(err)
			}

			logger.Warn(textOne)
			if err := checkLogLine(r, ""); err == nil {
				t.Errorf("No warning logs should be logged when log level is Error")
			}
		})

		t.Run("LogLevelWarn", func(*testing.T) {
			logger.SetSeverity(LogLevelWarn)
			if logger.Severity() != LogLevelWarn {
				t.Errorf("Log Level Wanted: %v Got: %v", LogLevelWarn, logger.Severity())
			}

			logger.Warn(textOne)
			textWanted := fmt.Sprintf("[WARN] %v\n", textOne)
			if err := checkLogLine(r, textWanted); err != nil {
				t.Error(err)
			}

			logger.Warnf("%v: %v", textOne, textTwo)
			textWanted = fmt.Sprintf("[WARN] %v: %v\n", textOne, textTwo)
			if err := checkLogLine(r, textWanted); err != nil {
				t.Error(err)
			}

			logger.Info(textOne)
			if err := checkLogLine(r, ""); err == nil {
				t.Errorf("No info logs should be logged when log level is Warn")
			}
		})

		t.Run("LogLevelInfo", func(*testing.T) {
			logger.SetSeverity(LogLevelInfo)
			if logger.Severity() != LogLevelInfo {
				t.Errorf("Log Level Wanted: %v Got: %v", LogLevelInfo, logger.Severity())
			}

			logger.Info(textOne)
			textWanted := fmt.Sprintf("[INFO] %v\n", textOne)
			if err := checkLogLine(r, textWanted); err != nil {
				t.Error(err)
			}

			logger.Infof("%v: %v", textOne, textTwo)
			textWanted = fmt.Sprintf("[INFO] %v: %v\n", textOne, textTwo)
			if err := checkLogLine(r, textWanted); err != nil {
				t.Error(err)
			}

			logger.Debug(textOne)
			if err := checkLogLine(r, ""); err == nil {
				t.Errorf("No debug logs should be logged when log level is info")
			}
		})

		t.Run("LogLevelDebug", func(*testing.T) {
			logger.SetSeverity(LogLevelDebug)
			if logger.Severity() != LogLevelDebug {
				t.Errorf("Log Level Wanted: %v Got: %v", LogLevelDebug, logger.Severity())
			}

			logger.Debug(textOne)
			textWanted := fmt.Sprintf("[DEBUG] %v\n", textOne)
			if err := checkLogLine(r, textWanted); err != nil {
				t.Error(err)
			}

			logger.Debugf("%v: %v", textOne, textTwo)
			textWanted = fmt.Sprintf("[DEBUG] %v: %v\n", textOne, textTwo)
			if err := checkLogLine(r, textWanted); err != nil {
				t.Error(err)
			}
		})
	})
}

func checkLogLine(r *bufio.Reader, textWanted string) error {
	textGot, err := r.ReadString('\n')
	if err != nil {
		return err
	}

	if strings.Compare(textWanted, textGot) != 0 {
		return fmt.Errorf("Expected \"%v\" Got \"%v\"", textWanted, textGot)
	}
	return nil
}
