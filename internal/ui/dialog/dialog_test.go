package dialog

import (
	"testing"
)

func TestNewOverlay(t *testing.T) {
	overlay := NewOverlay()
	if overlay == nil {
		t.Error("expected non-nil overlay")
	}
	if overlay.HasDialogs() {
		t.Error("expected no dialogs initially")
	}
}

func TestOverlay_OpenAndCloseDialog(t *testing.T) {
	overlay := NewOverlay()

	dialog := NewCommandsDialog(80, 24)
	overlay.OpenDialog(dialog)

	if !overlay.HasDialogs() {
		t.Error("expected dialog to be open")
	}

	overlay.CloseDialog(HelpID)

	if overlay.HasDialogs() {
		t.Error("expected dialog to be closed")
	}
}

func TestOverlay_CloseDialog_NotFound(t *testing.T) {
	overlay := NewOverlay()

	dialog := NewCommandsDialog(80, 24)
	overlay.OpenDialog(dialog)

	overlay.CloseDialog("non-existent")

	if !overlay.HasDialogs() {
		t.Error("dialog should still be open")
	}
}

func TestOverlay_View_Empty(t *testing.T) {
	overlay := NewOverlay()

	view := overlay.View(80, 24)
	if view != "" {
		t.Errorf("expected empty view, got %s", view)
	}
}

func TestOverlay_View_WithDialog(t *testing.T) {
	overlay := NewOverlay()

	dialog := NewCommandsDialog(80, 24)
	overlay.OpenDialog(dialog)

	view := overlay.View(80, 24)
	if view == "" {
		t.Error("expected non-empty view")
	}
}

func TestCommandsDialog_ID(t *testing.T) {
	dialog := NewCommandsDialog(80, 24)
	if dialog.ID() != HelpID {
		t.Errorf("expected %s, got %s", HelpID, dialog.ID())
	}
}

func TestCommandsDialog_HandleMsg_Escape(t *testing.T) {
	dialog := NewCommandsDialog(80, 24)

	action := dialog.HandleMsg(nil)
	if action != nil {
		t.Error("expected nil action for nil message")
	}
}

func TestCommandsDialog_View(t *testing.T) {
	dialog := NewCommandsDialog(80, 24)

	view := dialog.View(80, 24)
	if view == "" {
		t.Error("expected non-empty view")
	}
}

func TestCommandsDialog_View_InvalidWidth(t *testing.T) {
	dialog := NewCommandsDialog(80, 24)

	view := dialog.View(5, 24)
	if view == "" {
		t.Error("expected non-empty view even with small width")
	}
}
