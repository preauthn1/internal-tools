using target.Models;
using System;
using System.Collections.Generic;
using System.Text;
using target.Utility.Extensions;
using target.Models.Enums;
using target.Models.Attributes;

namespace target.Utility.Helper
{
    public class DataModelConvertHelper
    {

        public static double GetTimeZoneByAreaID(int timeZoneArea)
        {
            TimeZoneAreaEnum timeZone = (TimeZoneAreaEnum)timeZoneArea;

            if (!timeZone.IsDefinedValue())
                timeZone = TimeZoneAreaEnum.China_0800;

            return GetTimeZoneByEnum(timeZone);
        }
        

        public static double GetTimeZoneByEnum(TimeZoneAreaEnum timeZoneArea)
        {
            try
            {
                string timeZoneValue = timeZoneArea.GetAttributeFromEnum<ExtraAttributeValue>()?.Value;
                return double.Parse(timeZoneValue);
            }
            catch
            {

            }

            return 8.0;
        }

    }
}
