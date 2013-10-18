//
//  AppDelegate.m
//  TestWebsocketClient
//
//  Created by Johan Bilien on 10/17/13.
//  Copyright (c) 2013 Johan Bilien. All rights reserved.
//

#import "AppDelegate.h"

@interface AppDelegate()

@property(atomic, retain) NSNetServiceBrowser *serviceBrowser;
@property(atomic, retain) NSNetService *service;

@end

@implementation AppDelegate

@synthesize service;
@synthesize serviceBrowser;

- (void)applicationDidFinishLaunching:(NSNotification *)aNotification
{

    
    self.serviceBrowser = [[NSNetServiceBrowser alloc] init];
    self.serviceBrowser.delegate = self;
    [self.serviceBrowser searchForServicesOfType:@"_woven._tcp" inDomain:@"local."];
}

- (void)webSocketDidOpen:(SRWebSocket *)webSocket
{
    NSLog(@"webScoket didOpen");
}

- (void)webSocket:(SRWebSocket *)webSocket didReceiveMessage:(id)message
{
    NSLog(@"webScoket didReceiveMessage: %@", message);
}

- (void) netServiceBrowser:(NSNetServiceBrowser *)aNetServiceBrowser didFindService:(NSNetService *)aNetService moreComing:(BOOL)moreComing
{
    NSLog(@"netServiceBrowser found service %@", aNetService);
    self.service = aNetService;
    self.service.delegate = self;
    [self.service resolveWithTimeout:2.0];
}

- (void) netServiceBrowser:(NSNetServiceBrowser *)aNetServiceBrowser didRemoveService:(NSNetService *)aNetService moreComing:(BOOL)moreComing
{
    NSLog(@"netService stopped %@", aNetService);
}

- (void)netService:(NSNetService *)sender didNotResolve:(NSDictionary *)errorDict
{
    NSLog(@"netService didNotResolve %@", errorDict);
}

- (void)netServiceDidResolveAddress:(NSNetService *)sender
{
    NSLog(@"netServiceBrowser resolved service %@ on port %ld", sender, (long)sender.port);
    
    NSURL *url = [NSURL URLWithString:[NSString stringWithFormat:@"http://localhost:%ld/ws", sender.port]];
    SRWebSocket *socket = [[SRWebSocket alloc] initWithURL:url];
    socket.delegate = self;
    [socket open];
}

@end
