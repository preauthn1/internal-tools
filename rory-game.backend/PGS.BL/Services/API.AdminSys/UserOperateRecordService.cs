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

    public class UserOperateRecordService : BaseAdminSysService, IUserOperateRecordService
    {
        public UserOperateRecordRepository UserOperateRecordRepository { get; set; }

        public PlayerRepository PlayerRepository { get; set; }

        public MerchantRepository MerchantRepository { get; set; }

        public UserRepository UserRepository { get; set; }


        public PagedList<UserOperateRecordViewModel> GetByPagingFilter(int userID , QueryUserOperateRecordApiRequestParameterModel filter)
        {
            filter.MerchantIDs = filter.MerchantID.ToIntegerListCollection();
            filter.MerchantIDs = FilterAvailablePermissionMerchantIDs(userID , filter.MerchantIDs);

            filter.TargetIDs = filter.PlayerID.ToIntegerListCollection();
            filter.OperateTypes = new List<int>
                {
                     UserOperateTypeEnum.Player_UnlockStatus.ToInteger(),
                     UserOperateTypeEnum.Player_LockStatus.ToInteger(),
                };

            if (!string.IsNullOrWhiteSpace(filter.PlayerName))
            {
                filter.TargetIDs = PlayerRepository.GetPlayerIDsByPlayerName(filter.PlayerName);
            }

            var query = UserOperateRecordRepository.GetPagingByFilter(filter);

            var viewModels = ConvertToNewTypePagingObject<UserOperateRecordEntity, UserOperateRecordViewModel>(query);
            viewModels.Items = query.Items.Select(x =>
                {
                    var convert = ConvertToViewModels(x);
                    return convert;
                })
                .ToList();

            return viewModels;
        }

        private UserOperateRecordViewModel ConvertToViewModels(UserOperateRecordEntity queryUserOperateRecord)
        {
            var convert = queryUserOperateRecord.ReturnInjecteObject<UserOperateRecordViewModel>();
            convert.MerchantName = MerchantRepository.GetSingleByCache(queryUserOperateRecord.MerchantID)?.Name;
            convert.PlayerName = PlayerRepository.GetBriefByCache((int)queryUserOperateRecord.TargetID)?.PlayerName;
            convert.UpdateBy = UserRepository.GetSingleByCache(queryUserOperateRecord.UserID)?.UserName;
            convert.PlayerID = (int)queryUserOperateRecord.TargetID;

            return convert;
        }

    }
}
