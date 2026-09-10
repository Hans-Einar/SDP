# Kartlegging av eksisterende designspråk

Researchdato: **2026-09-10**. Omfang: **32 profiler** av språk, språkfamilier, notasjoner og uttrykkelig merkede tilgrensende formater/rammeverk. Dette er en bred, målrettet oversikt, ikke en uttømmende liste over alle språk i bruk.

Les [syntesen](../README.md) for SDP-behov, vesentlige skiller, åpne spørsmål og neste oppgave.

## Dekning

| Profil | Kategori |
| --- | --- |
| [UML 2.5.1](uml.md) | Modell-/designspråk |
| [SysML v2](sysml-v2.md) | Systemmodelleringsspråk |
| [ArchiMate](archimate.md) | Arkitekturmodell og grafisk notasjon |
| [AADL](aadl.md) | Architecture Description Language |
| [Structurizr DSL / C4](structurizr-dsl.md) | Arkitekturmodell som tekst |
| [LikeC4](likec4.md) | Utvidbar arkitektur-DSL |
| [PlantUML](plantuml.md) | Diagramnotasjon |
| [Mermaid](mermaid.md) | Diagramnotasjon |
| [D2](d2.md) | Diagram-DSL |
| [Context Mapper CML](context-mapper-cml.md) | DDD- og arkitektur-DSL |
| [Universal Variability Language (UVL)](uvl.md) | Feature-/variabilitetsmodell |
| [Clafer](clafer.md) | Struktur-/featuremodell med constraints |
| [BPMN 2.0.2](bpmn.md) | Prosessmodell og notasjon |
| [SCXML 1.0](scxml.md) | Tilstandsmaskin-/atferdsspråk |
| [TLA+ og PlusCal](tla-plus.md) | Formell atferdsspesifikasjon |
| [Alloy 6](alloy.md) | Relasjonelt formelt modelleringsspråk |
| [Object Constraint Language (OCL) 2.4](ocl.md) | Modellconstraints og queries |
| [OpenAPI](openapi.md) | HTTP API-kontrakt |
| [AsyncAPI 3.0](asyncapi.md) | Event-/meldingskontrakt |
| [Protocol Buffers](protobuf.md) | Data-/service-IDL |
| [Smithy 2.0](smithy.md) | Service-/data-IDL |
| [JSON Schema 2020-12](json-schema.md) | Datavalidering og schema |
| [CUE](cue.md) | Constraint-/konfigurasjonsspråk |
| [DBML](dbml.md) | Databasemodell-/diagram-DSL |
| [Ecore / Eclipse Modeling Framework](ecore.md) | Metamodel og modelleringsinfrastruktur |
| [ATL](atl.md) | Model-to-model-transformasjon |
| [QVT 1.3](qvt.md) | Query/View/Transformation-standard |
| [Epsilon-språkfamilien](epsilon.md) | Modellvalidering, transformasjon og generering |
| [Edapt](edapt.md) | Modellhistorikk-/migrasjonsrammeverk (tilgrensende) |
| [Rego / Open Policy Agent](rego.md) | Policyspråk |
| [ReqIF 1.2](reqif.md) | Kravutveksling (tilgrensende format) |
| [PROV-O med RDF/Turtle](prov-o.md) | Provenance-ontologi (tilgrensende) |

## Metode og evidensgrenser

Primærkilder er standardorganisasjonenes spesifikasjoner, utviklernes dokumentasjon og offisielle repositories. Hver profil skiller dokumentert formål/mekanismer fra vår vurdering av styrker, begrensninger og SDP-relevans. Negative konklusjoner om manglende livsløpsstøtte gjelder den undersøkte standard-/dokumentasjonsflaten, ikke alle mulige plugins eller forskningsutvidelser.

Alle eksempler er **egne, illustrative eksempler**, ikke kopierte normative eksempler. De er ikke parser-/runtime-testet. Fragmenter og pseudonotasjon er merket; de er ikke implementeringsklare fixtures. Full representasjon av Ponsse tilhører neste oppgave.

«Tilgjengelig dokumentasjon», «ikke arkivert» og en nylig push er forskjellige aktivitetssignaler. Ingen av dem alene beviser produksjonsmodenhet, utbredelse eller langsiktig vedlikehold. Oppførte standardversjoner er observerte referanser, ikke påstander om at alle integrasjoner støtter dem. Fremtidsdaterte releases behandles ikke som utgitt.

Lisensnavn fra GitHub er metadata-signaler, ikke full gjennomgang av alle filer eller en juridisk konklusjon. `NOASSERTION` betyr at metadata ikke ga et entydig svar. Standardtekst, parser, generator, runtime og hosted tjeneste kan ha forskjellige vilkår. Uavklarte vilkår er eksplisitte i profilene og skal avklares for konkrete kandidatverktøy før integrasjon.

## Viktige dekningshull

ArchiMate er bare delvis kontrollert: Archi-utviklernes håndbok var tilgjengelig, men The Open Groups normative sider lot seg ikke hente. Ikke bruk denne profilen som full normativ transition-vurdering.

Ikke undersøkt som egne profiler i denne runden: SysML v1, Capella/Arcadia, Modelica, SDL, Acme/Wright, DMN, Event-B, Z, SHACL, GraphQL, Avro, Thrift, RAML og TypeSpec. Xtext, MPS og Sirius er mulige språk-/editorbyggingsverktøy, men verktøyvalg er ikke gjort. Disse hullene betyr at ingen kandidat kan erklæres «best av alle språk» ut fra katalogen. Neste sammenligning kan ta inn én av disse dersom et konkret udekket behov krever det.

YAML og JSON er mulige serialiseringsformer, ikke i seg selv en semantisk SDP-modell. C4 er behandlet sammen med Structurizr DSL; EMF/Edapt, ReqIF og PROV-O er tatt med for avgrensede tilleggsbehov og telles ikke som fullverdige alternative systemdesignspråk.

## Hva som er kontrollert

Dokumentene er kontrollert for forventede seksjoner, eksempler, kildehenvisninger, unike filnavn og interne lenkemål. Nettkilder er lest som research; ingen komplett link-crawler eller installasjon av 32 verktøy er kjørt. Repositorydiff skal bare inneholde denne researchkatalogen. Ingen uavhengig verifikasjonsgodkjenning eller produkt-testpass påstås.
