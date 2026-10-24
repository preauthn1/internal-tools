using target.BL.Interface;
using target.BL.Services;
using target.DA.Repository;
using target.Models.CoreLibs;
using target.Models.Models;
using target.Models.Paging;
using target.Models.Table;
using target.Models.ViewModels;
using target.Utility.Extensions;
using target.Utility.Helper;
using System;
using System.Collections.Generic;
using System.Data;
using System.Data.SqlClient;
using System.Linq;

namespace target.BL.Services
{
    public class ExternalApiRequestLogService : BaseAdminSysService, IExternalApiRequestLogService
    {

        public ExternalApiRequestLogRepository ExternalApiRequestLogRepository { get; set; }



        public PagedList<ExternalApiRequestLogViewModel> GetByPagingFilter(int userID, QueryExternalApiRequestLogApiRequestParameterModel filter)
        {
            filter.IDs = filter.Id.ToIntegerListCollection();

            var cureentUser = GetCurrentOwnPermissionByCache(userID);

            var query = ExternalApiRequestLogRepository.GetPagingByFilter(filter);

            var viewModels = ConvertToNewTypePagingObject<ExternalApiRequestLogEntity, ExternalApiRequestLogViewModel>(query);
            viewModels.Items = query.Items.Select(x =>
                {
                    var convert = ConvertToViewModel(x);
                    return convert;
                })
                .ToList();

            return viewModels;
        }


        public ExternalApiRequestLogViewModel GetByID(long id)
        {
            var query = ExternalApiRequestLogRepository.GetById(id);
            return ConvertToViewModel(query);
        }


        public ExternalApiRequestLogViewModel GetByRequestID(string requestID)
        {
            var query = ExternalApiRequestLogRepository.GetByFilter(new QueryExternalApiRequestLogRequestViewModel
                {
                    RequestID = requestID
                })
                .FirstOrDefault();

            return ConvertToViewModel(query);
        }




        private ExternalApiRequestLogViewModel ConvertToViewModel(ExternalApiRequestLogEntity query)
        {
            //將model轉成viewModels
            var convert = query.ReturnInjecteObject<ExternalApiRequestLogViewModel>();
            return convert;
        }



    }

}
