using target.Models.CultureText;
using System.ComponentModel;

namespace target.Models.Models
{
    public enum ApiErrorCode
    {       
        Success = 0,
        SuccessOK = 1,

        [Description(i18nKeys.ErrorCode_ServerPlayerOutOfRange)]
        ServerPlayerOutOfRange = 400201,


        //401001~401099 會給前端拿來判斷強制導到登入頁，請謹慎使用
        //401001~401099 會給前端拿來判斷強制導到登入頁，請謹慎使用
        //401001~401099 會給前端拿來判斷強制導到登入頁，請謹慎使用
        //401001~401099 會給前端拿來判斷強制導到登入頁，請謹慎使用
        //401001~401099 會給前端拿來判斷強制導到登入頁，請謹慎使用
        //401001~401099 會給前端拿來判斷強制導到登入頁，請謹慎使用

        #region 401 身分授權類
        [Description(i18nKeys.ErrorCode_Authorize_AccessToken_NotValid)]
        Authorize_AccessToken_NotValid = 401001,

        [Description(i18nKeys.ErrorCode_Authorize_AccessToken_Expired)]
        Authorize_AccessToken_Expired = 401002,

        [Description(i18nKeys.ErrorCode_Authorize_TokenHasChanged)]
        Authorize_TokenHasChanged = 401003,

        [Description(i18nKeys.ErrorCode_Authorize_JwtTokenParsingFail)]
        Authorize_JwtTokenParsingFail = 401008,

        [Description(i18nKeys.ErrorCode_Authorize_JwtTokenVersionNotMatch)]
        Authorize_JwtTokenVersionNotMatch = 401009,

        #endregion



        #region 403 其他邏輯錯誤
        [Description(i18nKeys.ErrorCode_InvalidParameter)]
        InvalidParameter = 403001,

        [Description(i18nKeys.ErrorCode_InvalidParameter_TextEmpty)]
        InvalidParameter_TextEmpty = 403002,

        [Description(i18nKeys.ErrorCode_InvalidParameter_TextLength)]
        InvalidParameter_TextLength = 403003, 
        
        [Description(i18nKeys.ErrorCode_InvalidParameter_NotAccountFormat)]
        InvalidParameter_NotAccountFormat = 403004,

        [Description(i18nKeys.ErrorCode_InvalidParameter_NotSupportLanguageCodet)]
        InvalidParameter_NotSupportLanguageCodet = 403005,

        [Description(i18nKeys.ErrorCode_InvalidParameter_ConfirmCodeNotMatch)]
        InvalidParameter_ConfirmCodeNotMatch = 403009,

        [Description(i18nKeys.ErrorCode_InvalidParameter_NumberEmptyZero)]
        InvalidParameter_NumberEmptyZero = 403011,

        [Description(i18nKeys.ErrorCode_InvalidParameter_NumberMinus)]
        InvalidParameter_NumberMinus = 403012,

        [Description(i18nKeys.ErrorCode_InvalidParameter_NumberValueTooSmall)]
        InvalidParameter_NumberValueTooSmall = 403013,

        [Description(i18nKeys.ErrorCode_InvalidParameter_NumberValueTooLarge)]
        InvalidParameter_NumberValueTooLarge = 403014,

        [Description(i18nKeys.ErrorCode_InvalidParameter_NumberValueNotAllowFloat)]
        InvalidParameter_NumberValueNotAllowFloat = 403015,

        [Description(i18nKeys.ErrorCode_InvalidParameter_NumberValueNotAllowFloat)]
        InvalidParameter_NumberValueMustPlacesTwo = 403016,

        [Description(i18nKeys.ErrorCode_InvalidParameter_NumberValueNotAllowNull)]
        InvalidParameter_NumberValueNotAllowNull = 403018,

        [Description(i18nKeys.ErrorCode_InvalidParameter_Collection)]
        InvalidParameter_Collection = 403019,

        [Description(i18nKeys.ErrorCode_InvalidParameter_DateTimeRequired)]
        InvalidParameter_DateTimeRequired = 403021,

        [Description(i18nKeys.ErrorCode_InvalidParameter_StartTimeMustSmallThanEndTime)]
        InvalidParameter_StartTimeMustSmallThanEndTime = 403022,

        [Description(i18nKeys.ErrorCode_InvalidParameter_DateTimeRangeTooLarge)]
        InvalidParameter_DateTimeRangeTooLarge = 403023,

