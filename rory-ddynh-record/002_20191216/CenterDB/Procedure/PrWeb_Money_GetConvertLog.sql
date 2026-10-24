USE [CenterDB]
GO
/****** Object:  StoredProcedure [dbo].[PrWeb_Money_GetConvertLog]    Script Date: 2019/12/18 上午 11:54:17 ******/
SET ANSI_NULLS ON
GO
SET QUOTED_IDENTIFIER ON
GO
DROP PROCEDURE IF EXISTS dbo.PrWeb_Money_GetConvertLog
GO
-- =============================================
-- Author:  <xxx>
-- Create date: <xxx>
-- Description:	<获取兑换记录>
--	2018/03/22	Jason		調整為LINKEDSERVER取資料回來計算
--	2018/04/18	Jason		語法調整效能優化 & 未填入查詢日期時預設為一個月內
--	2018/04/24	Jason		新增@tnyConvertType參數及TotalCharge,TotalConvert,ChargeConvertRate計算
--	2018/05/16	Jason		#255 後台兌換功能調整 1. 預設列出審核通過的清單 2. 顯示已經打款完成的兌換金額總數與手續費總數
--	2018/05/18	Jason		#269 兌換頁面增加支付寶與銀行單獨呈現總兌換金額與手續費
--	2018/06/13	Jason		@chvPlatformIDs VARCHAR(32) change to VARCHAR(2000)
--	2018/06/27	jason		聯運需求, 開放後台給第三方查詢特定平台與渠道營運資訊 @chvPlatformIDs 改回@chvPlatformID & 新增@intAdminID
--	2018/08/10	alan		#660 檢查是否是嫌疑帳號
--	2019/05/31	minglung	新增@chvOperator參數
--  2019/07/15	minglung	新增[Memo],[GWID],[GatewayName],@lngMinAmount,@lngMaxAmount,@intGWID
--	2019/12/16	mars		增加返回BankCodeID
-- =============================================
CREATE PROCEDURE [dbo].[PrWeb_Money_GetConvertLog]
	@intAdminID				INT,
	@intPlatformID			INT = NULL,	--平台ID
	@intUserID				INT = NULL,			--用户ID
	@tnyStatus				TINYINT = NULL,		--状态
	@chvRealName			VARCHAR(64) = NULL,	--真名姓名
	@chvIPAddress			VARCHAR(15) = NULL,	--IP
	@dtbeginTime			DATETIME = NULL,		
	@dtendTime				DATETIME = NULL,
	@tnyConvertType			TINYINT = NULL,		-- 兌換類別 1:支付寶 2:銀聯	
	@chvOperator			VARCHAR(16) = NULL,	--操作人
	@lngMinAmount			BIGINT = 0,			--兑换金额
	@lngMaxAmount			BIGINT = 99999999,	--兑换金额
	@intGWID				INT = NULL,			--兑换商编号
	@intPageIndex			INT = 1,
	@intPageSize			INT = 20,
	@intPageCount			INT		OUTPUT
