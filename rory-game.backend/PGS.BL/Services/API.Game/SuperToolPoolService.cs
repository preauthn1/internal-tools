using target.BL.Interface;
using target.BL.Services;
using target.Models.CoreLibs;
using target.Models.Paging;
using target.Models.Table;
using target.Models.ViewModels;
using System;
using System.Collections.Generic;
using System.Data;
using System.Data.SqlClient;
using System.Linq;
using target.Utility.Extensions;
using target.Models.Models;
using target.DA.Repository;

namespace target.BL.Services
{
    public class SuperToolPoolService : BaseService, ISuperToolPoolService
    {
        public GameListRepository GameListRepository { get; set; }

        public SuperToolSystemPoolWinLossRepository SuperToolSystemPoolWinLossRepository { get; set; }


        public SuperToolPoolViewModel GetWinLossPool(SuperToolPoolRequestModel model)
        {
            model.RuleFor(x => x.GameID).Rule_NotZeroNotMinus();
            model.RuleFor(x => x.FlagID).Rule_NotMinus();
            
            if (!GameListRepository.IsGameExistByCache(model.GameID))
                throw new TerproException(ApiErrorCode.NotFoundGame);

            var query = SuperToolSystemPoolWinLossRepository.GetByGameID_And_FlagID(model.GameID , model.FlagID);
            
            if(query != null)
                return ConvertToViewModels(query);

            SuperToolSystemPoolWinLossRepository.UpdateOrCreateSystemPoolWinLoss(model.GameID, model.FlagID, 0, string.Empty);

            var queryAgain = SuperToolSystemPoolWinLossRepository.GetByGameID_And_FlagID(model.GameID, model.FlagID);

            if (queryAgain == null)
                throw new TerproException(ApiErrorCode.NotFoundData);

            return ConvertToViewModels(queryAgain);
        }


        public SuperToolPoolViewModel UpdateWinLossPool(SuperToolPoolUpdateModel model)
        {
            model.RuleFor(x => x.GameID).Rule_NotZeroNotMinus();
            model.RuleFor(x => x.FlagID).Rule_NotMinus();

            if (!GameListRepository.IsGameExistByCache(model.GameID))
                throw new TerproException(ApiErrorCode.NotFoundGame);

            SuperToolSystemPoolWinLossRepository.UpdateOrCreateSystemPoolWinLoss(model.GameID, model.FlagID, model.UpdatePoolAmount, model.GameData);

            var query = SuperToolSystemPoolWinLossRepository.GetByGameID_And_FlagID(model.GameID, model.FlagID);

            if (query == null)
                throw new TerproException(ApiErrorCode.NotFoundData);

            return ConvertToViewModels(query);
        }


        private SuperToolPoolViewModel ConvertToViewModels(SuperToolSystemPoolWinLossEntity entity)
        {
            var convert = entity.ReturnInjecteObject<SuperToolPoolViewModel>();
            return convert;
        }
        
    }

}
