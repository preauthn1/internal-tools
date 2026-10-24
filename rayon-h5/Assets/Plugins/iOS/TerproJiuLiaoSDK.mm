#import "JiuLiaoSDK.h"

/*
 加此段資訊加到plist
 
 <key>NSCameraUsageDescription</key>
 <string>此 App 需要您的同意才能使用相机</string>
 
 <key>NSPhotoLibraryUsageDescription</key>
 <string>此 App 需要您的同意才能读取媒体资料库</string>
 
 <key>LSApplicationQueriesSchemes</key>
 <array>
 <string>jiuliao</string>
 </array>
 
 */

extern "C" {	
	// 初始化
	// @param appID App Id
	// @param security 密钥
	void _init(const char* appID, const char* security)
	{
		[[JiuLiaoSDK shareInstance] initJiuLiao:[NSString stringWithUTF8String:appID] 
									security:[NSString stringWithUTF8String:security]];
	}
 
	// 将游戏分享到久聊
	// @param title 分享的题目(分享气泡上最大的字),长度不能超过1K
	// @param desc 描述内容,长度不能超过1K
	// @param appDownloadUrl 你的App的下载链接安卓
	// @param iOSDownloadUrl 你的App的下载链接iOS
	// @param androidBackURL 你的App的URL Schemes (安卓)
	// @param iOSBackUrl 你的App的URL Schemes (iOS)	 
	void _shareToJiuLiaoWithTitle(const char* title,
									const char* desc,
									const char* appDownloadUrl,
									const char* iOSDownloadUrl,
									const char* androidBackURL,
									const char* iOSBackUrl)
	{
	    [[JiuLiaoSDK shareInstance] shareToJiuLiaoWithTitle:[NSString stringWithUTF8String:title]
									description:[NSString stringWithUTF8String:desc]
									appDownloadUrl:[NSString stringWithUTF8String:appDownloadUrl]
									iOSDownloadUrl:[NSString stringWithUTF8String:iOSDownloadUrl]
									androidBackURL:[NSString stringWithUTF8String:androidBackURL]
									iOSBackUrl:[NSString stringWithUTF8String:iOSBackUrl]];
	}
	
	// 向久聊分享图片
    // @param images 需要分享的图片 
	void _shareImages(const char* images)
	{
        UIImage *myUIImage = [UIImage imageNamed:[NSString stringWithUTF8String:images]];
        
        [[JiuLiaoSDK shareInstance] shareImages:@[myUIImage]];
    }
	
	// 跳转到久聊支付
	// @param gameUserId 游戏玩家的账号
	// @param gameName 游戏包名称
	// @param amoutMoney 充值金额 单位：分
	// @param orderId 充值订单号，游戏为每笔充值订单生成的唯一订 单号
	void _goRecharge(const char* gameUserId, const char* gameName, int amoutMoney, const char* orderId)
	{
		[[JiuLiaoSDK shareInstance] goToJiuLiaoPay:[NSString stringWithUTF8String:gameUserId] 
									gameName:[NSString stringWithUTF8String:gameName] 
									amoutMoney:amoutMoney 
									orderId:[NSString stringWithUTF8String:orderId]];
	}	
	
	// 跳转到久聊转账
    // @param gameUserId 游戏玩家的账号
    // @param gameName 游戏包名称	
	void _goTransfer(const char* gameUserId, const char* gameName)
	{
	    [[JiuLiaoSDK shareInstance] goToJiuLiaoTransfer:[NSString stringWithUTF8String:gameUserId]
									gameName:[NSString stringWithUTF8String:gameName]];
	}	
}
