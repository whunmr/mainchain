package keeper_test

import (
	gocontext "context"
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"
	"github.com/unification-com/mainchain/app/test_helpers"
	"github.com/unification-com/mainchain/x/enterprise/types"
	"time"
)

func (suite *KeeperTestSuite) TestGRPCQueryParams() {
	app, ctx, queryClient, addrs := suite.app, suite.ctx, suite.queryClient, suite.addrs

	testParams := types.Params{
		EntSigners:        addrs[0].String(),
		Denom:             test_helpers.TestDenomination,
		MinAccepts:        1,
		DecisionTimeLimit: 600,
	}

	app.EnterpriseKeeper.SetParams(ctx, testParams)
	paramsResp, err := queryClient.Params(gocontext.Background(), &types.QueryParamsRequest{})

	suite.NoError(err)
	suite.Equal(testParams, paramsResp.Params)
}

func (suite *KeeperTestSuite) TestGRPCQueryEnterpriseUndPurchaseOrder() {
	app, ctx, queryClient, addrs := suite.app, suite.ctx, suite.queryClient, suite.addrs

	var (
		req   *types.QueryEnterpriseUndPurchaseOrderRequest
		expPo types.EnterpriseUndPurchaseOrder
	)

	testCases := []struct {
		msg      string
		malleate func()
		expPass  bool
	}{
		{
			"empty request",
			func() {
				req = &types.QueryEnterpriseUndPurchaseOrderRequest{}
				expPo = types.EnterpriseUndPurchaseOrder{}
			},
			false,
		},
		{
			"po id must not be 0",
			func() {
				req = &types.QueryEnterpriseUndPurchaseOrderRequest{
					PurchaseOrderId: 0,
				}
				expPo = types.EnterpriseUndPurchaseOrder{}
			},
			false,
		},
		{
			"po does not exist",
			func() {
				req = &types.QueryEnterpriseUndPurchaseOrderRequest{
					PurchaseOrderId: 99,
				}
				expPo = types.EnterpriseUndPurchaseOrder{}
			},
			false,
		},
		{
			"valid request",
			func() {

				poId := uint64(1)
				req = &types.QueryEnterpriseUndPurchaseOrderRequest{PurchaseOrderId: poId}

				expectedPo := types.EnterpriseUndPurchaseOrder{
					Id:             poId,
					Purchaser:      addrs[0].String(),
					Amount:         sdk.NewInt64Coin(test_helpers.TestDenomination, 100),
					Status:         types.StatusRaised,
					RaiseTime:      uint64(time.Now().Unix()),
					CompletionTime: 0,
					Decisions:      nil,
				}

				err := app.EnterpriseKeeper.SetPurchaseOrder(ctx, expectedPo)
				suite.Require().Nil(err)
				dbPo, found := app.EnterpriseKeeper.GetPurchaseOrder(ctx, poId)
				suite.Require().True(found)

				expPo = dbPo
			},
			true,
		},
	}

	for _, testCase := range testCases {
		suite.Run(fmt.Sprintf("Case %s", testCase.msg), func() {
			testCase.malleate()

			poRes, err := queryClient.EnterpriseUndPurchaseOrder(gocontext.Background(), req)

			if testCase.expPass {
				suite.Require().NoError(err)
				suite.Require().Equal(expPo, poRes.PurchaseOrder)
			} else {
				suite.Require().Error(err)
				suite.Require().Nil(poRes)
			}
		})
	}
}

