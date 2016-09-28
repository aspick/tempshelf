package tempshelf

import (
    "os"
    "encoding/json"
    "io/ioutil"
)

// Manifest represents manifest file
type Manifest struct {
    Meta    ManifestMeta `json:"meta"`
    Files   []FileRecord `json:"files"`
}


// ManifestMeta contain storage config
type ManifestMeta struct {
    Storage string  `json:"storage"`
    Bucket  string  `json:"bucket"`
    Region  string  `json:"region"`
    Token   string  `json:"token"`
    Secret  string  `json:"secret"`
    Prefix  string  `json:"prefix"`
}

// FileRecord represenets each file entity
type FileRecord struct {
    Name    string  `json:"name"`
    Expand  bool    `json:"expand"`
}

// ParseManifestFile create Manifest from manifest.json path
func ParseManifestFile(filepath string) Manifest {
    var manifest Manifest
    manifest.Load(filepath)

    return manifest
}

// Load assign data from manifest.json path
func (m *Manifest)Load(filepath string) {
    file, err := os.Open(filepath)
    if err != nil {
        panic(err)
    }
    defer file.Close()

    dec := json.NewDecoder(file)
    dec.Decode(m)
}

// Save write data to manifet.json path
func (m *Manifest)Save(filepath string) {
    manifestString, _ := json.MarshalIndent(m, "", "  ")
    err := ioutil.WriteFile(filepath, manifestString, 0644)
    if err != nil {
        panic(err)
    }
}
