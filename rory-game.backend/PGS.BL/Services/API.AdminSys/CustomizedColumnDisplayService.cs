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
using target.Models.ViewModels.API.AdminSys;
using Microsoft.Extensions.Caching.Memory;
using target.Models.Attributes;

namespace target.BL.Services
{

    public class CustomizedColumnDisplayService : BaseAdminSysService, ICustomizedColumnDisplayService
    {
        public MerchantRepository MerchantRepository { get; set; }

        public DynamicQueryCommandRepository DynamicQueryCommandRepository { get; set; }


        public DynamicQueryCommandViewModel GetDynamicQueryCommandSetting(int querySqlID)
        {           
            var querySqlItem = DynamicQueryCommandRepository.GetById(querySqlID);

            if (querySqlItem == null)
                throw new TerproException(ApiErrorCode.NotFoundData);

            var viewModel = querySqlItem.ReturnInjecteObject<DynamicQueryCommandViewModel>();

            if(!string.IsNullOrWhiteSpace(querySqlItem.FilterFormItems))
            {
                viewModel.AvailableFilterFormItems = querySqlItem.FilterFormItems.Split(',').ToList();
            }

            return viewModel;
        }



        /// <summary>
        /// 為了結果麻煩 這種類型的懶人查詢 不提供分頁
        /// 請控制查詢結果在 500筆以內
        /// </summary>
        public CustomizedColumnDynamicQueryResultViewModel GetDynamicQueryResult(int userID, CustomizedColumnDisplayRequestModel filter)
        {
            if(filter.BeginTime.HasValue && filter.EndTime.HasValue)
            {
                if((filter.EndTime.Value - filter.BeginTime.Value).TotalDays >= 62)
                {
                    throw new TerproException(ApiErrorCode.InvalidParameter_DateTimeRangeTooLarge);
                }
            }

            var querySqlItem = DynamicQueryCommandRepository.GetById(filter.QuerySqlID);

            if (querySqlItem == null)
                throw new TerproException(ApiErrorCode.NotFoundData);

            CheckRequiredColumns(querySqlItem , filter);

            var queryResult = DynamicQueryCommandRepository.QueryByTargetSql(querySqlItem.Sql, filter);

            if (queryResult.IsEmptyOrNull())
            {
                return new CustomizedColumnDynamicQueryResultViewModel
                {
                    Items = new List<Dictionary<string, object>>(),
                    ExtraObject = new CustomizedColumnExtraInfoViewModel
                    {
                        Columns = new List<CustomizedColumnDisplayColumnViewModel>()
                    }
                };
            }

            var firstItem = queryResult.FirstOrDefault();

            var columns = firstItem
                    .Select(z => z.Key)
                    .Select(z => InstanceDisplayColumnModel(z))
                    .ToList();

            string excelText = string.Empty;

            foreach (string column in columns.Select(z => z.ColumnName))
            {
                excelText += column + "\t";
            }

            excelText += Environment.NewLine;

            foreach (var item in queryResult)
            {
                foreach (string column in columns.Select(z => z.ColumnName))
                {
                    excelText += item[column] + "\t";
                }
                excelText += Environment.NewLine;
            }

            FormattingColumnValue(queryResult);

            return new CustomizedColumnDynamicQueryResultViewModel
            {
                Items = queryResult,
                ExtraObject = new CustomizedColumnExtraInfoViewModel
                {
                    Columns = columns,
                    CopyExcel = excelText
                }
            };
        }

        private static void FormattingColumnValue(List<Dictionary<string, object>> queryResult)
        {
            foreach (var eachRow in queryResult)
            {
                Dictionary<string, object> modifyValues = new Dictionary<string, object>();

                foreach (var eachColumn in eachRow)
                {
                    bool isAmount = eachColumn.Key.ToLower().Contains("amount");
                    bool isTime = eachColumn.Key.ToLower().Contains("time");

                    //暫存記住哪幾個欄位要被異動
                    if (isAmount)
                    {
                        modifyValues.Add(eachColumn.Key, ((decimal)eachColumn.Value).ToMoneyFormat());
                    }

                    if (isTime)
                    {
                        //反正是自己人用 強制轉+8
                        modifyValues.Add(eachColumn.Key, ((DateTime)eachColumn.Value).AddHours(8).ToString("yyyy-MM-dd HH:mm:ss"));
                    }
                }

                // 將amount欄位轉千分位符號
                foreach (var newValue in modifyValues)
                {
                    eachRow[newValue.Key] = newValue.Value;
                }
            }
        }

        private CustomizedColumnDisplayColumnViewModel InstanceDisplayColumnModel(string nameOfColumn)
        {
            return new CustomizedColumnDisplayColumnViewModel(nameOfColumn);
        }


        private void CheckRequiredColumns(DynamicQueryCommandEntity settingItem , CustomizedColumnDisplayRequestModel request)
        {
            if(settingItem.RequiredColumns.IsNullOrWhiteSpace())            
                return;            

            var columns = settingItem.RequiredColumns.Split(',').ToList();

            if (columns.IsEmptyOrNull())
                return;

            //檢查必填欄位是否有填寫
            foreach(string columnName in columns)
            {
                string columnNameLower = columnName.ToLower();

                foreach(var prop in request.GetType().GetProperties())
                {
                    string propName = prop.Name;
                    if (prop.Name.ToLower() == columnNameLower)
                    {
                        var value = prop.GetValue(request);

                        if(value == null || string.IsNullOrWhiteSpace(value.ToString()))
                        {
                            throw new TerproException(ApiErrorCode.InvalidParameter, new ErrorDetailModel(propName));
                        }
                    }
                }
            }
        }

    }
}