AS
BEGIN

	SET NOCOUNT ON

	DECLARE @lngTotalConvertAmount	BIGINT = 0
	DECLARE @lngTotalConvertIOSAmount	BIGINT = 0
	DECLARE @lngTotalConvertAndroidAmount	BIGINT = 0
	DECLARE @lngTotalTaxAmount		BIGINT = 0
	DECLARE @HeadLineJson			VARCHAR(MAX)
	DECLARE @AlipayAmount			BIGINT = 0
	DECLARE @UnionpayAmount			BIGINT = 0
	DECLARE @AlipayTaxAmount		BIGINT = 0
	DECLARE @UnionpayTaxAmount		BIGINT = 0
	DECLARE @TotalConvertUsers		INT = 0

	DECLARE @Temp_ConvertLog	TABLE
		(
			[Id]				INT,
			[UserID]			INT,
			[Amount]			BIGINT,
			[TaxAmount]			BIGINT,
			[RealTaxAmount]		BIGINT,
			[Status]			TINYINT,
			[ConvertType]		TINYINT,
			[RealName]			VARCHAR(64),
			[Account]			VARCHAR(64),
			[BankCodeID]		INT,
			[BankName]			VARCHAR(50),
			[CreateIP]			VARCHAR(15),
			[TransID]			VARCHAR(64),
			[TransData]			VARCHAR(MAX),
			[CreateTime]		DATETIME,
			[LastUpdateTime]	DATETIME,
			[CheckTime]			DATETIME,
			[PlatformID]		INT,
			[AreaID]			INT,
			[BankAmount]		BIGINT,
			[CashAmount]		BIGINT,
			[Operator]			VARCHAR(16),
			[Memo]				VARCHAR(200),
			[GWID]				INT,
			[GatewayName]		VARCHAR(20)
		)


	/*  暫存資料  */
	--------------------------------------------------------------------------------------
	INSERT INTO @Temp_ConvertLog
	EXEC [LOGDBLINKEDSERVER].[LogDB].[dbo].[PrWeb_Money_GetConvertLog]
		@intUserID				= @intUserID,
		@tnyStatus				= @tnyStatus,
		@chvRealName			= @chvRealName,
		@chvIPAddress			= @chvIPAddress,
		@dtbeginTime			= @dtbeginTime,	
		@dtendTime				= @dtendTime,
		@tnyConvertType			= @tnyConvertType,
		@chvOperator			= @chvOperator,
		@lngMinAmount			= @lngMinAmount,
		@lngMaxAmount			= @lngMaxAmount,
		@intGWID				= @intGWID


	DECLARE @Temp_PlatformID TABLE (PlatformID	INT)
	DECLARE @Temp_AreaID TABLE (AreaID	INT)

	--筛选平台
	INSERT INTO @Temp_PlatformID
	SELECT	DISTINCT PlatformID 
	FROM [dbo].[Fn_GetPlatformAccessList](@intAdminID)
	WHERE
		ISNULL(@intPlatformID,PlatformID) = PlatformID

	--筛选渠道
	IF EXISTS(Select 1 from Area.AreaAdminUser where AdminID=@intAdminID)
		INSERT INTO @Temp_AreaID Select AreaID From Area.AreaAdminUser where AdminID=@intAdminID
	ELSE
		INSERT INTO @Temp_AreaID Select AreaID From Area.Area where PlatformID in (Select PlatformID from @Temp_PlatformID)

	--------------------------------------------------------------------------------------


	/*  回傳資料  */
	--------------------------------------------------------------------------------------
	--取得總頁數
	SELECT	@intPageCount = COUNT(1)
	FROM @Temp_ConvertLog AS L
		INNER JOIN @Temp_PlatformID AS PID ON L.PlatformID = PID.PlatformID
		INNER JOIN Users.UserInfo AS U WITH (NOLOCK) ON L.UserID=U.UserID
	WHERE
		ISNULL(@intUserID,L.UserID) = L.UserID
		AND ISNULL(@tnyStatus,L.Status) = L.Status
		AND ISNULL(@chvRealName,L.RealName) = L.RealName
		AND ISNULL(@chvIPAddress,L.CreateIP) = L.CreateIP	
		AND ISNULL(@dtbeginTime,L.LastUpdateTime) <= L.LastUpdateTime
		AND ISNULL(@dtendTime,L.LastUpdateTime) >= L.LastUpdateTime
		AND ISNULL(@tnyConvertType,L.[ConvertType]) >= L.[ConvertType]
		AND ISNULL(@tnyStatus,L.Status) = L.Status
		AND L.AreaID IN (Select AreaID From @Temp_AreaID)

	--取得統計數
	SELECT	@lngTotalConvertAmount = SUM([Amount]),
			@lngTotalConvertIOSAmount = SUM(CASE WHEN U.HDType<>2 THEN ABS( ISNULL(L.Amount,0) ) ELSE 0 END),
			@lngTotalConvertAndroidAmount = SUM(CASE WHEN U.HDType=2 THEN ABS( ISNULL(L.Amount,0) ) ELSE 0 END),
			@lngTotalTaxAmount = SUM([TaxAmount]),
			@AlipayAmount = SUM(CASE WHEN ConvertType = 1 THEN Amount ELSE 0 END),
			@AlipayTaxAmount = SUM(CASE WHEN ConvertType = 1 THEN TaxAmount ELSE 0 END),
			@UnionpayAmount = SUM(CASE WHEN ConvertType = 2 THEN Amount ELSE 0 END),
			@UnionpayTaxAmount = SUM(CASE WHEN ConvertType = 2 THEN TaxAmount ELSE 0 END),
			@TotalConvertUsers = COUNT(DISTINCT L.UserID)
	FROM @Temp_ConvertLog AS L
		INNER JOIN @Temp_PlatformID AS Z ON L.PlatformID = Z.PlatformID
		LEFT JOIN Users.UserInfo U with (nolock) on L.UserID=U.UserID
	WHERE
		L.Status = 3
		AND L.AreaID IN (Select AreaID From @Temp_AreaID)
		
	SET @HeadLineJson	= '{"TotalConvertAmount":' + CAST(ISNULL(@lngTotalConvertAmount,0) AS VARCHAR) + 
							',"TotalConvertIOSAmount":' + CAST(ISNULL(@lngTotalConvertIOSAmount,0) AS VARCHAR) +
							',"TotalConvertAndroidAmount":' + CAST(ISNULL(@lngTotalConvertAndroidAmount,0) AS VARCHAR) +
							',"TotalTaxAmount":' + CAST(ISNULL(@lngTotalTaxAmount,0) AS VARCHAR) +
							',"Alipay":' + CAST(ISNULL(@AlipayAmount,0) AS VARCHAR) + 
							',"AlipayTax":'+ CAST(ISNULL(@AlipayTaxAmount,0) AS VARCHAR) +
							',"Unionpay":' + CAST(ISNULL(@UnionpayAmount,0) AS VARCHAR) +
							',"UnionpayTax":'+ CAST(ISNULL(@UnionpayTaxAmount,0) AS VARCHAR) +
							',"TotalConvertUsers":'+ CAST(ISNULL(@TotalConvertUsers,0) AS VARCHAR) + '}'

	;WITH tab AS(
			SELECT	L.*,
					'User.Id' = U.UserID,
					'User.Grade' = U.GradeID,
					'User.LoginName' = U.LoginName,
					'User.Area.Platform.Name' = P.Name,				
					TotalCharge = ISNULL(U.[SumCharge]+U.[SumAgentCharge],0),
					TotalConvert = ISNULL(U.[SumConvert],0),
					ChargeConvertRate = ISNULL(U.[SumConvert],0) - ISNULL(U.[SumCharge]+U.[SumAgentCharge],0),
					rowid = ROW_NUMBER() OVER (ORDER BY L.CreateTime DESC),
					CASE WHEN C.Class <> 99 THEN 1 ELSE 0 END as ClubAgent
					
			FROM @Temp_ConvertLog AS L
				INNER JOIN @Temp_PlatformID AS Z ON L.PlatformID = Z.PlatformID
				INNER JOIN [Area].[Platform] AS P WITH (NOLOCK) ON Z.PlatformID = P.PlatformID
				INNER JOIN Users.UserInfo AS U WITH (NOLOCK) ON L.UserID=U.UserID
				LEFT JOIN Club.dbo.AgentUsers C WITH (NOLOCK) ON L.UserID=C.UserID AND C.Class <> 99
			WHERE
				ISNULL(@intUserID,L.UserID) = L.UserID
				AND ISNULL(@tnyStatus,L.Status) = L.Status
				AND ISNULL(@chvRealName,L.RealName) = L.RealName
				AND ISNULL(@chvIPAddress,L.CreateIP) = L.CreateIP			
				AND ISNULL(@dtbeginTime,L.LastUpdateTime) <= L.LastUpdateTime
				AND ISNULL(@dtendTime,L.LastUpdateTime) >= L.LastUpdateTime				
				AND ISNULL(@tnyConvertType,L.[ConvertType]) >= L.[ConvertType]
				AND ISNULL(@tnyStatus,L.Status) = L.Status
				AND L.AreaID IN (Select AreaID From @Temp_AreaID)
		)
	SELECT 
	T.*, --兌換清單
	CASE WHEN S.AccountID IS NOT NULL THEN 1 ELSE 0 END AS Suspicious --嫌疑帳號
	FROM tab T
	LEFT JOIN Game.ConvertSuspicious S with(nolock) ON T.Account = S.AccountID
	WHERE T.rowid BETWEEN @intPageSize * (@intPageIndex -1) + 1 AND @intPageSize * @intPageIndex
	UNION ALL
	SELECT 0,0,0,0,0,0,0,'','',NULL,'','','',@HeadLineJson,GETDATE(),GETDATE(),GETDATE(),0,0,0,0,'','',0,'',
	0 as UserID,0 as GradeID,'' as LoginName,'' as Name,0,0,0,0,0 as ClubAgent,0 as Suspicious
	--------------------------------------------------------------------------------------
	ORDER BY T.Id DESC
END
