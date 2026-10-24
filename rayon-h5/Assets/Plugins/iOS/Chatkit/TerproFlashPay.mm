#import "ViewController.h"

extern "C" {
	void _showMessageChat(const char* iUrl,const char* app_key,const char* app_id,const char* user_id,const char* user_name,const char* device_id,const char* timestamp,int agentId,const char* sign,const char* key)
	{
		[[ViewController shareInstance]	showFPChat:[NSString stringWithUTF8String:iUrl]
										agentId:[NSString stringWithFormat:@"%d",agentId]
										appKey:[NSString stringWithUTF8String:app_key]
										appID:[NSString stringWithUTF8String:app_id]
										user_id:[NSString stringWithUTF8String:user_id]
										user_name:[NSString stringWithUTF8String:user_name]
										device_id:[NSString stringWithUTF8String:device_id]
										timestamp:[NSString stringWithUTF8String:timestamp]
										sign:[NSString stringWithUTF8String:sign]
										key:[NSString stringWithUTF8String:key]];
	}
}