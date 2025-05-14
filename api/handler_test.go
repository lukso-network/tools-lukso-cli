package api

import (
	"os"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"

	"github.com/lukso-network/tools-lukso-cli/api/errors"
	"github.com/lukso-network/tools-lukso-cli/api/types"
	"github.com/lukso-network/tools-lukso-cli/common"
	"github.com/lukso-network/tools-lukso-cli/common/file"
	"github.com/lukso-network/tools-lukso-cli/common/logger"
	"github.com/lukso-network/tools-lukso-cli/config"
	"github.com/lukso-network/tools-lukso-cli/dep/configs"
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
	loggerMock    logger.Logger // Could test log messages as well - this would mean expecting a lot of individual messages
	installerMock *mocks.MockInstaller
}

func (h *HandlerSuite) SetupSuite() {
	h.cfgMock = mocks.NewMockConfigurator(h.T())
	h.fileMock = mocks.NewMockManager(h.T())
	h.loggerMock = logger.ConsoleLogger{}
	h.installerMock = mocks.NewMockInstaller(h.T())

	configs.CfgFileManager = h.fileMock
	configs.CfgInstaller = h.installerMock

	h.handler = NewHandler(
		h.cfgMock,
		h.fileMock,
		h.loggerMock,
		h.installerMock,
	)
}

func (h *HandlerSuite) SetupTest() {
	// Flush expectations between every test to avoid confusion and expect everything explicitly.
	h.cfgMock.ExpectedCalls = nil
	h.fileMock.ExpectedCalls = nil
	h.installerMock.ExpectedCalls = nil
}

func (h *HandlerSuite) TestInit() {
	testCases := []testCase[types.InitRequest, types.InitResponse]{
		{
			"new init - missing IP in request",
			h.handler.Init,
			func(h *HandlerSuite) {
				h.cfgMock.EXPECT().Exists().Return(false).Once()
				h.fileMock.EXPECT().Exists(mock.Anything).Return(false) // assume any config file is not present
				h.installerMock.EXPECT().InstallFile(mock.Anything, mock.Anything, false).Return(nil)
				h.fileMock.EXPECT().Mkdir(file.SecretsDir, os.FileMode(common.ConfigPerms)).Return(nil).Once()
				h.fileMock.EXPECT().
					Write(
						file.JwtSecretPath,
						mock.MatchedBy(func(b []byte) bool { return len(b) == 64 }), // encoded JWT has twice the amount of bytes
						os.FileMode(common.ConfigPerms),
					).
					Return(nil).
					Once()
				h.fileMock.EXPECT().Mkdir(file.PidDir, os.FileMode(common.ConfigPerms)).Return(nil).Once()
				h.cfgMock.EXPECT().
					Create(
						config.NodeConfig{
							UseClients: config.UseClients{
								ExecutionClient: "",
								ConsensusClient: "",
								ValidatorClient: "",
							},
							Ipv4: "disabled",
						},
					).
					Return(nil).
					Once()
			},
			types.InitRequest{},
			types.InitResponse{},
		},
		{
			"new init - IP disabled explicitly",
			h.handler.Init,
			func(h *HandlerSuite) {
				h.cfgMock.EXPECT().Exists().Return(false).Once()
				h.fileMock.EXPECT().Exists(mock.Anything).Return(false) // assume any config file is not present
				h.installerMock.EXPECT().InstallFile(mock.Anything, mock.Anything, mock.Anything).Return(nil)
				h.fileMock.EXPECT().Mkdir(file.SecretsDir, os.FileMode(common.ConfigPerms)).Return(nil).Once()
				h.fileMock.EXPECT().
					Write(
						file.JwtSecretPath,
						mock.MatchedBy(func(b []byte) bool { return len(b) == 64 }), // encoded JWT has twice the amount of bytes
						os.FileMode(common.ConfigPerms),
					).
					Return(nil).
					Once()
				h.fileMock.EXPECT().Mkdir(file.PidDir, os.FileMode(common.ConfigPerms)).Return(nil).Once()
				h.cfgMock.EXPECT().
					Create(
						config.NodeConfig{
							UseClients: config.UseClients{
								ExecutionClient: "",
								ConsensusClient: "",
								ValidatorClient: "",
							},
							Ipv4: "disabled",
						},
					).
					Return(nil).
					Once()
			},
			types.InitRequest{Ip: "disabled"},
			types.InitResponse{},
		},
		{
			"new init - IP is invalid",
			h.handler.Init,
			func(h *HandlerSuite) {
				h.cfgMock.EXPECT().Exists().Return(false).Once()
				h.fileMock.EXPECT().Exists(mock.Anything).Return(false) // assume any config file is not present
				h.installerMock.EXPECT().InstallFile(mock.Anything, mock.Anything, mock.Anything).Return(nil)
				h.fileMock.EXPECT().Mkdir(file.SecretsDir, os.FileMode(common.ConfigPerms)).Return(nil).Once()
				h.fileMock.EXPECT().
					Write(
						file.JwtSecretPath,
						mock.MatchedBy(func(b []byte) bool { return len(b) == 64 }), // encoded JWT has twice the amount of bytes
						os.FileMode(common.ConfigPerms),
					).
					Return(nil).
					Once()
				h.fileMock.EXPECT().Mkdir(file.PidDir, os.FileMode(common.ConfigPerms)).Return(nil).Once()
				h.cfgMock.EXPECT().
					Create(
						config.NodeConfig{
							UseClients: config.UseClients{
								ExecutionClient: "",
								ConsensusClient: "",
								ValidatorClient: "",
							},
							Ipv4: "disabled",
						},
					).
					Return(nil).
					Once()
			},
			types.InitRequest{Ip: "invalidIp"},
			types.InitResponse{},
		},
		{
			"already initialized - reinit",
			h.handler.Init,
			func(h *HandlerSuite) {
				h.cfgMock.EXPECT().Exists().Return(true).Once()
				h.fileMock.EXPECT().Exists(mock.Anything).Return(false) // assume any config file is not present
				h.installerMock.EXPECT().InstallFile(mock.Anything, mock.Anything, mock.Anything).Return(nil)
				h.fileMock.EXPECT().Mkdir(file.SecretsDir, os.FileMode(common.ConfigPerms)).Return(nil).Once()
				h.fileMock.EXPECT().
					Write(
						file.JwtSecretPath,
						mock.MatchedBy(func(b []byte) bool { return len(b) == 64 }), // encoded JWT has twice the amount of bytes
						os.FileMode(common.ConfigPerms),
					).
					Return(nil).
					Once()
				h.fileMock.EXPECT().Mkdir(file.PidDir, os.FileMode(common.ConfigPerms)).Return(nil).Once()
				h.cfgMock.EXPECT().
					Create(
						config.NodeConfig{
							UseClients: config.UseClients{
								ExecutionClient: "",
								ConsensusClient: "",
								ValidatorClient: "",
							},
							Ipv4: "disabled",
						},
					).
					Return(nil).
					Once()
			},
			types.InitRequest{Ip: "invalidIp", Reinit: true},
			types.InitResponse{},
		},
		{
			"already initialized - no reinit",
			h.handler.Init,
			func(h *HandlerSuite) {
				h.cfgMock.EXPECT().Exists().Return(true).Once()
			},
			types.InitRequest{Reinit: false},
			types.InitResponse{
				Error: errors.ErrCfgExists,
			},
		},
	}

	for _, t := range testCases {
		h.Run(t.name, runCase(h, t))
	}
}

func TestHandlerSuite(t *testing.T) {
	suite.Run(t, new(HandlerSuite))
}
