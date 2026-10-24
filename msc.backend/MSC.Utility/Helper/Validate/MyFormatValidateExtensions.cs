using System;
using System.Collections.Generic;
using System.Security.Claims;
using System.Text;
using System.Linq;
using System.Security.Cryptography;
using MSC.Utility.Extensions;
using System.Linq.Expressions;
using System.Reflection;
using MSC.Models.Models;
using MSC.Utility.Helper;

namespace MSC.Utility.Extensions
{
    public static class MyFormatValidateExtensions
    {
        public static RuleModel<T, T2> RuleFor<T, T2>(this T source, Expression<Func<T, T2>> selector) where T : class
        {
            //這個擴充方法 禁止傳入null物件
            if (source == null)
                throw new MSCException(ErrorCode.InvalidParameter_ModelIsNull);

            var name = ((MemberExpression)selector.Body).Member?.Name;
            return new RuleModel<T, T2>
            {
                PropertyValue = (T2)source.GetType().GetProperty(name).GetValue(source),
                PropertyName = name
            };
        }
        

        public static RuleModel<T, DateTime> Rule_NotEmptyDateTime<T>(this RuleModel<T, DateTime> model)
        {
            if(model.PropertyValue == new DateTime())
                throw new MSCException(ErrorCode.InvalidParameter_DateTimeRequired, new ErrorDetailModel(model.PropertyName));

            return model;
        }

        public static RuleModel<T, DateTime> Rule_StartTimeEndTimeRange<T>(this RuleModel<T, DateTime> model , DateTime endTime , int daysRange)
        {
            if (model.PropertyValue > endTime)
                throw new MSCException(ErrorCode.InvalidParameter_StartTimeMustSmallThanEndTime);

            if ((endTime - model.PropertyValue).TotalDays > daysRange)
                throw new MSCException(ErrorCode.InvalidParameter_DateTimeRangeTooLarge);

            return model;
        }

        public static RuleModel<T, DateTime> Rule_DateTimeMustLargeThanPreviousDays<T>(this RuleModel<T, DateTime> model, int daysRange)
        {
            DateTime utcNow = DateTimeHelper.GetCurrentTime();
            DateTime compareDate = utcNow.AddDays(-1 * daysRange);

            if (compareDate > model.PropertyValue)
                throw new MSCException(ErrorCode.InvalidParameter_DateTimeRangeTooLarge);

            return model;
        }


        public static RuleModel<T, string> Rule_NotEmptyText<T>(this RuleModel<T, string> model)
        {
            if (string.IsNullOrWhiteSpace(model.PropertyValue))
                throw new MSCException(ErrorCode.InvalidParameter_TextEmpty , new ErrorDetailModel(model.PropertyName));

            return model;
        }
        
        public static RuleModel<T, string> Rule_IsSupportLanguageCode<T>(this RuleModel<T, string> model)
        {
            if (!RegexRuleHelper.IsSupportLanguage(model.PropertyValue))
                throw new MSCException(ErrorCode.InvalidParameter_NotSupportLanguageCode, new ErrorDetailModel(model.PropertyName));

            return model;
        }

        public static RuleModel<T, string> Rule_IsLetterAndNumberOnly<T>(this RuleModel<T, string> model)
        {
            if (!RegexRuleHelper.IsLetterAndNumberOnly(model.PropertyValue))
                throw new MSCException(ErrorCode.InvalidParameter_InvalidLetter, new ErrorDetailModel(model.PropertyName));

            return model;
        }

        public static RuleModel<T, string> Rule_IsNumeric<T>(this RuleModel<T, string> model)
        {
            if (!model.PropertyValue.IsNumeric())
                throw new MSCException(ErrorCode.InvalidParameter_InvalidLetter, new ErrorDetailModel(model.PropertyName));

            return model;
        }

        public static RuleModel<T, string> Rule_IsValidAccountFormat<T>(this RuleModel<T, string> model)
        {
            if (!RegexRuleHelper.IsAccountFormat(model.PropertyValue))
                throw new MSCException(ErrorCode.InvalidParameter_NotAccountFormat, new ErrorDetailModel(model.PropertyName));

            return model;
        }


