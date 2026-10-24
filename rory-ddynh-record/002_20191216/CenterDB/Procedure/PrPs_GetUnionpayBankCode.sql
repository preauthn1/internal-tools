USE [CenterDB]
GO
/****** Object:  StoredProcedure [dbo].[PrPs_GetUnionpayBankCode]    Script Date: 2019/12/18 上午 11:52:46 ******/
SET ANSI_NULLS ON
GO
SET QUOTED_IDENTIFIER ON
GO
DROP PROCEDURE IF EXISTS dbo.PrPs_GetUnionpayBankCode
GO
-- =============================================
-- Author:  alan
-- Create date: 2018-3-14
-- Description:	获取銀聯卡銀行代號
-- =============================================
CREATE PROCEDURE [dbo].[PrPs_GetUnionpayBankCode]
(
	@intUserID INT = NULL
)
AS
BEGIN
	SELECT B.Code, B.Bankname 
	FROM CenterDB.Money.UnionpayBankCode AS B WITH (NOLOCK)
END
