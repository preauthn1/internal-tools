USE [LogDB]
GO
/****** Object:  StoredProcedure [dbo].[PrPs_AddConvertLog]    Script Date: 2019/12/18 W 11:44:11 ******/
SET ANSI_NULLS ON
GO
SET QUOTED_IDENTIFIER ON
GO
DROP PROCEDURE IF EXISTS dbo.PrPs_AddConvertLog
GO
-- =============================================
-- Author:		Jason
-- Create date: 2018/03/22
-- Description:	gJϥΪ̧I
--	2018/04/24	jason	sWbankamount,cashamount,opreator gJ
--	2019/12/16		mars			addConvertLogW[BankCodeID
-- =============================================
CREATE PROCEDURE [dbo].[PrPs_AddConvertLog]
(
	@intUserID			INT,
	@lngAmount			BIGINT,
	@lngTaxAmount		BIGINT,
	@lngRealTaxAmount	BIGINT= 0,
	@tnyStatus			TINYINT,
	@chvRealName		VARCHAR(64),
	@tnyConvertType		TINYINT,
	@chvAccount			VARCHAR(64),
	@intBankCodeId		INT = NULL,
	@chvBankName		VARCHAR(64),
	@chvIP				VARCHAR(15),
	@dtStatTime			DATETIME,
	@intPlatformID		INT = NULL,
	@intAreaID			INT = NULL,
	@lngUserCashAmount	BIGINT = NULL,
	@lngUserBankAmount	BIGINT = NULL
)
AS
BEGIN
	SET NOCOUNT ON

	
	DECLARE @UserID				INT = @intUserID
	DECLARE @Amount				BIGINT = @lngAmount
	DECLARE @TaxAmount			BIGINT = @lngTaxAmount
	DECLARE @RealTaxAmount		BIGINT = @lngRealTaxAmount
	DECLARE @Status				TINYINT = @tnyStatus
	DECLARE @RealName			VARCHAR(64) = @chvRealName
	DECLARE @ConvertType		TINYINT = @tnyConvertType
	DECLARE @PayAccount			VARCHAR(64) = ISNULL(@chvAccount,' ')
	DECLARE @BankCodeId			INT = @intBankCodeId
	DECLARE @BankName			VARCHAR(64) = ISNULL(@chvBankName,' ')
	DECLARE @IP					VARCHAR(15) = @chvIP
	DECLARE @StatTime			DATETIME = @dtStatTime
	DECLARE @PlatformID			INT = @intPlatformID
	DECLARE @AreaID				INT = @intAreaID
	DECLARE @UserBankAmount		BIGINT = @lngUserBankAmount
	DECLARE @UserCashAmount		BIGINT = @lngUserCashAmount

	
	/*  gJIӽЬ  */
	------------------------------------------------------
	INSERT INTO [Money].[ConvertLog]
		(	
			UserID, Amount, TaxAmount, RealTaxAmount, Status, RealName, ConvertType,
			Account, BankName, CreateIP, CreateTime, LastUpdateTime, CheckTime, [PlatformID], [AreaID], [BankAmount], [CashAmount], 
			[BankCodeID]
		)
	VALUES
		(
			@UserID, @Amount, @TaxAmount, @RealTaxAmount, @Status, @RealName, @ConvertType,
			@PayAccount, @BankName,	@IP, @StatTime,	@StatTime, @StatTime, @PlatformID, @AreaID, @UserBankAmount, @UserCashAmount,
			@BankCodeId
		)
	------------------------------------------------------

	RETURN
END