        public static RuleModel<T, string> Rule_IsValidPasswordFormat<T>(this RuleModel<T, string> model)
        {
            if (!RegexRuleHelper.IsPasswordFormat(model.PropertyValue))
                throw new MSCException(ErrorCode.InvalidParameter_NotAvailablePassword, new ErrorDetailModel(model.PropertyName));

            return model;
        }

        public static RuleModel<T, List<string>> Rule_IsNotEmptyCollection<T>(this RuleModel<T, List<string>> model)
        {
            if (model.PropertyValue.IsEmptyOrNull())
                throw new MSCException(ErrorCode.InvalidParameter_NumberValueTooSmall, new ErrorDetailModel(model.PropertyName));

            return model;
        }

        public static RuleModel<T, List<int>> Rule_IsNotEmptyCollection<T>(this RuleModel<T, List<int>> model)
        {
            if (model.PropertyValue.IsEmptyOrNull())
                throw new MSCException(ErrorCode.InvalidParameter_NumberValueTooSmall, new ErrorDetailModel(model.PropertyName));

            return model;
        }

        public static RuleModel<T, int?> Rule_Rule_IsNotNull<T>(this RuleModel<T, int?> model)
        {
            if (model.PropertyValue == null)
                throw new MSCException(ErrorCode.InvalidParameter_NumberValueNotAllowNull, new ErrorDetailModel(model.PropertyName));

            return model;
        }


        public static RuleModel<T, string> Rule_MinLengthEnough<T>(this RuleModel<T, string> model, int minLength)
        {
            Rule_NotEmptyText(model);

            if(model.PropertyValue.Length < minLength)
                throw new MSCException(ErrorCode.InvalidParameter_TextLength, new ErrorDetailModel(model.PropertyName));

            return model;
        }

        public static RuleModel<T, string> Rule_MaxLengthOK<T>(this RuleModel<T, string> model, int maxLength)
        {
            model.PropertyValue = model.PropertyValue ?? string.Empty;

            if (model.PropertyValue.Length > maxLength)
                throw new MSCException(ErrorCode.InvalidParameter_TextLength, new ErrorDetailModel(model.PropertyName));

            return model;
        }


        public static RuleModel<T, int> Rule_NotZero<T>(this RuleModel<T, int> model) where T : class
        {
            if(model.PropertyValue == 0)
                throw new MSCException(ErrorCode.InvalidParameter_NumberEmptyZero, new ErrorDetailModel(model.PropertyName));

            return model;
        }

        public static RuleModel<T, int> Rule_NotZeroNotMinus<T>(this RuleModel<T, int> model) where T : class
        {
            if (model.PropertyValue == 0)
                throw new MSCException(ErrorCode.InvalidParameter_NumberEmptyZero, new ErrorDetailModel(model.PropertyName));

            if (model.PropertyValue < 0)
                throw new MSCException(ErrorCode.InvalidParameter_NumberMinus, new ErrorDetailModel(model.PropertyName));

            return model;
        }

        public static RuleModel<T, long> Rule_NotZeroNotMinus<T>(this RuleModel<T, long> model) where T : class
        {
            if (model.PropertyValue == 0)
                throw new MSCException(ErrorCode.InvalidParameter_NumberEmptyZero, new ErrorDetailModel(model.PropertyName));

            if (model.PropertyValue < 0)
                throw new MSCException(ErrorCode.InvalidParameter_NumberMinus, new ErrorDetailModel(model.PropertyName));

            return model;
        }


        public static RuleModel<T, long> Rule_IsValidTradeID<T>(this RuleModel<T, long> model) where T : class
        {
            if (model.PropertyValue < 2000000000000000)
                throw new MSCException(ErrorCode.InvalidParameter_NumberValueTooSmall);

            return model;
        }


