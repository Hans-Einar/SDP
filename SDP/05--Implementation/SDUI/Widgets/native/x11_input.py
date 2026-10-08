#!/usr/bin/env python3
"""Send actual X11 test input to an isolated native acceptance display.

Only the explicitly supplied display is opened. This is test tooling, not SDUI
runtime code. Coordinates are logical X11 window-relative pixels.
"""
import argparse
import ctypes as c
import ctypes.util
import json


def main():
    parser = argparse.ArgumentParser()
    parser.add_argument('--display', required=True)
    parser.add_argument('--title', required=True)
    parser.add_argument('--keep-focus', action='store_true', help='Do not force keyboard focus before input')
    parser.add_argument('action', choices=('locate', 'focus-state', 'move-window', 'close-window', 'click', 'right-click', 'double-click', 'wheel', 'drag', 'move', 'button-down', 'button-up', 'key', 'chord', 'key-down', 'key-up'))
    parser.add_argument('values', nargs='*')
    args = parser.parse_args()
    x = c.CDLL(ctypes.util.find_library('X11'))
    xt = c.CDLL(ctypes.util.find_library('Xtst'))
    dpy, win, atom = c.c_void_p, c.c_ulong, c.c_ulong
    def signature(lib, name, result, *parameters):
        fn = getattr(lib, name)
        fn.restype = result
        fn.argtypes = parameters
        return fn
    signature(x, 'XOpenDisplay', dpy, c.c_char_p)
    signature(x, 'XCloseDisplay', c.c_int, dpy)
    signature(x, 'XDefaultRootWindow', win, dpy)
    signature(x, 'XQueryTree', c.c_int, dpy, win, c.POINTER(win), c.POINTER(win), c.POINTER(c.POINTER(win)), c.POINTER(c.c_uint))
    signature(x, 'XFetchName', c.c_int, dpy, win, c.POINTER(c.c_void_p))
    signature(x, 'XFree', c.c_int, c.c_void_p)
    signature(x, 'XTranslateCoordinates', c.c_int, dpy, win, win, c.c_int, c.c_int, c.POINTER(c.c_int), c.POINTER(c.c_int), c.POINTER(win))
    signature(x, 'XSetInputFocus', c.c_int, dpy, win, c.c_int, c.c_ulong)
    signature(x, 'XGetInputFocus', c.c_int, dpy, c.POINTER(win), c.POINTER(c.c_int))
    signature(x, 'XMoveWindow', c.c_int, dpy, win, c.c_int, c.c_int)
    signature(x, 'XInternAtom', atom, dpy, c.c_char_p, c.c_int)
    signature(x, 'XSendEvent', c.c_int, dpy, win, c.c_int, c.c_long, c.c_void_p)
    signature(x, 'XStringToKeysym', atom, c.c_char_p)
    signature(x, 'XKeysymToKeycode', c.c_ubyte, dpy, atom)
    signature(x, 'XSync', c.c_int, dpy, c.c_int)
    signature(xt, 'XTestFakeMotionEvent', c.c_int, dpy, c.c_int, c.c_int, c.c_int, c.c_ulong)
    signature(xt, 'XTestFakeButtonEvent', c.c_int, dpy, c.c_uint, c.c_int, c.c_ulong)
    signature(xt, 'XTestFakeKeyEvent', c.c_int, dpy, c.c_uint, c.c_int, c.c_ulong)
    display = x.XOpenDisplay(args.display.encode())
    if not display:
        raise SystemExit('Cannot open explicit test display')
    try:
        root = x.XDefaultRootWindow(display)
        def find(window, depth=0):
            if depth > 32:
                return None
            name = c.c_void_p()
            if x.XFetchName(display, window, c.byref(name)) and name.value:
                raw_name = c.string_at(name.value)
                try:
                    text = raw_name.decode('utf-8')
                except UnicodeDecodeError:
                    # Legacy WM_NAME may use STRING/Latin-1 (e.g. middle dot).
                    text = raw_name.decode('latin-1')
                x.XFree(name)
                if text == args.title:
                    return window
            returned_root, parent, children, count = win(), win(), c.POINTER(win)(), c.c_uint()
            if not x.XQueryTree(display, window, c.byref(returned_root), c.byref(parent), c.byref(children), c.byref(count)):
                return None
            try:
                for i in range(count.value):
                    result = find(children[i], depth + 1)
                    if result:
                        return result
            finally:
                if children:
                    x.XFree(c.cast(children, c.c_void_p))
            return None
        target = find(root)
        if target is None:
            raise SystemExit('Named test window not found')
        left, top, child = c.c_int(), c.c_int(), win()
        if not x.XTranslateCoordinates(display, target, root, 0, 0, c.byref(left), c.byref(top), c.byref(child)):
            raise SystemExit('Cannot resolve test window position')
        if args.action == 'locate':
            print(json.dumps(dict(window=target, x=left.value, y=top.value)))
            return
        if args.action == 'focus-state':
            focused, revert = win(), c.c_int()
            x.XGetInputFocus(display, c.byref(focused), c.byref(revert))
            print(json.dumps(dict(window=target, focused=focused.value, matches=focused.value == target)))
            return
        if args.action == 'move-window':
            # Arrange isolated undecorated windows so both remain reachable;
            # this does not simulate any product command or focus decision.
            if len(args.values) != 2:
                raise SystemExit('Window move requires x y')
            x.XMoveWindow(display, target, int(args.values[0]), int(args.values[1]))
            x.XSync(display, 0)
            return
        if args.action == 'close-window':
            # Send the native WM protocol, without claiming a decoration click
            # or requiring a window manager on the isolated acceptance display.
            class ClientMessage(c.Structure):
                _fields_ = [('type', c.c_int), ('serial', c.c_ulong),
                            ('send_event', c.c_int), ('display', dpy),
                            ('window', win), ('message_type', atom),
                            ('format', c.c_int), ('data', c.c_long * 5)]
            class Event(c.Union):
                _fields_ = [('client', ClientMessage), ('pad', c.c_long * 24)]
            event = Event()
            event.client.type = 33  # ClientMessage
            event.client.send_event = 1
            event.client.display = display
            event.client.window = target
            event.client.message_type = x.XInternAtom(display, b'WM_PROTOCOLS', 0)
            event.client.format = 32
            event.client.data[0] = x.XInternAtom(display, b'WM_DELETE_WINDOW', 0)
            if not x.XSendEvent(display, target, 0, 0, c.byref(event)):
                raise SystemExit('Native close protocol could not be sent')
            x.XSync(display, 0)
            return
        if not args.keep_focus:
            x.XSetInputFocus(display, target, 1, 0)
        def button(number):
            xt.XTestFakeButtonEvent(display, number, 1, 0)
            xt.XTestFakeButtonEvent(display, number, 0, 40)
        if args.action in ('click', 'right-click', 'double-click', 'wheel', 'drag', 'move', 'button-down', 'button-up'):
            if len(args.values) < 2:
                raise SystemExit('Pointer action requires x y')
            px, py = (int(v) for v in args.values[:2])
            xt.XTestFakeMotionEvent(display, -1, left.value + px, top.value + py, 0)
            if args.action == 'move':
                pass
            elif args.action in ('button-down', 'button-up'):
                number = int(args.values[2]) if len(args.values) > 2 else 1
                if number not in (1, 2, 3):
                    raise SystemExit('Held pointer button must be 1, 2 or 3')
                xt.XTestFakeButtonEvent(display, number, args.action == 'button-down', 0)
            elif args.action == 'drag':
                if len(args.values) != 4:
                    raise SystemExit('Drag requires x1 y1 x2 y2')
                end_x, end_y = (int(v) for v in args.values[2:])
                xt.XTestFakeButtonEvent(display, 1, 1, 0)
                for step in range(1, 13):
                    move_x = round(px + (end_x - px) * step / 12)
                    move_y = round(py + (end_y - py) * step / 12)
                    xt.XTestFakeMotionEvent(display, -1, left.value + move_x, top.value + move_y, 20)
                xt.XTestFakeButtonEvent(display, 1, 0, 40)
            elif args.action == 'wheel':
                delta = int(args.values[2]) if len(args.values) > 2 else 1
                for _ in range(min(abs(delta), 100)):
                    button(5 if delta > 0 else 4)
            else:
                button(3 if args.action == 'right-click' else 1)
                if args.action == 'double-click':
                    button(1)
        else:
            pressed = []
            for key in args.values:
                code = x.XKeysymToKeycode(display, x.XStringToKeysym(key.encode()))
                if not code:
                    raise SystemExit('Unknown key: ' + key)
                xt.XTestFakeKeyEvent(display, code, args.action != 'key-up', 0)
                if args.action in ('key-down', 'key-up'):
                    pass
                elif args.action == 'chord':
                    pressed.append(code)
                else:
                    xt.XTestFakeKeyEvent(display, code, 0, 40)
            for code in reversed(pressed):
                xt.XTestFakeKeyEvent(display, code, 0, 40)
        x.XSync(display, 0)
    finally:
        x.XCloseDisplay(display)


if __name__ == '__main__':
    main()
