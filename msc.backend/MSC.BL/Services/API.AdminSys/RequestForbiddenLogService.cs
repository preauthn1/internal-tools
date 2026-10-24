using System;
using System.Collections.Generic;
using MSC.Models.CoreLibs;
using MSC.Models.Paging;
using MSC.Models.Table;
using MSC.Models.ViewModels;
using MSC.Utility.Extensions;
using MSC.BL.Interface;
using MSC.Models.Enums;
using MSC.Models.Models;
using System.Linq;
using MSC.Utility.Helper;
using MSC.Models;
using MSC.DA.Repository;

namespace MSC.BL.Services
{

    public class RequestForbiddenLogService : BaseAdminService, IRequestForbiddenLogService
    {

        public RequestForbiddenLogRepository RequestForbiddenLogRepository { get; set; }

        public MerchantRepository MerchantRepository { get; set; }


        public PagedList<RequestForbiddenLogViewModel> GetByPagingFilter(int userID, QueryRequestForbiddenLogApiRequestParameterModel filter)
        {
            filter.Types = filter.Type.ToIntegerListCollection();

            var query = RequestForbiddenLogRepository.GetPagingByFilter(filter);

            var paging = ConvertToNewTypePagingObject<RequestForbiddenLogEntity, RequestForbiddenLogViewModel>(query);

            paging.Items = query.Items
                .Select(x =>
                {
                    var convert = ConvertToViewModels(x);
                    return convert;
                })
                .ToList();

            return paging;
        }


        private RequestForbiddenLogViewModel ConvertToViewModels(RequestForbiddenLogEntity queryEntity)
        {
            var convert = queryEntity.ReturnInjecteObject<RequestForbiddenLogViewModel>();
            convert.TypeName = ((RequestForbiddenTypeEnum)queryEntity.Type).ToString();
            convert.MerchantName = MerchantRepository.GetSingleByCache(queryEntity.MerchantID)?.DisplayName;
            return convert;
        }


    }
}
