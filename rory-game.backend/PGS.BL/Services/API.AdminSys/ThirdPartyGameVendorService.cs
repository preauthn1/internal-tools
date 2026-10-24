using target.BL.Interface;
using target.BL.Services;
using target.DA.Repository;
using target.Models.CoreLibs;
using target.Models.Enums;
using target.Models.Models;
using target.Models.Paging;
using target.Models.Table;
using target.Models.ViewModels;
using target.Models.ViewModels.GameVendor;
using target.Utility.Extensions;
using target.Utility.Helper;
using System;
using System.Collections.Generic;
using System.Data;
using System.Data.SqlClient;
using System.Linq;

namespace target.BL.Services
{
    public class ThirdPartyGameVendorService : BaseAdminSysService, IThirdPartyGameVendorService
    {
        public BetRecordTransactionNoMappingRepository BetRecordTransactionNoMappingRepository { get; set; }

        public PerGameDetailLogRepository PerGameDetailLogRepository { get; set; }

        public GameListRepository GameListRepository { get; set; }


        public ThirdPartyGameVendorBetRecordDetailUrlViewModel GetDetailUrlByPegGameDetailLogID(int userID , int perGameDetailID)
        {
            var queryPerGame = PerGameDetailLogRepository.GetById(perGameDetailID);

            if (queryPerGame == null)
                throw new TerproException(ApiErrorCode.NotFoundBetRecord);

            CheckViewMerchantPermission(userID, queryPerGame.MerchantID);

            var queryGame = GameListRepository.GetById(queryPerGame.GameID);

            var vendorService = ServiceInstanceHelper.InstanceVendorService(queryGame.GameVendorType);

            if (queryPerGame.ThirdPartyBetRecordNo.IsNullOrWhiteSpace())
                return new ThirdPartyGameVendorBetRecordDetailUrlViewModel();

            ThirdPartyGameVendorBetRecordDetailUrlViewModel result = vendorService.GetExternalBetDetailInfo(new ExternalVendorBetRecordDetailInfoRequestModel
            {
                BetRecordNo = queryPerGame.ThirdPartyBetRecordNo,
                MerchantID = queryPerGame.MerchantID,
                ThirdPartyGameCode = queryGame.ThirdPartyGameCode
            });

            return result;
        }


        public ThirdPartyGameVendorBetRecordDetailUrlViewModel GetDetailUrlByTransactionNo(int userID , string pgsTransactionNo)
        {
            var query = BetRecordTransactionNoMappingRepository.GetByTransactionNo(pgsTransactionNo);

            if (query == null)
                throw new TerproException(ApiErrorCode.NotFoundBetRecord);

            var vendorService = ServiceInstanceHelper.InstanceVendorService(query.GameVendorID);

            var queryBetRecord = vendorService.GetBetRecordByPGSTransactionNo(pgsTransactionNo);

            CheckViewMerchantPermission(userID, queryBetRecord.MerchantID);

            ThirdPartyGameVendorBetRecordDetailUrlViewModel result = vendorService.GetExternalBetDetailInfo(new ExternalVendorBetRecordDetailInfoRequestModel
            {
                BetRecordNo = queryBetRecord.BetRecordNo,
                MerchantID = queryBetRecord.MerchantID,
                ThirdPartyGameCode = queryBetRecord.ThirdPartyGameCode
            });

            return result;
        }

    }

}