        [Description(i18nKeys.ErrorCode_InvalidParameter_DateTimeTooSmall)]
        InvalidParameter_DateTimeTooSmall = 403024,

        [Description(i18nKeys.ErrorCode_InvalidParameter_ParsingError)]
        InvalidParameter_ParsingError = 403059,

        [Description(i18nKeys.ErrorCode_InvalidParameter_TimestampError)]
        InvalidParameter_TimestampError = 403070,

        [Description(i18nKeys.ErrorCode_InvalidParameter_SignKeyError)]
        InvalidParameter_SignKeyError = 403071,

        [Description(i18nKeys.ErrorCode_InvalidParameter_SignKeyExpired)]
        InvalidParameter_SignKeyExpired = 403072,

        InvalidParameter_CollectionIsEmptyNull = 403098,

        [Description(i18nKeys.ErrorCode_InvalidParameter_ModelIsNull)]
        InvalidParameter_ModelIsNull = 403099,




        [Description(i18nKeys.ErrorCode_Forbidden_IPAddressNotAllow)]
        Forbidden_IPAddressNotAllow = 403201,

        [Description(i18nKeys.ErrorCode_Forbidden_DuplicateEnterGameServer)]
        Forbidden_DuplicateEnterGameServer = 403202,

        [Description(i18nKeys.ErrorCode_Forbidden_EnterGameError_WalletAmountNotEnough)]
        Forbidden_EnterGameError_WalletAmountNotEnough = 403203,

        [Description(i18nKeys.ErrorCode_Forbidden_RequiredApiCertificate)]
        Forbidden_RequiredApiCertificate = 403204,

        [Description(i18nKeys.ErrorCode_Forbidden_DoNotAllowWhenProduction)]
        Forbidden_DoNotAllowWhenProduction = 403222,
        
        //禁止編輯
        Forbidden_DoNotAllowEdit = 403223,

        //未定義的編號
        Forbidden_NotValidDefineValue = 403224,


        [Description(i18nKeys.ErrorCode_Forbidden_PagePermissionNotAllow)]
        Forbidden_PagePermissionNotAllow = 403271,

        [Description(i18nKeys.ErrorCode_Forbidden_WatchDataPermissionNotAllow)]
        Forbidden_WatchDataPermissionNotAllow = 403272,

        [Description(i18nKeys.ErrorCode_Forbidden_WhenPlayingGame_NotAllowTransfer)]
        Forbidden_WhenPlayingGame_NotAllowTransfer = 403273,

        Forbidden_NotSupportExternalLinkDisplay = 403288,


        SystemMaintain = 403501,

        GameUnderMaintain = 403508,


        [Description(i18nKeys.ErrorCode_InvalidStatus)]
        InvalidStatus = 403601,

        [Description(i18nKeys.ErrorCode_InvalidStatus_AccountStatus)]
        InvalidStatus_AccountStatus = 403602,

        [Description(i18nKeys.ErrorCode_InvalidStatus_OrganizationStatus)]
        InvalidStatus_OrganizationStatus = 403603,

        [Description(i18nKeys.ErrorCode_InvalidStatus_DisableNow)]
        InvalidStatus_DisableNow = 403609,

        [Description(i18nKeys.ErrorCode_InvalidStatus_GameStatus)]
        InvalidStatus_GameStatus = 403610,

        InvalidStatus_BetRecordDuplicateSettlement = 403611,



        [Description(i18nKeys.ErrorCode_WalletTranferFail_WalletAmountNotEnough)]
        WalletTranferFail_WalletAmountNotEnough = 403651,

        BetFail_WalletAmountNotEnough = 403652,


        [Description(i18nKeys.ErrorCode_ThirdPartyError_SettingDataNotCompleted)]
        ThirdPartyError_SettingDataNotCompleted = 403901,

        [Description(i18nKeys.ErrorCode_ThirdPartyApiFail_ResponseError)]
        ThirdPartyApiFail_ResponseError = 403911,

        [Description(i18nKeys.ErrorCode_ThirdPartyApiFail_ParsingFail)]
        ThirdPartyApiFail_ParsingFail = 403912,

        [Description(i18nKeys.ErrorCode_ThirdPartyApiFail_IsNotValid)]
        ThirdPartyApiFail_IsNotValid = 403913,

