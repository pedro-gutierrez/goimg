# Goimg

A simple Golang image resizing command line utility, based on Libvips and Bimg.

## Building

Make sure you have Libvips installed. Then simply clone this repo and run:

```
go get
go install
```

## Usage

The following command will download the original file at the specified url, and will produce the following 4 thumbnails:

* /tmp/sierra-original.jpg
* /tmp/sierra-big.jpg
* /tmp/sierra-regular.jpg
* /tmp/sierra-small.jpg

```
goimg url https://cdn-images-1.medium.com/max/2000/1*4BPoUrWcf67bdaKyZPARMg.png /tmp sierra jpg big,regular,small 
```

## Presets 

Goimg comes with the following thumbnail size presets:

```
sizes := map[string]int{
  "large":   1024,
  "big":     800,
  "regular": 480,
  "small":   320,
  "tiny":    128,
}
```

Goimg uses Libvips smartcrop in order to normalize aspect ratios that are bigger than 16:9. For smaller aspect ratios, the original value is kept.


## Related projects

* Libvips: https://github.com/jcupitt/libvips
* Bimg: https://github.com/h2non/bimg
* Imaginary: https://github.com/h2non/imaginary

