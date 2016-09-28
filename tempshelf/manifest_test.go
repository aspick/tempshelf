package tempshelf

import (
    "testing"
    "runtime"
    "path"

    "io/ioutil"
    "os"
    "reflect"
    "fmt"
)

func testManifetFilePath() string {
    _, filename, _, _ := runtime.Caller(0)
    return path.Join(path.Dir(path.Dir(filename)), "test", "manifest.json")
}

func assertManifest(m Manifest ,t *testing.T){
    if m.Meta.Storage != "s3" {
        t.Error("meta.storage should eq 's3'")
    }

    if len(m.Files) != 3 {
        t.Error("files count shoud eq 3")
    }

    if m.Files[0].Name != "file1" {
        t.Error("files[0].name should eq file1")
    }
}

func TestParseManifestFile(t *testing.T) {
    result := ParseManifestFile(testManifetFilePath())

    assertManifest(result, t)
}

func TestLoad(t *testing.T) {
    var m Manifest
    m.Load(testManifetFilePath())
    assertManifest(m, t)
}

func TestSave(t *testing.T) {
    var meta ManifestMeta
    meta.Storage = "s3"
    meta.Bucket = "test-bucket"
    meta.Region = "test-region"
    meta.Token = "test-token"
    meta.Secret = "test-secret"
    meta.Prefix = "app1"

    var file1 FileRecord
    file1.Name = "file1"
    file1.Expand = false

    var file2 FileRecord
    file2.Name = "file2"
    file2.Expand = false

    var file3 FileRecord
    file3.Name = "file3"
    file3.Expand = true

    var m Manifest
    m.Meta = meta
    m.Files = []FileRecord{}
    m.Files = append(m.Files, file1, file2, file3)

    tmpfile, err := ioutil.TempFile("", "example")
	if err != nil {
		panic(err)
	}
    defer os.Remove(tmpfile.Name())

    m.Save(tmpfile.Name())

    manifest2 := ParseManifestFile(tmpfile.Name())

    bytes, _ := ioutil.ReadFile(tmpfile.Name())
    fmt.Println(string(bytes))

    if !reflect.DeepEqual(m, manifest2) {
        fmt.Println(m)
        fmt.Println(manifest2)
        t.Error("manfest and manifest2 should equal")
    }

}
