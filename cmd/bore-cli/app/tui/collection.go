package tui

import (
	"fmt"
	"io"
	"strconv"
	"text/tabwriter"

	"go.trulyao.dev/bore/v2/database/models"
)

func (m *Manager) RenderCollectionsList(
	output io.Writer,
	collections models.Collections,
) error {
	writer := tabwriter.NewWriter(output, 0, 8, 4, ' ', 0)

	_, _ = fmt.Fprintln(writer, "ID\tNAME\tCREATED AT\tITEMS")

	for _, collection := range collections {
		createdAt := collection.CreatedAt.Format("Jan 02 2006 15:04")
		_, _ = fmt.Fprintf(
			writer,
			"%s\t%s\t%s\t%s\n",
			collection.ID,
			collection.Name,
			createdAt,
			strconv.Itoa(collection.ItemsCount),
		)
	}

	return writer.Flush()
}
