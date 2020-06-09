package grpc_test_test

import (
	"lib/grpc_test"
	"testing"
)

type s struct {
	grpc_test.Tester
}

func Test(t *testing.T) {
	grpc_test.RunSubTests(t, s{})
}

func (s) TestInfo(t *testing.T) {
	grpc_test.Info("Info", "message.")
}

func (s) TestInfoln(t *testing.T) {
	grpc_test.Infoln("Info", "message.")
}

func (s) TestInfof(t *testing.T) {
	grpc_test.Infof("%v %v.", "Info", "message")
}

func (s) TestInfoDepth(t *testing.T) {
	grpc_test.InfoDepth(0, "Info", "depth", "message.")
}

func (s) TestWarning(t *testing.T) {
	grpc_test.Warning("Warning", "message.")
}

func (s) TestWarningln(t *testing.T) {
	grpc_test.Warningln("Warning", "message.")
}

func (s) TestWarningf(t *testing.T) {
	grpc_test.Warningf("%v %v.", "Warning", "message")
}

func (s) TestWarningDepth(t *testing.T) {
	grpc_test.WarningDepth(0, "Warning", "depth", "message.")
}

func (s) TestError(t *testing.T) {
	const numErrors = 10
	grpc_test.TLogger.ExpectError("Expected error")
	grpc_test.TLogger.ExpectError("Expected ln error")
	grpc_test.TLogger.ExpectError("Expected formatted error")
	grpc_test.TLogger.ExpectErrorN("Expected repeated error", numErrors)
	grpc_test.Error("Expected", "error")
	grpc_test.Errorln("Expected", "ln", "error")
	grpc_test.Errorf("%v %v %v", "Expected", "formatted", "error")
	for i := 0; i < numErrors; i++ {
		grpc_test.Error("Expected repeated error")
	}
}
