package view

import (
	"fmt"
	"io"
	"strconv"
	"text/tabwriter"

	"go.trulyao.dev/bore/v2/database/models"
)

func (v *ViewManager) RenderCollectionsList(
	output io.Writer,
	collections models.Collections,
) error {
	writer := tabwriter.NewWriter(output, 0, 8, 4, ' ', 0)

	_, _ = fmt.Fprintln(writer, "AGGREGATE\tNAME\tCREATED AT\tITEMS")

	for _, collection := range collections {
		aggregate, err := collection.Aggregate()
		if err != nil {
			return err
		}

		createdAt := collection.CreatedAt.Format("Jan 02 2006 15:04")
		fmt.Fprintf(
			writer,
			"%s\t%s\t%s\t%s\n",
			aggregate.String(),
			collection.Name,
			createdAt,
			strconv.Itoa(collection.ItemsCount),
		)
	}

	return writer.Flush()
}