func (suite *KeeperTestSuite) TestGRPCQueryEnterpriseUndPurchaseOrders() {
	app, ctx, queryClient, addrs := suite.app, suite.ctx, suite.queryClient, suite.addrs

	testPos := []types.EnterpriseUndPurchaseOrder{}

	var (
		req    *types.QueryEnterpriseUndPurchaseOrdersRequest
		expRes *types.QueryEnterpriseUndPurchaseOrdersResponse
	)

	testCases := []struct {
		msg      string
		malleate func()
		expPass  bool
	}{
		{
			"empty request",
			func() {
				req = &types.QueryEnterpriseUndPurchaseOrdersRequest{}
			},
			true,
		},
		{
			"request pos with limit 3",
			func() {
				// create 5 test pos
				for i := 0; i < 5; i++ {
					newPo := types.EnterpriseUndPurchaseOrder{
						Purchaser: addrs[0].String(),
						Amount:    sdk.NewInt64Coin(test_helpers.TestDenomination, int64(i)+1),
					}

					poId, err := app.EnterpriseKeeper.RaiseNewPurchaseOrder(ctx, newPo)
					suite.Require().NoError(err)
					expectedPo, found := app.EnterpriseKeeper.GetPurchaseOrder(ctx, poId)
					suite.Require().True(found)
					testPos = append(testPos, expectedPo)
				}

				req = &types.QueryEnterpriseUndPurchaseOrdersRequest{
					Pagination: &query.PageRequest{Limit: 3},
				}

				expRes = &types.QueryEnterpriseUndPurchaseOrdersResponse{
					PurchaseOrders: testPos[:3],
				}
			},
			true,
		},
		{
			"request 2nd page with limit 4",
			func() {
				req = &types.QueryEnterpriseUndPurchaseOrdersRequest{
					Pagination: &query.PageRequest{Offset: 3, Limit: 3},
				}

				expRes = &types.QueryEnterpriseUndPurchaseOrdersResponse{
					PurchaseOrders: testPos[3:],
				}
			},
			true,
		},
		{
			"request with limit 2 and count true",
			func() {
				req = &types.QueryEnterpriseUndPurchaseOrdersRequest{
					Pagination: &query.PageRequest{Limit: 2, CountTotal: true},
				}

				expRes = &types.QueryEnterpriseUndPurchaseOrdersResponse{
					PurchaseOrders: testPos[:2],
				}
			},
			true,
		},
	}

	for _, testCase := range testCases {
		suite.Run(fmt.Sprintf("Case %s", testCase.msg), func() {
			testCase.malleate()

			pos, err := queryClient.EnterpriseUndPurchaseOrders(gocontext.Background(), req)

			if testCase.expPass {
				suite.Require().NoError(err)

				suite.Require().Len(pos.GetPurchaseOrders(), len(expRes.GetPurchaseOrders()))
				for i := 0; i < len(pos.GetPurchaseOrders()); i++ {
					suite.Require().Equal(pos.GetPurchaseOrders()[i].String(), expRes.GetPurchaseOrders()[i].String())
				}
			} else {
				suite.Require().Error(err)
				suite.Require().Nil(pos)
			}
		})
	}
}

func (suite *KeeperTestSuite) TestGRPCQueryLockedUndByAddress() {
	app, ctx, queryClient, addrs := suite.app, suite.ctx, suite.queryClient, suite.addrs

	var (
		req    *types.QueryLockedUndByAddressRequest
		expRes types.LockedUnd
	)

	testCases := []struct {
		msg      string
		malleate func()
		expPass  bool
	}{
		{
			"empty request",
			func() {
				req = &types.QueryLockedUndByAddressRequest{}
			},
			false,
		},
		{
			"invalid address",
			func() {
				req = &types.QueryLockedUndByAddressRequest{
					Owner: "rubbish",
				}
			},
			false,
		},
		{
			"valid request",
			func() {
				req = &types.QueryLockedUndByAddressRequest{
					Owner: addrs[0].String(),
				}

				l := types.LockedUnd{
					Owner:  addrs[0].String(),
					Amount: sdk.NewInt64Coin(test_helpers.TestDenomination, 1000),
				}

				err := app.EnterpriseKeeper.SetLockedUndForAccount(ctx, l)
				suite.Require().NoError(err)
				expRes = l
			},
			true,
		},
	}

	for _, testCase := range testCases {
		suite.Run(fmt.Sprintf("Case %s", testCase.msg), func() {
			testCase.malleate()

			lRes, err := queryClient.LockedUndByAddress(gocontext.Background(), req)

			if testCase.expPass {
				suite.Require().NoError(err)
				suite.Require().Equal(&expRes, lRes.LockedUnd)
			} else {
				suite.Require().Error(err)
				suite.Require().Nil(lRes)
			}
		})
	}
}

