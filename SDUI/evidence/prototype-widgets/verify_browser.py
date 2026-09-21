import json
from pathlib import Path
from playwright.sync_api import sync_playwright
root=Path(__file__).resolve().parents[2]
out=root/'evidence/prototype-widgets'
with sync_playwright() as p:
    browser=p.chromium.launch(executable_path='/usr/bin/google-chrome', headless=True,args=['--no-sandbox','--disable-dev-shm-usage'])
    page=browser.new_page(viewport={'width':1100,'height':850},device_scale_factor=1)
    requests=[]; errors=[]
    page.on('request',lambda r:requests.append(r.url))
    page.on('pageerror',lambda e:errors.append(str(e)))
    page.goto((root/'examples/prototype-controls.html').as_uri())
    page.get_by_label('Operatør').fill('Kari Nordmann')
    note='Kontroller barkmåling før neste kapp. Dette er en lengre utfylt merknad som skal brytes over flere linjer i utskriften uten at feltets bredde skjuler teksten.'
    page.get_by_label('Merknad').fill(note)
    button=page.get_by_role('button',name='OK',exact=True)
    button.hover();page.mouse.down();page.wait_for_timeout(100)
    assert button.evaluate('(el)=>getComputedStyle(el).transform')=='matrix(1, 0, 0, 1, 0, 2)'
    page.screenshot(path=str(out/'browser-pressed.png'),full_page=True)
    page.mouse.up()
    assert 'Lokalt knappetrykk 1: OK' in page.locator('#status').inner_text()
    button.focus();page.keyboard.press('Space')
    assert 'Lokalt knappetrykk 2: OK' in page.locator('#status').inner_text()
    assert page.get_by_role('button',name='Utilgjengelig').is_disabled()
    page.emulate_media(media='print')
    page.evaluate("window.dispatchEvent(new Event('beforeprint'))")
    assert page.get_by_label('Operatør').evaluate('(el)=>getComputedStyle(el).display')=='none'
    mirrors=page.locator('.print-value')
    assert mirrors.nth(0).is_visible() and mirrors.nth(0).inner_text()=='Kari Nordmann'
    assert mirrors.nth(2).inner_text()==note
    assert mirrors.nth(2).bounding_box()['height']>50
    page.screenshot(path=str(out/'print-preview.png'),full_page=True)
    page.emulate_media(media='screen')
    page.set_viewport_size({'width':390,'height':844})
    assert page.evaluate('document.documentElement.scrollWidth <= innerWidth')
    page.reload()
    assert page.get_by_label('Operatør').input_value()=='Ola Nordmann'
    assert not errors and all(u.startswith('file:') for u in requests)
    (out/'browser-verification.json').write_text(json.dumps({'browser':browser.version,'checks':['editing','mouse press appearance','keyboard activation','disabled button','print mirrors current values','print text wrapping','mobile no horizontal overflow','reload restores defaults','no external requests','no JavaScript errors'],'requests':requests,'printVerification':'CSS print media emulation; no physical printing or PDF verification'},indent=2)+'\n')
    browser.close()
print('10 browser checks passed')
