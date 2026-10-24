using Lck.Utility.Extensions.Helper;
using MSC.BL.Interface;
using MSC.DA.Repository;
using MSC.Models;
using MSC.Models.Attributes;
using MSC.Models.CoreLibs;
using MSC.Models.Enums;
using MSC.Models.Models;
using MSC.Models.Paging;
using MSC.Models.Table;
using MSC.Models.ViewModels;
using MSC.Utility.Extensions;
using MSC.Utility.Helper;
using target.Utility.Helper;
using System;
using System.Collections.Generic;
using System.Linq;

namespace MSC.BL.Services
{
    public class SystemSettingConfigService : BaseService, ISystemSettingConfigService
    {

        public SystemSettingConfigRepository SystemSettingConfigRepository { get; set; }



        public SystemConfigSettingPageViewModel GetGlobalSystemSettings()
        {
            var querySettings = SystemSettingConfigRepository.GetAll();

            var allEnums = ExtensionOfEnum.GetEnumAllValue<SystemSettingConfigEnum>().Select(x =>
                new SettingEnumStandardFormatConvertModel
                {
                    Group = x.GetAttributeFromEnum<DataCategoryAttribute>()?.Value ?? DataCategoryEnum.Common,
                    SettingKey = x.ToInteger(),
                    SettingKeyEnum = x
                })
                .ToList();

            var convertResult = SettingKeyEnumLogicHelper.ConvertToEditPageViewModel(allEnums);

            foreach (var item in convertResult.Groups.SelectMany(z => z.Items))
            {
                var dbSetting = querySettings.FirstOrDefault(x => x.SettingKey == item.SettingKey);

                if (dbSetting != null)
                {
                    //把db設值帶入
                    item.Value = dbSetting.SettingValue;
                }
            }

            return convertResult;
        }



        public void UpdateGlobalSystemSettings(int userID , List<SystemConfigSettingUpdateRequestItemModel> models)
        {
            var querySettings = SystemSettingConfigRepository.GetByFilter(new QuerySystemSettingConfigRequestViewModel
            {
                SettingKeyIDs = models.Select(z => z.SettingKey).ToList()
            });

            foreach (var formItem in models)
            {               
                var matchExist = querySettings.FirstOrDefault(x => x.SettingKey == formItem.SettingKey);

                if (matchExist == null)
                {
                    SystemSettingConfigRepository.Create(new SystemSettingConfigEntity
                    {
                        SettingKey = formItem.SettingKey,
                        SettingValue = formItem.Value,
                    });
                }
                else
                {
                    matchExist.SettingValue = formItem.Value;
                    SystemSettingConfigRepository.Update(matchExist);
                }

                CacheRemoveHelper.WhenUpdateSystemSettingConfig(formItem.SettingKey);
            }
        }


    }
}
