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
using target.Models.ViewModels.API.AdminSys;
using target.DA.Repository;
using Lck.Utility.Extensions.Helper;

namespace target.BL.Services
{

    public class OrganizationService : BaseAdminSysService, IOrganizationService
    {

        public OrganizationRepository OrganizationRepository { get; set; }

        public OrganizationMerchantMappingRepository OrganizationMerchantMappingRepository { get; set; }

        public OrganizationRoleMappingRepository OrganizationRoleMappingRepository { get; set; }

        public UserRepository UserRepository { get; set; }



        public PagedList<OrganizationViewModel> GetPagingByFilter(QueryOrganizationApiRequestParameterModel filter)
        {
            PagedList<OrganizationEntity> query = OrganizationRepository.GetPagingByFilter(filter);

            var paging = ConvertToNewTypePagingObject<OrganizationEntity,OrganizationViewModel>(query);

            paging.Items = query.Items
                .Select(x => 
                {
                    var convert = ConvertToViewModels(x);
                    return convert;
                })
                .ToList();

            SetExtraProperty(paging.Items);
            return paging;
        }


        public OrganizationViewModel GetByID(int orgID)
        {
            var queryOrganization = OrganizationRepository.GetById(orgID);

            if (queryOrganization == null)
                throw new TerproException(ApiErrorCode.NotFoundOrganization);

            var convert = ConvertToViewModels(queryOrganization);

            convert.Roles = OrganizationRoleMappingRepository.GetByOrgID(orgID).Select(x => x.RoleID).ToList();
            convert.Merchants = OrganizationMerchantMappingRepository.GetByOrgID(orgID).Select(x => x.MerchantID).ToList();
            return convert;
        }


        public void Create(OrganizationCreateRequestModel model)
        {
            model.RuleFor(x => x.OrganizationName).Rule_NotEmptyText().Rule_MaxLengthOK(50);
            model.RuleFor(x => x.OrganizationCode).Rule_NotEmptyText().Rule_MaxLengthOK(5);
            model.RuleFor(x => x.Memo).Rule_MaxLengthOK(50); 
            model.RuleFor(x => x.Status).Rule_NotMinus();
            model.RuleFor(x => x.Lang).Rule_NotEmptyText().Rule_IsSupportLanguageCode();
                       
            var newEntity = model.ReturnInjecteObject<OrganizationEntity>();

            var currentOrganization = GetCurrentUser();
            newEntity.OrganizationType = OrganizationTypeEnum.Normal.ToInteger(); //暫時不允許編輯
            newEntity.CreateTime = DateTimeHelper.GetCurrentTime();
            newEntity.UpdateTime = DateTimeHelper.GetCurrentTime();
            newEntity.UpdateBy = currentOrganization.LoginName;
            newEntity.CreateBy = currentOrganization.LoginName;

            var queryDuplicate = OrganizationRepository.GetByFilter(new QueryOrganizationRequestViewModel
            {
                OrganizationCode = model.OrganizationCode
            });

            if (queryDuplicate.IsNotEmptyNull())
                throw new TerproException(ApiErrorCode.BusinessLogicDeny_DuplicateDataExist);

            //step 1. 新增org
            var newOrganizationID = OrganizationRepository.Create(newEntity);            

            //step 2. 寫入 mapping
            OrganizationMerchantMappingRepository.DeleteAndInsertNew(newOrganizationID, model.Merchants);

            //step 3. 寫入 role mapping
            OrganizationRoleMappingRepository.DeleteAndInsertNew(newOrganizationID, model.Roles);

            //更新cache
            CacheRemoveHelper.RemoveCache_WhenUpdateOrganization();
        }


        public void Update(OrganizationUpdateRequestModel model)
        {
            model.RuleFor(x => x.OrganizationID).Rule_NotZeroNotMinus();
            model.RuleFor(x => x.OrganizationName).Rule_NotEmptyText();
            model.RuleFor(x => x.Memo).Rule_MaxLengthOK(50);
            model.RuleFor(x => x.Status).Rule_NotMinus();
            model.RuleFor(x => x.Lang).Rule_NotEmptyText().Rule_IsSupportLanguageCode();

            var queryExist = OrganizationRepository.GetById(model.OrganizationID);

            if (queryExist == null)
                throw new TerproException(ApiErrorCode.NotFoundOrganization);

            queryExist.OrganizationName = model.OrganizationName;
            queryExist.Memo = model.Memo;
            queryExist.Lang = model.Lang;
            queryExist.Status = model.Status;
            queryExist.UpdateTime = DateTimeHelper.GetCurrentTime();
            queryExist.UpdateBy = GetCurrentUser().LoginName;


            //step 1. 編輯
            OrganizationRepository.Update(queryExist);

            //step 2. 更新 mapping
            OrganizationMerchantMappingRepository.DeleteAndInsertNew(model.OrganizationID , model.Merchants);

            //step 3. 更新 role mapping
            OrganizationRoleMappingRepository.DeleteAndInsertNew(model.OrganizationID , model.Roles);


            //更新cache
            CacheRemoveHelper.RemoveCache_WhenUpdateOrganization();
        }


        public void Delete(int targetID)
        {
            var queryExist = OrganizationRepository.GetById(targetID);

            if (queryExist == null)
                throw new TerproException(ApiErrorCode.NotFoundOrganization);

            queryExist.Status = NormalActiveStatus.PermanentDelete.ToInteger();
            OrganizationRepository.Update(queryExist);
        }

        

        private OrganizationViewModel ConvertToViewModels(OrganizationEntity queryOrganization)
        {
            var convert = queryOrganization.ReturnInjecteObject<OrganizationViewModel>();
            return convert;
        }

        private void SetExtraProperty(List<OrganizationViewModel> items)
        {
            var queryUsers = UserRepository.SelectSpecificColumnsByFilter(new QueryUserRequestViewModel(), " UserID , OrganizationID ");

            foreach (var item in items)
            {
                item.UsersTotalCount = queryUsers.Where(z => z.OrganizationID == item.OrganizationID).Count();
            }
        }

    }
}
