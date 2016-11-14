// package provides text retrieval backend for the ditta
package manager

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/gotk3/gotk3/gdk"
	"github.com/gotk3/gotk3/gtk"
)

type Manager struct {
	id        string
	text      string
	fileName  Option
	canrw     bool
	x, y      int
	modifiers modifiers
	selection SelectionBounds
}

type modifiers struct {
	alt   bool
	ctrl  bool
	shift bool
	super bool
}

type SelectionBounds struct {
	X, Y        int
	InSelection bool
}

// NewManager returns a new Manager struct with members initialised
func NewManager(id string) *Manager {
	manager := &Manager{
		id:       id,
		text:     "",
		fileName: NewOption(),
		x:        0,
		y:        0,
		modifiers: modifiers{
			alt:   false,
			ctrl:  false,
			super: false,
			shift: false,
		},
		selection: SelectionBounds{0, 0, false},
	}
	return manager
}

func (m *Manager) ReadKey(win *gtk.ApplicationWindow, ev *gdk.Event) {
	keyEvent := &gdk.EventKey{ev}
	key := keyEvent.KeyVal()
	switch {
	case key >= 32 && key <= 127: // typable characters
		if m.modifiers.alt || m.modifiers.ctrl || m.modifiers.super {
			m.checkShortcut(win, string(key))
			break
		}
		m.addChar(string(key))
	case key == gdk.KEY_BackSpace: // backspace
		if m.x == 0 {
			lines := strings.Split(m.text, "\n")
			newx := 0
			if m.y > 0 {
				newx = len(lines[m.y-1])
			}
			m.bkspChar()
			m.x = newx
		} else {
			m.bkspChar()
		}
	case key == gdk.KEY_Delete: // delete
		m.delChar()
	case key == gdk.KEY_Return: // enter
		m.addChar("\n")
		m.y++
		m.x = 0
	case key == gdk.KEY_Left: // arrow left
		if m.modifiers.shift {
			if !m.selection.InSelection {
				m.selection.InSelection = true
				m.selection.X, m.selection.Y = m.x, m.y
			}
		} else {
			m.selection.InSelection = false
		}
		if m.modifiers.ctrl {
			m.cursorLeftWord()
			break
		}
		m.cursorLeft()
	case key == gdk.KEY_Right: // arrow right
		if m.modifiers.shift {
			if !m.selection.InSelection {
				m.selection.InSelection = true
				m.selection.X, m.selection.Y = m.x, m.y
			}
		} else {
			m.selection.InSelection = false
		}
		if m.modifiers.ctrl {
			m.cursorRightWord()
			break
		}
		m.cursorRight()
	case key == gdk.KEY_Up: // arrow up
		if m.modifiers.shift {
			if !m.selection.InSelection {
				m.selection.InSelection = true
				m.selection.X, m.selection.Y = m.x, m.y
			}
		} else {
			m.selection.InSelection = false
		}
		m.cursorUp()
	case key == gdk.KEY_Down: // arrow down
		if m.modifiers.shift {
			if !m.selection.InSelection {
				m.selection.InSelection = true
				m.selection.X, m.selection.Y = m.x, m.y
			}
		} else {
			m.selection.InSelection = false
		}
		m.cursorDown()
	case key == gdk.KEY_Alt_L, key == gdk.KEY_Alt_R:
		m.modifiers.alt = true
	case key == gdk.KEY_Control_L, key == gdk.KEY_Control_R:
		m.modifiers.ctrl = true
	case key == gdk.KEY_Super_L, key == gdk.KEY_Super_R:
		m.modifiers.super = true
	case key == gdk.KEY_Shift_L, key == gdk.KEY_Shift_R:
		m.modifiers.shift = true
	}
	win.QueueDraw()
}

func (m *Manager) CheckModifier(win *gtk.ApplicationWindow, ev *gdk.Event) {
	keyEvent := &gdk.EventKey{ev}
	key := keyEvent.KeyVal()
	switch {
	case key == gdk.KEY_Alt_L, key == gdk.KEY_Alt_R:
		m.modifiers.alt = false
	case key == gdk.KEY_Control_L, key == gdk.KEY_Control_R:
		m.modifiers.ctrl = false
	case key == gdk.KEY_Super_L, key == gdk.KEY_Super_R:
		m.modifiers.super = false
	case key == gdk.KEY_Shift_L, key == gdk.KEY_Shift_R:
		m.modifiers.shift = false
	}
}

func (m *Manager) checkShortcut(win *gtk.ApplicationWindow, key string) {
	scsession := shortcut{m.modifiers, key}
	for name, val := range bindings {
		if val == scsession {
			switch name {
			case "ditta.Quit":
				app, err := win.GetApplication()
				if err != nil {
					fmt.Println(err)
					return
				}
				app.Quit()
			case "ditta.Save":
				m.Save()
			case "ditta.Paste":
				cb, err := gtk.ClipboardGet(gdk.SELECTION_CLIPBOARD)
				if err != nil {
					fmt.Println(err)
					return
				}
				stuff, err := cb.WaitForText()
				if err != nil {
					fmt.Println(err)
					return
				}
				m.addChar(stuff)
			case "ditta.Copy":
				cb, err := gtk.ClipboardGet(gdk.SELECTION_CLIPBOARD)
				if err != nil {
					fmt.Println(err)
					return
				}
				copied := m.GetSelectionText()
				if copied != "" {
					cb.SetText(copied)
				}
				m.selection.InSelection = false
			case "ditta.Cut":
				cb, err := gtk.ClipboardGet(gdk.SELECTION_CLIPBOARD)
				if err != nil {
					fmt.Println(err)
					return
				}
				cut := m.CutStuff()
				if cut != "" {
					cb.SetText(cut)
				}
				m.selection.InSelection = false
			}
		}
	}
}

