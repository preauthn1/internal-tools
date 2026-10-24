using System;
using System.Collections.Generic;
using System.Text;
using System.Linq;
using MSC.Models.ViewModels;
using MSC.Utility.Helper;
using MSC.Utility.Extensions;
using MSC.Models.Models;
using MSC.Models.Attributes;

namespace target.Utility.Helper
{
    public class SettingKeyEnumLogicHelper
    {


        public static SystemConfigSettingPageViewModel ConvertToEditPageViewModel(List<SettingEnumStandardFormatConvertModel> settings)
        {
            var allCategory = settings.Select(z=> z.Group).Distinct().ToList();

            List<SystemConfigSettingGroupViewModel> groups = new List<SystemConfigSettingGroupViewModel>();

            //編輯介面是以分類來批次編輯的
            foreach (var eachGroup in allCategory)
            {
                string groupName = MultiLanguageTextHelper.GetEnums_i18nText(eachGroup);

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
                SettingKeyDisplayName = MultiLanguageTextHelper.GetEnums_i18nText(enumValue),
                DefineType = GetReferenceDataTypeFromEnum(enumValue),
                IsDisableEdit = !IsEdtiEnable(enumValue),
                DefaultValue = GetDefinedDefaultValue(enumValue)
            };

            newInstance.DataTypeName = newInstance.DefineType.TypeName;
            newInstance.Value = newInstance.DefaultValue;
            return newInstance;
        }



        private static DefineTypeCheckResultModel GetReferenceDataTypeFromEnum(System.Enum enumValue)
        {
            var attribute = enumValue.GetAttributeFromEnum<ReferenceType>();
            var defineType = GetReferenceDataTypeFromEnum(attribute);
            return defineType;
        }


        public static DefineTypeCheckResultModel GetReferenceDataTypeFromEnum(ReferenceType attribute)
        {
            Type defineType = typeof(string);

            if (attribute != null)
                defineType = attribute.Type;

            string typeName = defineType.Name.ToLower();

            if (typeName.Contains("enum"))
            {
                return new DefineTypeCheckResultModel
                {
                    Type = defineType,
                    IsEnums = true,
                    TypeName = typeName
                };
            }

            return new DefineTypeCheckResultModel
            {
                Type = defineType,
                IsEnums = false,
                TypeName = typeName
            };
        }

        public static string GetDefinedDefaultValue(System.Enum parseToEnum)
        {
            var attr = parseToEnum.GetAttributeFromEnum<SettingDefaultValue>();

            if (attr == null)
                return "";

            return attr.Value;
        }


        public static bool IsEdtiEnable(System.Enum parseToEnum)
        {
            var attr = parseToEnum.GetAttributeFromEnum<IsDisabledEditAttribute>();

            if (attr == null)
                return true;

            return false;
        }


        public static void CheckIsValidSettingValueInput<T>(List<SystemConfigSettingUpdateRequestItemModel> settingItems) where T: System.Enum
        {
            foreach (var item in settingItems)
            {
                item.RuleFor(x => x.SettingKey).Rule_NotZeroNotMinus();

                System.Enum settingKeyEnum = (T)(object)item.SettingKey;

                if (!IsEdtiEnable(settingKeyEnum))
                {
                    throw new MSCException(ErrorCode.Forbidden_DoNotAllowEdit);
                }

                if (!ExtensionOfEnum.IsDefinedValue(settingKeyEnum))
                {
                    throw new MSCException(ErrorCode.Forbidden_NotValidDefineValue);
                }

                var dataTypeModel = InstanceSettingKeyItemViewModel(settingKeyEnum);

                //檢查資料型別是否合法
                if (dataTypeModel.DefineType.IsEnums)
                {
                    //若為列舉、只吃數字格式
                    if (!item.Value.IsNumeric())
                    {
                        throw new MSCException(ErrorCode.InvalidParameter_NumberEmptyZero);
                    }
                }
                else
                {
                    //若為數字格式
                    if (dataTypeModel.DefineType.Type == typeof(decimal))
                    {
                        if (!item.Value.IsNumeric())
                        {
                            throw new MSCException(ErrorCode.InvalidParameter_ParsingError);
                        }
                    }
                }
            }
        }



    }
}
