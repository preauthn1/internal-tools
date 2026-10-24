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

    public class UserOperateLogService : BaseAdminService, IUserOperateLogService
    {

        public UserOperateLogRepository UserOperateLogRepository { get; set; }
        

        public PagedList<UserOperateLogViewModel> GetByPagingFilter(int userID, QueryUserOperateLogApiRequestParameterModel filter)
        {
            var query = UserOperateLogRepository.GetPagingByFilter(filter);

            var paging = ConvertToNewTypePagingObject<UserOperateLogEntity, UserOperateLogViewModel>(query);

            paging.Items = query.Items
                .Select(x =>
                {
                    var convert = ConvertToViewModels(x);
                    return convert;
                })
                .ToList();

            return paging;
        }


        public void WriteOperateLog(int userID, UserOperateLogCreateRequestModel model)
        {
            UserOperateLogRepository.CreateByParameter(userID, (UserOperateTypeEnum)model.OperationType, null, model.DebugInfo);
        }


        private UserOperateLogViewModel ConvertToViewModels(UserOperateLogEntity queryEntity)
        {
            var convert = queryEntity.ReturnInjecteObject<UserOperateLogViewModel>();
            convert.OperationTypeText = ((UserOperateTypeEnum)queryEntity.OperationType).ToString();
            return convert;
        }


    }
}