func (suite *KeeperTestSuite) TestGRPCQueryTotalLocked() {
	app, ctx, queryClient := suite.app, suite.ctx, suite.queryClient

	req := &types.QueryTotalLockedRequest{}
	locked := sdk.NewInt64Coin(test_helpers.TestDenomination, 1000)
	expectedRes := &types.QueryTotalLockedResponse{
		Amount: locked,
	}

	err := app.EnterpriseKeeper.SetTotalLockedUnd(ctx, locked)
	suite.Require().NoError(err)

	lRes, err := queryClient.TotalLocked(gocontext.Background(), req)

	suite.Require().NoError(err)
	suite.Require().Equal(expectedRes, lRes)
}

func (suite *KeeperTestSuite) TestGRPCQueryTotalUnlocked() {
	app, ctx, queryClient, addrs := suite.app, suite.ctx, suite.queryClient, suite.addrs

	toLock := sdk.NewInt64Coin(test_helpers.TestDenomination, 1000)
	toUnock := sdk.NewInt64Coin(test_helpers.TestDenomination, 100)
	err := app.EnterpriseKeeper.MintCoinsAndLock(ctx, addrs[0], toLock)
	suite.Require().NoError(err)

	err = app.EnterpriseKeeper.UnlockCoinsForFees(ctx, addrs[0], sdk.Coins{toUnock})
	suite.Require().NoError(err)

	req := &types.QueryTotalUnlockedRequest{}

	expectedUnlocked := app.EnterpriseKeeper.GetTotalUnLockedUnd(ctx)

	expectedRes := &types.QueryTotalUnlockedResponse{
		Amount: expectedUnlocked,
	}

	lRes, err := queryClient.TotalUnlocked(gocontext.Background(), req)

	suite.Require().NoError(err)
	suite.Require().Equal(expectedRes, lRes)
}

func (suite *KeeperTestSuite) TestGRPCQueryEnterpriseSupply() {
	app, ctx, queryClient, addrs := suite.app, suite.ctx, suite.queryClient, suite.addrs

	toLock := sdk.NewInt64Coin(test_helpers.TestDenomination, 1000)
	toUnlock := sdk.NewInt64Coin(test_helpers.TestDenomination, 100)

	supply := app.BankKeeper.GetSupply(ctx).GetTotal().AmountOf(test_helpers.TestDenomination)
	baseSupply := sdk.NewCoin(test_helpers.TestDenomination, supply)
	locked := toLock.Sub(toUnlock)
	unlocked := baseSupply.Add(toUnlock)
	total := baseSupply.Add(toLock)

	expectedTotalSupply := types.UndSupply{
		Denom:  test_helpers.TestDenomination,
		Locked: locked.Amount.Uint64(),
		Amount: unlocked.Amount.Uint64(),
		Total:  total.Amount.Uint64(),
	}

	err := app.EnterpriseKeeper.MintCoinsAndLock(ctx, addrs[0], toLock)
	suite.Require().NoError(err)

	err = app.EnterpriseKeeper.UnlockCoinsForFees(ctx, addrs[0], sdk.Coins{toUnlock})
	suite.Require().NoError(err)

	req := &types.QueryEnterpriseSupplyRequest{}

	expectedRes := &types.QueryEnterpriseSupplyResponse{
		Supply: expectedTotalSupply,
	}

	lRes, err := queryClient.EnterpriseSupply(gocontext.Background(), req)

	suite.Require().NoError(err)
	suite.Require().Equal(expectedRes, lRes)
}

