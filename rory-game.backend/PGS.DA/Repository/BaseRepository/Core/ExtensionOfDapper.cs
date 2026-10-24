using Dapper;

namespace target.Core.Repository
{
    public static class ExtensionOfDapper
    {
        public static int GetDynamicReturnCode(this DynamicParameters para)
        {
            int returnCode = para.Get<int>("@RetVal");
            return returnCode;
        }

        public static T GetOutputVariable<T>(this DynamicParameters para , string parameterName)
        {
            var result = para.Get<T>(parameterName);
            return result;
        }

    }

}
