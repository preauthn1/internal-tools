
using MSC.DA.Repository;

namespace MSC.Core.Repository
{
    /// <summary>
    /// YǸRepositoryd/@άd iHgbo
    /// </summary>
    public class RepositoryHelper
    {
        public static string GetUserNameByUserID(int userID)
        {
            var rep = new UserRepository();
            return rep.GetUserNameByUserID(userID);
        }

    }
}
