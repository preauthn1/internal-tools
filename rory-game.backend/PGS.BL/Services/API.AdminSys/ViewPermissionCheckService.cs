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

    public class ViewPermissionCheckService : BaseAdminSysService, IViewPermissionCheckService
    {

        public UserRepository UserRepository { get; set; }


        public PlayerRepository PlayerRepository { get; set; }


        public void CheckViewPlayerPermission(int currentUserID, int targerPlayerID)
        {
            CheckValueIsZeroOrMinus(currentUserID);
            CheckValueIsZeroOrMinus(targerPlayerID);

            var userPermission = GetCurrentOwnPermissionByCache(currentUserID);

            //超級管理員
            if (userPermission.IsSuperAdminOrg)
                return;

            var queryPlayer = PlayerRepository.GetById(targerPlayerID);

            if (queryPlayer == null)
                throw new TerproException(ApiErrorCode.NotFoundPlayer);
            
            var isExist = userPermission.AvailableMerchantIDs.Contains(queryPlayer.MerchantID);

            if(!isExist)
                throw new TerproException(ApiErrorCode.Forbidden_WatchDataPermissionNotAllow);

        }



        public void CheckViewUserPermission(int currentUserID , int targetUserID)
        {
            CheckValueIsZeroOrMinus(currentUserID);
            CheckValueIsZeroOrMinus(targetUserID);

            var userPermission = GetCurrentOwnPermissionByCache(currentUserID);

            //超級管理員
            if (userPermission.IsSuperAdminOrg)
                return;

            var query = UserRepository.GetByIds(new List<int>
            {
                currentUserID,
                targetUserID,
            });

            var matchCurrentUser = query.FirstOrDefault(x => x.UserID == currentUserID);
            var matchTargetUser = query.FirstOrDefault(x => x.UserID == targetUserID);

            if (matchCurrentUser == null || matchTargetUser == null)
                throw new TerproException(ApiErrorCode.NotFoundUser);

            //如果同組織OK
            if (matchCurrentUser.OrganizationID == matchTargetUser.OrganizationID)
                return;

            throw new TerproException(ApiErrorCode.Forbidden_WatchDataPermissionNotAllow);
        }


        public void CheckViewMerchantPermission(int currentUserID, int targetMerchantID)
        {
            CheckValueIsZeroOrMinus(currentUserID);
            CheckValueIsZeroOrMinus(targetMerchantID);

            var userPermission = GetCurrentOwnPermissionByCache(currentUserID);

            //超級管理員
            if (userPermission.IsSuperAdminOrg)
                return;

            if(!userPermission.AvailableMerchantIDs.Contains(targetMerchantID))
                throw new TerproException(ApiErrorCode.Forbidden_WatchDataPermissionNotAllow);
        }


        private void CheckValueIsZeroOrMinus(int myValue)
        {
            if (myValue == 0)
                throw new TerproException(ApiErrorCode.InvalidParameter_NumberEmptyZero);

            if (myValue <= 0)
                throw new TerproException(ApiErrorCode.InvalidParameter_NumberMinus);
        }



    }
}
