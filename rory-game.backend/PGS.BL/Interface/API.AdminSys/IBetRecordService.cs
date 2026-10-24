using target.Models.CoreLibs;
using target.Models.Paging;
using target.Models.Table;
using target.Models.ViewModels;
using System;
using System.Collections.Generic;

namespace target.BL.Interface
{
    public interface IBetRecordService
    {
        PagedList<UnionBetRecordViewModel> GetUnionBetRecordsByPagingFilter(int userID , QueryUnionBetRecordRequestViewModel filter);

        BetRecordGeneralInfoViewModel GetByTransactionNo(int userID, BetRecordQueryRequestModel model);

        BetRecordNotifyQueueDetailInfoViewModel GetNotifyQueueInfo(int userID, long queueID);

        void ManualRefundBetRecord(int userID, BetRecordManualRefundRequestModel model);

        void ManualSettleBetRecord(int userID, BetRecordManualSettleRequestModel model);
    }
}
