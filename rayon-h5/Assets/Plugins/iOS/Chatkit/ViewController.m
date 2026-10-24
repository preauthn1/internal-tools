//
//  ViewController.m
//  Chat
//
//  Created by User on 2019/3/5.
//  Copyright © 2019年 User. All rights reserved.
//

#import "ViewController.h"
#import <ChatKit/ChatSDK.h>
#import <CommonCrypto/CommonDigest.h>

#define FileHashDefaultChunkSizeForReadingData 1024*1

static const char encodingTable[] = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789+/";
static NSString *secret = @"9d849053b3dbe374685af44f525e8f52";

// @interface ViewController ()

// @end

static ViewController * _instance = nil;
@implementation ViewController

+ (instancetype)shareInstance {
    static dispatch_once_t onceToken ;
    dispatch_once(&onceToken, ^{
        _instance = [[super allocWithZone:NULL] init] ;
    }) ;
    
    return _instance ;
}

+ (id)allocWithZone:(struct _NSZone *)zone {
    return [ViewController shareInstance] ;
}

- (id)copyWithZone:(struct _NSZone *)zone {
    return [ViewController shareInstance] ;
}

- (void)viewDidLoad {
    [super viewDidLoad];
    // Do any additional setup after loading the view, typically from a nib.
}

- (void)showFPChat:(NSString *)iUrl agentId:(NSString *)agentId appKey:(NSString *)app_key appID:(NSString *)app_id user_id:(NSString *)user_id user_name:(NSString *)user_name device_id:(NSString *)device_id timestamp:(NSString *)timestamp sign:(NSString *)sign key:(NSString *)key {
    secret = [REDACTED];
    
    NSDictionary * iData = @{@"id":agentId
                             ,@"app_key":app_key
                             ,@"app_id":app_id
                             ,@"user_id":user_id
                             ,@"user_name":user_name
                             ,@"device_id":device_id
                             ,@"timestamp":timestamp
                             ,@"sign":sign
                             };
    NSMutableDictionary *params = [iData mutableCopy];
    if ([timestamp  isEqual: @"0"]) {
        params =                        [@{
                                           @"id":agentId
                                           ,@"app_key":app_key
                                           ,@"app_id":app_id
                                           ,@"user_id":user_id
                                           ,@"user_name":user_name
                                           ,@"device_id":device_id
                                           ,@"timestamp": [self getNowTimeTimestamp]
                                           } mutableCopy];
        
        params[@"sign"] = [self presignString:params];
    }
    
    //     NSLog(@"params.id           = %@", params[@"id"]);
    //     NSLog(@"params.app_key      = %@", params[@"app_key"]);
    //     NSLog(@"params.app_id       = %@", params[@"app_id"]);
    //     NSLog(@"params.user_id      = %@", params[@"user_id"]);
    //     NSLog(@"params.user_name    = %@", params[@"user_name"]);
    //     NSLog(@"params.device_id    = %@", params[@"device_id"]);
    //     NSLog(@"params.timestamp    = %@", params[@"timestamp"]);
    //     NSLog(@"params.sign         = %@", params[@"sign"]);
    
    [ChatSDK showMessageChat:iUrl params:params isLandscape:TRUE];
}

- (IBAction)goChat1:(id)sender {
    //测试参数
    NSMutableDictionary *params = [@{
                                     @"id":@"193"
                                     ,@"app_key":@"10003"
                                     ,@"app_id":@"59"
                                     ,@"user_id":@"133132"
                                     ,@"user_name":@"AZLDC"
                                     ,@"device_id":@"iPhoneX"
                                     ,@"timestamp": [self getNowTimeTimestamp]
                                     } mutableCopy];
    
    //客户端测试签名，实际场景需要在服务端执行签名操作
    //先对所有参数的key进行字符串ASCII码升序排序
    //然后拼接成带签名的字符串 k1=v1&k2=v2
    //在待签名的字符串后面拼接上密钥（一个app_key 对应一个密钥）进行 md5 得到小写 sign
    params[@"sign"] = [self presignString:params];
    
    //唤起聊天消息界面
    [ChatSDK showMessageChat:@"https://example.invalid/index/login/auto_login" params:params isLandscape:TRUE];
}

- (IBAction)goChat2:(id)sender {
    
    //测试参数
    NSMutableDictionary *params = [@{
                                     @"id":@"-1"
                                     ,@"app_key":@"10003"
                                     ,@"app_id":@"59"
                                     ,@"user_id":@"133132"
                                     ,@"user_name":@"AZLDC"
                                     ,@"device_id":@"iPhoneX"
                                     ,@"timestamp": [self getNowTimeTimestamp]
                                     } mutableCopy];
    
    params[@"sign"] = [self presignString:params];
    
    //唤起聊天消息界面
    [ChatSDK showMessageChat:@"https://example.invalid/index/login/auto_login" params:params isLandscape:TRUE];
}

-(NSString *)getNowTimeTimestamp{
    
    NSDateFormatter *formatter = [[NSDateFormatter alloc] init] ;
    
    [formatter setDateStyle:NSDateFormatterMediumStyle];
    
    [formatter setTimeStyle:NSDateFormatterShortStyle];
    
    [formatter setDateFormat:@"YYYY-MM-dd HH:mm:ss"]; // ----------设置你想要的格式,hh与HH的区别:分别表示12小时制,24小时制
    
    //设置时区,这个对于时间的处理有时很重要
    
    NSTimeZone* timeZone = [NSTimeZone timeZoneWithName:@"Asia/Beijing"];
    
    [formatter setTimeZone:timeZone];
    
    NSDate *datenow = [NSDate date];//现在时间,你可以输出来看下是什么格式
    
    NSString *timeSp = [NSString stringWithFormat:@"%ld", (long)[datenow timeIntervalSince1970]];
    
    return timeSp;
    
}

-(NSString*)presignString:(NSMutableDictionary*)params{
    NSMutableArray *allKeys = [params.allKeys mutableCopy];
    if(allKeys.count==0){
        return @"";
    }
    [allKeys sortUsingSelector:@selector(compare:)];
    NSMutableArray *keyValue = [NSMutableArray array];
    for (NSString* key in allKeys) {
//        if([key isEqualToString:@"id"] && [params[key] isEqualToString:@"-1"]){
//            continue;
//        }
        [keyValue addObject:[NSString stringWithFormat:@"%@=%@", key, params[key]]];
    }
    NSString* param = [[keyValue componentsJoinedByString:@"&"] stringByTrimmingCharactersInSet:[NSCharacterSet whitespaceCharacterSet]] ;
    
    NSString *saltString = [NSString stringWithFormat:@"%@%@", param, secret];
    NSString *presignString = [self MD5:saltString];
    return presignString;
}

- (NSString*) MD5:(NSString*)str {
    
    const char *cStr = [str UTF8String];
    unsigned char digest[CC_MD5_DIGEST_LENGTH];
    CC_MD5( cStr, strlen(cStr), digest );
    
    NSMutableString *output = [NSMutableString stringWithCapacity:CC_MD5_DIGEST_LENGTH * 2];
    
    for(int i = 0; i < CC_MD5_DIGEST_LENGTH; i++){
        [output appendFormat:@"%02x", digest[i]];
    }
    return output;
}

@end
