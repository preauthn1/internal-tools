using System;
using System.Collections.Generic;
using target.Models.CoreLibs;
using target.Models.Paging;
using target.Models.Table;
using target.Models.ViewModels;
using target.Utility.Extensions;
using target.BL.Interface;
using target.Models.Enums;
using target.Models.Models;
using System.Linq;
using target.DA.Repository;
using target.Utility.Helper;
using target.Utility;

namespace target.BL.Services
{

    public class ThirdPartyGameLogicService : BaseService, IThirdPartyGameLogicService
    {

        public ThirdPartySlot88MD5RecordRepository ThirdPartySlot88MD5RecordRepository { get; set; }
     

        public Slot88MD5InfoResultModel GetPlayerInformationByMD5(QuerySlot88MD5RequestModel model)
        {
            model.RuleFor(x => x.MD5).Rule_NotEmptyText();

            var queryMD5 = ThirdPartySlot88MD5RecordRepository.GetByFilter(new QueryThirdPartySlot88MD5RecordRequestViewModel
                {
                    MD5Text = model.MD5
                })
                .FirstOrDefault();

            if (queryMD5 == null)
                throw new TerproException(ApiErrorCode.NotFoundData);

            var parsingJson = queryMD5.OriginalJson.JsonConvertToModel<Slot88MD5InfoResultModel>();

            if (parsingJson == null)
                throw new TerproException(ApiErrorCode.InvalidParameter_ParsingError);

            DateTime now = DateTimeHelper.GetCurrentTime();

            if(now >= parsingJson.ExpiredTime)
                throw new TerproException(ApiErrorCode.InvalidParameter_SignKeyExpired);

            return parsingJson;
        }
    }
}
