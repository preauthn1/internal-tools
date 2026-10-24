using target.BL.Interface;
using target.BL.Services;
using target.DA.Repository;
using target.Models.CoreLibs;
using target.Models.Enums;
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
    public class VendorMaintenanceService : BaseAdminSysService, IVendorMaintenanceService
    {

        public VendorMaintenanceRepository VendorMaintenanceRepository { get; set; }



        public PagedList<VendorMaintenanceViewModel> GetByPagingFilter(int userID, QueryVendorMaintenanceApiRequestParameterModel filter)
        {
            var cureentUser = GetCurrentOwnPermissionByCache(userID);

            var query = VendorMaintenanceRepository.GetPagingByFilter(filter);

            var viewModels = ConvertToNewTypePagingObject<VendorMaintenanceEntity, VendorMaintenanceViewModel>(query);
            viewModels.Items = query.Items.Select(x =>
                {
                    var convert = ConvertToViewModel(x);
                    return convert;
                })
                .ToList();

            return viewModels;
        }




        public CRUDResultModel Create(int userID, VendorMaintenanceCreateViewModel model)
        {
            model.RuleFor(x => x.VendorID).Rule_NotZeroNotMinus();
            model.RuleFor(x => x.StartTime).Rule_NotEmptyDateTime();
            model.RuleFor(x => x.EndTime).Rule_NotEmptyDateTime();
            model.RuleFor(x => x.StartTime).Rule_StartTimeMustSmallThanEndTime(model.EndTime);


            var newEntity = model.ReturnInjecteObject<VendorMaintenanceEntity>();

            DateTime now = DateTimeHelper.GetCurrentTime();

            var currentVendor = GetCurrentUser();

            newEntity.StartTime = SecondsToZero(newEntity.StartTime);
            newEntity.EndTime = SecondsToZero(newEntity.EndTime);
            newEntity.Status = NormalActiveStatus.Enable.ToInteger();
            newEntity.CreateTime = now;
            newEntity.CreateBy = currentVendor.LoginName;
            newEntity.UpdateTime = now;
            newEntity.UpdateBy = currentVendor.LoginName;

            var resultId = VendorMaintenanceRepository.Create(newEntity);

            return new CRUDResultModel
            {
                ResultId = resultId,
                IsSuccess = true,
                Message = "新增成功",
            };
        }

        public CRUDResultModel Delete(int userID, int targetID)
        {
            var queryExsit = VendorMaintenanceRepository.GetById(targetID);

            if (queryExsit == null)
                throw new TerproException(ApiErrorCode.NotFoundData);

            queryExsit.Status = NormalActiveStatus.PermanentDelete.ToInteger();

            VendorMaintenanceRepository.Delete(targetID);
            return new CRUDResultModel(true);
        }


        private VendorMaintenanceViewModel ConvertToViewModel(VendorMaintenanceEntity query)
        {
            DateTime now = DateTimeHelper.GetCurrentTime();

            //將model轉成viewModels
            var convert = query.ReturnInjecteObject<VendorMaintenanceViewModel>();
            convert.VendorName = ((GameVendorTypeEnum)query.VendorID).ToString();
            convert.StatusText = MultiLanguageTextHelper.GetEnumsCultureText((NormalActiveStatus)query.Status);
            convert.IsExpired = now > query.EndTime;

            return convert;
        }


        private DateTime SecondsToZero(DateTime time)
        {
            return new DateTime(time.Year, time.Month, time.Day, time.Hour, time.Minute, 0);
        }


    }

}
