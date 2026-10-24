using System;
using System.Text.Json;
using System.Text.Json.Serialization;
using System.Text.RegularExpressions;
using target.Models.Extensions;

namespace target.Utility.Helper
{
    public class ApiAlwaysTimestampConverter : JsonConverter<DateTime>
    {
        public override DateTime Read(ref Utf8JsonReader reader, Type typeToConvert, JsonSerializerOptions options)
        {
            string myValue = reader.GetString();

            if (myValue.IsNumeric())
            {
                return Double.Parse(myValue).UnixTimeStampToDateTime();
            }

            return DateTime.Parse(reader.GetString());
        }

        public override void Write(Utf8JsonWriter writer, DateTime value, JsonSerializerOptions options)
        {
            writer.WriteStringValue(value.ToUniversalTime().ToString("yyyy-MM-ddTHH:mm:ssZ")); //IOS-8601 +0 UTC
            //writer.WriteStringValue(value.ToUniversalTime().CovertToUnixTimestamp().ToString());  //timestamp
        }
    }


}
