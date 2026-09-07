package stringutil

import (
	"fmt"
	"strings"

	"charm.land/lipgloss/v2"
)

type Builder struct {
	sb    strings.Builder
	style *lipgloss.Style
}

func NewBuilder() *Builder {
	return &Builder{}
}

func (b *Builder) WithStyle(style lipgloss.Style) *Builder {
	b.style = &style
	return b
}

func (b *Builder) Write(value string) *Builder {
	if b.style != nil {
		value = b.style.Render(value)
	}
	return b.write(value)
}

func (b *Builder) Writef(format string, a ...any) *Builder {
	value := fmt.Sprintf(format, a...)
	return b.Write(value)
}

func (b *Builder) WriteLn(value string) *Builder {
	return b.Write(value).Ln()
}

func (b *Builder) Ln() *Builder {
	return b.write("\n")
}

func (b *Builder) Tab() *Builder {
	return b.write("\t")
}

func (b *Builder) WriteStyle(value string, style lipgloss.Style) *Builder {
	return b.write(style.Render(value))
}

func (b *Builder) WriteStyleRanges(value string, ranges ...lipgloss.Range) *Builder {
	return b.write(lipgloss.StyleRanges(value, ranges...))
}

func (b *Builder) WriteRepeat(value string, count int) *Builder {
	return b.Write(strings.Repeat(value, count))
}

func (b *Builder) IsNotBlank() bool {
	return b.sb.Len() > 0 && strings.TrimSpace(b.sb.String()) != ""
}

func (b *Builder) Len() int {
	return b.sb.Len()
}

func (b *Builder) String() string {
	return b.sb.String()
}

func (b *Builder) StringStyle(style lipgloss.Style) string {
	return style.Render(b.sb.String())
}

func (b *Builder) write(value string) *Builder {
	b.sb.WriteString(value)
	return b
}
