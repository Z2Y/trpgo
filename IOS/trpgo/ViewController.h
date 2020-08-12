//
//  ViewController.h
//  trpgo
//
//  Created by Z2Y on 2020/8/5.
//  Copyright © 2020 Z2Y. All rights reserved.
//

#import <UIKit/UIKit.h>
#import <GLKit/GLKit.h>

#define GLES_SILENCE_DEPRECATION true

@interface EngoViewController : UIViewController

@property (strong, nonatomic) EAGLContext *context;

- (void)suspendGame;

- (void)resumeGame;

@end

