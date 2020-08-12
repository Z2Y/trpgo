//
//  ViewController.m
//  trpgo
//
//  Created by Z2Y on 2020/8/5.
//  Copyright © 2020 Z2Y. All rights reserved.
//

#import "ViewController.h"
#import "Bridge/Bridge.h"

#define TOUCH_TYPE_BEGIN 0 // touch.TypeBegin
#define TOUCH_TYPE_MOVE  1 // touch.TypeMove
#define TOUCH_TYPE_END   2 // touch.TypeEnd

@interface EngoViewController ()

@end

@implementation EngoViewController {
    GLKView* glkView_;
    bool     started_;
    bool     active_;
}

- (GLKView*) glkView {
    if (!glkView_) {
        glkView_ = [[GLKView alloc] init];
        glkView_.multipleTouchEnabled = YES;
        glkView_.userInteractionEnabled = YES;
    }
  return glkView_;
}


- (void)viewDidLoad {
    [super viewDidLoad];
    
    if (!started_) {
      @synchronized(self) {
        active_ = true;
      }
      started_ = true;
    }
    
    self.context = [[EAGLContext alloc] initWithAPI:kEAGLRenderingAPIOpenGLES2];
    [self glkView].delegate = (id<GLKViewDelegate>)(self);
    [self glkView].context = self.context;
    [self.view addSubview: self.glkView];
    [EAGLContext setCurrentContext:self.context];
    // self.glkView.enableSetNeedsDisplay = NO;

    CADisplayLink *displayLink = [CADisplayLink displayLinkWithTarget:self selector:@selector(drawFrame)];
    [displayLink addToRunLoop:[NSRunLoop currentRunLoop] forMode:NSDefaultRunLoopMode];
}

- (void)viewWillLayoutSubviews {
  CGRect viewRect = [[self view] frame];
  [[self glkView] setFrame:viewRect];
}

- (void)viewDidLayoutSubviews {
    [super viewDidLayoutSubviews];
    CGRect viewRect = [[self view] frame];
    if (!BridgeIsRunning()) {
          BridgeStart(viewRect.size.width, viewRect.size.height);
    }
}

- (void)drawFrame{
  @synchronized(self) {
    if (!active_) {
      return;
    }
    [[self glkView] setNeedsDisplay];
  }
}

- (void)glkView:(GLKView *)view drawInRect:(CGRect)rect {
    BridgeUpdate();
}

- (void)updateTouches:(NSSet*)touches {
    NSInteger index = 0;
    for (UITouch* touch in touches) {
        if (touch.view != [self glkView]) {
          continue;
        }
        CGPoint location = [touch locationInView:touch.view];
        BridgeTouch(location.x, location.y, index, touch.phase);
        index++;
    }
}

- (void)touchesBegan:(NSSet*)touches withEvent:(UIEvent*)event {
  [self updateTouches:touches];
}
- (void)touchesMoved:(NSSet*)touches withEvent:(UIEvent*)event {
  [self updateTouches:touches];
}
- (void)touchesEnded:(NSSet*)touches withEvent:(UIEvent*)event {
  [self updateTouches:touches];
}
- (void)touchesCancelled:(NSSet*)touches withEvent:(UIEvent*)event {
  [self updateTouches:touches];
}

- (void)suspendGame {
    @synchronized(self) {
      active_ = false;
    }
}

- (void)resumeGame {
    @synchronized(self) {
      active_ = true;
    }
}

@end
