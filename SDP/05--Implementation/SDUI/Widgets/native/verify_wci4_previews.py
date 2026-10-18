#!/usr/bin/env python3
"""Real XTest preview geometry, immutable preparation and connected SDL workflows."""
import argparse
import time
import verify_wci3_text as text
import verify_wci3_values as scalar

TITLE='SDUI WCI4 Previews'
TABS='page/view/body/workspace'
SPLIT='page/view/body'
ART=TABS+'/art/viewport'
REPORTS=TABS+'/reports/reportViewport'
HERO,PROSE,TAIL=[ART+'/'+n for n in ('hero','prose','tail')]
RA,RB,RSVG=[REPORTS+'/'+n for n in ('reportA','reportB','reportSVG')]
DETAIL='page/detail'
DSVG,DDOC=[DETAIL+'/detailViewport/'+n for n in ('detailSVG','detailDoc')]
NOTE='page/view/body/aside/note'
DNOTE=DETAIL+'/detailNote'
LOAD='page/view/body/aside/loadStatus'
MARK='page/view/body/aside/markStatus'
BAR='page/view/header/toolbar/'

class Trial(text.Trial):
    def native(self,action,*values,title=TITLE,keep_focus=False):
        if action in ('drag','wheel'):
            scalar.Trial.native(self,'move',*values[:2],title=title,keep_focus=keep_focus)
            time.sleep(.1) # Native pointer delivery, not a publication barrier.
        return scalar.Trial.native(self,action,*values,title=title,keep_focus=keep_focus)

    def center(self,rect,action='click',*args,title=TITLE):
        return self.native(action,round(rect['X']+rect['W']/2),round(rect['Y']+rect['H']/2),*args,title=title)

    def tab(self,name):
        s=self.state();page=next(p for p in s['tabs'][TABS]['pages'] if p['id']==name)
        cursor=self.center(page['clip'])
        return self.changed(lambda d:d['tabs'][TABS]['selected']==name,cursor)

    def preview(self,path,state=None):
        return (state or self.state())['previews'][path]

    def unchanged(self,before,after):
        return before['publishedIdentity']==after['publishedIdentity'] and before['source']==after['source'] and before['previews']==after['previews'] and before['actionCalls']==after['actionCalls']

    def geometry(self):
        initial=self.state();p=self.preview(HERO,initial);baseline=initial['rendererCalls']
        self.record('observed-preparation-renderer-baseline',baseline)
        self.check('direct SVG and bounded prose are rendered with no SDL action',p['status']=='rendered' and self.preview(PROSE,initial)['status']=='rendered' and not any(initial['actionCalls'].values()))
        self.check('actual mounted Accessible describes rendered image separately from caption',p['mounted'] and p['accessibleLabel']=='Output by site (rendered)' and p['accessibleRole']=='text' and p['caption']['H']>0)
        self.check('full image preserves 2:1 aspect and ends before caption',abs(p['image']['W']/p['image']['H']-2)<1e-7 and p['image']['Y']+p['image']['H']<=p['caption']['Y'])
        self.check('inactive and closed content is prepared without fabricated mounted objects',all(not self.preview(x,initial)['mounted'] for x in (RA,RB,RSVG,DSVG,DDOC)))
        self.screenshot('native-shape-montage-and-prose')
        self.command('mutate-input Hero');s=self.state()
        self.check('mutating caller loan leaves live image digest and geometry unchanged',s['loanSHA256'][HERO]!=p['sha256'] and self.preview(HERO,s)==p and s['publishedIdentity']==initial['publishedIdentity'])
        self.command('caption Hero long');self.command('resize 900 780');s=self.state()
        self.check('caption resize retains full accessibility and frozen renderer baseline',self.preview(HERO,s)['accessibleLabel']==p['accessibleLabel'] and s['rendererCalls']==baseline)
        self.screenshot('long-caption-relative-fit')
        start=self.preview(HERO);clip=self.state()['viewports'][ART]['Clip']
        cursor=self.center(clip,'wheel',8)
        s=self.changed(lambda d:d['viewports'][ART]['Offset']['Y']>120,cursor);after=self.preview(HERO,s)
        offset=s['viewports'][ART]['Offset']['Y']
        self.check('outer scroll translates full image without refitting visible fragment',after['rect']['Y']<0 and abs(after['image']['W']-start['image']['W'])<.01 and abs(after['image']['Y']-(start['image']['Y']-offset))<.01 and after['clip']['Y']+.01>=s['viewports'][ART]['Clip']['Y'])
        self.check('preview never overlaps outer scrollbar gutter',after['clip']['X']+after['clip']['W']<=s['viewports'][ART]['VerticalGutter']['X']+.01)
        self.screenshot('negative-box-under-shared-clip')
        s=self.tab('reports')
        self.check('equal Markdown paths keep independent descriptions and partial outcomes',self.preview(RA,s)['status']=='partial' and self.preview(RB,s)['status']=='partial' and self.preview(RA,s)['accessibleLabel']!=self.preview(RB,s)['accessibleLabel'])
        self.screenshot('per-diagram-fallback-with-surrounding-prose')
        s=self.tab('art')
        self.check('tab reveal preserves outer offset and renderer baseline',abs(s['viewports'][ART]['Offset']['Y']-offset)<.01 and s['rendererCalls']==baseline)
        r=s['splits'][SPLIT]['DividerClip'];cursor=self.center(r)
        self.native('chord','Control_L','Home');s=self.observed(lambda d:d['snapshot']['Splits'][SPLIT]['Collapsed']=='first')
        self.check('split collapse detaches visible preview geometry',not self.preview(HERO,s)['visible'] and s['splits'][SPLIT]['First']['W']==0)
        self.native('key','space');s=self.observed(lambda d:d['snapshot']['Splits'][SPLIT]['Collapsed']=='none')
        self.check('split restore preserves content and invokes no renderer or domain action',s['rendererCalls']==baseline and s['actionCalls']==initial['actionCalls'] and self.preview(HERO,s)['sha256']==p['sha256'])

    def fallback(self):
        for slot in ('wide','tall','negative-origin'):
            self.command('resource Hero '+slot);self.command('reload');p=self.preview(HERO)
            ratio={'wide':4,'tall':.25,'negative-origin':2}[slot]
            self.check(slot+' SVG is rendered at its actual intrinsic aspect',p['status']=='rendered' and abs(p['image']['W']/p['image']['H']-ratio)<1e-7)
            self.screenshot('resource-'+slot)
        self.command('policy Hero label');self.command('resource Hero malformed');self.command('reload');p=self.preview(HERO)
        self.check('malformed content becomes a visible explicit accessible fallback',p['status']=='label' and p['diagnostic'] and p['statusRect']['H']>0 and p['image']['W']==0 and 'unavailable' in p['accessibleLabel'].lower())
        self.screenshot('visible-svg-fallback')
        for alias,slot in [('Hero','unsupported'),('Hero','missing'),('Hero','bad-dimensions')]:
            self.command('resource '+alias+' '+slot);self.command('reload');self.check(slot+' obeys declared label',self.preview(HERO)['status']=='label')
        before=self.state();self.command('binding Hero digest');self.command('reload',expected='error');after=self.state()
        self.check('digest identity failure remains fatal under label and retains old bundle',self.unchanged(before,after))
        self.command('resource Hero good');self.command('policy Hero reject');self.command('reload')
        for alias,slot in [('Hero','malformed'),('DetailSVG','missing')]:
            before=self.state();self.command('resource '+alias+' '+slot);self.command('reload',expected='error')
            self.check('reject on '+alias+' retains actual published bundle',self.unchanged(before,self.state()))
            self.command('resource '+alias+' good');self.command('reload')
        before=self.state();self.command('policy ReportB reject');self.command('reload',expected='error')
        self.check('hidden same-text reject path cannot borrow another path label policy',self.unchanged(before,self.state()))
        self.command('policy ReportB label');self.command('renderer ReportA error');self.command('reload');s=self.tab('reports')
        self.check('failed diagram renderer keeps surrounding Markdown as partial',self.preview(RA,s)['status']=='partial' and self.preview(RB,s)['status']=='partial')
        self.screenshot('renderer-error-keeps-prose')
        self.command('document ReportA html');self.command('reload');s=self.state()
        self.check('unsupported whole document replaces only its own path with label',self.preview(RA,s)['status']=='label' and self.preview(RB,s)['status']=='partial')
        self.screenshot('whole-document-fallback-independent-path')
        before=self.state();self.command('binding ReportB revision');self.command('reload',expected='error')
        self.check('renderer identity failure retains independent published outcomes',self.unchanged(before,self.state()))

    def lifetime(self):
        initial=self.state();baseline=initial['rendererCalls'];p=self.preview(HERO,initial)
        self.command('prepare');prepared=self.state();self.check('preparation holds detached candidate without publication',prepared['candidateHeld'] and prepared['publishedIdentity']==initial['publishedIdentity'])
        counts=prepared['rendererCalls'];self.command('resource Hero alternate');self.command('publish',expected='error');s=self.state()
        self.check('stale resource guard rejects publication without replay or live replacement',not s['candidateHeld'] and s['publishedIdentity']==initial['publishedIdentity'] and self.preview(HERO,s)==p and s['rendererCalls']==counts)
        self.command('reload');s=self.state();self.check('explicit fresh preparation publishes changed digest',self.preview(HERO,s)['sha256']!=p['sha256'])
        self.command('prepare');before=self.state();self.command('renderer-revision ReportA r2');self.command('publish',expected='error');s=self.state()
        self.check('stale renderer revision rejects held candidate and retains current bundle',not s['candidateHeld'] and self.unchanged(before,s) and s['rendererCalls']==before['rendererCalls'])
        self.command('reload');self.command('prepare');before=self.state();self.command('abandon');s=self.state()
        self.check('abandon disposes only detached candidate without altering live content',not s['candidateHeld'] and self.unchanged(before,s) and s['rendererCalls']==before['rendererCalls'])
        self.command('fail resource-ticket');before=self.state();self.command('caption Hero short',expected='error');s=self.state()
        self.check('resource ticket failure preserves published image and zero renderer calls',self.unchanged(before,s) and s['rendererCalls']==before['rendererCalls'])
        self.command('caption Hero short');s=self.state();self.check('successful presentation ticket uses frozen resources without rerender',s['rendererCalls']==before['rendererCalls'] and self.preview(HERO,s)['sha256']==self.preview(HERO,before)['sha256'])
        self.screenshot('replacement-after-guard-recovery')

    def forms(self):
        baseline=self.state()['rendererCalls']
        cursor=self.control(BAR+'loadButton');s=self.changed(lambda d:self.calls(d,'Load')==1,cursor)
        self.check('actual SDL Load publishes into extended readonly receiver',self.accepted(s,LOAD)!='Not loaded')
        cursor=self.control(BAR+'markButton');s=self.changed(lambda d:self.calls(d,'Mark')==1,cursor)
        self.check('shared command publishes extended readonly result exactly once',self.accepted(s,MARK)!='No mark')
        cursor=self.native('chord','Control_L','m');s=self.changed(lambda d:self.calls(d,'Mark')==2,cursor)
        self.check('shared shortcut invokes one SDL action with no automatic replay',self.calls(s,'Mark')==2)
        self.paste(NOTE,'Saved from preview');cursor=self.key_for(NOTE,'key','Return');s=self.changed(lambda d:self.calls(d,'SaveNote')==1,cursor)
        self.check('ordinary text commit coexists with frozen previews',self.accepted(s,NOTE)=='Saved from preview' and s['rendererCalls']==baseline)
        cursor=self.control(BAR+'openButton');s=self.changed(lambda d:d['snapshot']['Surfaces'][DETAIL]['Open'],cursor)
        p=self.preview(DSVG,s)
        self.check('dialog mounts its own visible accessible artwork',p['mounted'] and p['visible'] and p['accessibleLabel']=='Detail artwork (rendered)' and p['canvas']==(DETAIL if self.args.nonmodal else 'main'))
        self.screenshot('independent-detail-preview')
        self.paste(DNOTE,'Committed child');cursor=self.key_for(DNOTE,'key','Return');s=self.changed(lambda d:self.calls(d,'SaveDetail')==1,cursor)
        cursor=self.control(DETAIL+'/detailActions/cancelButton');s=self.changed(lambda d:not d['snapshot']['Surfaces'][DETAIL]['Open'],cursor)
        self.check('Cancel detaches dialog preview and preserves prior child domain commit',not self.preview(DSVG,s)['mounted'] and s['domain']['values'].get('SaveDetail')=='Committed child' and s['domain']['commits']['SaveDetail']==1 and self.accepted(s,DNOTE)=='Committed child')
        cursor=self.control(BAR+'openButton');self.changed(lambda d:d['snapshot']['Surfaces'][DETAIL]['Open'],cursor)
        self.paste(DNOTE,'Accepted detail');self.command('accept false');cursor=self.control(DETAIL+'/detailActions/acceptButton');s=self.changed(lambda d:self.calls(d,'AcceptDetail')==1,cursor)
        self.check('rejected SDL Accept retains dialog draft and image',s['snapshot']['Surfaces'][DETAIL]['Open'] and self.draft(s,DNOTE)=='Accepted detail' and self.preview(DSVG,s)['visible'])
        self.command('accept true');cursor=self.control(DETAIL+'/detailActions/acceptButton');s=self.changed(lambda d:not d['snapshot']['Surfaces'][DETAIL]['Open'],cursor)
        self.check('successful SDL Accept closes artwork with exactly two explicit calls',self.calls(s,'AcceptDetail')==2 and self.accepted(s,DNOTE)=='Accepted detail' and not self.preview(DSVG,s)['mounted'])
        cursor=self.control(BAR+'openButton');self.changed(lambda d:d['snapshot']['Surfaces'][DETAIL]['Open'],cursor)
        self.command('parent-hide');s=self.state();self.check('parent hide closes child preview without rerender',not s['snapshot']['Surfaces'][DETAIL]['Open'] and not self.preview(DSVG,s)['mounted'] and s['rendererCalls']==baseline)
        self.command('parent-show');cursor=self.control(BAR+'openButton');self.changed(lambda d:d['snapshot']['Surfaces'][DETAIL]['Open'],cursor)
        self.command('reload');s=self.state();self.check('accepted reload closes old dialog and publishes a fresh main preview',not s['snapshot']['Surfaces'][DETAIL]['Open'] and not self.preview(DSVG,s)['mounted'] and self.preview(HERO,s)['visible'])
        cursor=self.control(BAR+'openButton');self.changed(lambda d:d['snapshot']['Surfaces'][DETAIL]['Open'],cursor)
        cursor=self.native('close-window');self.wait(lambda e:e['event']=='closed',cursor)
        results=[e['data'] for e in self.events if e['event']=='dialog-result']
        self.check('five dialog openings each yield exactly one terminal receipt',len(results)==5 and len({(r['Surface']['Handle']['Path'],r['Surface']['ModelRevision'],r['Surface']['OpenGeneration']) for r in results})==5)


def main():
    p=argparse.ArgumentParser()
    p.add_argument('--binary',required=True);p.add_argument('--display',required=True);p.add_argument('--out',required=True)
    p.add_argument('--nonmodal',action='store_true');p.add_argument('--variant',choices=('geometry','fallback','lifetime','forms'),required=True)
    a=p.parse_args();t=None
    try:
        t=Trial(a);getattr(t,a.variant)();t.record('result','passed')
    except Exception as e:
        if t:t.record('result',{'failed':str(e)})
        raise
    finally:
        if t:t.close()

if __name__=='__main__':main()