func (suite *KeeperTestSuite) TestGRPCQueryTotalSupply() {
	app, ctx, queryClient, addrs := suite.app, suite.ctx, suite.queryClient, suite.addrs

	toLock := sdk.NewInt64Coin(test_helpers.TestDenomination, 1000)
	toUnlock := sdk.NewInt64Coin(test_helpers.TestDenomination, 100)

	supply := app.BankKeeper.GetSupply(ctx).GetTotal().AmountOf(test_helpers.TestDenomination)
	baseSupply := sdk.NewCoin(test_helpers.TestDenomination, supply)
	expectedTotalSupply := baseSupply.Add(toUnlock)

	expectedResponse := &types.QueryTotalSupplyResponse{
		Supply: sdk.NewCoins(
			expectedTotalSupply,
		),
	}

	err := app.EnterpriseKeeper.MintCoinsAndLock(ctx, addrs[0], toLock)
	suite.Require().NoError(err)

	err = app.EnterpriseKeeper.UnlockCoinsForFees(ctx, addrs[0], sdk.Coins{toUnlock})
	suite.Require().NoError(err)

	req := &types.QueryTotalSupplyRequest{}

	lRes, err := queryClient.TotalSupply(gocontext.Background(), req)

	suite.Require().NoError(err)
	suite.Require().Equal(expectedResponse, lRes)
}

func (suite *KeeperTestSuite) TestGRPCQuerySupplyOf() {
	app, ctx, queryClient, addrs := suite.app, suite.ctx, suite.queryClient, suite.addrs

	toLock := sdk.NewInt64Coin(test_helpers.TestDenomination, 1000)
	toUnlock := sdk.NewInt64Coin(test_helpers.TestDenomination, 100)

	supply := app.BankKeeper.GetSupply(ctx).GetTotal().AmountOf(test_helpers.TestDenomination)
	baseSupply := sdk.NewCoin(test_helpers.TestDenomination, supply)
	expectedTotalSupply := baseSupply.Add(toUnlock)

	expectedResponse := &types.QuerySupplyOfResponse{
		Amount: expectedTotalSupply,
	}

	err := app.EnterpriseKeeper.MintCoinsAndLock(ctx, addrs[0], toLock)
	suite.Require().NoError(err)

	err = app.EnterpriseKeeper.UnlockCoinsForFees(ctx, addrs[0], sdk.Coins{toUnlock})
	suite.Require().NoError(err)

	req := &types.QuerySupplyOfRequest{Denom: test_helpers.TestDenomination}

	lRes, err := queryClient.SupplyOf(gocontext.Background(), req)

	suite.Require().NoError(err)
	suite.Require().Equal(expectedResponse, lRes)
}

func (suite *KeeperTestSuite) TestGRPCQueryWhitelist() {
	app, ctx, queryClient, addrs := suite.app, suite.ctx, suite.queryClient, suite.addrs

	var whitelistedAddrs []string

	for _, addr := range addrs {
		err := app.EnterpriseKeeper.AddAddressToWhitelist(ctx, addr)
		suite.Require().NoError(err)
		whitelistedAddrs = append(whitelistedAddrs, addr.String())
	}

	req := &types.QueryWhitelistRequest{}
	res, err := queryClient.Whitelist(gocontext.Background(), req)
	suite.Require().NoError(err)
	suite.Require().Equal(whitelistedAddrs, res.Addresses)
}

func (suite *KeeperTestSuite) TestGRPCQueryWhitelisted() {
	app, ctx, queryClient, addrs := suite.app, suite.ctx, suite.queryClient, suite.addrs
	notListed := test_helpers.GenerateRandomTestAccounts(10)

	for _, addr := range addrs {
		err := app.EnterpriseKeeper.AddAddressToWhitelist(ctx, addr)
		suite.Require().NoError(err)
	}

	for _, addr := range addrs {
		req := &types.QueryWhitelistedRequest{Address: addr.String()}
		res, err := queryClient.Whitelisted(gocontext.Background(), req)
		suite.Require().NoError(err)
		suite.Require().True(res.Whitelisted)
		suite.Require().Equal(addr.String(), res.Address)
	}

	for _, addr := range notListed {
		req := &types.QueryWhitelistedRequest{Address: addr.String()}
		res, err := queryClient.Whitelisted(gocontext.Background(), req)
		suite.Require().NoError(err)
		suite.Require().False(res.Whitelisted)
		suite.Require().Equal(addr.String(), res.Address)
	}
}
