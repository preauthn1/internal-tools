using Lck.Utility.Extensions.Helper;
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
    public class LanguageTextService : BaseAdminSysService, ILanguageTextService
    {

        public LanguageTextRepository LanguageTextRepository { get; set; }

        public GameListRepository GameListRepository { get; set; }


        public PagedList<LanguageTextViewModel> GetByPagingFilter(int userID, QueryLanguageTextApiRequestParameterModel filter)
        {
            var cureentUser = GetCurrentOwnPermissionByCache(userID);

            var query = LanguageTextRepository.GetGroupByPagingByFilter(filter);

            var viewModels = ConvertToNewTypePagingObject<LanguageTextGroupModel, LanguageTextViewModel>(query);
            viewModels.Items = query.Items.Select(x =>
                {
                    var convert = ConvertToViewModel(x);
                    return convert;
                })
                .ToList();

            SetExtraProperty(viewModels.Items);

            return viewModels;
        }


        public LanguageTextSingleInfoViewModel GetAllTextByLanguageCode(string languageCode)
        {
            var query = LanguageTextRepository.GetByFilter(new QueryLanguageTextRequestViewModel
            {
                LanguageCode = languageCode
            });

            if (query.IsEmptyOrNull())
            {
                return new LanguageTextSingleInfoViewModel
                {
                    Items = new List<LanguageTextItemViewModel>()
                };
            }

            var allLang = ExtensionOfEnum.GetEnumAllValue<LanguageEnum>();

            var items = allLang.Select(z =>
                {
                    string language = z.ToString();
                    var convert = new LanguageTextItemViewModel
                    {
                        LanguageCode = languageCode,
                        Language = language,
                        LanguageText = MultiLanguageTextHelper.GetEnumsCultureText(z)
                    };

                    var matchLang = query.FirstOrDefault(z => z.Language == language);

                    if(matchLang != null)
                    {
                        convert.Id = matchLang.Id;
                        convert.Text = matchLang.Text;
                    }

                    return convert;
                })
                .ToList();


            string insertSQL = GenerateInsertSQL(query);

            return new LanguageTextSingleInfoViewModel
            {
                InsertSQL = insertSQL,
                Items = items
            };
        }


        private string GenerateInsertSQL(List<LanguageTextEntity> query)
        {
            string insertSql = @"INSERT INTO [dbo].[LanguageText] ([LanguageCode] , [Language] ,[Text]) VALUES ";

            List<string> sqls = new List<string>();

            foreach(var item in query.Where(z=> z.Text.IsNotNullEmptyText()))
            {
                sqls.Add($" ('{item.LanguageCode}' , '{item.Language}' , '{item.Text}')");
            }

            insertSql = insertSql + string.Join(",", sqls);
            return insertSql;
        }


        public CRUDResultModel Create(int userID, LanguageTextUpdateRequestModel model)
        {
            model.RuleFor(x => x.LanguageCode).Rule_NotEmptyText();
            model.RuleFor(x => x.Language).Rule_NotEmptyText();
            model.RuleFor(x => x.Text).Rule_NotEmptyText();

            var queryExsit = LanguageTextRepository.GetByFilter(new QueryLanguageTextRequestViewModel
            {
                LanguageCode = model.LanguageCode,
                Language = model.Language
            });

            if (queryExsit.IsNotEmptyNull())
                throw new TerproException(ApiErrorCode.BusinessLogicDeny_DuplicateDataExist);


            var insert = model.ReturnInjecteObject<LanguageTextEntity>();

            LanguageTextRepository.Create(insert);
            return new CRUDResultModel(true);
        }


        public CRUDResultModel BatchUpdate(int userID, LanguageTextBatchUpdateRequestModel model)
        {
            model.RuleFor(x => x.LanguageCode).Rule_NotEmptyText();
                        
            var queryAllExist = LanguageTextRepository.GetByFilter(new QueryLanguageTextRequestViewModel
            {
                LanguageCode = model.LanguageCode
            });


            foreach(var item in model.Items)
            {
                item.RuleFor(z => z.Language).Rule_NotEmptyText();

                var matchExist = queryAllExist.FirstOrDefault(z => z.Language == item.Language);

                if (string.IsNullOrWhiteSpace(item.Text))
                    continue;

                if(matchExist == null)
                {
                    LanguageTextRepository.Create(new LanguageTextEntity
                    {
                        Language = item.Language,
                        Text = item.Text,
                        LanguageCode = model.LanguageCode
                    });
                }
                else
                {
                    matchExist.Text = item.Text;
                    LanguageTextRepository.Update(matchExist);
                }
            }

            //更新redis
            CacheRemoveHelper.RemoveCache_WhenUpdateLanguageCode(model.LanguageCode);
            return new CRUDResultModel(true);
        }


        /// <summary>
        /// 檢查是否有哪幾款遊戲沒建立languageText，建立預設值
        /// </summary>
        public void CreateNotExistGameNameLanguageText()
        {
            var queryGames = GameListRepository.GetByFilter(new QueryGameListRequestViewModel
            {
                IsNotDelete = true
            });

            List<string> languageCodes = queryGames
                .Where(z=> AppHelper.IsExternalVendorGame(z.GameID))
                .Where(z=> z.LanguageCode.IsNotNullEmptyText())
                .Select(z => z.LanguageCode)
                .ToList();

            List<string> queryNotExistCodes = LanguageTextRepository.GetNotExistLanguageCode(languageCodes);

            List<LanguageTextEntity> insertItems = new List<LanguageTextEntity>();

            foreach(string langCode in queryNotExistCodes)
            {
                var matchGame = queryGames.Where(z => z.LanguageCode == langCode).FirstOrDefault();

                insertItems.Add(new LanguageTextEntity
                {
                    LanguageCode = langCode,
                    Text = matchGame.GameName,
                    Language = LanguageEnum.zh_cn.ToString()
                });
            }

            LanguageTextRepository.CreateMany(insertItems);
        }


        private LanguageTextViewModel ConvertToViewModel(LanguageTextGroupModel query)
        {
            //將model轉成viewModels
            var convert = query.ReturnInjecteObject<LanguageTextViewModel>();
            return convert;
        }

        private void SetExtraProperty(List<LanguageTextViewModel> viewModels)
        {
            List<string> codes = viewModels.Select(z => z.LanguageCode).ToList();

            var queryTexts = LanguageTextRepository.GetByFilter(new QueryLanguageTextRequestViewModel
            {
                LanguageCodeList = codes
            });

            foreach(var item in viewModels)
            {
                var matchText = queryTexts.Where(z => z.LanguageCode == item.LanguageCode).ToList();

                if(matchText.IsNotEmptyNull())
                {
                    item.AllTexts = string.Join("、", matchText.Select(z => z.Text));
                    item.AllTextList = matchText.Select(z => z.Text).ToList();
                }
            }

        }

    }

}
