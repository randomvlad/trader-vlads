package toast

import (
	"fmt"

	"github.com/randomvlad/trader-vlads/internal/appstyle"
)

type Toast struct {
	message string
	Show    bool
}

func (m *Toast) Message(text string, a ...any) {
	m.message = fmt.Sprintf(text, a...)
	m.Show = true
}

func (m *Toast) Clear() {
	m.message = ""
	m.Show = false
}

func (m *Toast) Render() string {
	return appstyle.StyleToast.Render(m.message)
}
