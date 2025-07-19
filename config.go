package bore

type Config struct {
	/*
		// DataPath is the path to the storage directory.
		DataPath string `toml:"data_path" json:"data_path"`
	*/

	// SystemClipboardPassthrough enables passing the bore clipboard data to the native clipboard on copy
	SystemClipboardPassthrough bool `toml:"system_clipboard_passthrough" json:"system_clipboard_passthrough"`

	// DeleteOnPaste deletes the content after it has been pasted
	DeleteOnPaste bool `toml:"delete_on_paste" json:"delete_on_paste"`
}
