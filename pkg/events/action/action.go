package action

//go:generate go tool github.com/abice/go-enum --marshal

// ENUM(create_item,bump_item,delete_item,create_collection,delete_collection)
type Action string
