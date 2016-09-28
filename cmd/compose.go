// Copyright © 2016 NAME HERE <EMAIL ADDRESS>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	// "github.com/aspick/tempshelf/ziptool"
	"github.com/pierrre/archivefile/zip"
	"github.com/aspick/tempshelf/tempshelf"
	"errors"
	"path"
	"path/filepath"
	"io/ioutil"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
    "github.com/aws/aws-sdk-go/service/s3"
)

// composeCmd represents the compose command
var composeCmd = &cobra.Command{
	Use:   "compose [flags] <manifest path>",
	Short: "compose files and upload to storage",
	Long: `compose files and upload to storage`,
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("compose called")
		if len(args) == 0 {
			return errors.New("no input specified")
		}

		manifetPath := args[0]
		fmt.Println("target", manifetPath)

		// listup files
		targetAbspath, e1 := filepath.Abs(manifetPath)
		if e1 != nil {
			return e1
		}
		targetBaseDir := path.Dir(targetAbspath)
		fmt.Println("target", targetBaseDir)

		files, e2 := ioutil.ReadDir(targetBaseDir)
		if e2 != nil {
			return e2
		}

		manifest := tempshelf.ParseManifestFile(manifetPath)

		records := []tempshelf.FileRecord{}

		// update manifest
		for _, file := range files {
			if file.Name() == "manifest.json" {
				continue
			}
			if file.Name() == ".tmp" && file.IsDir() == true {
				continue
			}

			var fileRecord tempshelf.FileRecord
			if file.IsDir() {
				fileRecord.Expand = true
				fileRecord.Name = file.Name() + ".zip"
			}else{
				fileRecord.Expand = false
				fileRecord.Name = file.Name()
			}

			records = append(records, fileRecord)

			fmt.Println("file:", file.Name(), "dir:", file.IsDir())
		}

		manifest.Files = records
		manifest.Save(manifetPath)

		// archie dir to zip
		// archived := []string{}
		tempDirPath := path.Join(targetBaseDir, ".tmp")
		os.MkdirAll(tempDirPath, 0777)

		defer os.RemoveAll(tempDirPath)
		
		for _, record := range records {
			if record.Expand == false {
				continue
			}

			var targetBasename = strings.Replace(record.Name, ".zip", "", 1)
			var targetPath = path.Join(targetBaseDir, targetBasename)
			var destPath = path.Join(tempDirPath, record.Name)

			// fmt.Println("zipping", destPath, targetPath)

			ziperr := zip.ArchiveFile(targetPath, destPath, nil)
			if ziperr != nil {
				fmt.Println("ziperr", ziperr)
			}
		}

		// upload files
	    s3cli := tempshelf.S3Client(manifest)

		for _, record := range records {
			var tpath = path.Join(targetBaseDir, record.Name)
			if record.Expand == true {
				tpath = path.Join(tempDirPath, record.Name)
			}
			key := manifest.Meta.Prefix + "/" + record.Name

			// fmt.Println("src", tpath, "dest", key)

			targetFile, openErr := os.Open(tpath)
			if openErr != nil {
				fmt.Println(openErr)
				continue
			}

			defer targetFile.Close()

			_, s3err := s3cli.PutObject(&s3.PutObjectInput{
		        Bucket: aws.String(manifest.Meta.Bucket),
		        Key:    aws.String(key),
		        Body:   targetFile,
		    })

			if s3err != nil {
				fmt.Println(s3err)
			}
		}

		// clean .tmp dir


		return nil
	},
}

func init() {
	RootCmd.AddCommand(composeCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// composeCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// composeCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

}
