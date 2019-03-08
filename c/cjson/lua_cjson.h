//
//  lua-cjson.h
//  ares
//
//  Created by liubin ouyang on 12-4-21.
//  Copyright (c) 2012年 playcrab. All rights reserved.
//

#ifndef ares_lua_cjson_h
#define ares_lua_cjson_h

#define ENABLE_CJSON_GLOBAL 1

#ifdef __cplusplus
extern "C" {
#endif
    
int luaopen_cjson(lua_State *l);
    
#ifdef __cplusplus
}
#endif

#endif
