package main

import (
	"fmt"
	"github.com/spf13/cobra"
	"gopkg.in/h2non/bimg.v1"
	"io/ioutil"
	"net/http"
	"path/filepath"
	"strings"
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

type ResizeOpts struct {
	Height int
	Force  bool
}

func GetResizeOpts(width int, w int, h int) ResizeOpts {
	if h > w {
		return GetResizeOpts(width, w, w)
	} else {
		if width > 128 {
			actual := (float64(w) / float64(h) * 100) / 100
			if actual > 1.77 {
				return ResizeOpts{
					Height: int(float64(width) / float64(1.77)),
					Force:  false,
				}
			} else {
				return ResizeOpts{
					Height: h * width / w,
					Force:  false,
				}
			}
		} else {
			return ResizeOpts{
				Height: width,
				Force:  true,
			}
		}
	}
}

func WithImageOptions(s string, m bimg.ImageMetadata, next func(o bimg.Options)) {

	sizes := map[string]int{
		"large":   1024,
		"big":     800,
		"regular": 480,
		"small":   320,
		"tiny":    128,
		"avatar":  64,
	}

	width, ok := sizes[s]
	if ok {

		resizeOpts := GetResizeOpts(width, m.Size.Width, m.Size.Height)

		//fmt.Println(strings.Replace(fmt.Sprintf("%#v", resizeOpts), ", ", "\n", -1))

		opts := bimg.Options{
			Width:          width,
			Height:         resizeOpts.Height,
			Force:          resizeOpts.Force,
			Crop:           true,
			SmartCrop:      true,
			Embed:          true,
			Type:           bimg.JPEG,
			Gravity:        bimg.GravitySmart,
			Interpretation: bimg.InterpretationSRGB,
		}

		//fmt.Println(strings.Replace(fmt.Sprintf("%#v", opts), ", ", "\n", -1))

		next(opts)
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
