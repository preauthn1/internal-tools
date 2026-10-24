USE [LogDB]
GO

/****** Object:  Table [dbo].[SlotGameHistorySeqence]    Script Date: 2019/12/13 W 09:53:48 ******/
DROP TABLE [dbo].[SlotGameHistorySeqence]
GO

/****** Object:  Table [dbo].[SlotGameHistorySeqence]    Script Date: 2019/12/13 W 09:53:48 ******/
SET ANSI_NULLS ON
GO

SET QUOTED_IDENTIFIER ON
GO

CREATE TABLE [Game].[SlotGameHistorySeqence](
	[ID] [int] IDENTITY(1,1) NOT FOR REPLICATION NOT NULL,
 CONSTRAINT [PK_SlotGameHistorySeqence] PRIMARY KEY CLUSTERED 
(
	[ID] ASC
)WITH (PAD_INDEX = OFF, STATISTICS_NORECOMPUTE = OFF, IGNORE_DUP_KEY = OFF, ALLOW_ROW_LOCKS = ON, ALLOW_PAGE_LOCKS = ON) ON [PRIMARY]
) ON [PRIMARY]
GO


