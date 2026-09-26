package install

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func fixture() Descriptor {
	b := []byte("Toolkit-managed instructions\n")
	h := Hash(b)
	return Descriptor{ReleaseSchema, "dev-gip-1", strings.Repeat("a", 40), Protocol, "sdp-five-phase/0.1", "sdp-project-management/0.2", []string{Protocol}, []File{{"AGENTS.md", "file", "managed", &h, b}}, []string{}, []string{}, []Asset{}}
}
func TestStrictRecords(t *testing.T) {
	d := fixture()
	b, _ := Canonical(d)
	var got Descriptor
	if e := Decode(b, MetadataLimit, &got); e != nil {
		t.Fatal(e)
	}
	if e := ValidateDescriptor(got); e != nil {
		t.Fatal(e)
	}
	for _, bad := range []string{strings.Replace(string(b), `"release":`, `"extra":1,"release":`, 1), strings.Replace(string(b), `"release":`, `"release":"x","release":`, 1), strings.Replace(string(b), `"release":`, `"Release":"x","release":`, 1), strings.Replace(string(b), `"retired":[]`, `"retired":null`, 1), "schemaVersion: &x foo\nrelease: *x\n", "x: !!str foo\n", "x: 1.5\n", "x: 1\n---\nx: 2\n"} {
		if e := Decode([]byte(bad), MetadataLimit, &got); e == nil {
			t.Fatalf("accepted %s", bad)
		}
	}
	if e := Decode(b, 10, &got); e == nil {
		t.Fatal("bound ignored")
	}
	d.Files[0].SHA256 = new(string)
	if e := ValidateDescriptor(d); e == nil {
		t.Fatal("bad hash")
	}
}
func TestRootAndReadOnlyInspection(t *testing.T) {
	root := t.TempDir()
	before, _ := os.ReadDir(root)
	r, e := Root(root)
	if e != nil || r != root {
		t.Fatal(r, e)
	}
	s, e := Inspect(root, []string{"AGENTS.md"})
	if e != nil || s.Snapshot["AGENTS.md"].Type != "absent" {
		t.Fatal(s, e)
	}
	after, _ := os.ReadDir(root)
	if len(before) != len(after) {
		t.Fatal("inspection wrote files")
	}
	os.Mkdir(filepath.Join(root, "SDP"), 0700)
	if r, e = Root(filepath.Join(root, "SDP")); e != nil || r != root {
		t.Fatal(r, e)
	}
	os.Symlink(t.TempDir(), filepath.Join(root, "SDP", "escape"))
	if _, e = Inspect(root, nil); e == nil {
		t.Fatal("symlink accepted")
	}
	for _, p := range []string{"../escape", "a/../b", "/root", "a\\b", "a:b", "a/CON.txt", "a/trailing.", "a//b"} {
		if e = Relative(p); e == nil {
			t.Fatal(p)
		}
	}
}
func TestCanonicalGolden(t *testing.T) {
	b, _ := Canonical(map[string]any{"z": 1, "a": "<ø>"})
	if string(b) != "{\"a\":\"<ø>\",\"z\":1}\n" {
		t.Fatal(string(b))
	}
	p := Plan{SchemaVersion: PlanSchema}
	a := planHash(p)
	p.PlanDigest = "ignored"
	if a != planHash(p) {
		t.Fatal("recursive digest")
	}
}
func TestDescriptorCollisions(t *testing.T) {
	d := fixture()
	f := d.Files[0]
	f.Path = "agents.md"
	d.Files = append(d.Files, f)
	if ValidateDescriptor(d) == nil {
		t.Fatal("case collision")
	}
}

func TestRecordShapes(t *testing.T) {
 d:=fixture();raw,_:=Canonical(d);input:=Input{Path:"/fixture/release.json",SHA256:Hash(raw),Bytes:raw,Provenance:"local-development"}
 a:=Adoption{SchemaVersion:AdoptionSchema,ProjectRoot:"/fixture/project",Baseline:"manual",TargetDigest:input.SHA256,Inventory:map[string]Observation{},Moves:[]Move{},RefreshManaged:[]string{}}
 r:=Receipt{ReceiptSchema,d.Release,input.SHA256,d.SourceCommit,d.ProcessProfile,d.ManagementProfile,d.Capabilities,"local-development","install-0123456789abcdef01234567","2026-09-27T00:00:00Z"}
 p:=Plan{SchemaVersion:PlanSchema,ProjectRoot:a.ProjectRoot,Release:input,Snapshot:map[string]Observation{},Actions:[]Action{},Preserved:[]string{},Conflicts:[]string{},Warnings:[]string{}}
 p.PlanDigest=planHash(p)
 j:=Journal{SchemaVersion:JournalSchema,Plan:p,Steps:[]Action{}}
 cases:=[]struct{name string;value any;target any}{{"release",d,&Descriptor{}},{"adoption",a,&Adoption{}},{"receipt",r,&Receipt{}},{"plan",p,&Plan{}},{"journal",j,&Journal{}}}
 for _,c:=range cases{t.Run(c.name,func(t *testing.T){b,e:=Canonical(c.value);if e!=nil{t.Fatal(e)};if e=Decode(b,RecordLimit,c.target);e!=nil{t.Fatal(e)};gold:=filepath.Join("testdata",c.name+".json");if os.Getenv("SDP_UPDATE_GOLDEN")=="1"{if e=os.WriteFile(gold,b,0600);e!=nil{t.Fatal(e)}};want,e:=os.ReadFile(gold);if e!=nil||string(want)!=string(b){t.Fatalf("golden differs: %v",e)};missing:=strings.Replace(string(b),`"schemaVersion":"`+map[string]string{"release":ReleaseSchema,"adoption":AdoptionSchema,"receipt":ReceiptSchema,"plan":PlanSchema,"journal":JournalSchema}[c.name]+`",`,"",1);if Decode([]byte(missing),RecordLimit,c.target)==nil{t.Fatal("missing schema accepted")}})}
}
