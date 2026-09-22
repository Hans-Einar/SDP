package viewpoint

type Spec struct {
	ID, Title, Note string
	Levels          []string
}

var Catalog = []Spec{
	{"VP01", "Bruksmål og sporbarhet", "pursues, supports og contributes-to; modellens omfang, uten oppdiktet System-grense.", []string{"A0", "A1"}},
	{"VP02", "Arkitektur og logisk inndeling", "Container/Unit og contains; bibliotekstruktur er ikke en deployment-allokering.", []string{"A2", "A3"}},
	{"VP03", "Ansvar og kapabiliteter over arkitekturen", "owns, realizes og provides. Capability er ikke Feature.", []string{"A1", "A2", "A3"}},
	{"VP04", "Grensesnitt og samarbeid", "consumes viser bruk; ingen tilbyder, Channel eller kjørbar meldingsflyt utledes.", []string{"A2", "A3"}},
	{"VP05", "Avhengigheter per modus", "requires in mode; modi har ingen implisitt arv.", []string{"A2", "A3"}},
	{"VP06", "Aktiviteter og leveranseplan", "refines, addresses, delivers og depends-on; planstatus er en eksplisitt kildepåstand.", []string{}},
	{"VP07", "Features over arkitekturen", "contributes-to, owns og eksplisitt allocated-to per modus. Uspesifisert allokering vises som hull.", []string{"A1", "A2", "A3"}},
	{"VP08", "Channel-kontrakter og sekvenser", "Eksplisitte scenario-steg validert mot permits, deltakelse, modus og request/resultat-korrelasjon.", []string{"A1", "A4"}},
	{"VP09", "Dataset, Datagram og persistent Database", "Eksplisitte holdere, kilde, kontrakter, varianter, felt og projeksjoner.", []string{"A1", "A4"}},
	{"VP10", "Datagram-koding og packet", "Kun closed kontrakt med validert Encoding og eksplisitte bitplasseringer.", []string{"A4"}},
	{"VP11", "Egenskaper, sporbarhet og modellhull", "Deklarasjoner og alle fakta med kildeposisjoner; støttegrenser beholdes.", []string{"A0", "A1", "A2", "A3", "A4", "A5"}},
}

func SpecFor(id string) (Spec, bool) {
	for _, s := range Catalog {
		if s.ID == id {
			return s, true
		}
	}
	return Spec{}, false
}
