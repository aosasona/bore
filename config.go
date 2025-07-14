package bore

// TODO: move path to the CLI itself (can be overriden by using the -c flag)
// TODO: automatically detect the data directory for each platform and remove this config option
type Config struct {
	// Path is the path to the configuration file
	Path string `toml:"-" json:"config_path"`

	// DataDir is the path to the directory where the application stores its data
	DataDir string `toml:"data_dir" json:"data_dir"`

	// SystemClipboardPassthrough enables passing the bore clipboard data to the native clipboard on copy
	SystemClipboardPassthrough bool `toml:"system_clipboard_passthrough" json:"system_clipboard_passthrough"`

	// DeleteOnPaste deletes the content after it has been pasted
	DeleteOnPaste bool `toml:"delete_on_paste" json:"delete_on_paste"`
}
