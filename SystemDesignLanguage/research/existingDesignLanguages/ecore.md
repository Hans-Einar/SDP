# Ecore / Eclipse Modeling Framework

[Katalog](README.md) · Kategori: **Metamodel og modelleringsinfrastruktur** · Research: **2026-09-10**

## Formål og abstraksjonsnivå

Ecore definerer modellens typer, attributter og referanser. EMF er rammeverket rundt modellen, blant annet serialisering og kodegenerering; ikke ett selvstendig designspråk.

## Modellmekanismer

Identitet: EPackage namespace URI og modelelementer; instans-ID-policy må defineres. Relasjoner: EReference med cardinality/containment. Contracts: struktur; rikere constraints via andre språk/verktøy. Utvidelse: egen metamodel. Views/editorer krever verktøy.

## Styrker og begrensninger — vår vurdering

**Styrke:** Et etablert grunnlag dersom SDP trenger egen semantikk og integrasjon med modelltransformasjoner.

**Begrensning:** En egen Ecore-metamodel er fortsatt vårt eget språkansvar. JVM-/Eclipse-infrastruktur og model-persistence kan være mer enn små prosjekter trenger.

## Historikk, endring og transitions

Ecore beskriver typer, ikke automatisk deres evolusjon. Edapt og Epsilon Flock er separate mekanismer for migrasjon; Git kan supplere med revisjoner.

## Illustrativt eksempel

Illustrativ Ecore-klasse uten attributter; ikke en SDP-metamodelbeslutning. Eksemplet er ikke parser-/runtime-testet.

```xml
<ecore:EPackage xmi:version="2.0"
 xmlns:xmi="http://www.omg.org/XMI"
 xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance"
 xmlns:ecore="http://www.eclipse.org/emf/2002/Ecore"
 name="demo" nsURI="urn:example:demo" nsPrefix="demo">
  <eClassifiers xsi:type="ecore:EClass" name="Feature"/>
</ecore:EPackage>
```

## Verktøy, vedlikehold og vilkår

Eclipse-prosjektet og dokumentasjonen er tilgjengelige; prosjektsiden oppgir EPL-2.0. Konkrete bundle-NOTICE-filer og versjonskompatibilitet må kontrolleres før bruk. Ingen toolchain er installert eller valgt.

## Primærkilder

Alle kilder kontrollert 2026-09-10; se katalogens metode for evidens- og lisensbegrensninger.

- [Eclipse EMF prosjekt](https://projects.eclipse.org/projects/modeling.emf)
