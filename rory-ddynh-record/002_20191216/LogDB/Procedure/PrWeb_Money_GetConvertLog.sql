USE [LogDB]
GO
/****** Object:  StoredProcedure [dbo].[PrWeb_Money_GetConvertLog]    Script Date: 2019/12/18 上午 11:50:18 ******/
SET ANSI_NULLS ON
GO
SET QUOTED_IDENTIFIER ON
GO
DROP PROCEDURE IF EXISTS dbo.PrWeb_Money_GetConvertLog
GO
-- =============================================
-- Author:  Jason
-- Create date: 2018/03/22
-- Description:	<获取兑换记录>
--	2018/04/18	Jason		語法優化及調整時間限制未指定日期則由變數替換
--	2018/04/24	Jason		新增[BankAmount],[CashAmount],[Operator]
--	2018/05/16	Jason		#255 後台兌換功能調整 1. 預設列出審核通過的清單 2. 顯示已經打款完成的兌換金額總數與手續費總數
--	2018/05/19	alan		調整統計數據, 改用到帳時間判斷
--	2018-06-13	Jason		調整@dtbeginTime,@dtendTime 30日限制改為全範圍
--	2019/05/31	minglung	新增@chvOperator參數
--  2019/07/15	minglung	新增[Memo],[GWID],[GatewayName],@lngMinAmount,@lngMaxAmount,@intGWID
--	2019/12/16	mars		增加返回BankCodeID
-- =============================================
CREATE PROCEDURE [dbo].[PrWeb_Money_GetConvertLog]
(
	@intUserID		INT = NULL,			--用户ID
	@tnyStatus		TINYINT = NULL,		--状态
	@chvRealName	VARCHAR(64) = NULL,	--真名姓名
	@chvIPAddress	VARCHAR(15) = NULL,	--IP
	@dtbeginTime	DATETIME = NULL,		
	@dtendTime		DATETIME = NULL,
	@tnyConvertType	TINYINT = NULL,		--兌換類別 1:支付寶 2:銀聯
	@chvOperator	VARCHAR(16) = NULL,	--操作人
	@lngMinAmount	BIGINT = 0,			--兑换金额
	@lngMaxAmount	BIGINT = 99999999,	--兑换金额
	@intGWID		INT = NULL			--兑换商编号
)
AS
BEGIN
	
	DECLARE @UserID			INT = @intUserID
	DECLARE @Status			TINYINT = @tnyStatus
	DECLARE @RealName		VARCHAR(64) = @chvRealName
	DECLARE @IPAddress		VARCHAR(15) = @chvIPAddress
	DECLARE @beginTime		DATETIME = @dtbeginTime		
	DECLARE @endTime		DATETIME = @dtendTime
	DECLARE @ConvertType	TINYINT = @tnyConvertType
	DECLARE @Operator		VARCHAR(16) = @chvOperator

	IF (@beginTime IS NULL AND @endTime IS NULL)
	BEGIN
		SET @beginTime	=	DATEADD(DAY,-30,DATEDIFF(DAY,0,GETDATE()))
		SET @endTime	=	DATEADD(DAY,0,DATEDIFF(DAY,0,GETDATE()))
	END
	ELSE IF (@beginTime IS NULL AND @endTime IS NOT NULL)
	BEGIN
		--SET @beginTime	=	DATEADD(DAY,-30,DATEDIFF(DAY,0,@endTime))
		SET @beginTime	=	'2018-03-01 00:00:00.000'
	END
	ELSE IF (@beginTime IS NOT NULL AND @endTime IS NULL)
	BEGIN
		--SET @endTime	=	DATEADD(DAY,30,DATEDIFF(DAY,0,@beginTime))
		SET @endTime	=	GETDATE()
	END
	
	SET @lngMinAmount = -@lngMinAmount * 100
	SET @lngMaxAmount = -@lngMaxAmount * 100

	SELECT	L.[Id],
			L.[UserID],
			L.[Amount],
			L.[TaxAmount],
			L.[RealTaxAmount],
			L.[Status],
			L.[ConvertType],
			L.[RealName],
			L.[Account],
			L.[BankCodeID],
			L.[BankName],
			L.[CreateIP],
			L.[TransID],
			L.[TransData],
			L.[CreateTime],
			L.[LastUpdateTime],
			L.[CheckTime],
			L.[PlatformID],
			L.[AreaID],
			L.[BankAmount],
			L.[CashAmount],
			L.[Operator],
			L.[Memo],
			L.[GWID],
			C.[GatewayName]
	FROM [Money].[ConvertLog] AS L WITH (NOLOCK)
	LEFT JOIN CenterDB.Payment.PayGateway AS C WITH (NOLOCK) ON L.GWID = C.ID	
	WHERE	
		ISNULL(@UserID,L.UserID) = L.UserID
		AND ISNULL(@RealName,L.RealName) = L.RealName
		AND ISNULL(@IPAddress,L.CreateIP) = L.CreateIP
		AND ISNULL(@beginTime,L.LastUpdateTime) <= L.LastUpdateTime
		AND ISNULL(@endTime,L.LastUpdateTime) >= L.LastUpdateTime				
		AND ISNULL(@ConvertType,L.ConvertType) = L.ConvertType
		AND (@Operator is null or @Operator = L.Operator)
		AND Amount Between @lngMaxAmount AND @lngMinAmount
		AND (@intGWID is null or @intGWID = L.GWID)
END
