//
//  lua-cjson.h
//  ares
//
//  Created by huachangmiao.
//  Copyright (c) 2016年 playcrab. All rights reserved.
//

#ifndef ares_lua_hjson_h
#define ares_lua_hjson_h

#define ENABLE_HJSON_GLOBAL 1

#ifdef __cplusplus
extern "C" {
#endif
    
int luaopen_hjson(lua_State *l);
    
#ifdef __cplusplus
}
#endif

#endif
