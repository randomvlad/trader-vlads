package util

import (
	"fmt"
	"strings"

	"charm.land/lipgloss/v2"
)

type StringBuilder struct {
	sb    strings.Builder
	style *lipgloss.Style
}

func (b *StringBuilder) WithStyle(style lipgloss.Style) *StringBuilder {
	b.style = &style
	return b
}

func (b *StringBuilder) Write(value string) *StringBuilder {
	if b.style != nil {
		value = b.style.Render(value)
	}
	return b.write(value)
}

func (b *StringBuilder) Writef(format string, a ...any) *StringBuilder {
	value := fmt.Sprintf(format, a...)
	return b.Write(value)
}

func (b *StringBuilder) WriteLn(value string) *StringBuilder {
	return b.Write(value).Ln()
}

func (b *StringBuilder) Ln() *StringBuilder {
	return b.write("\n")
}

func (b *StringBuilder) Tab() *StringBuilder {
	return b.write("\t")
}

func (b *StringBuilder) WriteStylized(value string, style lipgloss.Style) *StringBuilder {
	return b.write(style.Render(value))
}

func (b *StringBuilder) WriteRepeat(value string, count int) *StringBuilder {
	return b.Write(strings.Repeat(value, count))
}

func (b *StringBuilder) String() string {
	return b.sb.String()
}

func (b *StringBuilder) write(value string) *StringBuilder {
	b.sb.WriteString(value)
	return b
}
