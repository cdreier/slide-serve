# slide-serve

i really like the idea behind [slide-html](https://github.com/trikita/slide-html) and the [Takahashi method](https://en.wikipedia.org/wiki/Takahashi_method)

for really, really fast slide-development i created this small server

## features

* just write plaintext files, no need to think about html or indents
  * you can even write more than one plaintext file!
* auto reload the browser on that slide you just edited
  * this is just awesome!
* add styles and background images next to your content
  * no hardcoded slide-numbers in css
  * don't count your slides, i'll do it for you ;)
* use local images

## usage

quickstart with a pre-build binary for linux, windows or mac: on the [release page](https://github.com/cdreier/slide-serve/releases)

put the file in a place you really like - and start it!  
the `-dir` flag should point to the folder with your presentation

``` bash
./slide-serve -dir how-to-overengineer-3000
```

### your presentation

* your slide-files should end with `.md` 
* `styles.css` is automaticly added to your presentation
* your presentation directory is routed to `http://localhost:8080/static/` so you can access local images or fonts - this is important for your custom stylesheets!

### syntax additions

`@img/imagename.gif` and `@css/stylesheet.css`

the path is relative to your presentation directory, if you like to create an image folder it would change to `@img/images/imagename.gif`

note: instead of writing the harcoded slide number, you should use the SLIDENUMBER placeholder

``` css
.slide-SLIDENUMBER h1 { 
  color: #fff; 
  text-shadow: 1px 1px 3px #333; 
}
```

example slide with background image and stylesheet:

```
# AWESOME
.
@img/example_bg.jpg
@css/whiteHeadline.css
```

example presentation: https://github.com/cdreier/slide-serve/tree/master/example

## cli flags

-dev  
* start with `-dev` or `-dev true` to enable auto-reloading (default false)
* still much awesome feature! 

-dir  
* this is the directory of your presentation (default "example")

-port 
* http port the server is starting on (default "8080")

-title 
* html title in the browser (default "Slide")