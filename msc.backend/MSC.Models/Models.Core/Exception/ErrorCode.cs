using MSC.Models.Attributes;
using MSC.Models.CultureText;
using MSC.Models.Enums;
using System.ComponentModel;

namespace MSC.Models.Models
{
    public enum ErrorCode
    {


        Success = 1,


       

        [ReferenceText(i18nKeys.ErrorCode_ServerPlayerOutOfRange)]
        ServerPlayerOutOfRange = 400201,


        //401001~401099 會給前端拿來判斷強制導到登入頁，請謹慎使用
        //401001~401099 會給前端拿來判斷強制導到登入頁，請謹慎使用
        //401001~401099 會給前端拿來判斷強制導到登入頁，請謹慎使用
        //401001~401099 會給前端拿來判斷強制導到登入頁，請謹慎使用
        //401001~401099 會給前端拿來判斷強制導到登入頁，請謹慎使用
        //401001~401099 會給前端拿來判斷強制導到登入頁，請謹慎使用


        //401 身分授權
        [ReferenceText(i18nKeys.ErrorCode_Authorize_AccessToken_NotValid)]
        Authorize_AccessToken_NotValid = 401001,

        [ReferenceText(i18nKeys.ErrorCode_Authorize_AccessToken_Expired)]
        Authorize_AccessToken_Expired = 401002,

        [ReferenceText(i18nKeys.ErrorCode_Authorize_TokenHasChanged)]
        Authorize_TokenHasChanged = 401003,

        [ReferenceText(i18nKeys.ErrorCode_Authorize_JwtTokenParsingFail)]
        Authorize_JwtTokenParsingFail = 401008,

        [ReferenceText(i18nKeys.ErrorCode_Authorize_JwtTokenVersionNotMatch)]
        Authorize_JwtTokenVersionNotMatch = 401009,

        [ReferenceText(i18nKeys.ErrorCode_Authorize_AccessToken_Expired)]
        Authorize_PleaseLogin = 401099,

        [ReferenceText(i18nKeys.ErrorCode_InvalidParameter)]
        InvalidParameter = 403001,

        [ReferenceText(i18nKeys.ErrorCode_InvalidParameter_TextEmpty)]
        InvalidParameter_TextEmpty = 403002,

        [ReferenceText(i18nKeys.ErrorCode_InvalidParameter_TextLength)]
        InvalidParameter_TextLength = 403003, 
        
        [ReferenceText(i18nKeys.ErrorCode_InvalidParameter_NotAccountFormat)]
        InvalidParameter_NotAccountFormat = 403004,

        [ReferenceText(i18nKeys.ErrorCode_InvalidParameter_NotSupportLanguageCode)]
        InvalidParameter_NotSupportLanguageCode = 403005,

        [ReferenceText(i18nKeys.ErrorCode_InvalidParameter_NotAccountFormat)]
        InvalidParameter_NotAvailablePassword = [REDACTED],

        [ReferenceText(i18nKeys.ErrorCode_InvalidParameter_ConfirmCodeNotMatch)]
        InvalidParameter_ConfirmCodeNotMatch = 403009,

        InvalidParameter_OriginalPasswordNotMatch = 403010,

        [ReferenceText(i18nKeys.ErrorCode_InvalidParameter_NumberEmptyZero)]
        InvalidParameter_NumberEmptyZero = 403011,

        [ReferenceText(i18nKeys.ErrorCode_InvalidParameter_NumberMinus)]
        InvalidParameter_NumberMinus = 403012,

        [ReferenceText(i18nKeys.ErrorCode_InvalidParameter_NumberValueTooSmall)]
        InvalidParameter_NumberValueTooSmall = 403013,

        [ReferenceText(i18nKeys.ErrorCode_InvalidParameter_NumberValueTooLarge)]
        InvalidParameter_NumberValueTooLarge = 403014,

        [ReferenceText(i18nKeys.ErrorCode_InvalidParameter_NumberValueNotAllowFloat)]
        InvalidParameter_NumberValueNotAllowFloat = 403015,

        [ReferenceText(i18nKeys.ErrorCode_InvalidParameter_NumberValueNotAllowNull)]
        InvalidParameter_NumberValueNotAllowNull = 403018,

        [ReferenceText(i18nKeys.ErrorCode_InvalidParameter_Collection)]
        InvalidParameter_Collection = 403019,

        [ReferenceText(i18nKeys.ErrorCode_InvalidParameter_DateTimeRequired)]
        InvalidParameter_DateTimeRequired = 403021,

        [ReferenceText(i18nKeys.ErrorCode_InvalidParameter_StartTimeMustSmallThanEndTime)]
        InvalidParameter_StartTimeMustSmallThanEndTime = 403022,

        [ReferenceText(i18nKeys.ErrorCode_InvalidParameter_DateTimeRangeTooLarge)]
        InvalidParameter_DateTimeRangeTooLarge = 403023,

        [ReferenceText(i18nKeys.ErrorCode_InvalidParameter_DateTimeTooSmall)]
        InvalidParameter_DateTimeTooSmall = 403024,

