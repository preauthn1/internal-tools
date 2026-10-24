USE [LogDB]
GO

/****** Object:  StoredProcedure [dbo].[PrWeb_Slot_GetSequence]    Script Date: 2019/12/13 W 10:09:44 ******/
SET ANSI_NULLS ON
GO

SET QUOTED_IDENTIFIER ON
GO





-- =============================================
-- Author:		Jason
-- Create date: 2019/05/24
-- Description:	SlotGameHistoryLog
-- 20191213	Louis	קSlotGameHistorySeqenceSchemaName
-- =============================================
ALTER PROCEDURE [dbo].[PrWeb_Slot_GetSequence]
(
	@intLogID INT OUTPUT
)
AS
BEGIN
	SET NOCOUNT ON

	
	BEGIN TRAN
		INSERT [Game].[SlotGameHistorySeqence] WITH (TABLOCKX) DEFAULT VALUES;
    ROLLBACK TRAN

    SELECT @intLogID = SCOPE_IDENTITY()

	RETURN
END
GO


