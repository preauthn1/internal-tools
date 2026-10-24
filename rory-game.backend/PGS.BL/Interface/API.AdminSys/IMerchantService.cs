using System;
using System.Collections.Generic;
using target.Models.CoreLibs;
using target.Models.Paging;
using target.Models.Table;
using target.Models.ViewModels;

namespace target.BL.Interface
{

    public interface IMerchantService
    {

        PagedList<MerchantViewModel> GetByPagingFilter(QueryMerchantApiRequestParameterModel filter);

        MerchantViewModel GetByID(int merchantID);

        MerchantCopyInfoPreviewViewModel GetCopyPreviewInfo(int merchantID);

        void Create(MerchantCreateRequestModel model);

        void Update(MerchantUpdateRequestModel model);

        void Delete(int targetID);

    }
}
