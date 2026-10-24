using Dapper;
using game.Entity;
using Lck.Utility.Extensions.Helper;
using target.Core.Repository;
using target.Models.Models;
using target.Models.Paging;
using target.Models.Table;
using target.Models.ViewModels;
using target.Utility.Extensions;
using System;
using System.Collections.Generic;
using System.Linq;

namespace target.DA.Repository
{
    public class LanguageTextRepository : FormatIntegerPrimaryKeyRepository<LanguageTextEntity>
    {
        protected override string ConnectionString
        {
            get { return GetConnectionStringByDbName(DBConnType.MasterDB); }
        }

        protected override string TableName
        {
            get { return "[dbo].[LanguageText]"; }
        }

        protected override string PrimaryKeyColumnName
        {
            get { return "Id"; }
        }

        public void CreateNotExist(List<LanguageTextEntity> insertItems)
        {
            List<string> codes = insertItems.Select(z => z.LanguageCode).ToList();

            var filter = new QueryLanguageTextRequestViewModel
            {
                LanguageCodeList = codes
            };

            var queryExistDate = GetByFilter(filter);

            var groupBy = queryExistDate.GroupBy(z => z.LanguageCode).ToDictionary(x => x.First().LanguageCode, z => z.ToList());

            List<LanguageTextEntity> readyToInsert = new List<LanguageTextEntity>();

            foreach (var item in insertItems)
            {
                var isExist = groupBy.ContainsKey(item.LanguageCode);

                if (!isExist)
                {
                    readyToInsert.Add(item);
                }
                else
                {
                    List<LanguageTextEntity> match = groupBy[item.LanguageCode];
                    bool isExistLanguage = match.Any(x => x.Language == item.Language);

                    if (!isExistLanguage)
                    {
                        readyToInsert.Add(item);
                    }
                }
            }

            CreateMany(readyToInsert);
        }


        public List<LanguageTextEntity> GetByLanguageCodeByCache(string languageCode)
        {
            string cacheKey = CacheKey.DbEntity_LanguageTexts(languageCode);

            var queryCache = CacheHelper.GetCacheOrQuery(cacheKey, () =>
               {
                   var filter = new QueryLanguageTextRequestViewModel
                   {
                       LanguageCode = languageCode
                   };

                   var query = GetByFilter(filter);
                   return query;

               }, 30);

            return queryCache;
        }

        /// <summary>
        /// 查出db不存在的languageCodes
        /// </summary>
        public List<string> GetNotExistLanguageCode(List<string> languageCodeItems)
        {
            string sql = $@"SELECT distinct [LanguageCode] 
                            FROM {TableName} 
                            WHERE LanguageCode in @languageCodeItems";

            var queryExistCodes = QueryByDapper<LanguageTextBriefModel>(sql, new { languageCodeItems = languageCodeItems });

            return languageCodeItems.Where(z=> !queryExistCodes.Any(x=> x.LanguageCode == z)).ToList();
        }

        public List<LanguageTextEntity> GetByFilter(QueryLanguageTextRequestViewModel filter)
        {
            string sql = GenerateWhereSQL(filter);
            var query = QueryByDapper<LanguageTextEntity>(sql, filter);
            return query;
        }

        public PagedList<LanguageTextGroupModel> GetGroupByPagingByFilter(QueryLanguageTextRequestViewModel filter)
        {
            string whereSQL = " WHERE 1=1 ";

            if (filter.LanguageCode.IsNotNullEmptyText())
            {
                whereSQL += " AND LanguageCode = @LanguageCode ";
            }

            if (filter.LanguageCodeLikeQuery.IsNotNullEmptyText())
            {
                string likeQuery = GenerateLikeQuery("LanguageCode", "@LanguageCodeLikeQuery");
                whereSQL += $" AND {likeQuery} ";
            }

            string sql = $@"SELECT innerGroup.* FROM 
                         (
                           SELECT LanguageCode as Id , LanguageCode, count(*) as totalCount from [dbo].[LanguageText]
                           {whereSQL}
                           GROUP BY LanguageCode

                         ) AS innerGroup ";

            var query = GetPagingEntitiesBySQLScript<LanguageTextGroupModel>(sql, filter, filter);
            return query;
        }


        public PagedList<LanguageTextEntity> GetPagingByFilter(QueryLanguageTextRequestViewModel filter)
        {
            string sql = GenerateWhereSQL(filter);

            var query = GetPagingEntitiesBySQLScript<LanguageTextEntity>(sql, filter, filter);
            return query;
        }


        private string GenerateWhereSQL(QueryLanguageTextRequestViewModel filter)
        {
            string sql = $@"SELECT * FROM {TableName} 
                           WHERE 1=1 ";

            if (filter.LanguageCode.IsNotNullEmptyText())
                sql += " AND LanguageCode = @LanguageCode ";

            if (filter.LanguageCodeList.IsNotEmptyNull())
                sql += " AND LanguageCode in @LanguageCodeList ";

            if (filter.Language.IsNotEmptyNull())
                sql += " AND Language = @Language ";

            return sql;
        }



    }


}
