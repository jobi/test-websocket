//
//  AppDelegate.h
//  TestWebsocketClient
//
//  Created by Johan Bilien on 10/17/13.
//  Copyright (c) 2013 Johan Bilien. All rights reserved.
//

#import <Cocoa/Cocoa.h>
#import <SocketRocket/SRWebSocket.h>

@interface AppDelegate : NSObject <NSApplicationDelegate, SRWebSocketDelegate, NSNetServiceBrowserDelegate, NSNetServiceDelegate>

@property (assign) IBOutlet NSWindow *window;

@end
