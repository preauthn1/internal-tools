using target.Models.CoreLibs;
using target.Models.Paging;
using target.Models.Table;
using target.Models.ViewModels;
using System;
using System.Collections.Generic;

namespace target.BL.Interface
{
    public interface ILanguageTextService
    {

        PagedList<LanguageTextViewModel> GetByPagingFilter(int userID, QueryLanguageTextApiRequestParameterModel filter);


        LanguageTextSingleInfoViewModel GetAllTextByLanguageCode(string languageCode);


        CRUDResultModel Create(int userID, LanguageTextUpdateRequestModel model);


        CRUDResultModel BatchUpdate(int userID, LanguageTextBatchUpdateRequestModel model);


        void CreateNotExistGameNameLanguageText();

    }
}
