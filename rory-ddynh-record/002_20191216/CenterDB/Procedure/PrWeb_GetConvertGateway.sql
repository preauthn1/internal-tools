USE [CenterDB]
GO
/****** Object:  StoredProcedure [dbo].[PrWeb_GetConvertGateway]    Script Date: 2019/12/18 上午 11:53:54 ******/
SET ANSI_NULLS ON
GO
SET QUOTED_IDENTIFIER ON
GO
DROP PROCEDURE IF EXISTS dbo.PrWeb_GetConvertGateway
GO
-- =============================================
-- Author:  Marshall
-- Create date: 2019-12-02
-- Description:	获取代付廠商
--	2019/12/17	Mars			增加關聯Pay.ConvertGatewayBankSetting/Payment.PayGatewayConvertTypeMapping取得該筆兌換當時所能選擇的代付渠道
-- =============================================
CREATE PROCEDURE [dbo].[PrWeb_GetConvertGateway]	
(
	@intAdminId INT,
	@intConvertType INT = NULL,
	@intUserId INT = NULL,
	@intBankCodeID INT = NULL
)
AS
BEGIN

	SET NOCOUNT ON
	
	IF (@intUserId IS NULL)
	BEGIN
		SELECT DISTINCT G.ID,
				G.GatewayName,
				1 as ConvertType,
				1 as [Status],
				1 as Company
		FROM Payment.PayGateway as G 
		WHERE G.SupportPayout = 'true'
	END
	ELSE
	BEGIN
		SELECT DISTINCT G.ID,
				G.GatewayName,
				C.ConvertType,
				S.[Status],
				1 as Company
		FROM Payment.PayGateway as G 
		JOIN Payment.PayGatewayConvertTypeMapping as C 
			ON C.GWID = G.ID 
			AND (@intConvertType IS NULL OR C.ConvertType = @intConvertType)
		JOIN Pay.ConvertGatewayBankSetting as S
			ON S.ConvertGatewayID = G.ID 
			AND S.Status = 1 
			AND (@intBankCodeID IS NULL OR S.UnionPayBankCodeId = @intBankCodeID) 
		JOIN Payment.MerchantInfo AS M 
			ON M.GWID = G.ID
			AND M.PlatformID IN (SELECT PlatformID FROM dbo.Fn_GetPlatformInfoByUserID(@intUserId))
	END
END


