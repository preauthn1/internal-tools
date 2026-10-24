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
using target.Utility.Helper;
using target.Models;
using target.DA.Repository;

namespace target.BL.Services
{

    public class VendorService : BaseAdminSysService, IVendorService
    {
        public VendorRepository VendorRepository { get; set; }


        public PagedList<VendorViewModel> GetByPagingFilter(int userID , QueryVendorApiRequestParameterModel filter)
        {
            var query = VendorRepository.GetPagingByFilter(filter);

            var viewModels = ConvertToNewTypePagingObject<VendorEntity, VendorViewModel>(query);
            viewModels.Items = query.Items.Select(x =>
                {
                    var convert = ConvertToViewModels(x);
                    return convert;
                })
                .ToList();

            return viewModels;
        }


        public VendorViewModel GetByID(int targetID)
        {
            var queryVendor = VendorRepository.GetById(targetID);

            if (queryVendor == null)
                throw new TerproException(ApiErrorCode.NotFoundData);

            var convert = ConvertToViewModels(queryVendor);
            return convert;
        }


        public void Create(VendorCreateRequestModel model)
        {           
            model.RuleFor(x => x.VendorID).Rule_NotZeroNotMinus();
            model.RuleFor(x => x.VendorCode).Rule_NotEmptyText();
            model.RuleFor(x => x.Name).Rule_NotEmptyText();
            

            var newEntity = model.ReturnInjecteObject<VendorEntity>();

            var currentVendor = GetCurrentUser();
            newEntity.CreateTime = DateTimeHelper.GetCurrentTime();
            newEntity.UpdateTime = DateTimeHelper.GetCurrentTime();
            newEntity.UpdateBy = currentVendor.LoginName;
            newEntity.CreateBy = currentVendor.LoginName;


            var newVendorID = VendorRepository.Create(newEntity);
        }


        public void Update(VendorUpdateRequestModel model)
        {
            model.RuleFor(x => x.VendorID).Rule_NotZeroNotMinus();
            model.RuleFor(x => x.VendorCode).Rule_NotEmptyText();
            model.RuleFor(x => x.Name).Rule_NotEmptyText();
            model.RuleFor(x => x.Status).Rule_NotZeroNotMinus();

            var queryExist = VendorRepository.GetById(model.VendorID);

            if (queryExist == null)
                throw new TerproException(ApiErrorCode.NotFoundData);

            queryExist.UpdateTime = DateTimeHelper.GetCurrentTime();
            queryExist.UpdateBy = GetCurrentUser().LoginName;
            queryExist.VendorCode = model.VendorCode;
            queryExist.Status = model.Status;
            queryExist.Name = model.Name;

            VendorRepository.Update(queryExist);
        }


        public void Delete(int targetID)
        {
            var queryExist = VendorRepository.GetById(targetID);

            if (queryExist == null)
                throw new TerproException(ApiErrorCode.NotFoundData);

            queryExist.Status = NormalActiveStatus.PermanentDelete.ToInteger();
            VendorRepository.Update(queryExist);
        }


        private VendorViewModel ConvertToViewModels(VendorEntity queryVendor)
        {
            var convert = queryVendor.ReturnInjecteObject<VendorViewModel>();
            return convert;
        }

    }
}
