package hashtable

import (
	"fmt"
	"strings"
	"text/tabwriter"
)

func (t *Table) Render() string {
	var builder strings.Builder
	writer := tabwriter.NewWriter(&builder, 2, 2, 2, ' ', 0)
	writeHeader(writer)
	for idx, slot := range t.slots {
		writeRow(writer, idx, slot)
	}
	_ = writer.Flush()
	return builder.String()
}

func writeHeader(writer *tabwriter.Writer) {
	fmt.Fprintln(writer, "Idx\tID\tC\tU\tT\tL\tD\tP0\tPi\tV\tH")
}

func writeRow(writer *tabwriter.Writer, idx int, slot Entry) {
	present := slot.Flags.Occupied || slot.Flags.Deleted
	fmt.Fprintf(
		writer,
		"%d\t%s\t%d\t%d\t%d\t%d\t%d\t%s\t%s\t%s\t%s\n",
		idx,
		slot.Key,
		boolToInt(slot.Flags.Collision),
		boolToInt(slot.Flags.Occupied),
		boolToInt(slot.Flags.Terminal),
		boolToInt(slot.Flags.Link),
		boolToInt(slot.Flags.Deleted),
		formatNumber(slot.Next, present),
		slot.Value,
		formatNumber(slot.V, present),
		formatNumber(slot.Home, present),
	)
}

func boolToInt(flag bool) int {
	if flag {
		return 1
	}
	return 0
}

func formatNumber(value int, ok bool) string {
	if !ok {
		return ""
	}
	return fmt.Sprintf("%d", value)
}
