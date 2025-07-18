package bore

type Config struct {
	// SystemClipboardPassthrough enables passing the bore clipboard data to the native clipboard on copy
	SystemClipboardPassthrough bool `toml:"system_clipboard_passthrough" json:"system_clipboard_passthrough"`

	// DeleteOnPaste deletes the content after it has been pasted
	DeleteOnPaste bool `toml:"delete_on_paste" json:"delete_on_paste"`
}
