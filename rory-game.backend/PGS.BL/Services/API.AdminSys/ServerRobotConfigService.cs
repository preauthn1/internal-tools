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

    public class ServerRobotConfigService : BaseAdminSysService, IServerRobotConfigService
    {

        public ServerRobotConfigRepository ServerRobotConfigRepository { get; set; }


        public ServerListRepository ServerListRepository { get; set; }



        public PagedList<ServerRobotConfigViewModel> GetByPagingFilter(QueryServerRobotConfigApiRequestParameterModel filter)
        {
            filter.GameIDs = filter.GameID.HasValue ? filter.GameID.Value.ToListCollection() : new List<int>();
            var query = ServerRobotConfigRepository.GetPagingByFilter(filter);

            var viewModels = ConvertToNewTypePagingObject<ServerRobotConfigEntity, ServerRobotConfigViewModel>(query);
            viewModels.Items = query.Items.Select(x =>
            {
                var convert = ConvertToViewModels(x);
                return convert;
            })
            .ToList();

            return viewModels;
        }


        public ServerRobotConfigViewModel GetByID(int serverID)
        {
            var queryServerRobotConfig = ServerRobotConfigRepository.GetById(serverID);

            if (queryServerRobotConfig == null)
                throw new TerproException(ApiErrorCode.NotFoundData);

            var convert = ConvertToViewModels(queryServerRobotConfig);
            return convert;
        }


        public void Create(ServerRobotConfigCreateRequestModel model)
        {
            model.RuleFor(x => x.ServerID).Rule_NotZeroNotMinus();
            model.RuleFor(x => x.MaxLevel).Rule_NotMinus().Rule_NotTooLarge(10);
            model.RuleFor(x => x.MinLevel).Rule_NotMinus().Rule_NotTooLarge(10);
            model.RuleFor(x => x.MaxRobotCount).Rule_NotMinus();
            model.RuleFor(x => x.MinRobotCount).Rule_NotMinus();
            model.RuleFor(x => x.MaxLevel).Rule_NotMinus();
            model.RuleFor(x => x.BeginTime).Rule_NotEmptyText();
            model.RuleFor(x => x.EndTime).Rule_NotEmptyText();
                       
            bool isServerDataExist = ServerListRepository.IsExist(model.ServerID);

            if(!isServerDataExist)
                throw new TerproException(ApiErrorCode.NotFoundServer);

            bool isDataExist = ServerRobotConfigRepository.IsExist(model.ServerID);

            if(isDataExist)
                throw new TerproException(ApiErrorCode.BusinessLogicDeny_DuplicateDataExist);

            var newEntity = model.ReturnInjecteObject<ServerRobotConfigEntity>();

            var currentServerRobotConfig = GetCurrentUser();
            newEntity.CreateTime = DateTimeHelper.GetCurrentTime();
            newEntity.UpdateTime = DateTimeHelper.GetCurrentTime();
            newEntity.IsRobotEnable = true;
            newEntity.UpdateBy = currentServerRobotConfig.LoginName;
            newEntity.CreateBy = currentServerRobotConfig.LoginName;


            var newServerRobotConfigID = ServerRobotConfigRepository.Create(newEntity);
        }


        public void Update(ServerRobotConfigUpdateRequestModel model)
        {
            model.RuleFor(x => x.ServerID).Rule_NotZeroNotMinus();
            model.RuleFor(x => x.MaxLevel).Rule_NotMinus().Rule_NotTooLarge(10);
            model.RuleFor(x => x.MinLevel).Rule_NotMinus().Rule_NotTooLarge(10);
            model.RuleFor(x => x.MaxRobotCount).Rule_NotMinus();
            model.RuleFor(x => x.MinRobotCount).Rule_NotMinus();
            model.RuleFor(x => x.MaxLevel).Rule_NotMinus();
            model.RuleFor(x => x.BeginTime).Rule_NotEmptyText();
            model.RuleFor(x => x.EndTime).Rule_NotEmptyText();

            var queryExist = ServerRobotConfigRepository.GetById(model.ServerID);

            if (queryExist == null)
                throw new TerproException(ApiErrorCode.NotFoundData);

            queryExist.MaxLevel = model.MaxLevel;
            queryExist.MinLevel = model.MinLevel;
            queryExist.MaxRobotCount = model.MaxRobotCount;
            queryExist.MinRobotCount = model.MinRobotCount;
            queryExist.BeginTime = model.BeginTime;
            queryExist.EndTime = model.EndTime;
            queryExist.IsRobotEnable = model.IsRobotEnable;
            queryExist.UpdateTime = DateTimeHelper.GetCurrentTime();
            queryExist.UpdateBy = GetCurrentUser().LoginName;

            ServerRobotConfigRepository.Update(queryExist);

        }


        public void Delete(int targetID)
        {
            var queryExist = ServerRobotConfigRepository.GetById(targetID);

            if (queryExist == null)
                throw new TerproException(ApiErrorCode.NotFoundData);

            ServerRobotConfigRepository.Delete(targetID);
        }


        private ServerRobotConfigViewModel ConvertToViewModels(ServerRobotConfigEntity queryServerRobotConfig)
        {
            var convert = queryServerRobotConfig.ReturnInjecteObject<ServerRobotConfigViewModel>();
            convert.ServerName = ServerListRepository.GetBriefByCache(queryServerRobotConfig.ServerID)?.ServerName;
            return convert;
        }


    }
}
