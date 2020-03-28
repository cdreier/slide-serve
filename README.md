# slide-serve

![Drone Build Status](https://github-drone.d.drailing.net/api/badges/cdreier/slide-serve/status.svg)

![demogif](https://github.com/cdreier/slide-serve/blob/master/demogif/slide-serve-demo.gif)

i really like the idea behind [slide-html](https://github.com/trikita/slide-html) and the [Takahashi method](https://en.wikipedia.org/wiki/Takahashi_method)

you only need to learn 5 pseudo-markdown rules and a few annotations to build a nice presentation.

for really, really fast slide-creation and iteration and a good *developer*-experience i created this server

and for a good talk, you can control the presentation from the presenter view, including the next slide and your speaker notes

![presenter-view](https://github.com/cdreier/slide-serve/blob/master/demogif/presenter-view.png)

please take a look at the [extended example presentation](http://htmlpreview.github.io/?https://github.com/cdreier/slide-serve/blob/master/example_html/Slide.html) from slide-html

## quickstart

i added a [quickstart package](https://github.com/cdreier/slide-serve/tree/master/quickstart) in the repository, with my own styles and short demos from all features

just download the latest binary for your platform from the [release page](https://github.com/cdreier/slide-serve/releases/latest), place it somewhere and start

```
./slide-serve -dir path/to/quickstart
```
and visit [http://127.0.0.1:8080](http://127.0.0.1:8080)

## features

* just write plaintext files, no need to think about html or indents
  * you can even write more than one plaintext file, they are added in alphanumerical order
* start with the `-dev` flag to auto reload the browser on that slide you just edited
  * this is just awesome!
* just go to http://127.0.0.1:8080/presenter and controll the presentation from the presenter view!
  * see the current and the next slide
  * see the speaker-@notes
* add styles and background images next to your content
  * no hardcoded slide-numbers in css
  * don't count your slides, i'll do it for you ;)
* syntax highlighter with [chroma](https://github.com/alecthomas/chroma)
* use local images
* add small javascripts for the last bits of awesomeness
  * take a look at [useful snippets](#useful-snippets) for an example

## your presentation

* your slide-files must end with `.md` to gets loaded
* `styles.css` in your presentation-root dir is automaticly added to your presentation
* your presentation directory is routed to `http://localhost:8080/static/` so you can access local images or fonts - this is important for your custom stylesheets!

## syntax additions

`@img/imagename.gif`  
`@css/stylesheet.css`  
`@js/jsfile.js`  
`@code/language`  
`@append`  
`@hidden`  
`@note`  

the path is relative to your presentation directory, if you like to create an image folder it would change to `@img/images/imagename.gif`

with die `@code` annotation, the code-formatting is set for this slide 

note: instead of writing the harcoded slide number, you should use the SLIDENUMBER placeholder

``` css
@code/css
.slide-SLIDENUMBER h1 { 
  color: #fff; 
  text-shadow: 1px 1px 3px #333; 
}
```

example slide with background image and stylesheet:

``` md
# AWESOME
.
@img/example_bg.jpg
@css/whiteHeadline.css
```

with the `@js` annotation, you can add javascript to your slide - the code will run in the moment you enter the slide.

``` md
# javascript!
@js/jsdemo.js
```

2 variables are available: `_slide` is your current slide-div, and `_slideIndex` - is your slide index... obviously

the jsdemo looks like

```js
_slide.style.backgroundColor = "red"
setTimeout(() => {
  _slide.style.backgroundColor = "inherit"
}, 1000)
```

`@append` appends the current slide to the previous, so you don't have to repeat the whole content.

following example still produces 2 slides, but the second with `#this is` prepended

``` md
# this is

# awesome
@append
```

example presentation: https://github.com/cdreier/slide-serve/tree/master/example

> to see the example presentation, just run slide-serve without any flags

## syntax highlighter

i used [chroma](https://github.com/alecthomas/chroma) for syntax highlightning. 

just tell me what language you are using on your slide (with `@code/lang`) and you are good to go.

with the `-syntaxhl` start flag you can set the highlighter theme (there is a [list](https://github.com/alecthomas/chroma/tree/master/styles) with all the themes). the default is monokai

example (the position of the `@code` annotation does not matter)

``` md
# CSS backgrounds
@code/css
  .slide-12 {
    background: url("icon.png"),
                url("bg.png");
    background-repeat: no-repeat, repeat;
    background-size: 40%, auto;
    background-position: bottom right;
  }
```

# useful snippets

## list items popping up

first, lets make the paragraphs looks like list-items

```css 
.slide-SLIDENUMBER .slide-content p {
  list-style-type: disc; 
  display: list-item;
  text-align: left;
}
```

and then use append for every item

```md
# lists
@css/list.css

new list item 
@append

next item
@append

item 3
@append
```

## restarting gifs

as gifs are looping in the background, this script reloads the gif every time you enter the slide

```js
// restartGif.js
function getBackgroundImageUrl(element){
  var styles = getComputedStyle(element)
  return styles.backgroundImage.replace("url(\"", "").replace("\")", "")
}
var url = getBackgroundImageUrl(_slide)
var reloadUrl = url + "?" + Math.random()
_slide.style.backgroundImage = "url('" + reloadUrl + "')"
```

and in an example slide

```
.
@img/funnyCat.gif
@js/restartGif.js
```

## no image scaling

in the past, i needed a down-scaling of the images - if you want to reset this behaviour, add this to your `styles.css`

```css
.slide {
  background-size: unset !important;
}
```

## add a logo on every slide

i added an empty div at the end of the presentation HTML, so we can just use it in our global `styles.css`

```css
#logo {
  position: fixed;
  right: 10px;
  bottom: 10px;
  font-family: Vollkorn;
  background-color: #e2e4e6;
  background-image: url("/static/images/your_logo.png");
  background-repeat: no-repeat;
}
```

# printing

the default print styles from [slide-html](https://github.com/trikita/slide-html) are included per default, so you can print a summary of your presentation.

with the `-pdf` flag, the styles are changed for printing the full slides, or to export to a pdf.

if something wont fit with your presentation styles, you can overwrite the print styles:

```css
@media print {
  @page {size: landscape}

  /* ... */

  .slide h1 {
		font-size: 8rem !important;
	}
	.slide p {
		font-size: 2rem !important;
	}

}
```

# cli flags

-dev  
* start with `-dev` or `-dev true` to enable auto-reloading (default false)
* still much awesome feature! 

-dir  
* this is the directory of your presentation (default "example")

-port 
* http port the server is starting on (default "8080")

-syntaxhl
* string with the highlighter theme for your code (default "monokai")

-ratio
* slide ratio, possible values: 4x3, 16x9 or 16x10 (default "16x9")

-title 
* html title in the browser (default "Slide")

-pdf 
* overwrite the summary print styles with full slide prints (default false)

-click 
* on default you only navigate with arrow keys, this enabled 'next slide' on click

# build on your own

i used [pkger](https://github.com/markbates/pkger) to embedd the assets in the binary, you should check it out!

# deployment

You can deploy on [Zeit Now](https://zeit.co/now) for example by doing the following:

- `npm install -g now`
- Run your presentation locally without development mode
- Create an `html` directory: `mkdir html`
- Inside of the directory do a curl: `curl localhost:8080 -o index.html`
- Create a `now.json` file with something like this:
  ```
  {
    "version": 2,
    "name": "name-of-the-presentation",
    "alias": [
      "name-of-the-presentation"
    ]
  } ```
- If you have pictures / fonts, create a symlink a folder called `static`. f.e. `ln -s <place-with-images> static/images`
- Run `now`
- Now will deploy your presentation.
- To make it available on https://name-of-the-presentation.now.sh you have to run: `now alias`
