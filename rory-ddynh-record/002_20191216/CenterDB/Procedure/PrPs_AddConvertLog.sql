USE [CenterDB]
GO
/****** Object:  StoredProcedure [dbo].[PrPs_AddConvertLog]    Script Date: 2019/12/18 上午 11:51:13 ******/
SET ANSI_NULLS ON
GO
SET QUOTED_IDENTIFIER ON
GO
DROP PROCEDURE IF EXISTS dbo.PrPs_AddConvertLog
GO
-- =============================================
-- Author:  <xxx>
-- Create date: <xxx>
-- Description:	<提交兑换申请>
--			{0:"兑换完成",1:"审核中",2:"审核通过",3:"已打款",4:"订单异常",5 :"异常帐号",6:"黑名单",7:"账户错误",8:"拒绝兑换",9:"疑似洗钱",10:"补单完成",11:"退还金币"};
-- Changelog:
	
--		2018/03/22		Jason			調整語法改為呼叫LOGDB SP執行LOG寫入、新增@dtStatTime & @lngRealTaxAmount、新增@chvUnionpayBankName
--		2018/03/27		alan			增加錢包最低金額 <您的现金不足，无法兑换,兑换数量与手续费一共需要'+dbo.Fn_ShowMoney(ABS(@lngDeductAmount) + @lngConvertRetainAmount)>
--		2018/04/10		Jason			新增UpdateUserCash時寫入AreaID & PlatformID
--		2018/04/17		Jason			新增PrPs_AddConvertLog時寫入AreaID & PlatformID
--		2018/04/24		Jason			新增PrPs_AddConvertLog時寫入CashAmount & BankAmount
--		2018/05/03		alan			判斷如果大於自動審核金額, 需進行人工審核
--		2018/05/16		alan			修正銀行名稱帶入bankcode錯誤, 應該記錄銀行名稱
--		2019/08/21		mars			增加兌換前的稽核流水判斷
--		2019/10/29		mars			修正系统配置for platform #148
--		2019/10/29		marshall		稽核流水有個判斷式移除
--		2019/12/17		mars			addConvertLog增加BankCodeID
-- =============================================
CREATE PROCEDURE [dbo].[PrPs_AddConvertLog]
@intUserID			INT,					--用户ID
@chvCer				VARCHAR(32),			--登录凭证
@lngAmount			BIGINT,					--兑换数量,100表示1.00
@tnyConvertType		TINYINT,			--選擇兌換類型(支付寶1, 銀聯卡2)
@chvIP				VARCHAR(15),			--客户端IP
@chvErrMsg			VARCHAR(256) OUTPUT		--返回消息
AS
BEGIN

	SET NOCOUNT ON
	Declare @chvLoginName varchar(32)=''
	DECLARE @tnyUserType tinyint
	DECLARE @chvRealName			VARCHAR(64)				--支付宝真实姓名
	DECLARE @chvAlipayAccount		VARCHAR(64)				--支付宝账号
	DECLARE @chvUnionpayAccount		VARCHAR(64)				--銀聯卡帳號
	DECLARE @tnyStatus				TINYINT					--兑换状态(在兑换开启的情况，同时没有超过限制配置，则直接为审核通过)
	DECLARE @tnyConvertStatus		TINYINT					--兑换开关状态 1开启 0未开启
	DECLARE @chvPromptMsg			VARCHAR(256)			--未开启时需要给客户端的提示消息
	DECLARE @lngAutoCheckLimit		BIGINT					--单次兑换限制金额	
	DECLARE @lngConvertRetainAmount	BIGINT					--兑换保留金币数
	DECLARE @lngConvertMinAmount	BIGINT					--兑换最小金币数
	DECLARE @lngConvertMaxAmount	BIGINT					--兑换最大金币数
	DECLARE @lngTaxAmount			BIGINT					--兑换手续费
	DECLARE @lngDeductAmount		BIGINT					--用户需要扣除的数量,是负数
	DECLARE @chvUnionpayBankCode	VARCHAR(5)				--銀行代號
	DECLARE @chvBankname			varchar(50) = ''		--銀行名稱
	DECLARE	@chvUnionpayRealName	VARCHAR(64)				--銀行帳戶名稱
	DECLARE @dtStatTime				DATETIME = GETDATE()	--用戶提出申請時間
	DECLARE @lngRealTaxAmount		BIGINT = 0				--預設手續費為0
	DECLARE @intGameID				INT						--遊戲ID 目前為空值
	DECLARE @intAreaID				INT						--渠道
	DECLARE @intPlatformID			INT
	DECLARE	@lngCurrentAmount		BIGINT
	DECLARE @UserCashAmount			BIGINT					--現金
	DECLARE @UserBankAmount			BIGINT					--銀行金額
		
	SET @lngAmount = ABS(@lngAmount)
	SET @lngTaxAmount = dbo.Fn_GetConvertTax(@lngAmount)
	SET @lngDeductAmount = -1*(@lngAmount+@lngTaxAmount)
	
	---從上次提款到現在的流水限制確認開始
	DECLARE	@intLastConvertFlowDiff int = 0,
			@intLastConvertFlowAmountLimit int = 0,
			@intSumFlowAmountLimit int = 0,
			@dateLastConvertTime date,
			@dateNOW date,
			@bitLastConvertFlowAmountCheck bit = 0,
			@SumFlowAmount BIGINT = 0,
			@FlowAmountLimit BIGINT = 0
	
	--取得用戶渠道/總流水
	SELECT --@intUserAreaID = AreaID,
			@SumFlowAmount = SumFlowAmount
	FROM Users.UserInfo WITH (NOLOCK) WHERE UserID = @intUserID

	EXEC [CenterDB].[dbo].[PrWeb_GetLastConvertFlowAmountAndLastConvertTime]
		 @intUserID = @intUserID,
		 @intLastConvertFlowDiff = @intLastConvertFlowDiff				OUTPUT,	--從最後一次提領到現在遊戲輸贏的流水(預設總流水 - 上次提領時的流水)
		 @intLastConvertFlowAmountLimit = @intLastConvertFlowAmountLimit		OUTPUT,	--此欄位已作廢、沒有意義了
		 @intSumFlowAmountLimit = @intSumFlowAmountLimit						OUTPUT,	--從最後一次提領到現在的稽核流水限額(mars)
		 @dateLastConvertTime = @dateLastConvertTime							OUTPUT,	--最後一次提款的放款時間
		 @dateNOW = @dateNOW													OUTPUT	--目前的時間

	---首存活動的返利資格確認開始
	DECLARE @ChargeBonusInfo TABLE (FlowAmountNow BIGINT,FlowAmountLimit BIGINT,BonusRate INT,FlowAmountRate INT,EventStatus INT)
	DECLARE @UnionPlatform INT
	
	SET @SumFlowAmount = ISNULL(@SumFlowAmount,0)
	
	INSERT INTO @ChargeBonusInfo
	EXEC [dbo].[PrPs_Event_ChargeBonusAmount_Get]
		@intUserID = @intUserID,
		@chvMsg = ''
	
	--Mars 稽核流水
	IF (@intLastConvertFlowDiff < @intSumFlowAmountLimit)
	BEGIN
		SET @chvErrMsg = '未达所需稽核流水，请联系客服';
		RETURN -1
	END

	SELECT @FlowAmountLimit = FlowAmountLimit FROM @ChargeBonusInfo
	IF (@SumFlowAmount < @FlowAmountLimit)
	BEGIN
		SET @chvErrMsg = '未达所需稽核流水，请联系客服';
		RETURN -1
	END
	
	SELECT	
			@tnyUserType = UserType,
			@chvLoginName = LoginName,
			@chvRealName = RealName,
			@chvAlipayAccount = AlipayAccount,
			@chvUnionpayAccount = UnionpayAccount,
			@chvUnionpayBankCode = [UnionpayBank],
			@chvUnionpayRealName = [UnionpayRealName],
			@intAreaID = [AreaID]
	FROM Users.UserInfo WITH (NOLOCK)
	WHERE
		UserID = @intUserID

	SELECT @intPlatformID = [PlatformID]
	FROM [Area].[Area] WITH (NOLOCK)
	WHERE
		[AreaID] = @intAreaID

	SET @UserCashAmount	= ISNULL((SELECT [Amount] FROM [Money].[Cash] WITH (NOLOCK) WHERE [UserID] = @intUserID),0)
	SET @UserBankAmount	= ISNULL((SELECT [Amount] FROM [Money].[Bank] WITH (NOLOCK) WHERE [UserID] = @intUserID),0)


	--检查是否在游戏房间中
	IF EXISTS(SELECT 1 FROM Game.GameUserList WITH(NOLOCK) WHERE UserID=@intUserID)
	BEGIN
		SET @chvErrMsg = '操作失败:帐号在游戏房间中'
		RETURN -1
	END

	IF @tnyUserType = 1
	BEGIN
		--檢查是否被鎖定
		IF EXISTS(select 1 from Game.ForbidAccount with (nolock) where AccountID = @chvLoginName)
		BEGIN
			SET @chvErrMsg = '您的手机号兑换功能已被封锁, 请与客服联系'
			RETURN -1
		END

	END

	IF @tnyConvertType = 1
	BEGIN

		--如果没有绑定支付宝账号返回提示信息
		IF (LEN(ISNULL(@chvRealName,''))<1 OR LEN(ISNULL(@chvAlipayAccount,''))<1)
		BEGIN
			SET @chvErrMsg = '您的账号未绑定支付宝账号，请绑定后再操作'
			RETURN -1
		END

		--是否正式帳號
		IF @tnyUserType = 1
		BEGIN
			--檢查是否被鎖定
			IF EXISTS(select 1 from Game.ForbidAccount with (nolock) where IdType=2 and AccountID = @chvAlipayAccount)
				BEGIN
					SET @chvErrMsg = '您的支付宝帐号兑换功能已被封锁, 请与客服联系'
					RETURN -1
				END
		END

		SET @chvBankname = ''

	END
	ELSE IF @tnyConvertType = 2
	BEGIN

		--如果沒有綁定銀行卡帳號返回提示信息
		IF (LEN(ISNULL(@chvUnionpayAccount,''))<1 OR LEN(ISNULL(@chvUnionpayBankCode,''))<1 OR LEN(ISNULL(@chvUnionpayRealName,''))<1)
		BEGIN
			SET @chvErrMsg = '您的账号未绑定銀行卡账号，请绑定后再操作'
			RETURN -1
		END
		
		--是否正式帳號
		IF @tnyUserType = 1
		BEGIN
			--檢查是否被鎖定
			IF EXISTS(select 1 from Game.ForbidAccount with (nolock) where IdType=3 and AccountID = @chvUnionpayAccount)
				BEGIN
					SET @chvErrMsg = '您的银行帐号兑换功能已被封锁, 请与客服联系'
					return -1
				END
		END

		SELECT @chvBankname = Bankname from  [Money].[UnionpayBankCode] with (nolock) where Code=@chvUnionpayBankCode
		IF @@Rowcount <> 1 
		BEGIN
			SET @chvErrMsg = '银行名称错误'
			RETURN -1
		END
	END
	ELSE
	BEGIN
			SET @chvErrMsg = '兌換類型發生錯誤'
			RETURN -1
	END
	
	SELECT TOP 1	@tnyConvertStatus =	ConvertStatus,
					@chvPromptMsg = PromptMsg,
					@lngAutoCheckLimit = AutoCheckLimit,
					@lngConvertRetainAmount = ConvertRetainAmount,
					@lngConvertMinAmount = ConvertMinAmount,
					@lngConvertMaxAmount = ConvertMaxAmount
	FROM ServerInfo.SystemConfig WITH (NOLOCK)
	WHERE PlatformID in (0, @intPlatformID)
	ORDER BY PlatformID DESC
	
	IF @tnyConvertStatus <> 1
	BEGIN
		SET @chvErrMsg = @chvPromptMsg
		RETURN -1
	END
	
	IF @lngAmount%5000 <>0
	BEGIN
		SET @chvErrMsg = '兑换金额必须是' + dbo.Fn_ShowMoney(5000) + '的倍数'
		RETURN -1
	END
	
	IF @lngConvertMinAmount>@lngAmount
	BEGIN
		SET @chvErrMsg = '兑换金额不能低于'+dbo.Fn_ShowMoney(@lngConvertMinAmount)
		RETURN -1
	END
	
	IF @lngConvertMaxAmount<@lngAmount
	BEGIN
		SET @chvErrMsg = '兑换金额不能高于'+dbo.Fn_ShowMoney(@lngConvertMaxAmount)
		RETURN -1
	END
	
	IF EXISTS(SELECT 1 FROM Money.Cash WHERE Amount+@lngDeductAmount<0 AND UserID=@intUserID)
	BEGIN
		SET @chvErrMsg = '您的现金不足，无法兑换,兑换数量与手续费一共需要'+dbo.Fn_ShowMoney(ABS(@lngDeductAmount) + @lngConvertRetainAmount)
		RETURN -1
	END
	
	IF EXISTS(SELECT 1 FROM Money.Cash WITH(NOLOCK) WHERE UserID=@intUserID AND Amount<@lngConvertRetainAmount+ABS(@lngDeductAmount))
	BEGIN
		SET @chvErrMsg = '您的保留现金不足' + dbo.Fn_ShowMoney(@lngConvertRetainAmount) + '，无法兑换'
		RETURN -1
	END	

	SET @tnyStatus = 2	--默认为审核通过

	--如果大於自動審核金額, 需進行人工審核
	IF @lngAmount>@lngAutoCheckLimit
	BEGIN
		SET @tnyStatus = 0
	END
	
	--扣玩家金幣
	SET @lngAmount = -1 * ABS(@lngAmount)
	EXEC dbo.PrPsWeb_UpdateUserCash 
		@intUserID			= @intUserID,
		@sintSourceType		= 14000,
		@lngUpdateAmount	= @lngDeductAmount,
		@chvIPAddress		= @chvIP,
		@chvMemo			= '兑换扣金币',
		@lngCurrentAmount	= @lngCurrentAmount	OUTPUT,
		@lngTaxAmount		= 0,
		@intAreaID			= @intAreaID,
		@intPlatformID		= @intPlatformID

	/*  --写兑换申请记录  */
	-------------------------------------------------------------------------------------------
	--INSERT INTO Log.ConvertLog(UserID,Amount,TaxAmount,RealTaxAmount,Status,RealName,ConvertType,UnionpayAccount, AlipayAccount,CreateIP,CreateTime,LastUpdateTime,CheckTime)
	--SELECT @intUserID,@lngAmount,@lngTaxAmount,0,@tnyStatus,@chvRealName,
	--@tnyConvertType,@chvUnionpayAccount,
	--@chvAlipayAccount, @chvIP,GETDATE(),GETDATE(),GETDATE()
	--FROM Users.UserInfo WITH(NOLOCK) 
	--WHERE UserID = @intUserID

	DECLARE @RealName	VARCHAR(64)
	DECLARE @Account	VARCHAR(64)
	
	SET @RealName = CASE WHEN @tnyConvertType = 1 THEN @chvRealName ELSE @chvUnionpayRealName END
	SET @Account = CASE WHEN @tnyConvertType = 1 THEN @chvAlipayAccount ELSE @chvUnionpayAccount END
	
	DECLARE @intBankCodeID INT = CONVERT(int, @chvUnionpayBankCode)
	
	EXEC LogDB.dbo.PrPs_AddConvertLog
		@intUserID			= @intUserID,
		@lngAmount			= @lngAmount,
		@lngTaxAmount		= @lngTaxAmount,
		@lngRealTaxAmount	= @lngRealTaxAmount,
		@tnyStatus			= @tnyStatus,
		@chvRealName		= @RealName,
		@tnyConvertType		= @tnyConvertType,
		@chvAccount			= @Account,
		@intBankName		= @intBankCodeID,
		@chvBankName		= @chvBankname,
		@chvIP				= @chvIP,
		@dtStatTime			= @dtStatTime,
		@intPlatformID		= @intPlatformID,
		@intAreaID			= @intAreaID,
		@lngUserCashAmount	= @UserCashAmount,
		@lngUserBankAmount	= @UserBankAmount
	-------------------------------------------------------------------------------------------
	--更新用户数据
	UPDATE Users.UserInfo
	SET SumConvert = SumConvert + ABS(@lngAmount)
		,SumConvertTimes = SumCOnvertTimes + 1
		,SumConvertTax = SumConvertTax + @lngTaxAmount
	WHERE
		UserID = @intUserID
	
	SET @chvErrMsg = '兑换成功，请等待3-5分钟后到账'

	--检查输赢是否过大警报提醒
	DECLARE @lngWinLost				BIGINT					--玩家总输赢 总兑换减总充值
	DECLARE @lngWinLostLimit		BIGINT					--系统总输赢警报值
	DECLARE @fltMultiple			FLOAT					--输赢倍率 
	DECLARE @lngSumCharge			BIGINT
	DECLARE @lngSumConvert			BIGINT
		
	SELECT
	@lngWinLost = SumConvert-(SumCharge+SumAgentCharge),
	@lngSumCharge = SumCharge+SumAgentCharge,
	@lngSumConvert = SumConvert	
	FROM Users.UserInfo WITH (NOLOCK)
	WHERE UserID = @intUserID

	--计算倍数
	IF @lngSumConvert>=@lngSumCharge
		SET @fltMultiple=CAST(@lngSumConvert as FLOAT)/CAST(CASE WHEN @lngSumCharge=0 THEN 1 ELSE @lngSumCharge END as FLOAT)
	ELSE
		SET @fltMultiple=CAST(@lngSumCharge as FLOAT)/CAST(CASE WHEN @lngSumConvert=0 THEN 1 ELSE @lngSumConvert END as FLOAT)

	--保留小数点两位
	SET @fltMultiple=convert(numeric(17,2),ROUND(@fltMultiple,2))
	
	SELECT TOP 1 @lngWinLostLimit = ConvertWinLostAlert
	FROM ServerInfo.SystemConfig WITH (NOLOCK)
	WHERE PlatformID in (0, @intPlatformID)
	ORDER BY PlatformID DESC

	IF ABS(@lngWinLost) >= @lngWinLostLimit
	BEGIN
		--支付宝
		IF @tnyConvertType = 1
		BEGIN
			INSERT INTO Notify.Notify.ConvertAlert(UserID,UserName,WinLost,Status,Multiple) values(@intUserID, @chvRealName, @lngWinLost, 1, @fltMultiple)
		END
		--银行
		IF @tnyConvertType = 2
		BEGIN
			INSERT INTO Notify.Notify.ConvertAlert(UserID,UserName,WinLost,Status,Multiple) values(@intUserID, @chvUnionpayRealName, @lngWinLost, 1, @fltMultiple)
		END
	END
	
	RETURN 1
END