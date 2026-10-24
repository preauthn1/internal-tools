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
using Lck.Utility.Extensions.Helper;

namespace target.BL.Services
{

    public class ServerListService : BaseAdminSysService, IServerListService
    {
        public ServerListRepository ServerListRepository { get; set; }

        public ServerMerchantMappingRepository ServerMerchantMappingRepository { get; set; }

        public PlayerCurrentServerRepository PlayerCurrentServerRepository { get; set; }


        public PagedList<ServerListViewModel> GetByPagingFilter(QueryServerListApiRequestParameterModel filter)
        {
            filter.GameIDs = filter.GameID.HasValue ? filter.GameID.Value.ToListCollection() : filter.GameIDs;

            var query = ServerListRepository.GetPagingByFilter(filter);

            var queryMappings = ServerMerchantMappingRepository.GetAll();

            var viewModels = ConvertToNewTypePagingObject<ServerListEntity, ServerListViewModel>(query);
            viewModels.Items = query.Items.Select(x =>
                {
                    var convert = ConvertToViewModels(x , queryMappings);
                    return convert;
                })
                .ToList();

            return viewModels;
        }


        public ServerListViewModel GetByID(int serverID)
        {
            var queryServerList = ServerListRepository.GetById(serverID);

            if (queryServerList == null)
                throw new TerproException(ApiErrorCode.NotFoundServer);

            var queryMappings = ServerMerchantMappingRepository.GetByServerID(serverID);
            var convert = ConvertToViewModels(queryServerList , queryMappings);
            return convert;
        }


        public void Create(ServerListCreateRequestModel model)
        {
            model.RuleFor(x => x.ServerID).Rule_NotZeroNotMinus();
            model.RuleFor(x => x.ServerName).Rule_NotEmptyText();
            model.RuleFor(x => x.GameID).Rule_NotZeroNotMinus();
            model.RuleFor(x => x.MinAvailableMoney).Rule_NotMinus();
            model.RuleFor(x => x.Status).Rule_NotMinus();
            model.RuleFor(x => x.FlagLevel).Rule_NotZeroNotMinus();
            model.RuleFor(x => x.Currency).Rule_NotEmptyText();

            bool isExist = ServerListRepository.IsExist(model.ServerID);

            if (isExist)
                throw new TerproException(ApiErrorCode.BusinessLogicDeny_DuplicateDataExist);

            var querySameName = ServerListRepository.GetByFilter(new QueryServerListRequestViewModel
            {
                ServerNameFullMatch = model.ServerName,
                GameIDs = model.GameID.ToListCollection()
            });

            if(querySameName.IsNotEmptyNull())
                throw new TerproException(ApiErrorCode.BusinessLogicDeny_DuplicateDataExist);


            ////請實作參數檢查
            var newEntity = model.ReturnInjecteObject<ServerListEntity>();
            var currentServerList = GetCurrentUser();

            newEntity.CreateTime = DateTimeHelper.GetCurrentTime();
            newEntity.UpdateTime = DateTimeHelper.GetCurrentTime();
            newEntity.UpdateBy = currentServerList.LoginName;

            var newServerListID = ServerListRepository.Create(newEntity);
        }


        public void Update(ServerListUpdateRequestModel model)
        {
            model.RuleFor(x => x.ServerID).Rule_NotZeroNotMinus();
            model.RuleFor(x => x.ServerName).Rule_NotEmptyText();
            model.RuleFor(x => x.GameID).Rule_NotZeroNotMinus();
            model.RuleFor(x => x.MinAvailableMoney).Rule_NotMinus();
            model.RuleFor(x => x.Status).Rule_NotMinus();
            model.RuleFor(x => x.FlagLevel).Rule_NotZeroNotMinus();
            model.RuleFor(x => x.Currency).Rule_NotEmptyText();

            var queryExist = ServerListRepository.GetById(model.ServerID);

            if (queryExist == null)
                throw new TerproException(ApiErrorCode.NotFoundServer);

            queryExist.ServerName = model.ServerName;
            queryExist.GameID = model.GameID;
            queryExist.FlagLevel = model.FlagLevel;
            queryExist.ClientAddr = model.ClientAddr;
            queryExist.HttpAddr = model.HttpAddr;
            queryExist.NotifyAddr = model.NotifyAddr;
            queryExist.Status = model.Status;
            queryExist.MinAvailableMoney = model.MinAvailableMoney;
            queryExist.UpdateTime = DateTimeHelper.GetCurrentTime();
            queryExist.UpdateBy = GetCurrentUser().LoginName;
            queryExist.Currency = model.Currency;
            queryExist.RoomFlagLevel = model.RoomFlagLevel;

            ServerListRepository.Update(queryExist);

            CacheRemoveHelper.RemoveCache_WhenUpdateServer(model.ServerID);
        }


        public void Delete(int serverID)
        {
            var queryExist = ServerListRepository.GetById(serverID);

            if (queryExist == null)
                throw new TerproException(ApiErrorCode.NotFoundServer);

            var now = DateTimeHelper.GetCurrentTime();

            if ((now - queryExist.UpdateTime).TotalSeconds >= 600)
                throw new TerproException(ApiErrorCode.BusinessLogicDeny_NotAllowDeleteData_ModifyRecently);

            var queryCurretPlayer = PlayerCurrentServerRepository.GetByFilter(new QueryPlayerCurrentServerRequestViewModel
            {
                ServerIDs = serverID.ToListCollection()
            });

            if (queryCurretPlayer.IsNotEmptyNull())
                throw new TerproException(ApiErrorCode.BusinessLogicDeny_NotAllowDeleteData_CurrentServer);

            ServerListRepository.Delete(serverID);
            ServerListRepository.DeleteRelationDataWhenServerDelete(serverID);


            CacheRemoveHelper.RemoveCache_WhenUpdateServer(serverID);
        }


        private ServerListViewModel ConvertToViewModels(ServerListEntity queryServerList , List<ServerMerchantMappingEntity> mappings)
        {
            var convert = queryServerList.ReturnInjecteObject<ServerListViewModel>();
            convert.AssignMerchantTotalCount = mappings.Where(x => x.ServerID == queryServerList.ServerID).Count();
            convert.GameName = DbMultiLanguageTextHelper.GetGameDisplayName(queryServerList.GameID);
            convert.RoomFlagLevelText = MultiLanguageTextHelper.GetEnumsCultureText((GameRoomFlagLevelEnum)queryServerList.RoomFlagLevel);
            return convert;
        }

    }
}
