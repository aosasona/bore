//go:build ios

#import <UIKit/UIKit.h>

const char* readClipboard() {
	@autoreleasepool {
		NSString *clipboard = [UIPasteboard generalPasteboard].string;
		return clipboard ? strdup([clipboard UTF8String]) : NULL;
	}
}

void writeClipboard(char *text) {
	@autoreleasepool {
		NSString *nsText = [NSString stringWithUTF8String:text];
		[UIPasteboard generalPasteboard].string = nsText;
	}
}
