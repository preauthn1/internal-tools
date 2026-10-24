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

    public class ExternalApiRequestLogService : BaseAdminService, IExternalApiRequestLogService
    {
        public ExternalApiRequestLogRepository ExternalApiRequestLogRepository { get; set; }


        public PagedList<ExternalApiRequestLogViewModel> GetByPagingFilter(int userID, QueryExternalApiRequestLogApiRequestParameterModel filter)   
        {
            var query = ExternalApiRequestLogRepository.GetPagingByFilter(filter);

            var viewModels = ConvertToNewTypePagingObject<ExternalApiRequestLogEntity, ExternalApiRequestLogViewModel>(query);
            viewModels.Items = query.Items.Select(x =>
                {
                    var convert = ConvertToViewModels(x);
                    return convert;
                })
                .ToList();

            return viewModels;
        }

        private ExternalApiRequestLogViewModel ConvertToViewModels(ExternalApiRequestLogEntity model)
        {
            var convert = model.ReturnInjecteObject<ExternalApiRequestLogViewModel>();
            return convert;
        }


    }
}
