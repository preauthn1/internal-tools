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
    public class GamePayTableImageService : BaseAdminSysService, IGamePayTableImageService
    {

        public GamePayTableImageRepository GamePayTableImageRepository { get; set; }



        public PagedList<GamePayTableImageViewModel> GetByPagingFilter(int userID, QueryGamePayTableImageApiRequestParameterModel filter)
        {
            var cureentUser = GetCurrentOwnPermissionByCache(userID);
            filter.GameIDs = filter.GameID.ToIntegerListCollection();

            var query = GamePayTableImageRepository.GetPagingByFilter(filter);

            var viewModels = ConvertToNewTypePagingObject<GamePayTableImageEntity, GamePayTableImageViewModel>(query);
            viewModels.Items = query.Items.Select(x =>
                {
                    var convert = ConvertToViewModel(x);
                    return convert;
                })
                .ToList();

            return viewModels;
        }



        public CRUDResultModel CreateOrUpdate(int userID, GamePayTableImageUpdateRequestModel model)
        {
            model.RuleFor(x => x.Code).Rule_NotEmptyText();
            model.RuleFor(x => x.GameID).Rule_NotZeroNotMinus();
            model.RuleFor(x => x.ImageUrl).Rule_NotEmptyText();

            var newEntity = model.ReturnInjecteObject<GamePayTableImageEntity>();


            var queryExist = GamePayTableImageRepository.GetByFilter(new QueryGamePayTableImageRequestViewModel
                {
                    GameIDs = model.GameID.ToListCollection(),
                    Code = model.Code
                })
                .FirstOrDefault();


            //若資料以存在 update
            if(queryExist != null)
            {
                model.Id = queryExist.Id;
                return Update(userID , model);
            }

            newEntity.CreateTime = DateTimeHelper.GetCurrentTime();
            newEntity.Code = model.Code.Trim();

            var resultId = GamePayTableImageRepository.Create(newEntity);

            return new CRUDResultModel
            {
                ResultId = resultId,
                IsSuccess = true,
                Message = "新增成功",
            };
        }


        public CRUDResultModel Update(int userID, GamePayTableImageUpdateRequestModel model)
        {
            model.RuleFor(x => x.Id).Rule_NotZero().Rule_NotMinus();
            model.RuleFor(x => x.Code).Rule_NotEmptyText();
            model.RuleFor(x => x.ImageUrl).Rule_NotEmptyText();

            var queryExsit = GamePayTableImageRepository.GetById(model.Id);

            if (queryExsit == null)
                throw new TerproException(ApiErrorCode.NotFoundData);

            queryExsit.Code = model.Code;
            queryExsit.Memo = model.Memo;
            queryExsit.ImageUrl = model.ImageUrl;


            GamePayTableImageRepository.Update(queryExsit);
            return new CRUDResultModel(true);
        }


        public CRUDResultModel Delete(int userID, int targetID)
        {
            var queryExsit = GamePayTableImageRepository.GetById(targetID);

            if (queryExsit == null)
                throw new TerproException(ApiErrorCode.NotFoundData);

            GamePayTableImageRepository.Delete(targetID);
            return new CRUDResultModel(true);
        }


        private GamePayTableImageViewModel ConvertToViewModel(GamePayTableImageEntity query)
        {
            //將model轉成viewModels
            var convert = query.ReturnInjecteObject<GamePayTableImageViewModel>();
            convert.GameName = DbMultiLanguageTextHelper.GetGameDisplayName(query.GameID);

            return convert;
        }



    }

}
