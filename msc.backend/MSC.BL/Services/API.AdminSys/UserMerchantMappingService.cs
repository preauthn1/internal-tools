using Lck.Utility.Extensions.Helper;
using MSC.BL.Interface;
using MSC.DA.Repository;
using MSC.Models;
using MSC.Models.CoreLibs;
using MSC.Models.Enums;
using MSC.Models.Models;
using MSC.Models.Paging;
using MSC.Models.Table;
using MSC.Models.ViewModels;
using MSC.Utility.Extensions;
using MSC.Utility.Helper;
using System;
using System.Collections.Generic;
using System.Linq;

namespace MSC.BL.Services
{
    public class UserMerchantMappingService : BaseService, IUserMerchantMappingService
    {
        public UserMerchantMappingRepository UserMerchantMappingRepository { get; set; }

        public MerchantRepository MerchantRepository { get; set; }


        public UserMerchantMappingViewModel GetUserMerchantsEnableFlag(int userID)
        {
            var queryAllMerchants = MerchantRepository.GetAll();

            var query = UserMerchantMappingRepository.GetByFilter(new QueryUserMerchantMappingRequestViewModel
            {
                UserIDs = userID.ToListCollection()
            });

            return new UserMerchantMappingViewModel
            {
                Items = queryAllMerchants.Select(z => new UserMerchantMappingItemViewModel
                {
                    MerchantID = z.MerchantID,
                    MerchantName = z.DisplayName,
                    IsEnable = query.Any(x=> x.MerchantID == z.MerchantID)
                })
                .ToList()
            };
        }

        public void UpdateUserMerchants(int currentUserID, UserMerchantMappingRequestModel model)
        {
            model.RuleFor(z => z.UserID).Rule_NotZeroNotMinus();
            model.RuleFor(z => z.MerchantIDs).Rule_IsNotEmptyCollection(); //至少要有一筆

            // Step 1. 清除全部
            UserMerchantMappingRepository.DeleteByUserID(model.UserID);

            var inserts = model.MerchantIDs.Select(z => new UserMerchantMappingEntity
                {
                    UserID = model.UserID,
                    MerchantID = z
                })
                .ToList();

            // Step 2. 寫入
            UserMerchantMappingRepository.CreateMany(inserts);

            // Step 3. 更新cache
            CacheRemoveHelper.WhenUpdateUserMerchantMappings(model.UserID);
        }

    }
}
