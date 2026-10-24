//
//  ViewController.h
//  Chat
//
//  Created by User on 2019/3/5.
//  Copyright © 2019年 User. All rights reserved.
//
#import <UIKit/UIKit.h>

@interface ViewController : UIViewController

+ (instancetype) shareInstance;

- (void)showFPChat:(NSString *)iUrl agentId:(NSString *)agentId appKey:(NSString *)app_key appID:(NSString *)app_id user_id:(NSString *)user_id user_name:(NSString *)user_name device_id:(NSString *)device_id timestamp:(NSString *)timestamp sign:(NSString *)sign key:(NSString *)key;

@end

