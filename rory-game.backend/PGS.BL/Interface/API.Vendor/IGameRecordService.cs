using System;
using System.Collections.Generic;
using target.Models.CoreLibs;
using target.Models.Paging;
using target.Models.Table;
using target.Models.ViewModels;

namespace target.BL.Interface
{

    public interface IGameRecordService
    {

        PagedList<GameRecordViewModel> GetBetRecords(GameRecordRequestModel model);

        BetRecordDetailInfoLinkGenerateViewModel GenerateDetailInfoLink(BetRecordDetailInfoLinkGenerateRequestModel model);

        SingleWalletGameRecordViewModel GetBetRecordBySingleWalletTransactionNo(SingleWalletGameRecordRequestModel model);

        MerchantReportSummaryViewModel GetReportSummary(VendorReportSummaryRequestModel model);

    }
}
