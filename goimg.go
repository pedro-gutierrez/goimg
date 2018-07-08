package main

import (
	"fmt"
	"github.com/spf13/cobra"
	"gopkg.in/h2non/bimg.v1"
	"io/ioutil"
	"net/http"
	"path/filepath"
	"strings"
	//	"os"
)

func main() {

	var urlCmd = &cobra.Command{
		Use:   "url [url] [dir] [basename] [type] [formats]",
		Short: "Convert an image",
		Long:  "Convert an image",
		Args:  cobra.MinimumNArgs(5),
		Run: func(cmd *cobra.Command, args []string) {
			WithBytesFromUrl(args[0], func(data []byte) {
				WithImageName("original", args[2], args[3], func(n1 string) {
					WriteImage(args[1], n1, data, func() {
						WithImageFromBytes(data, func(i *bimg.Image, m bimg.ImageMetadata) {
							WithEachSize(args[4], func(s string) {
								WithImageOptions(s, m, func(opts bimg.Options) {
									WithNewImage(i, opts, func(i2 []byte) {
										WithImageName(s, args[2], args[3], func(n string) {
											WriteImage(args[1], n, i2, func() {})
										})
									})
								})
							})
						})
					})
				})
			})
		},
	}

	var rootCmd = &cobra.Command{Use: "goimg"}
	rootCmd.AddCommand(urlCmd)
	rootCmd.Execute()
}

func WriteImage(dir string, name string, data []byte, next func()) {
	filename := filepath.Join(dir, name)
	bimg.Write(filename, data)
	fmt.Println(fmt.Sprintf("Written %s", filename))
	next()
}

func WithImageName(size string, basename string, ext string, next func(n string)) {
	next(fmt.Sprintf("%s-%s.%s", basename, size, ext))
}

func WithNewImage(i *bimg.Image, o bimg.Options, next func(i2 []byte)) {
	i2, err := i.Process(o)
	if err != nil {
		fmt.Println(fmt.Sprintf("Unable to convert image"))
		fmt.Println(o)
		fmt.Println(err)
	} else {
		next(i2)
	}
}

func WithEachSize(sizes string, next func(s string)) {
	for _, s := range strings.Split(sizes, ",") {
		next(s)
	}
}

func WithImageOptions(s string, m bimg.ImageMetadata, next func(o bimg.Options)) {

	sizes := map[string]int{
		"large":   1024,
		"big":     800,
		"regular": 480,
		"small":   320,
		"tiny":    128,
	}

	width, ok := sizes[s]
	if ok {

		next(bimg.Options{
			Width:          width,
			Height:         m.Size.Height * width / m.Size.Width,
			Crop:           true,
			Embed:          true,
			Type:           bimg.JPEG,
			Interpretation: bimg.InterpretationSRGB,
		})
	} else {
		fmt.Println(fmt.Sprintf("Size not supported: %s", s))
	}
}

func WithImageFromBytes(data []byte, next func(i *bimg.Image, m bimg.ImageMetadata)) {
	i := bimg.NewImage(data)
	m, err := i.Metadata()
	if err != nil {
		fmt.Println("Error getting image metadata")
		return
	}
	next(i, m)
}

func WithBytesFromUrl(url string, next func(d []byte)) {
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println(fmt.Sprintf("Error getting %s", url))
		return
	}

	defer resp.Body.Close()
	body, err2 := ioutil.ReadAll(resp.Body)
	if err2 != nil {
		fmt.Println("Error converting original file to a byte array")
		return
	}
	next(body)
}
