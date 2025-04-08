package api

import (
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/lukso-network/tools-lukso-cli/api/types"
	"github.com/lukso-network/tools-lukso-cli/test/mocks"
)

type testCase[Rq types.Request, Rs types.Response] struct {
	name string
	// expector runs expects on mocks for the current testcase
	handlerFn    HandlerFunc[Rq, Rs]
	expector     func(*HandlerSuite)
	req          Rq
	expectedResp Rs
}

// runCase returns a func used in suite.Suite.Run execution.
func runCase[Rq types.Request, Rs types.Response](h *HandlerSuite, t testCase[Rq, Rs]) func() {
	return func() {
		t.expector(h)

		actual := t.handlerFn(t.req)

		h.Equal(t.expectedResp, actual)
	}
}

type HandlerSuite struct {
	suite.Suite

	handler Handler

	cfgMock       *mocks.MockConfigurator
	fileMock      *mocks.MockManager
	loggerMock    *mocks.MockLogger
	installerMock *mocks.MockInstaller
}

func (h *HandlerSuite) SetupSuite() {
	h.cfgMock = mocks.NewMockConfigurator(h.T())
	h.fileMock = mocks.NewMockManager(h.T())
	h.loggerMock = mocks.NewMockLogger(h.T())
	h.installerMock = mocks.NewMockInstaller(h.T())

	h.handler = NewHandler(
		h.cfgMock,
		h.fileMock,
		h.loggerMock,
		h.installerMock,
	)
}

func (h *HandlerSuite) SetupTest() {
	return
}

func (h *HandlerSuite) TestInit() {
	testCases := []testCase[types.InitRequest, types.InitResponse]{
		{
			"successful init",
			h.handler.Init,
			func(h *HandlerSuite) {
			},
			types.InitRequest{},
			types.InitResponse{},
		},
	}

	for _, t := range testCases {
		h.Run(t.name, runCase(h, t))
	}
}

func TestHandlerSuite(t *testing.T) {
	suite.Run(t, new(HandlerSuite))
}
