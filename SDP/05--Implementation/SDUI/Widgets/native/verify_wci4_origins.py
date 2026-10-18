#!/usr/bin/env python3
"""Verify real X11 screenshot pixels against SVG-origin geometry, not aspect alone."""
import argparse
import json
import math
from pathlib import Path
import subprocess
import sys
import verify_wci4_previews as previews

GREEN = bytes((35, 130, 72))

class Trial(previews.Trial):
    def pixels(self, name, slot):
        s = self.state()
        p = self.preview(previews.HERO, s)
        self.check(name + ' mounted rendered image', p['status'] == 'rendered' and p['visible'])
        window = json.loads(subprocess.check_output([sys.executable, str(Path(__file__).with_name('x11_input.py')), '--display', self.args.display, '--title', previews.TITLE, 'locate']))
        self.check(name + ' mapped native window', window['viewable'])
        self.screenshot(name)
        path = self.out / (name + '.png')
        width, height = map(int, subprocess.check_output(['identify', '-format', '%w %h', str(path)]).split())
        rgb = subprocess.check_output(['convert', str(path), '-alpha', 'off', '-depth', '8', 'RGB:-'])
        assert len(rgb) == width * height * 3
        im, clip = p['image'], p['clip']
        # This recorded native profile fixes FYNE_SCALE=1. Coordinates are canvas
        # physical pixels; the actual X window origin is added exactly once.
        x, y = window['x'], window['y']
        # Scan three physical pixels outside the clip too: these fixture
        # margins contain no other #238248 artwork and expose leaked paint.
        roi = [max(0, math.floor(x+clip['X'])-3), max(0, math.floor(y+clip['Y'])-3), min(width, math.ceil(x+clip['X']+clip['W'])+3), min(height, math.ceil(y+clip['Y']+clip['H'])+3)]
        rect = [x+im['X'], y+im['Y'], x+im['X']+im['W'], y+im['Y']+im['H']]
        if slot == 'root-transform':
            rect = [x+im['X']+.1*im['W'], y+im['Y']+.1*im['H'], x+im['X']+.3*im['W'], y+im['Y']+.3*im['H']]
        expected = [max(rect[0],x+clip['X'],0), max(rect[1],y+clip['Y'],0), min(rect[2],x+clip['X']+clip['W'],width), min(rect[3],y+clip['Y']+clip['H'],height)]
        def green_coverage(offset):
            r,g,b=rgb[offset:offset+3]
            alpha=(255-r)/220
            return alpha if abs(g-(255-125*alpha))<=2 and abs(b-(255-183*alpha))<=2 else 0
        # Fractional native image placement adds linear filtering to SVG edge
        # antialiasing. Measure the 50%-coverage contour for geometry; exact
        # green remains mandatory throughout the separately inset interior.
        points = [(px,py) for py in range(roi[1],roi[3]) for px in range(roi[0],roi[2]) if green_coverage((py*width+px)*3)>=.5]
        bbox = [min(q[0] for q in points),min(q[1] for q in points),max(q[0] for q in points)+1,max(q[1] for q in points)+1] if points else None
        interior = [math.ceil(expected[0]+1),math.ceil(expected[1]+1),math.floor(expected[2]-1),math.floor(expected[3]-1)]
        solid = all(rgb[(py*width+px)*3:(py*width+px)*3+3] == GREEN for py in range(interior[1],interior[3]) for px in range(interior[0],interior[2]))
        self.record('pixel-oracle', dict(name=name, slot=slot, window=window, profile={'FYNE_SCALE':'1','FYNE_THEME':'light'}, image=im, clip=clip, expected=expected, observed=bbox, halfCoveragePixels=len(points), coverageThreshold=.5, interiorSolid=solid))
        self.check(name + ' actual half-coverage green bounds within one physical pixel', bbox is not None and all(abs(a-b)<=1.01 for a,b in zip(bbox,expected)))
        self.check(name + ' solid interior and no green outside expected bounds', solid and all(expected[0]-1<=px<expected[2]+1 and expected[1]-1<=py<expected[3]+1 for px,py in points))
        return im['W']/400

    def origins(self):
        scales=[]
        for slot in ('negative-origin','positive-origin','mixed-origin','root-transform'):
            self.command('resource Hero '+slot)
            self.command('reload')
            for size in (500,1100):
                self.command('resize '+str(size)+' 850')
                self.observed(lambda d:self.preview(previews.HERO,d)['visible'])
                scales.append(self.pixels(slot+'-'+str(size),slot))
        self.check('actual native fit includes below and above intrinsic scale',min(scales)<1<max(scales))
        self.command('resource Hero negative-origin');self.command('reload')
        before=self.preview(previews.HERO)
        clip=self.state()['viewports'][previews.ART]['Clip']
        cursor=self.center(clip,'wheel',8)
        self.changed(lambda d:d['viewports'][previews.ART]['Offset']['Y']>120,cursor)
        after=self.preview(previews.HERO)
        self.check('shared clip translates image without refitting',after['image']['Y']<before['image']['Y'] and abs(after['image']['W']-before['image']['W'])<.01)
        self.pixels('negative-origin-shared-clip','negative-origin')


def main():
    p=argparse.ArgumentParser()
    p.add_argument('--binary',required=True);p.add_argument('--display',required=True);p.add_argument('--out',required=True)
    p.add_argument('--reproduce-old',action='store_true')
    a=p.parse_args();a.variant='origins';a.nonmodal=False;t=None
    import os
    os.environ['FYNE_SCALE']='1';os.environ['FYNE_THEME']='light'
    try:
        t=Trial(a)
        if a.reproduce_old:
            t.command('resource Hero negative-origin');t.command('reload')
            t.observed(lambda d:t.preview(previews.HERO,d)['visible'])
            t.pixels('negative-origin-old','negative-origin')
        else:t.origins()
        t.record('result','passed')
    except Exception as e:
        if t:t.record('result',{'failed':str(e)})
        raise
    finally:
        if t:t.close()

if __name__=='__main__':main()