func (m *Manager) GetText(w, h int) string {
	return m.getNLines(int(h/20.0) + 1)
}

func (m *Manager) getNLines(n int) string {
	lines := strings.Split(m.text, "\n")
	text := ""
	for i, line := range lines {
		text += line + "\n"
		if i+1 == n {
			break
		}
	}
	text = text[:len(text)-1]
	return text
}

func (m *Manager) Save() {
	filename, err := m.fileName.Unexpect()
	if err != nil {
		fmt.Println(err)
		return
	}
	err = ioutil.WriteFile(filename, []byte(m.text), 0664)
	if err != nil {
		fmt.Println(err)
		return
	}
}

func (m *Manager) SetFileName(filename string) {
	m.fileName.Set(filename)
	_, err := os.Stat(filename)
	if err == nil || os.IsNotExist(err) {
		b, err := ioutil.ReadFile(filename)
		if err != nil {
			m.canrw = false
			return
		}
		m.text = string(b)
	}
}

func (m Manager) GetCursorXY() (int, int) {
	return m.x, m.y
}

func (m Manager) IsSelection() bool {
	return m.selection.InSelection
}

func (m Manager) GetSelectionStart() (int, int) {
	return m.selection.X, m.selection.Y
}

func (m *Manager) CutStuff() string {
	if !m.selection.InSelection {
		return ""
	}
	offset := m.getOffset()
	selectionOffset := m.getOffsetFor(m.selection.X, m.selection.Y)
	m.x, m.y = m.selection.X, m.selection.Y
	if selectionOffset > offset {
		selection := m.text[offset:selectionOffset]
		m.text = m.text[:offset] + m.text[selectionOffset:]
		return selection
	}
	selection := m.text[selectionOffset:offset]
	m.text = m.text[:selectionOffset] + m.text[offset:]
	return selection
}

func (m Manager) GetSelectionText() string {
	if !m.selection.InSelection {
		return ""
	}
	offset := m.getOffset()
	selectionOffset := m.getOffsetFor(m.selection.X, m.selection.Y)
	if selectionOffset > offset {
		return m.text[offset:selectionOffset]
	}
	return m.text[selectionOffset:offset]
}

func (m *Manager) cursorLeft() {
	if m.x != 0 {
		m.x--
		return
	}
	if m.y == 0 {
		m.x = 0
		return
	}
	m.y--
	lines := strings.Split(m.text, "\n")
	if len(lines) >= m.y+1 {
		m.x = len(lines[m.y])
	}
}

func (m *Manager) cursorRightWord() {
	offset := m.getOffset()
	first := true
	for {
		if string(m.text[offset]) == " " {
			if !first {
				break
			}
			first = false
		}
		if offset+1 > len(m.text) {
			break
		}
		m.cursorRight()
		offset++
	}
}

func (m *Manager) cursorLeftWord() {
	offset := m.getOffset()
	first := true
	for {
		if string(m.text[offset]) == " " {
			if !first {
				break
			}
			first = false
		}
		if offset-1 < 0 {
			break
		}
		m.cursorLeft()
		offset--
	}
}

func (m *Manager) cursorRight() {
	lines := strings.Split(m.text, "\n")
	if len(lines) < m.y+1 {
		return
	}
	if len(lines[m.y]) == m.x {
		if len(lines) != m.y+1 {
			m.y++
			m.x = 0
		}
	} else {
		m.x++
	}
}

func (m *Manager) cursorUp() {
	lines := strings.Split(m.text, "\n")
	if m.y > 0 {
		m.y--
		if len(lines[m.y]) < m.x {
			m.x = len(lines[m.y])
		}
	}
}

func (m *Manager) cursorDown() {
	lines := strings.Split(m.text, "\n")
	if m.y+1 < len(lines) {
		m.y++
		if len(lines[m.y]) < m.x {
			m.x = len(lines[m.y])
		}
	}
}

func (m *Manager) addChar(c string) {
	offset := m.getOffset()
	m.text = m.text[:offset] + c + m.text[offset:]
	m.x += len(c)
}

func (m *Manager) bkspChar() {
	offset := m.getOffset()
	if offset > 0 {
		m.text = m.text[:offset-1] + m.text[offset:]
		m.cursorLeft()
	}
}

func (m *Manager) delChar() {
	offset := m.getOffset()
	if offset+1 < len(m.text) {
		m.text = m.text[:offset] + m.text[offset+1:]
	} else {
		m.text = ""
	}
}

func (m Manager) getOffset() int {
	return m.getOffsetFor(m.x, m.y)
}

func (m Manager) getOffsetFor(x, y int) int {
	lines := strings.Split(m.text, "\n")
	offset := 0
	for i, val := range lines {
		if y > i {
			offset += len(val) + 1
		}
		if y == i {
			offset += x
			break
		}
	}
	return offset
}
