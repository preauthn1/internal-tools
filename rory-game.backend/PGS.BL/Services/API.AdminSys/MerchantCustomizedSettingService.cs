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
using target.Models.Attributes;
using target.Models.ViewModels.API.AdminSys;

namespace target.BL.Services
{

    public class MerchantCustomizedSettingService : BaseAdminSysService, IMerchantCustomizedSettingService
    {

        public MerchantRepository MerchantRepository { get; set; }

        public MerchantCustomizedSettingRepository MerchantCustomizedSettingRepository { get; set; }


        public SystemConfigSettingPageViewModel GetMerchantAdvanceSettings(int userID , int merchantID)
        {
            var queryMerchantSettings = MerchantCustomizedSettingRepository.GetByFilter(new QueryMerchantCustomizedSettingRequestViewModel
            {
                MerchantIDs = merchantID.ToListCollection()
            });

            var allEnums = ExtensionOfEnum.GetEnumAllValue<MerchantCustomizedSettingKeyEnum>()
                .Select(x =>
                new SettingEnumStandardFormatConvertModel
                {
                    Group = x.GetAttributeFromEnum<DataCategoryAttribute>()?.Value ?? DataCategoryEnum.Common,
                    SettingKey = x.ToInteger(),
                    SettingKeyEnum = x
                })
                .OrderBy(z=> z.Group == DataCategoryEnum.Common ? 999999 : z.Group.ToInteger())
                .ToList();

            var convertResult = SettingKeyEnumLogicHelper.ConvertToEditPageViewModel(allEnums);

            foreach (var item in convertResult.Groups.SelectMany(z => z.Items))
            {
                var dbSetting = queryMerchantSettings.FirstOrDefault(x => x.CustomizedSettingKey == item.SettingKey);

                if (dbSetting != null)
                {
                    //把db設值帶入
                    item.Value = dbSetting.SettingValue;
                }

                //特定的settingKey 只能用定義好的下拉選單來設置
                if (item.SettingKey == (int)MerchantCustomizedSettingKeyEnum.Style_LoadingStyleMode)
                {
                    item.DataDropdowns = CommonLogicHelper.GetDropdownValuesByEnumType<MerchantLoadingStyleModeEnum>();
                }
                if (item.SettingKey == (int)MerchantCustomizedSettingKeyEnum.System_ClientLogEnable)
                {
                    item.DataDropdowns = CommonLogicHelper.GetDropdownValuesByEnumType<NormalActiveStatus>();
                }
                if (SettingKeyEnumLogicHelper.IsDataTypeYesOrNoEnum((MerchantCustomizedSettingKeyEnum)item.SettingKey))
                {
                    item.DataDropdowns = CommonLogicHelper.GetDropdownValuesByEnumType<NormalYesOrNoEnum>();
                }
            }

            return convertResult;
        }


        public void UpdateSettings(int userID, SystemConfigSettingUpdateRequestModel model)
        {
            model.RuleFor(x => x.MerchantID).Rule_NotZeroNotMinus();
            model.RuleFor(x => x.Items).Rule_NotEmptyNull();

            //檢查格式合法
            SettingKeyEnumLogicHelper.CheckIsValidSettingValueInput<MerchantCustomizedSettingKeyEnum>(model.Items);

            foreach (var item in model.Items)
            {
                //迴圈編輯 (這功能不常使用 應該還好)
                var queryExist = MerchantCustomizedSettingRepository.GetByFilter(new QueryMerchantCustomizedSettingRequestViewModel
                    {
                        MerchantIDs = model.MerchantID.ToListCollection(),
                        CustomizedSettingKeyIDs = item.SettingKey.ToListCollection()
                    })
                    .FirstOrDefault();   

                if (queryExist == null)
                {
                    MerchantCustomizedSettingRepository.Create(new MerchantCustomizedSettingEntity
                    {
                        CustomizedSettingKey = item.SettingKey,
                        MerchantID = model.MerchantID,
                        CreateTime = DateTimeHelper.GetCurrentTime(),
                        SettingValue = item.Value
                    });
                }
                else
                {
                    queryExist.SettingValue = item.Value;
                    MerchantCustomizedSettingRepository.Update(queryExist);
                }

                CacheRemoveHelper.RemoveCache_WhenUpdateMerchantSetting(model.MerchantID, item.SettingKey);
            }
        }



    }
}
