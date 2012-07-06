package logorithm

import (
	"bytes"
	"fmt"
	"regexp"
	"testing"
)

const (
	SOFTWARE = "death-ray"
	VERSION  = "1.0.0"
	PROGRAM  = "energy-pump"
	PID      = 1234
	VERBOSE  = false
)

func TestNew(t *testing.T) {
	var buffer bytes.Buffer
	logger := New(&buffer, VERBOSE, SOFTWARE, VERSION, PROGRAM, PID)

	if logger.Flags() != 0 {
		t.Error("Logger flags should be zero")
	}
	if logger.Prefix() != "" {
		t.Error("Logger prefix should be empty")
	}
	if logger.sequenceNumber != 0 {
		t.Error("Initial sequenceNumber must be zero")
	}
	if logger.verbose != VERBOSE {
		t.Errorf("verbose value should be %v", VERBOSE)
	}
	expectedFmtStr := fmt.Sprintf(BASE_FORMAT_STR, SOFTWARE, VERSION, PROGRAM, PID)
	if logger.formatString != expectedFmtStr {
		t.Error("formatString value should be " + expectedFmtStr)
	}
}

func TestLog(t *testing.T) {
	var buffer bytes.Buffer
	logger := New(&buffer, VERBOSE, SOFTWARE, VERSION, PROGRAM, PID)

	oldSeqN := logger.sequenceNumber
	logger.log("TOKEN", "")
	if logger.sequenceNumber != oldSeqN+1 {
		t.Error("sequenceNumber wasn't incremented")
	}
	buffer.Reset()

	logger.log("TOKEN", "MSG")
	match, _ := regexp.MatchString("^TOKEN.*MSG\n$", buffer.String())
	if !match {
		t.Error("Log format mismatch")
	}
	buffer.Reset()

	logger.log("TOKEN", "%s", "MSG")
	match, _ = regexp.MatchString("^TOKEN.*MSG\n$", buffer.String())
	if !match {
		t.Error("Interpolation didn't work")
	}
	buffer.Reset()

	logger.log("TOKEN", "")
	match, _ = regexp.MatchString("x-timestamp=\"\\d{4}-\\d{2}-\\d{2}T\\d{2}:\\d{2}:\\d{2}\\.\\d+\\+\\d{2}:\\d{2}\"", buffer.String())
	if !match {
		t.Errorf("Timestamp format should be time.RFC3339Nano\nGot: %s", buffer.String())
	}
	buffer.Reset()

	logger.log("TOKEN", "")
	match, _ = regexp.MatchString("x-severity=\"TOKEN\"", buffer.String())
	if !match {
		t.Errorf("Severity mismatch\nGot: %s", buffer.String())
	}
	buffer.Reset()

	logger.Emerg("")
	match, _ = regexp.MatchString("^EMERG.*x-severity=\"EMERG\"", buffer.String())
	if !match {
		t.Errorf("Emerg severity should be EMERG\nGot: %s", buffer.String())
	}
	buffer.Reset()

	logger.Alert("")
	match, _ = regexp.MatchString("^ALERT.*x-severity=\"ALERT\"", buffer.String())
	if !match {
		t.Errorf("Alert severity should be ALERT\nGot: %s", buffer.String())
	}
	buffer.Reset()

	logger.Critical("")
	match, _ = regexp.MatchString("^CRIT.*x-severity=\"CRIT\"", buffer.String())
	if !match {
		t.Errorf("Critical severity should be CRIT\nGot: %s", buffer.String())
	}
	buffer.Reset()

	logger.Error("")
	match, _ = regexp.MatchString("^ERR.*x-severity=\"ERR\"", buffer.String())
	if !match {
		t.Errorf("Error severity should be ERR\nGot: %s", buffer.String())
	}
	buffer.Reset()

	logger.Warning("")
	match, _ = regexp.MatchString("^WARN.*x-severity=\"WARN\"", buffer.String())
	if !match {
		t.Errorf("Warning severity should be WARN\nGot: %s", buffer.String())
	}
	buffer.Reset()

	logger.Notice("")
	match, _ = regexp.MatchString("^NOTICE.*x-severity=\"NOTICE\"", buffer.String())
	if !match {
		t.Errorf("Notice severity should be NOTICE\nGot: %s", buffer.String())
	}
	buffer.Reset()

	logger.Info("")
	match, _ = regexp.MatchString("^INFO.*x-severity=\"INFO\"", buffer.String())
	if !match {
		t.Errorf("Info severity should be INFO\nGot: %s", buffer.String())
	}
	buffer.Reset()

	logger.verbose = true
	logger.Debug("")
	match, _ = regexp.MatchString("^DEBUG.*x-severity=\"DEBUG\"", buffer.String())
	if !match {
		t.Errorf("Debug severity should be DEBUG\nGot: %s", buffer.String())
	}
	buffer.Reset()

	logger.verbose = false
	logger.Debug("")
	if buffer.String() != "" {
		t.Error("Debug should not log if verbose flag is false")
	}
	buffer.Reset()
}

// Should be run with GOMAXPROCS greater than 1 like `go test -cpu 4`
func TestThreadSafeSequenceNumber(t *testing.T) {
	// chosen to run quickly with cpu=1 and give enough chance for a race
	const factor = 200
	var buf bytes.Buffer
	join := make(chan struct{})

	logger := New(&buf, VERBOSE, SOFTWARE, VERSION, PROGRAM, PID)

	for i := 0; i < factor; i++ {
		go func() {
			for j := 0; j < factor; j++ {
				logger.Info("go go go speedracer")
			}
			join <- struct{}{}
		}()
	}

	for i := 0; i < factor; i++ {
		<-join
	}

	if logger.sequenceNumber != factor*factor {
		t.Error("sequenceNumber has race condition")
	}
}
