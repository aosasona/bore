//go:build ios && cgo

package clipboard

/*
#cgo CFLAGS: -x objective-c
#cgo LDFLAGS: -framework Foundation -framework UIKit

#import <stdlib.h>
char *readClipboard();
void writeClipboard(char *text);
*/
import "C"

import (
	"context"
	"unsafe"
)

type iosClipboard struct{}

func NewNativeClipboard() (NativeClipboard, error) {
	return &iosClipboard{}, nil
}

// Available implements NativeClipboard.
func (i *iosClipboard) Available() bool {
	return true
}

// Binaries implements NativeClipboard.
func (i *iosClipboard) Binaries() Binaries {
	return Binaries{}
}

// Clear implements NativeClipboard.
func (i *iosClipboard) Clear(ctx context.Context) error {
	return i.Write(ctx, []byte{})
}

// Read implements NativeClipboard.
func (i *iosClipboard) Read(ctx context.Context) ([]byte, error) {
	cStr := C.readClipboard()
	if cStr == nil {
		return nil, nil // or an error if you prefer
	}
	defer C.free(unsafe.Pointer(cStr))
	return []byte(C.GoString(cStr)), nil
}

// Write implements NativeClipboard.
func (i *iosClipboard) Write(ctx context.Context, data []byte) error {
	rawText := C.CString(string(data))
	defer C.free(unsafe.Pointer(rawText))
	C.writeClipboard(rawText)
	return nil
}

var _ NativeClipboard = (*iosClipboard)(nil)