        public static RuleModel<T, long> Rule_IsValidTimestamp<T>(this RuleModel<T, long> model) where T : class
        {
            DateTime nowUTC = DateTimeHelper.GetCurrentTime();

            var ts_5_min = nowUTC.AddMinutes(-5).CovertToUnixTimestamp();
            var ts_10_min_future = nowUTC.AddMinutes(10).CovertToUnixTimestamp();

            //時間太久之前: 不合法
            if (model.PropertyValue < ts_5_min)
                throw new MSCException(ErrorCode.InvalidParameter_TimestampError);

            //時間是未來日期: 不合法
            if (model.PropertyValue > ts_10_min_future)
                throw new MSCException(ErrorCode.InvalidParameter_TimestampError);

            return model;
        }


        public static RuleModel<T, decimal> Rule_NotZero<T>(this RuleModel<T, decimal> model) where T : class
        {
            if (model.PropertyValue == 0)
                throw new MSCException(ErrorCode.InvalidParameter_NumberEmptyZero, new ErrorDetailModel(model.PropertyName));

            return model;
        }

        public static RuleModel<T, decimal> Rule_IsIntegerValue<T>(this RuleModel<T, decimal> model) where T : class
        {
            if (!model.PropertyValue.IsIntegerValue())
                throw new MSCException(ErrorCode.InvalidParameter_NumberValueNotAllowFloat, new ErrorDetailModel(model.PropertyName));

            return model;
        }

        public static RuleModel<T, decimal> Rule_NotTooSmall<T>(this RuleModel<T, decimal> model , decimal minValue) where T : class
        {
            if (minValue > model.PropertyValue)
                throw new MSCException(ErrorCode.InvalidParameter_NumberValueTooSmall, new ErrorDetailModel(model.PropertyName));

            return model;
        }

        public static RuleModel<T, decimal> Rule_NotTooLarge<T>(this RuleModel<T, decimal> model, decimal maxValue) where T : class
        {
            if (maxValue < model.PropertyValue)
                throw new MSCException(ErrorCode.InvalidParameter_NumberValueTooLarge, new ErrorDetailModel(model.PropertyName));

            return model;
        }

        public static RuleModel<T, int> Rule_NotTooLarge<T>(this RuleModel<T, int> model, int maxValue) where T : class
        {
            if (maxValue < model.PropertyValue)
                throw new MSCException(ErrorCode.InvalidParameter_NumberValueTooLarge, new ErrorDetailModel(model.PropertyName));

            return model;
        }

        public static RuleModel<T, decimal> Rule_NotZeroNotMinus<T>(this RuleModel<T, decimal> model) where T : class
        {
            if (model.PropertyValue == 0)
                throw new MSCException(ErrorCode.InvalidParameter_NumberEmptyZero, new ErrorDetailModel(model.PropertyName));

            if (model.PropertyValue < 0)
                throw new MSCException(ErrorCode.InvalidParameter_NumberMinus, new ErrorDetailModel(model.PropertyName));

            return model;
        }


        public static RuleModel<T, decimal> Rule_NotMinus<T>(this RuleModel<T, decimal> model) where T : class
        {
            if (model.PropertyValue < 0)
                throw new MSCException(ErrorCode.InvalidParameter_NumberMinus, new ErrorDetailModel(model.PropertyName));

            return model;
        }
        
        public static RuleModel<T, int> Rule_NotMinus<T>(this RuleModel<T, int> model) where T : class
        {
            if (model.PropertyValue < 0)
                throw new MSCException(ErrorCode.InvalidParameter_NumberMinus, new ErrorDetailModel(model.PropertyName));

            return model;
        }


        public static RuleModel<T, T2> Rule_ModelNotNull<T , T2>(this RuleModel<T, T2> model) 
            where T : class 
            where T2 : class
        {
            if (model.PropertyValue == null)
                throw new MSCException(ErrorCode.InvalidParameter_ModelIsNull, new ErrorDetailModel(model.PropertyName));

            return model;
        }

    }





}
