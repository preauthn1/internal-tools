using game.Entity.Game;
using Microsoft.VisualStudio.TestTools.UnitTesting;
using NSubstitute;
using target.BL.Interface;
using target.BL.Services;
using target.DA.Repository;
using target.Models.CoreLibs;
using target.Models.Enums;
using target.Models.Models;
using target.Models.ViewModels;
using target.Utility;
using target.Utility.Extensions;
using target.Utility.Helper;
using target.Utility.IOC;
using target.Utility.Utility;
using System;
using System.Collections.Generic;
using System.Linq;

namespace target.Testing
{
    [TestClass]
    public class SingleWalletPlatformApiRequestRefServiceTests : BaseTest
    {
        public SingleWalletMerchantSessionTokenRepository SingleWalletMerchantSessionTokenRepository { get; set; }

        public ISingleWalletPlatformApiRequestRefService ServiceInstance { get; set; }


        public class MockNotFoundRecordPlatformApiRequestRefService : SingleWalletPlatformApiRequestRefService            
        {
            protected override ApiRequestResultModel<T> RequestApiResult<T>(RequestParameters para, object body = null)
            {
                dynamic res = new ApiRequestResultModel<BaseSingleWalletModeResponseModel<SingleWalletModeCancelApiResponseModel>>
                {
                    ResponseText = "{ 'code' : '404016'  ,  'message' : '查無注單' }",
                    ResultObject = new BaseSingleWalletModeResponseModel<SingleWalletModeCancelApiResponseModel>
                    {
                        Code = ApiErrorCode.NotFoundBetRecord.ToInteger(),
                        Data = null
                    }
                };

                return res;
            }
        }

        public class MockBalanceNotEnoughPlatformApiRequestRefService : SingleWalletPlatformApiRequestRefService
        {
            protected override ApiRequestResultModel<T> RequestApiResult<T>(RequestParameters para, object body = null)
            {
                dynamic res = new ApiRequestResultModel<BaseSingleWalletModeResponseModel<SingleWalletModePlaceBetApiResponseModel>>
                {
                    ResponseText = "{ 'code' : '403652' ,  'message' : '餘額不足' }",
                    ResultObject = new BaseSingleWalletModeResponseModel<SingleWalletModePlaceBetApiResponseModel>
                    {
                        Code = ApiErrorCode.BetFail_WalletAmountNotEnough.ToInteger(),
                        Data = null
                    }
                };

                return res;
            }
        }


        public class MockTimeoutErrorPlatformApiRequestRefService : SingleWalletPlatformApiRequestRefService
        {
            protected override ApiRequestResultModel<T> RequestApiResult<T>(RequestParameters para, object body = null)
            {
                dynamic res = new ApiRequestResultModel<BaseSingleWalletModeResponseModel<SingleWalletModeCancelApiResponseModel>>
                {
                    Exception = new Exception("timeout"),
                    IsTimeoutException = true,
                    ResponseText = "Error!",
                };

                return res;
            }
        }


        [TestMethod]
        public void SingleWallet_PlaceBet_NotEnoughBalance_Test()
        {
            var mockService = new MockBalanceNotEnoughPlatformApiRequestRefService();

            try
            {
                var result = mockService.PlaceBet(new SingleWalletModePlaceBetParameter
                {
                    SessionToken = [REDACTED],
                    BetRecordNoPGS = _TEST_DATA.SingleWallet.BetRecordNo,
                    BetAmount = 10,
                    GameID = _TEST_DATA.Default.GameID,
                    PlayerID = _TEST_DATA.Default.PlayerID,
                });
            }
            catch(Exception ex)
            {
                AssertExpectException(ex, new List<ApiErrorCode>
                {
                    ApiErrorCode.BetFail_WalletAmountNotEnough
                });
            }            
           
            return;
        }


        [TestMethod]
        public void SingleWallet_Cancel_NotFoundBetRecord_Test()
        {
            var mockService = new MockNotFoundRecordPlatformApiRequestRefService();

            var result = mockService.Cancel(new SingleWalletModeCancelApiRequestModel
            {
                GameID = _TEST_DATA.Default.GameID,
                MerchantID = _TEST_DATA.SingleWallet.MerchantID,
                PlayerID = _TEST_DATA.SingleWallet.PlayerID,
                TransactionNo = RandomHelper.GetOrderNo("UNI")
            });

            Assert.IsNotNull(result);

            return;
        }


        [TestMethod]
        public void SingleWallet_Cancel_Timeout_Test()
        {
            var mockService = new MockTimeoutErrorPlatformApiRequestRefService();

            try
            {
                var result = mockService.Cancel(new SingleWalletModeCancelApiRequestModel
                {
                    GameID = _TEST_DATA.Default.GameID,
                    MerchantID = _TEST_DATA.SingleWallet.MerchantID,
                    PlayerID = _TEST_DATA.SingleWallet.PlayerID,
                    TransactionNo = RandomHelper.GetOrderNo("UNI")
                });

            }
            catch(Exception ex)
            {
                AssertExpectException(ex, new List<ApiErrorCode>
                {
                    ApiErrorCode.ThirdPartyApiFail_Timeout
                });
            }

            return;
        }


    }
}

