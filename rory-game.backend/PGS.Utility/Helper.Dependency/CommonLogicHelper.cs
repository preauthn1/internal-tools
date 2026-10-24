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

namespace target.Utility.Helper
{
    public class CommonLogicHelper
    {
        public static bool IsSingleWalletMode(MerchantWalletModeEnum walletMode)
        {
            return walletMode == MerchantWalletModeEnum.SingleWalletMode;
        }

        public static bool IsSingleWalletMode(int walletMode)
        {
            return IsSingleWalletMode((MerchantWalletModeEnum)walletMode);
        }

        public static decimal CalculateRTP_Percentage(decimal betAmount , decimal payoutAmount)
        {
            return Math.Round(betAmount == 0 ? 0 : payoutAmount / betAmount * 100 , 2);
        }


        public static List<DropdownValueViewModel> GetDropdownValuesByEnumType<T>() where T : IComparable, IFormattable, IConvertible
        {
            var allEnums = ExtensionOfEnum.GetEnumAllValue<T>();

            var dropdowns = allEnums.Select(x =>
            {
                dynamic convert = x;
                return new DropdownValueViewModel
                {
                    Value = ((int)convert).ToString(),
                    Text_CN = MultiLanguageTextHelper.GetEnumsCultureText(convert)
                };
            })
                .ToList();

            //這種enums 會排除部分選項
            if (typeof(T) == typeof(NormalActiveStatus))
            {
                return dropdowns
                    .Where(x => x.Value != NormalActiveStatus.PermanentDelete.ToIntegerString())
                    .ToList();
            }

            return dropdowns;

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

    }
}
