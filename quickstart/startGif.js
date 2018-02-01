function getBackgroundImageUrl(element){
  var styles = getComputedStyle(element)
  return styles.backgroundImage.replace("url(\"", "").replace("\")", "")
}
var url = getBackgroundImageUrl(_slide)
var reloadUrl = url + "?" + Math.random()
_slide.style.backgroundImage = "url('" + reloadUrl + "')"