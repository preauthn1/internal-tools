using MSC.Models.Paging;
using MSC.Utility.IOC;
using System;
using System.Collections.Generic;
using System.Linq;
using System.Reflection;
using MSC.Utility.Extensions;
using MSC.Utility.Helper;
using MSC.Core.Repository;
using MSC.DA.Repository;
using Lck.Utility.Extensions.Helper;
using MSC.Models.Models;

namespace MSC.BL.Services
{
    public class BaseAdminService : BaseService
    {

        public UserRepository UserRepository { get; set; }


        public BaseAdminService()
        {
            InstanceAllRepositoryImplement();
        }

        protected void CheckHasMerchantViewPermission(int userID, int merchantID)
        {
            List<int> availableMerchantIDs = GetUserMerchantIDsByCache(userID);

            if (availableMerchantIDs.Contains(merchantID))
                return;

            throw new MSCException(ErrorCode.Forbidden_WatchDataPermissionNotAllow);
        }


        protected List<int> FilterAvailableMerchants(int userID , List<int> currentMerchantIDs)
        {
            List<int> availableMerchantIDs = GetUserMerchantIDsByCache(userID);

            if (currentMerchantIDs.IsEmptyOrNull())
                return availableMerchantIDs;

            return currentMerchantIDs.Where(z=> availableMerchantIDs.Contains(z)).ToList();
        }


        protected List<int> GetUserMerchantIDsByCache(int userID)
        {
            string cacheKey = CacheKey.AdminSysAPI_UserMerchantIDs(userID);

            var queryCache = CacheHelper.GetCacheOrQuery(cacheKey, () =>
               {
                   var rep = new UserMerchantMappingRepository();

                   var query = rep.GetByFilter(new Models.ViewModels.QueryUserMerchantMappingRequestViewModel
                   {
                       UserIDs = userID.ToListCollection()
                   });

                   return query.Select(z => z.MerchantID).ToList();
               }, 30);

            return queryCache;
        }


        /// <summary>
        /// 轉換成另一種class的 PagedList 物件
        /// </summary>
        protected PagedList<T_New> ConvertToNewTypePagingObject<T_Original , T_New>(PagedList<T_Original> model)
        {
            var resulk = model.ReturnInjecteObject<PagedList<T_New>>();
            return resulk;
        }


        protected string GetUserNameByUserID(int userID)
        {
            var queryUser = UserRepository.GetSingleByCache(userID);
            return queryUser?.UserName;
        }


    }
}
