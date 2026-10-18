#!/usr/bin/env python3
"""Actual X11 command/surface acceptance; fixture commands only arrange conditions."""
import argparse
import subprocess
import sys
import time
from pathlib import Path
import verify_wci1 as base

TITLE = "SDUI WCI2 Commands"
CHILD_TITLE = "SDUI WCI2 Settings"
base.TITLE = TITLE
TOOLBAR = "page/view/header/toolbar/"
DIALOG = "page/settings"
NAME = DIALOG + "/form/body/name"
NOTE = DIALOG + "/form/body/note"
ACCEPT = DIALOG + "/form/footer/acceptButton"
CANCEL = DIALOG + "/form/footer/cancelButton"
CLOSE = DIALOG + "/form/footer/closeButton"
EDIT_MENU = DIALOG + "/editMenu"
TREE = "page/view/body/items"
PREVIEW = "page/view/footer/preview"
ITEM_PREVIEW = "page/view/footer/itemPreview"
BAR = "page/view/header/actions"
CONTEXT = "page/itemMenu"

class Trial(base.Trial):
    def __init__(self, args):
        args.binary_args = ["--nonmodal"] if args.nonmodal else []
        variant = args.variant
        args.variant = "wci2-m2-" + variant
        try:
            super().__init__(args)
        finally:
            args.variant = variant

    def native(self, action, *values, title=TITLE, keep_focus=False):
        cursor = len(self.events)
        self.record("x11-input", {"title": title, "action": action,
                                 "values": values, "keep_focus": keep_focus})
        command = [sys.executable, str(Path(__file__).with_name("x11_input.py")),
                   "--display", self.args.display, "--title", title]
        if keep_focus:
            command.append("--keep-focus")
        subprocess.run(command + [action, *map(str, values)], check=True, timeout=10)
        return cursor

    def input(self, action, *values):
        return self.native(action, *values)

    def screenshot(self, name):
        # Dialog animation/desktop painting can trail accepted input by more
        # than the base capture delay. This is only visual capture pacing;
        # all interaction assertions use correlated runtime/action events.
        time.sleep(1)
        super().screenshot(name)

    def calls(self, state, name):
        return state["actionCalls"].get(name, 0)

    def field(self, state, path):
        return next(w for w in state["widgets"] if w["InstancePath"] == path)

    def opened(self, state):
        return state["snapshot"]["Surfaces"][DIALOG]["Open"]

    def surface_title(self):
        return CHILD_TITLE if self.args.nonmodal else TITLE

    def control(self, path, action="click"):
        state = self.state()
        info = state["controls"][path]
        rect = info["clip"]
        assert info["visible"] and rect["W"] > 0 and rect["H"] > 0, path
        title = info["title"]
        return self.native(action, round(rect["X"] + rect["W"]/2),
                           round(rect["Y"] + rect["H"]/2), title=title)

    def type_field(self, path, text):
        self.control(path)
        title = self.surface_title() if path.startswith(DIALOG + "/") else TITLE
        cursor = self.native("chord", "Control_L", "a", title=title)
        self.native("key", *list(text), title=title)
        self.changed(lambda d: self.field(d, path)["Draft"] == text, cursor)

    def open_dialog(self):
        cursor = self.control(TOOLBAR + "openButton")
        return self.changed(self.opened, cursor)

    def receipt(self, cursor, kind, path=DIALOG):
        return self.wait(lambda e: e["event"] == "dialog-result"
                         and e["data"]["Kind"] == kind
                         and e["data"]["Surface"]["Handle"]["Path"] == path, cursor)["data"]

    def menu_item(self, root, path, action="click"):
        deadline = time.monotonic() + 3
        while True:
            state = self.state()
            menu = state.get("menus", {}).get(root)
            item = next((i for i in menu["items"] if i["path"] == path), None) if menu else None
            if item and item["clip"]["W"] > 0 and item["clip"]["H"] > 0:
                r = item["clip"]
                return self.native(action, round(r["X"] + r["W"]/2),
                                   round(r["Y"] + r["H"]/2), title=menu["title"])
            if time.monotonic() >= deadline:
                raise AssertionError("Native menu item not visible: " + path)
            time.sleep(.05)  # driver layout/hover, never provider or callback completion

    def tree_row(self, item, action="click"):
        row = next(r for r in self.state()["rows"][TREE] if r["item"] == item)
        rect = row["clip"]
        return self.native(action, round(rect["X"] + min(100, rect["W"]/2)),
                           round(rect["Y"] + rect["H"]/2))

    def commands(self):
        initial = self.state()
        self.check("preparation executes no SDL action", not any(initial["actionCalls"].values()))
        cursor = self.control(TOOLBAR + "runButton")
        self.changed(lambda d: self.calls(d, "Run") == 1, cursor)
        self.check("toolbar invokes shared Run once", self.calls(self.state(), "Run") == 1)
        cursor = self.native("chord", "Control_L", "r")
        self.changed(lambda d: self.calls(d, "Run") == 2, cursor)
        self.check("key invokes same Run once", self.calls(self.state(), "Run") == 2)
        self.control(BAR)
        cursor = self.menu_item(BAR, BAR + "/runItem")
        self.changed(lambda d: self.calls(d, "Run") == 3, cursor)
        self.check("menu invokes same Run once after native dismissal", self.calls(self.state(), "Run") == 3)
        cursor = self.control(TOOLBAR + "flagButton")
        state = self.changed(lambda d: self.calls(d, "Toggle") == 1, cursor)
        self.check("pointer toggle publishes typed checked state", state["snapshot"]["Commands"]["page/flag"]["Checked"])
        cursor = self.native("chord", "Control_L", "k")
        state = self.changed(lambda d: self.calls(d, "Toggle") == 2, cursor)
        self.check("key toggles canonical state once", not state["snapshot"]["Commands"]["page/flag"]["Checked"])
        cursor = self.control(TOOLBAR + "betaButton")
        state = self.changed(lambda d: d["snapshot"]["Commands"]["page/beta"]["Checked"], cursor)
        self.check("exclusive selection clears its peer atomically", not state["snapshot"]["Commands"]["page/alpha"]["Checked"] and self.calls(state, "Toggle") == 2)
        self.control(TOOLBAR + "betaButton")
        self.check("already checked exclusive command emits no domain call", self.calls(self.state(), "Toggle") == 2)
        self.control(BAR)
        self.menu_item(BAR, BAR + "/more", "move")
        cursor = self.menu_item(BAR, BAR + "/more/alphaItem")
        state = self.changed(lambda d: d["snapshot"]["Commands"]["page/alpha"]["Checked"], cursor)
        self.check("nested pointer menu updates shared exclusive state", not state["snapshot"]["Commands"]["page/beta"]["Checked"])
        self.command("disable Run")
        self.control(TOOLBAR + "runButton")
        self.native("chord", "Control_L", "r")
        cursor = self.native("chord", "Control_L", "k")
        self.changed(lambda d: self.calls(d, "Toggle") == 3, cursor)
        self.check("disabled toolbar and key cannot invoke Run", self.calls(self.state(), "Run") == 3)
        self.control(BAR)
        cursor = self.native("key", "Home", "Return")
        self.changed(lambda d: not d["snapshot"]["Menus"][BAR]["Open"], cursor)
        self.check("disabled native menu keyboard activation cannot invoke Run", self.calls(self.state(), "Run") == 3)
        self.command("enable Run")
        cursor = self.native("chord", "Control_L", "r")
        self.changed(lambda d: self.calls(d, "Run") == 4, cursor)
        self.check("re-enabled shared command remains usable", self.calls(self.state(), "Run") == 4)
        info = self.state()["controls"][TOOLBAR + "runButton"]["clip"]
        self.native("move", round(info["X"] + info["W"]/2), round(info["Y"] + info["H"]/2))
        self.screenshot("shared-commands")

    def context(self):
        cursor = self.tree_row("alpha")
        self.changed(lambda d: d["snapshot"]["Collections"][TREE]["Selected"] == "alpha", cursor)
        self.tree_row("beta", "right-click")
        self.check("secondary click does not replace selected item", self.state()["snapshot"]["Collections"][TREE]["Selected"] == "alpha")
        cursor = self.menu_item(CONTEXT, CONTEXT + "/inspectItem")
        state = self.changed(lambda d: self.calls(d, "Inspect") == 1, cursor)
        self.check("item command uses captured row without retargeting selection", state["snapshot"]["Collections"][TREE]["Selected"] == "alpha" and "beta" in self.field(state, ITEM_PREVIEW)["Value"])
        self.tree_row("beta", "right-click")
        menu = self.state()["menus"][CONTEXT]
        item = next(i for i in menu["items"] if i["path"] == CONTEXT + "/inspectItem")
        rect = item["clip"]
        self.command("replace-items")
        self.native("click", round(rect["X"] + rect["W"]/2), round(rect["Y"] + rect["H"]/2))
        cursor = self.native("chord", "Control_L", "r")
        self.changed(lambda d: self.calls(d, "Run") == 1, cursor)
        self.check("obsolete context cannot invoke after collection replacement", self.calls(self.state(), "Inspect") == 1)
        self.screenshot("context-recovery")

    def dialog(self):
        self.open_dialog()
        self.type_field(NAME, "draft")
        self.screenshot("dialog-draft")
        self.check("typing does not persist dialog", self.calls(self.state(), "Save") == 0)
        cursor = self.control(NAME, "right-click")
        self.changed(lambda d: d["snapshot"]["Menus"][EDIT_MENU]["Open"], cursor)
        self.screenshot("dialog-context")
        cursor = self.native("key", "Escape", title=self.surface_title())
        state = self.changed(lambda d: not d["snapshot"]["Menus"][EDIT_MENU]["Open"], cursor)
        self.check("Escape dismisses menu without reverting draft or closing surface", self.opened(state) and self.field(state, NAME)["Draft"] == "draft" and self.field(state, NAME)["Dirty"])
        cursor = self.native("key", "Escape", title=self.surface_title())
        self.changed(lambda d: not self.field(d, NAME)["Dirty"], cursor)
        self.check("next Escape reverts draft and keeps surface", self.opened(self.state()))
        cursor = self.native("key", "Escape", title=self.surface_title())
        result = self.receipt(cursor, "cancel")
        self.check("clean Escape cancels once without Save", self.calls(self.state(), "Save") == 0 and result["Reason"] == "user")
        self.open_dialog()
        self.type_field(NAME, "saved")
        self.command("accept false")
        cursor = self.control(ACCEPT)
        state = self.changed(lambda d: self.calls(d, "Save") == 1, cursor)
        self.check("rejected Accept retains dirty form", self.opened(state) and self.field(state, NAME)["Dirty"])
        self.command("accept true")
        cursor = self.control(ACCEPT)
        result = self.receipt(cursor, "accept")
        state = self.state()
        self.check("successful Accept commits captured fields and closes once", not self.opened(state) and self.calls(state, "Save") == 2 and self.field(state, NAME)["Value"] == "saved" and result["Domain"] == "succeeded")
        self.open_dialog()
        cursor = self.control(CLOSE)
        result = self.receipt(cursor, "close")
        self.check("Close is distinct from Cancel and emits no Save", self.calls(self.state(), "Save") == 2 and result["Reason"] == "user")
        self.screenshot("dialog-completed")

    def keyboard(self):
        self.control(BAR)
        cursor = self.native("key", "Home", "Return")
        self.changed(lambda d: self.calls(d, "Run") == 1, cursor)
        self.check("native menu Home Enter invokes once", self.calls(self.state(), "Run") == 1)
        self.control(BAR)
        cursor = self.native("key", "End", "Right", "End", "space")
        state = self.changed(lambda d: d["snapshot"]["Commands"]["page/beta"]["Checked"], cursor)
        self.check("nested keyboard End Space selects exclusive peer", not state["snapshot"]["Commands"]["page/alpha"]["Checked"])
        self.control(BAR)
        cursor = self.native("key", "Escape")
        self.changed(lambda d: not d["snapshot"]["Menus"][BAR]["Open"], cursor)
        self.check("menu Escape dismisses without action", self.calls(self.state(), "Run") == 1)
        cursor = self.tree_row("alpha")
        self.changed(lambda d: d["snapshot"]["Collections"][TREE]["Selected"] == "alpha", cursor)
        cursor = self.native("chord", "Shift_L", "F10")
        self.changed(lambda d: d["snapshot"]["Menus"][CONTEXT]["Open"], cursor)
        cursor = self.native("key", "Home", "Return")
        state = self.changed(lambda d: self.calls(d, "Inspect") == 1, cursor)
        self.check("Shift F10 captures focused item for keyboard action", "alpha" in self.field(state, ITEM_PREVIEW)["Value"])
        self.type_field("page/view/body/mainDraft", "keep")
        cursor = self.native("chord", "Control_L", "r")
        state = self.changed(lambda d: self.calls(d, "Run") == 2, cursor)
        self.check("declared shortcut works from Entry without committing its draft", self.field(state, "page/view/body/mainDraft")["Draft"] == "keep" and self.field(state, "page/view/body/mainDraft")["Dirty"])
        self.screenshot("keyboard-commands")

    def focus(self):
        state = self.open_dialog()
        if self.args.nonmodal:
            self.native("move-window", 880, 180, title=CHILD_TITLE, keep_focus=True)
        self.check("opening focuses first eligible form control", state["focused"] == NAME)
        self.screenshot("surface-open")
        for path in (NOTE, ACCEPT, CANCEL, CLOSE, NAME):
            cursor = self.native("key", "Tab", title=self.surface_title())
            self.changed(lambda d: d["focused"] == path, cursor)
            self.check("surface Tab reaches " + path.rsplit("/", 1)[-1], self.state()["focused"] == path)
        self.type_field(NAME, "retained")
        cursor = self.control(TOOLBAR + "runButton")
        if self.args.nonmodal:
            state = self.changed(lambda d: self.calls(d, "Run") == 1, cursor)
            self.check("nonmodal allows parent action and preserves child draft", self.opened(state) and self.field(state, NAME)["Draft"] == "retained")
        else:
            # Positive child editing event follows the parent click on the same
            # native event queue, proving that the blocked click was processed.
            self.type_field(NAME, "retainedx")
            self.check("modal prevents parent action and retains child editing", self.calls(self.state(), "Run") == 0)
        cursor = self.control(CANCEL)
        self.receipt(cursor, "cancel")
        state = self.state()
        self.check("Cancel restores surviving opener focus", state["focused"] == TOOLBAR + "openButton")
        cursor = self.native("key", "space", keep_focus=True)
        self.changed(self.opened, cursor)
        self.check("restored opener receives actual keyboard input", self.opened(self.state()))
        cursor = self.control(CANCEL)
        self.receipt(cursor, "cancel")
        if self.args.nonmodal:
            self.open_dialog()
            self.native("move-window", 880, 180, title=CHILD_TITLE, keep_focus=True)
            cursor = self.control(TOOLBAR + "runButton")
            state = self.changed(lambda d: self.calls(d, "Run") == 2, cursor)
            cursor = self.native("close-window", title=CHILD_TITLE, keep_focus=True)
            self.receipt(cursor, "close")
            self.check("closing unfocused nonmodal preserves active parent control", self.state()["focused"] == TOOLBAR + "runButton")
            cursor = self.native("key", "space", keep_focus=True)
            self.changed(lambda d: self.calls(d, "Run") == 3, cursor)
            self.check("closing unfocused child does not steal OS keyboard focus", self.calls(self.state(), "Run") == 3)
        self.screenshot("focus-restored")

    def command_failures(self):
        receiver = "page/view/footer/flagPreview"
        for index, failure in enumerate(("error", "malformed", "draft-conflict"), 1):
            self.command("action Toggle " + failure)
            self.control(BAR)
            cursor = self.menu_item(BAR, BAR + "/options/flagItem")
            error = self.wait(lambda e: e["event"] == "error" and "domain=" in e["data"], cursor)["data"]
            state = self.state()
            expected = "succeeded" if failure == "draft-conflict" else "unknown"
            self.check("menu " + failure + " invokes once and reports domain outcome", self.calls(state, "Toggle") == index and "domain=" + expected in error)
            self.check("menu " + failure + " dismisses without speculative checked state", not state["snapshot"]["Menus"][BAR]["Open"] and not state["snapshot"]["Commands"]["page/flag"]["Checked"])
            if failure == "draft-conflict":
                self.check("accepted reentrant receiver draft survives failed publication", self.field(state, receiver)["Dirty"] and self.field(state, receiver)["Draft"] == "New draft during SDL action")
        # A new explicit gesture is distinct from automatically replaying the
        # failed invocation. The existing accepted reentrant draft stays dirty.
        cursor = self.control(TOOLBAR + "runButton")
        self.changed(lambda d: self.calls(d, "Run") == 1, cursor)
        self.check("later unrelated action cannot replay failed toggle", self.calls(self.state(), "Toggle") == 3)
        self.screenshot("command-failures")

    def nested(self):
        assert self.args.nonmodal, "Nested parent regression requires --nonmodal"
        self.open_dialog()
        self.native("move-window", 880, 180, title=CHILD_TITLE, keep_focus=True)
        self.type_field(NAME, "parentdraft")

        def open_child():
            cursor = self.control(NAME, "right-click")
            self.changed(lambda d: d["snapshot"]["Menus"][EDIT_MENU]["Open"], cursor)
            cursor = self.menu_item(EDIT_MENU, EDIT_MENU + "/childItem")
            return self.changed(lambda d: d["snapshot"]["Surfaces"]["page/x"]["Open"], cursor)

        state = open_child()
        parent, child = state["surfaces"][DIALOG], state["surfaces"]["page/x"]
        self.check("shorter source sibling uses actual nonmodal parent canvas", child["canvas"] == DIALOG and child["title"] == CHILD_TITLE and state["snapshot"]["Surfaces"]["page/x"]["ParentSurface"] == parent["target"])
        self.check("child body scales against actual parent body", all(abs(child["size"][axis] - parent["size"][axis] * .4) < .0001 for axis in ("W", "H")))
        self.check("child opening focuses its own first field", state["focused"] == "page/x/childInput")
        self.control("page/x/childInput")
        self.native("chord", "Control_L", "a", title=CHILD_TITLE)
        cursor = self.native("key", *list("childdraft"), title=CHILD_TITLE)
        self.changed(lambda d: self.field(d, "page/x/childInput")["Draft"] == "childdraft", cursor)
        self.screenshot("nested-modal-parent")
        cursor = self.native("key", "Escape", title=CHILD_TITLE)
        state = self.changed(lambda d: not self.field(d, "page/x/childInput")["Dirty"], cursor)
        self.check("child Escape reverts only its own draft", state["snapshot"]["Surfaces"]["page/x"]["Open"] and self.field(state, NAME)["Draft"] == "parentdraft")
        cursor = self.native("key", "Escape", title=CHILD_TITLE)
        self.receipt(cursor, "cancel", "page/x")
        state = self.state()
        self.check("clean child Escape preserves parent and restores its focus", self.opened(state) and state["focused"] == NAME and self.field(state, NAME)["Dirty"])
        open_child()
        cursor = self.control("page/x/childClose")
        self.receipt(cursor, "close", "page/x")
        self.check("explicit child Close preserves parent draft without Save", self.opened(self.state()) and self.calls(self.state(), "Save") == 0 and self.field(self.state(), NAME)["Draft"] == "parentdraft")
        open_child()
        cursor = self.native("close-window", title=CHILD_TITLE, keep_focus=True)
        child_result = self.receipt(cursor, "close", "page/x")
        parent_result = self.receipt(cursor, "close", DIALOG)
        state = self.state()
        self.check("native parent close revokes dynamic child before teardown", child_result["Reason"] == "parent-closed" and parent_result["Reason"] == "parent-closed" and child_result["Sequence"] == 0 and parent_result["Sequence"] == 0 and not self.opened(state) and not state["snapshot"]["Surfaces"]["page/x"]["Open"] and self.calls(state, "Save") == 0)
        self.screenshot("nested-parent-closed")

    def accept_failures(self):
        prior_generation = 0
        for failure in ("error", "malformed", "draft-conflict", "resource-conflict"):
            state = self.open_dialog()
            surface = state["snapshot"]["Surfaces"][DIALOG]
            generation = surface["Target"]["OpenGeneration"]
            self.check("fresh " + failure + " opening has a new token", generation > prior_generation)
            prior_generation = generation
            self.type_field(NAME, "pending")
            before = self.state()
            calls = self.calls(before, "Save")
            self.command("accept " + failure)
            cursor = self.control(ACCEPT)
            state = self.changed(lambda d: self.calls(d, "Save") == calls + 1
                                 and d["snapshot"]["Surfaces"][DIALOG]["AcceptBlocked"], cursor)
            expected = "succeeded" if failure.endswith("conflict") else "unknown"
            surface = state["snapshot"]["Surfaces"][DIALOG]
            self.check(failure + " preserves outcome and blocks replay", surface["Open"]
                       and surface["Domain"] == expected and surface["AcceptSequence"] > 0)
            attempt = surface["AcceptSequence"]
            self.screenshot("accept-" + failure)
            self.control(ACCEPT)
            self.type_field(NAME, "recovery")
            self.check("blocked " + failure + " remains editable without another Save", self.calls(self.state(), "Save") == calls + 1)
            committed = self.state()["domain"]["Commits"]
            cursor = self.control(CANCEL)
            result = self.receipt(cursor, "cancel")
            state = self.state()
            self.check("Cancel after " + failure + " retains attempt and does not roll back domain", result["Domain"] == expected
                       and result["AcceptSequence"] == attempt and not self.opened(state)
                       and state["domain"]["Commits"] == committed)
        self.screenshot("accept-failures-recovered")

    def lifecycle(self):
        self.open_dialog()
        self.type_field(NAME, "unsaved")
        before = self.state()
        target = before["snapshot"]["Surfaces"][DIALOG]["Target"]
        receipts = sum(e["event"] == "dialog-result" for e in self.events)
        for stage in ("profile", "binding", "resource", "guard", "stale"):
            result = self.command("fail " + stage)["data"]
            state = self.state()
            self.check("failed " + stage + " reload preserves open surface and draft", result["status"] == "error"
                       and self.opened(state) and state["snapshot"]["Surfaces"][DIALOG]["Target"] == target
                       and self.field(state, NAME)["Draft"] == "unsaved"
                       and sum(e["event"] == "dialog-result" for e in self.events) == receipts)
        cursor = len(self.events)
        self.command("reload")
        result = self.receipt(cursor, "close")
        state = self.state()
        self.check("successful reload closes old opening once without Save", result["Reason"] == "reload"
                   and result["Surface"] == target and not self.opened(state) and self.calls(state, "Save") == 0)
        state = self.open_dialog()
        self.check("reopened successor has no unaccepted old form draft", not self.field(state, NAME)["Dirty"])
        cursor = len(self.events)
        self.command("parent-hide")
        result = self.receipt(cursor, "close")
        self.check("parent hide closes published child once", result["Reason"] == "parent-hidden")
        self.command("parent-show")
        self.check("parent show does not revive child opening", not self.opened(self.state()))
        self.open_dialog()
        if self.args.nonmodal:
            cursor = self.native("close-window", title=CHILD_TITLE)
            result = self.receipt(cursor, "close")
            self.check("native child close protocol produces one lifecycle Close", result["Reason"] == "parent-closed" and result["Sequence"] == 0 and self.calls(self.state(), "Save") == 0)
            self.open_dialog()
        cursor = self.native("close-window")
        result = self.receipt(cursor, "close")
        self.wait(lambda e: e["event"] == "closed", cursor)
        self.check("native parent close drains child result before teardown", result["Reason"] == "parent-closed")


def main():
    p = argparse.ArgumentParser()
    p.add_argument("--binary", required=True)
    p.add_argument("--display", required=True)
    p.add_argument("--out", required=True)
    p.add_argument("--variant", choices=("commands", "context", "keyboard", "focus", "dialog", "nested", "command_failures", "accept_failures", "lifecycle"), default="commands")
    p.add_argument("--nonmodal", action="store_true")
    args = p.parse_args()
    trial = None
    try:
        trial = Trial(args)
        getattr(trial, args.variant)()
        trial.record("result", "passed")
    except Exception as error:
        if trial:
            trial.record("result", {"failed": str(error)})
        raise
    finally:
        if trial:
            trial.close()

if __name__ == "__main__":
    main()
