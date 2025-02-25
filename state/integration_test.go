package state

import (
	"context"
	"encoding/json"
	"os"
	"testing"

	abci "github.com/cometbft/cometbft/abci/types"
	rpcclient "github.com/cometbft/cometbft/rpc/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	sdk "github.com/cosmos/cosmos-sdk/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"google.golang.org/grpc"

	libhead "github.com/celestiaorg/go-header"
	"github.com/sunriselayer/sunrise/test/util/genesis"
	"github.com/sunriselayer/sunrise/test/util/testfactory"
	"github.com/sunriselayer/sunrise/test/util/testnode"
	blobtypes "github.com/sunriselayer/sunrise/x/blob/types"

	"github.com/sunriselayer/sunrise-da/core"
	"github.com/sunriselayer/sunrise-da/header"
)

func TestIntegrationTestSuite(t *testing.T) {
	suite.Run(t, new(IntegrationTestSuite))
}

type IntegrationTestSuite struct {
	suite.Suite

	cleanups []func() error
	accounts []genesis.Account
	cctx     testnode.Context

	accessor *CoreAccessor
}

func (s *IntegrationTestSuite) SetupSuite() {
	if testing.Short() {
		s.T().Skip("skipping test in unit-tests")
	}
	s.T().Log("setting up integration test suite")

	cfg := core.DefaultTestConfig()
	s.cctx = core.StartTestNodeWithConfig(s.T(), cfg)
	s.accounts = cfg.Genesis.Accounts()

	signer := blobtypes.NewKeyringSigner(s.cctx.Keyring, s.accounts[0].Name, s.cctx.ChainID)
	accessor := NewCoreAccessor(signer, localHeader{s.cctx.RpcClient}, "", "", "")
	setClients(accessor, s.cctx.GRPCClient, s.cctx.Client)
	s.accessor = accessor

	// required to ensure the Head request is non-nil
	_, err := s.cctx.WaitForHeight(3)
	require.NoError(s.T(), err)
}

func setClients(ca *CoreAccessor, conn *grpc.ClientConn, abciCli rpcclient.ABCIClient) {
	ca.coreConn = conn
	// create the query client
	queryCli := banktypes.NewQueryClient(ca.coreConn)
	ca.queryCli = queryCli
	// create the staking query client
	stakingCli := stakingtypes.NewQueryClient(ca.coreConn)
	ca.stakingCli = stakingCli

	ca.rpcCli = abciCli
}

func (s *IntegrationTestSuite) TearDownSuite() {
	s.T().Log("tearing down integration test suite")
	require := s.Require()
	require.NoError(s.accessor.Stop(s.cctx.GoContext()))
	for _, c := range s.cleanups {
		err := c()
		require.NoError(err)
	}
}

func (s *IntegrationTestSuite) getAddress(acc string) sdk.Address {
	rec, err := s.cctx.Keyring.Key(acc)
	require.NoError(s.T(), err)

	addr, err := rec.GetAddress()
	require.NoError(s.T(), err)

	return addr
}

type localHeader struct {
	client rpcclient.Client
}

func (l localHeader) Head(
	ctx context.Context,
	_ ...libhead.HeadOption[*header.ExtendedHeader],
) (*header.ExtendedHeader, error) {
	latest, err := l.client.Block(ctx, nil)
	if err != nil {
		return nil, err
	}
	h := &header.ExtendedHeader{
		RawHeader: latest.Block.Header,
	}
	return h, nil
}

func (s *IntegrationTestSuite) TestGetBalance() {
	require := s.Require()
	for _, acc := range s.accounts {
		bal, err := s.accessor.BalanceForAddress(context.Background(), Address{s.getAddress(acc.Name)})
		require.NoError(err)
		require.True(bal.IsPositive(), acc.Name)
	}
}

// This test can be used to generate a json encoded block for other test data,
// such as that in share/availability/light/testdata
func (s *IntegrationTestSuite) TestGenerateJSONBlock() {
	t := s.T()
	t.Skip("skipping testdata generation test")
	resp, err := s.cctx.FillBlock(4, s.accounts[0].Name, flags.BroadcastSync)
	require := s.Require()
	require.NoError(err)
	require.Equal(abci.CodeTypeOK, resp.Code)
	require.NoError(s.cctx.WaitForNextBlock())

	// download the block that the tx was in
	res, err := testfactory.QueryWithoutProof(s.cctx.Context, resp.TxHash)
	require.NoError(err)

	block, err := s.cctx.Client.Block(s.cctx.GoContext(), &res.Height)
	require.NoError(err)

	pBlock, err := block.Block.ToProto()
	require.NoError(err)

	file, err := os.OpenFile("sample-block.json", os.O_CREATE|os.O_RDWR, os.ModePerm)
	defer file.Close() //nolint: staticcheck
	require.NoError(err)

	err = json.NewEncoder(file).Encode(pBlock)
	require.NoError(err)
}
