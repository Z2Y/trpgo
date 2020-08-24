// +build darwin
// +build ios

#include "_cgo_export.h"

#import <Foundation/Foundation.h>

const char* nsstring2cstring(NSString *s) {
    if (s == NULL) { return NULL; }

    const char *cstr = [s UTF8String];
    return cstr;
}


const char* mainBundlePath(void) {
  return nsstring2cstring([[NSBundle mainBundle] bundlePath]);
}