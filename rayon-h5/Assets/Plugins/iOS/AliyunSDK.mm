#import "YunCeng.h"

#define MakeStringCopy( _x_ ) ( _x_ != NULL && [_x_ isKindOfClass:[NSString class]] ) ? strdup( [_x_ UTF8String] ) : NULL

extern "C" {
    
    int _initEx(const char* appkey, const char* token)
    {
        int ret = [YunCeng initEx:appkey:token];
        return ret;
    }
    
    char* _getProxyTcpByIp(const char* token, const char* groupname, const char* ip, const char* port)
    {
        char target_ip[128]= {0};
        char target_port[32] = {0};
        int ret = [YunCeng getProxyTcpByIp : token : [REDACTED] : ip : port : target_ip : 128 : target_port: 32];
        
        NSString *str = [NSString stringWithFormat:@"%i\n%s\n%s", ret, target_ip, target_port];
        
        return MakeStringCopy(str);
    }
    
    char* _getProxyTcpByDomain(const char* token, const char* groupname, const char* domain, const char* port)
    {
        char target_ip[128]= {0};
        char target_port[32] = {0};
        int ret = [YunCeng getProxyTcpByDomain : token : [REDACTED] : domain : port : target_ip : 128 : target_port : 32];
        
        NSString *str = [NSString stringWithFormat:@"%i\n%s\n%s", ret, target_ip, target_port];
        
        return MakeStringCopy(str);
    }
    
    int _reportUserData(int type, const char* msg, int sync)
    {
        int ret = [YunCeng reportUserData: type : msg : sync];
        return ret;
    }
    
    char* _GetSession()
    {
        char session[384] = {0};
        int ret = [YunCeng GetSession:session :384];
        
        NSString *str = [NSString stringWithFormat:@"%i\n%s", ret, session];
        
        return MakeStringCopy(str);
    }
}
