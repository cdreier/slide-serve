// var _slide = document.getElementById("")
_slide.style.backgroundColor = "red"
setTimeout(() => {
  _slide.style.backgroundColor = "blue"
  setTimeout(() => {
    _slide.style.backgroundColor = "inherit"
  }, 1000)
}, 1000)