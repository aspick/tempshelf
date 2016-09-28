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
	"errors"
	"path"
	"path/filepath"
	"github.com/aspick/tempshelf/tempshelf"
	"os"

	"github.com/spf13/cobra"
	"github.com/aws/aws-sdk-go/aws"
    "github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/pierrre/archivefile/zip"
)

// pullCmd represents the pull command
var pullCmd = &cobra.Command{
	Use:   "pull [flags] <manifet path>",
	Short: "pull files binary with manifest.json file",
	Long: ``,
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("pull called")

		if len(args) == 0 {
			return errors.New("no input specified")
		}

		manifetPath := args[0]
		fmt.Println("target", manifetPath)

		targetAbspath, e1 := filepath.Abs(manifetPath)
		if e1 != nil {
			return e1
		}
		targetBaseDir := path.Dir(targetAbspath)
		fmt.Println("target", targetBaseDir)

		manifest := tempshelf.ParseManifestFile(manifetPath)

		// create temp dir
		tempDirPath := path.Join(targetBaseDir, ".tmp")
		os.MkdirAll(tempDirPath, 0777)
		defer os.RemoveAll(tempDirPath)

		// download
		s3cli := tempshelf.S3Client(manifest)
		downloader := s3manager.NewDownloaderWithClient(s3cli)

		for _, record := range manifest.Files {
			var localPath = path.Join(targetBaseDir, record.Name)
			if record.Expand == true {
				localPath = path.Join(tempDirPath, record.Name)
			}

			key := manifest.Meta.Prefix + "/" + record.Name

			file, ferr := os.Create(localPath)
			if ferr != nil {
				fmt.Println(ferr)
				continue
			}
			defer file.Close()

			_, s3err := downloader.Download(file, &s3.GetObjectInput{
				Bucket: aws.String(manifest.Meta.Bucket),
				Key:	aws.String(key),
			})
			if s3err != nil {
				fmt.Println(s3err)
			}
		}

		// unarchive
		for _, record := range manifest.Files {
			if record.Expand != true {
				continue
			}

			// var targetBasename = strings.Replace(record.Name, ".zip", "", 1)
			var dstPath = targetBaseDir
			var srcPath = path.Join(tempDirPath, record.Name)

			zipErr := zip.UnarchiveFile(srcPath, dstPath, nil)
			if zipErr != nil {
				fmt.Println("ZipErr", zipErr)
			}
		}

		return nil
	},
}

func init() {
	RootCmd.AddCommand(pullCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// pullCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// pullCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

}
