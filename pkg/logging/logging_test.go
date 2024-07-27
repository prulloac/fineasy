package logging

import (
	"testing"
)

func TestNewLogger(t *testing.T) {
	logger := NewLogger()
	if logger.level != Info {
		t.Errorf("Expected log level to be Info, got %v", logger.level)
	}
}

func TestNewLoggerWithPrefix(t *testing.T) {
	logger := NewLoggerWithPrefix("Test")
	if logger.prefix != "Test" {
		t.Errorf("Expected prefix to be Test, got %v", logger.prefix)
	}
}

func TestLogger_Debugf(t *testing.T) {
	logger := NewLoggerWithLevel(Debug)
	logger.Debugf("Test")
}

func TestLogger_Infof(t *testing.T) {
	logger := NewLoggerWithLevel(Info)
	logger.Infof("Test")
}

func TestLogger_Warningf(t *testing.T) {
	logger := NewLoggerWithLevel(Warning)
	logger.Warningf("Test")
}

func TestLogger_Errorf(t *testing.T) {
	logger := NewLoggerWithLevel(Error)
	logger.Errorf("Test")
}

func TestLogger_Fatalf(t *testing.T) {
	logger := NewLoggerWithLevel(Fatal)
	logger.Fatalf("Test")
}

func TestLogger_Println(t *testing.T) {
	Println("Test")
}

func TestLogger_Errorln(t *testing.T) {
	Errorln("Test")
}

func TestLogger_SetDepth(t *testing.T) {
	logger := NewLogger()
	logger.SetDepth(3)
	if logger.depth != 3 {
		t.Errorf("Expected depth to be 3, got %v", logger.depth)
	}
}

func TestLogger_fmtMessage(t *testing.T) {
	logger := NewLogger()
	message := logger.fmtMessage("Test")
	if message != "[INFO] Test" {
		t.Errorf("Expected message to be Test, got %v", message)
	}
}

func TestLogger_fmtMessageWithPrefix(t *testing.T) {
	logger := NewLoggerWithPrefix("<Test>")
	message := logger.fmtMessage("Test")
	if message != "[INFO] <Test> Test" {
		t.Errorf("Expected message to be Test: Test, got %v", message)
	}
}

func TestLogger_fmtMessageWithArgs(t *testing.T) {
	logger := NewLogger()
	message := logger.fmtMessage("Test %v", "Test")
	if message != "[INFO] Test Test" {
		t.Errorf("Expected message to be Test Test, got %v", message)
	}
}

func TestLogger_fmtMessageWithPrefixAndArgs(t *testing.T) {
	logger := NewLoggerWithPrefix("<Test>")
	message := logger.fmtMessage("Test %v", "Test")
	if message != "[INFO] <Test> Test Test" {
		t.Errorf("Expected message to be Test: Test Test, got %v", message)
	}
}
