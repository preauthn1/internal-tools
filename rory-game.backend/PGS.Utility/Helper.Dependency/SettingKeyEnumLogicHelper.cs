using target.Models;
using System;
using System.Collections.Generic;
using System.Text;
using target.Utility.Extensions;
using target.Models.Models;
using target.Models.Attributes;
using target.Models.ViewModels.API.AdminSys;
using target.Models.Enums;
using System.Linq;
using target.Models.Table;
using target.Models.ViewModels;

namespace target.Utility.Helper
{
    public class SettingKeyEnumLogicHelper
    {

        /// <summary>
        /// 根據type不同 設置 MerchantCustomizedSettingKeyEnum 各自的預設值
        /// </summary>
        public static MerchantCustomizedSettingModel GetMerchantSettingOrDefaultValueByEntity(MerchantCustomizedSettingKeyEnum type, MerchantCustomizedSettingEntity setting)
        {
            if (setting == null)
                return GetMerchantSettingOrDefaultValue(type, null);

            var convert = new MerchantCustomizedSettingModel
            {
                MerchantID = setting.MerchantID,
                SettingValue = setting.SettingValue
            };

            return GetMerchantSettingOrDefaultValue(type, convert);
        }


        /// <summary>
        /// 根據type不同 設置 MerchantCustomizedSettingKeyEnum 各自的預設值
        /// </summary>
        public static MerchantCustomizedSettingModel GetMerchantSettingOrDefaultValue(MerchantCustomizedSettingKeyEnum type, MerchantCustomizedSettingModel setting)
        {
            if (setting != null)
                return setting;

            setting = new MerchantCustomizedSettingModel();
            setting.CustomizedSettingKey = type.ToInteger();
            setting.SettingValue = GetMerchantSettingDefaultValueByType(type);
            return setting;
        }


        public static string GetMerchantSettingDefaultValueByType(MerchantCustomizedSettingKeyEnum type)
        {
            //根據type不同 塞預設值
            switch (type)
            {
                case MerchantCustomizedSettingKeyEnum.DefaultWalletBalance:
                    return "0";

                case MerchantCustomizedSettingKeyEnum.Style_LoadingStyleMode:
                    return MerchantLoadingStyleModeEnum.Default.ToIntegerString();

                case MerchantCustomizedSettingKeyEnum.System_ClientLogEnable:
                    return "0";

                case MerchantCustomizedSettingKeyEnum.System_IsAlwaysTrialMode:
                    return NormalYesOrNoEnum.No.ToIntegerString();

                default:
                    break;
            }

            return "";
        }


        public static SystemConfigSettingPageViewModel ConvertToEditPageViewModel(List<SettingEnumStandardFormatConvertModel> settings)
        {
            var allCategory = settings.Select(z=> z.Group).Distinct().ToList();

            List<SystemConfigSettingGroupViewModel> groups = new List<SystemConfigSettingGroupViewModel>();

            //編輯介面是以分類來批次編輯的
            foreach (var eachGroup in allCategory)
            {
                string groupName = MultiLanguageTextHelper.GetEnumsCultureText(eachGroup);

                var matchSettingsKeys = settings.Where(x => x.Group == eachGroup).ToList();

                var myGroup = new SystemConfigSettingGroupViewModel
                {
                    GroupID = eachGroup.ToInteger(),
                    GroupName = groupName,
                    Items = new List<SystemConfigSettingViewModel>()
                };

                foreach (var item in matchSettingsKeys)
                {
                    var newInstance = InstanceSettingKeyItemViewModel(item.SettingKeyEnum);
                    myGroup.Items.Add(newInstance);
                }

                myGroup.IsAllDisableEdit = myGroup.Items.All(x => x.IsDisableEdit);
                groups.Add(myGroup);
            }

            var viewModel = new SystemConfigSettingPageViewModel
            {
                Groups = groups
            };

            return viewModel;
        }


        public static SystemConfigSettingViewModel InstanceSettingKeyItemViewModel(System.Enum enumValue)
        {
            var newInstance = new SystemConfigSettingViewModel
            {
                SettingKey = enumValue.ToInteger(),
                SettingKeyDisplayName = MultiLanguageTextHelper.GetEnumsCultureText(enumValue),
                DefineType = GetReferenceDataTypeFromEnum(enumValue),
                IsDisableEdit = !IsEdtiAttributeEnable(enumValue),
                Value = enumValue.GetAttributeFromEnum<SettingDefaultValue>()?.Value
            };

            newInstance.DataTypeName = newInstance.DefineType.TypeName;
            return newInstance;
        }



        private static DefineTypeCheckResultModel GetReferenceDataTypeFromEnum(System.Enum enumValue)
        {
            var attribute = enumValue.GetAttributeFromEnum<ReferenceType>();
            var defineType = CommonLogicHelper.GetReferenceDataTypeFromEnum(attribute);
            return defineType;
        }


        public static bool IsEdtiAttributeEnable(System.Enum parseToEnum)
        {
            var attr = parseToEnum.GetAttributeFromEnum<DisableEdit>();

            if (attr == null)
                return true;

            return false;
        }

        public static bool IsDataTypeYesOrNoEnum(System.Enum parseToEnum)
        {
            var attr = parseToEnum.GetAttributeFromEnum<ReferenceType>();

            if (attr == null)
                return false;

            return attr.Type == typeof(NormalYesOrNoEnum);
        }


        public static void CheckIsValidSettingValueInput<T>(List<SystemConfigSettingUpdateRequestItemModel> settingItems) where T: System.Enum
        {
            foreach (var item in settingItems)
            {
                item.RuleFor(x => x.SettingKey).Rule_NotZeroNotMinus();

                System.Enum settingKeyEnum = (T)(object)item.SettingKey;

                if (!IsEdtiAttributeEnable(settingKeyEnum))
                {
                    throw new TerproException(ApiErrorCode.Forbidden_DoNotAllowEdit);
                }

                if (!ExtensionOfEnum.IsDefinedValue(settingKeyEnum))
                {
                    throw new TerproException(ApiErrorCode.Forbidden_NotValidDefineValue);
                }

                var dataTypeModel = InstanceSettingKeyItemViewModel(settingKeyEnum);

                //檢查資料型別是否合法
                if (dataTypeModel.DefineType.IsEnums)
                {
                    //若為列舉、只吃數字格式
                    if (!item.Value.IsNumeric())
                    {
                        throw new TerproException(ApiErrorCode.InvalidParameter_NumberEmptyZero);
                    }
                }
                else
                {
                    //若為數字格式
                    if (dataTypeModel.DefineType.Type == typeof(decimal))
                    {
                        if (!item.Value.IsNumeric())
                        {
                            throw new TerproException(ApiErrorCode.InvalidParameter_ParsingError);
                        }
                    }
                }
            }
        }

    }
}