        [ReferenceText(i18nKeys.ErrorCode_InvalidParameter_InvalidLetter)]
        InvalidParameter_InvalidLetter = 403025,

        [ReferenceText(i18nKeys.ErrorCode_InvalidParameter_ParsingError)]
        InvalidParameter_ParsingError = 403059,

        [ReferenceText(i18nKeys.ErrorCode_InvalidParameter_TimestampError)]
        InvalidParameter_TimestampError = 403070,

        [ReferenceText(i18nKeys.ErrorCode_InvalidParameter_SignKeyError)]
        InvalidParameter_SignKeyError = 403071,

        [ReferenceText(i18nKeys.ErrorCode_InvalidParameter_SignKeyExpired)]
        InvalidParameter_SignKeyExpired = 403072,

        [ReferenceText(i18nKeys.ErrorCode_InvalidParameter_ModelIsNull)]
        InvalidParameter_ModelIsNull = 403099,




        [ReferenceText(i18nKeys.ErrorCode_Forbidden_IPAddressNotAllow)]
        Forbidden_IPAddressNotAllow = 403201,

        [ReferenceText(i18nKeys.ErrorCode_Forbidden_DuplicateEnterGameServer)]
        Forbidden_DuplicateEnterGameServer = 403202,

        [ReferenceText(i18nKeys.ErrorCode_Forbidden_EnterGameError_WalletAmountNotEnough)]
        Forbidden_EnterGameError_WalletAmountNotEnough = 403203,

        [ReferenceText(i18nKeys.ErrorCode_Forbidden_RequiredApiCertificate)]
        Forbidden_RequiredApiCertificate = 403204,

        [ReferenceText(i18nKeys.ErrorCode_Forbidden_DoNotAllowWhenProduction)]
        Forbidden_DoNotAllowWhenProduction = 403222,

        [ReferenceText(i18nKeys.ErrorCode_Forbidden_PagePermissionNotAllow)]
        Forbidden_PagePermissionNotAllow = 403271,

        [ReferenceText(i18nKeys.ErrorCode_Forbidden_WatchDataPermissionNotAllow)]
        Forbidden_WatchDataPermissionNotAllow = 403272,

        [ReferenceText(i18nKeys.ErrorCode_Forbidden_WhenPlayingGame_NotAllowTransfer)]
        Forbidden_WhenPlayingGame_NotAllowTransfer = 403273,

        Forbidden_OrderAmount_NotAllow = 403281,

        Forbidden_CardPoint_OrderAmount_NotMatch = 403282,

        //卡號已銷毀
        Forbidden_Card_CardStatusIsDestroy = 403283,

        Forbidden_Card_CardCurrencyNotSupport = 403284,

        Forbidden_NotValidDefineValue = 403331,

        Forbidden_DoNotAllowEdit = 403332,

        Forbidden_BalanceNotEnough = 403455,


        [ReferenceText(i18nKeys.ErrorCode_InvalidStatus)]
        InvalidStatus = 403601,

        [ReferenceText(i18nKeys.ErrorCode_InvalidStatus_AccountStatus)]
        InvalidStatus_AccountStatus = 403602,

        [ReferenceText(i18nKeys.ErrorCode_InvalidStatus_OrganizationStatus)]
        InvalidStatus_OrganizationStatus = 403603,

        [ReferenceText(i18nKeys.ErrorCode_InvalidStatus_DisableNow)]
        InvalidStatus_DisableNow = 403609,

        [ReferenceText(i18nKeys.ErrorCode_WalletTranferFail_WalletAmountNotEnough)]
        WalletTranferFail_WalletAmountNotEnough = 403651,



        [ReferenceText(i18nKeys.ErrorCode_ParsingError_EncryptOrDecryptError)]
        ParsingError_EncryptOrDecryptError = 403701,

        [ReferenceText(i18nKeys.ErrorCode_ParsingError_ParseTokenObjectFail)]
        ParsingError_ParseTokenObjectFail = 403702,

        [ReferenceText(i18nKeys.ErrorCode_ThirdPartyError_SettingDataNotCompleted)]
        ThirdPartyError_SettingDataNotCompleted = 403901,

        [ReferenceText(i18nKeys.ErrorCode_ThirdPartyApiFail_ResponseError)]
        ThirdPartyApiFail_ResponseError = 403911,

        [ReferenceText(i18nKeys.ErrorCode_ThirdPartyApiFail_ParsingFail)]
        ThirdPartyApiFail_ParsingFail = 403912,

        [ReferenceText(i18nKeys.ErrorCode_ThirdPartyApiFail_IsNotValid)]
        ThirdPartyApiFail_IsNotValid = 403913,

        [ReferenceText(i18nKeys.ErrorCode_ThirdPartyApiFail_SignCheckNotMatch)]
        ThirdPartyApiFail_SignCheckNotMatch = 403914,

        [ReferenceText(i18nKeys.ErrorCode_ThirdPartyApiFail_SignCheckNotMatch)]
        ThirdPartyApiFail_TimestampInvalid = 403915,

