package manager

type shortcut struct {
	modifiers modifiers
	key       string
}

type shortcuts map[string]shortcut

var bindings shortcuts = shortcuts{
	// File Menu Shortcuts
	"ditta.Save": shortcut{modifiers: modifiers{ctrl: true}, key: "s"},
	"ditta.Quit": shortcut{modifiers: modifiers{ctrl: true}, key: "q"},
	// Edit Menu Shortcuts
	"ditta.Paste": shortcut{modifiers: modifiers{ctrl: true}, key: "v"},
	"ditta.Copy":  shortcut{modifiers: modifiers{ctrl: true}, key: "c"},
	"ditta.Cut":   shortcut{modifiers: modifiers{ctrl: true}, key: "x"},
	// Tools menu shortcuts
	"ditta.FullScreen": shortcut{modifiers: modifiers{ctrl: true, alt: true}, key: "f"},
}