        [Description(i18nKeys.ErrorCode_ThirdPartyApiFail_SignCheckNotMatch)]
        ThirdPartyApiFail_SignCheckNotMatch = 403914,

        ThirdPartyApiFail_ResponseDataUnExpected = 403915,

        ThirdPartyApiFail_Timeout = 403918,

        #endregion

        #region 404 查無資料類

        [Description(i18nKeys.ErrorCode_NotFoundData)]
        NotFoundData = 404001,

        [Description(i18nKeys.ErrorCode_NotFoundMerchant)]
        NotFoundMerchant = 404002,

        [Description(i18nKeys.ErrorCode_NotFoundPlayer)]
        NotFoundPlayer = 404003,

        [Description(i18nKeys.ErrorCode_NotFoundUser)]
        NotFoundUser = 404004,

        [Description(i18nKeys.ErrorCode_NotFoundUserLogin)]
        NotFoundUserLogin = 404005,

        [Description(i18nKeys.ErrorCode_NotFoundOrganization)]
        NotFoundOrganization = 404006,

        [Description(i18nKeys.ErrorCode_NotFoundServer)]
        NotFoundServer = 404008,
        
        [Description(i18nKeys.ErrorCode_NotFoundGame)]
        NotFoundGame = 404009,

        [Description(i18nKeys.ErrorCode_NotFoundGameRecord)]
        NotFoundGameRecord = 404010,

        [Description(i18nKeys.ErrorCode_NotFoundPlayerWallet)]
        NotFoundPlayerWallet = 404011,

        [Description(i18nKeys.ErrorCode_NotFoundServerCurrencyMatch)]
        NotFoundServerCurrencyMatch = 404012,

        [Description(i18nKeys.ErrorCode_NotFoundRole)]
        NotFoundRole = 404013,

        NotFoundRecord = 404014,

        NotFoundSessionToken = [REDACTED],

        NotFoundBetRecord = 404016,

        [Description(i18nKeys.ErrorCode_NotFoundAvailableServer_MerchantApiSelected)]
        NotFoundAvailableServer_MerchantApiSelected = 404051,

        [Description(i18nKeys.ErrorCode_NotFoundAvailableServer_FlagLevelNotMatch)]
        NotFoundAvailableServer_FlagLevelNotMatch = 404052,


        NotFoundPerGameLogButBetRecordExist = 404888,

        #endregion




        [Description(i18nKeys.ErrorCode_InternalServerError)]
        InternalServerError = 500001,

        [Description(i18nKeys.ErrorCode_NotImplement)]
        NotImplement = 500002,

        [Description(i18nKeys.ErrorCode_BusinessLogicDeny)]
        BusinessLogicDeny = 500018,

        [Description(i18nKeys.ErrorCode_BusinessLogicDeny_SettingDataNotCompleted)]
        BusinessLogicDeny_SettingDataNotCompleted = 500019,

        [Description(i18nKeys.ErrorCode_BusinessLogicDeny_NotAllowDeleteData_ModifyRecently)]
        BusinessLogicDeny_NotAllowDeleteData_ModifyRecently = 500020,

        [Description(i18nKeys.ErrorCode_BusinessLogicDeny_NotAllowDeleteData_CurrentServer)]
        BusinessLogicDeny_NotAllowDeleteData_CurrentServer = 500021,

        BusinessLogicDeny_DuplicateOrderNo = 500028,

        [Description(i18nKeys.ErrorCode_BusinessLogicDeny_DuplicateDataExist)]
        BusinessLogicDeny_DuplicateDataExist = 500029,

        [Description(i18nKeys.ErrorCode_BusinessLogicDeny_CurrencyNotMatch)]
        BusinessLogicDeny_CurrencyNotMatch = 500030,


        BusinessLogicDeny_MatchDataIsNotSingleOne = 500035,


        FunkyApi_InvalidInput = 700400,
        FunkyApi_PlayerIsNotLogin = 700401,
        FunkyApi_InsufficientBalance = 700402,
        FunkyApiBet_Already_Exists_ = 700403,
        FunkyApi_Bet_Was_not_found = 700404,
        FunkyApi_API_Suspended = 700405,
        FunkyApi_over_The_max_winning = 700406,
        FunkyApi_Over_The_max_loss = 700407,
        FunkyApi_Bet_Was_Already_Settled = 700409,
        FunkyApi_Bet_Was_Already_Cancelled = 700410,
        FunkyApi_GameCodeNotDefine = 700499,


    }
}