        [ReferenceText(i18nKeys.ErrorCode_ThirdPartyApiFail_TokenExpired)]
        ThirdPartyApiFail_TokenExpired = 403916,




        [ReferenceText(i18nKeys.ErrorCode_NotFoundData)]
        NotFoundData = 404001,

        [ReferenceText(i18nKeys.ErrorCode_NotFoundMerchant)]
        NotFoundMerchant = 404002,

        [ReferenceText(i18nKeys.ErrorCode_NotFoundPlayer)]
        NotFoundPlayer = 404003,

        [ReferenceText(i18nKeys.ErrorCode_NotFoundUser)]
        NotFoundUser = 404004,

        [ReferenceText(i18nKeys.ErrorCode_NotFoundUserLogin)]
        NotFoundUserLogin = 404005,

        [ReferenceText(i18nKeys.ErrorCode_NotFoundOrganization)]
        NotFoundOrganization = 404006,

        [ReferenceText(i18nKeys.ErrorCode_NotFoundServer)]
        NotFoundServer = 404008,
        
        [ReferenceText(i18nKeys.ErrorCode_NotFoundGame)]
        NotFoundGame = 404009,

        [ReferenceText(i18nKeys.ErrorCode_NotFoundGameRecord)]
        NotFoundGameRecord = 404010,

        [ReferenceText(i18nKeys.ErrorCode_NotFoundPlayerWallet)]
        NotFoundPlayerWallet = 404011,

        [ReferenceText(i18nKeys.ErrorCode_NotFoundServerCurrencyMatch)]
        NotFoundServerCurrencyMatch = 404012,

        [ReferenceText(i18nKeys.ErrorCode_NotFoundRole)]
        NotFoundRole = 404013,

        [ReferenceText(i18nKeys.ErrorCode_NotFoundAvailableServer_MerchantApiSelected)]
        NotFoundAvailableServer_MerchantApiSelected = 404051,

        [ReferenceText(i18nKeys.ErrorCode_NotFoundAvailableServer_FlagLevelNotMatch)]
        NotFoundAvailableServer_FlagLevelNotMatch = 404052,        

        [ReferenceText(i18nKeys.ErrorCode_NotFoundRequestData)]
        NotFoundRequestData = 404053,

        NotFoundCardNo_OrPassword = [REDACTED],




        [ReferenceText(i18nKeys.ErrorCode_NotFoundCardNo_FormatIncorrect)]
        NotFoundCardNo_FormatIncorrect = 404056,

        [ReferenceText(i18nKeys.Common_Undefined)]
        NotFoundNotifyQueue = 404061,

        [ReferenceText(i18nKeys.Common_Undefined)]
        NotFoundCardConsume = 404062,




        [ReferenceText(i18nKeys.ErrorCode_InternalServerError)]
        InternalServerError = 500001,

        [ReferenceText(i18nKeys.ErrorCode_NotImplement)]
        NotImplement = 500002,

        [ReferenceText(i18nKeys.ErrorCode_BusinessLogicDeny)]
        System_SystemMaintain = 500015,

        [ReferenceText(i18nKeys.ErrorCode_BusinessLogicDeny)]
        BusinessLogicDeny = 500018,

        [ReferenceText(i18nKeys.ErrorCode_BusinessLogicDeny_SettingDataNotCompleted)]
        BusinessLogicDeny_SettingDataNotCompleted = 500019,

        [ReferenceText(i18nKeys.ErrorCode_BusinessLogicDeny_NotAllowDeleteData_ModifyRecently)]
        BusinessLogicDeny_NotAllowDeleteData_ModifyRecently = 500020,

        [ReferenceText(i18nKeys.ErrorCode_BusinessLogicDeny_NotAllowDeleteData_CurrentServer)]
        BusinessLogicDeny_NotAllowDeleteData_CurrentServer = 500021,

        [ReferenceText(i18nKeys.ErrorCode_BusinessLogicDeny_DuplicateDataExist)]
        BusinessLogicDeny_DuplicateDataExist = 500029,

        [ReferenceText(i18nKeys.ErrorCode_BusinessLogicDeny_CurrencyNotMatch)]
        BusinessLogicDeny_CurrencyNotMatch = 500030,

        [ReferenceText(i18nKeys.ErrorCode_BusinessLogicDeny_UndefinedValue)]
        BusinessLogicDeny_UndefinedValue = 500031,

        [ReferenceText(i18nKeys.ErrorCode_BusinessLogicDeny_UndefinedColumnNameValueInRawEntity)]
        BusinessLogicDeny_UndefinedColumnNameValueInRawEntity = 500032,


        BusinessLogicDeny_RequestTooManyTimes = 500033,

        CommonError_Timeout = 500351,



        [ReferenceText(i18nKeys.ErrorCode_BusinessLogicDeny_UndefinedValue)]
        StoreProcedure_CheckIsValidCard_NotFoundData = 550001,

        [ReferenceText(i18nKeys.ErrorCode_BusinessLogicDeny_UndefinedValue)]
        StoreProcedure_CheckIsValidCard_CardAlreadyUse = 550002,



    }
}
