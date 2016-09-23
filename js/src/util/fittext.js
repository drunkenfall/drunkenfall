/* global getComputedStyle */
/*!
* FitText.js 1.0 jQuery free version
*
* Copyright 2011, Dave Rupert http://daverupert.com
* Released under the WTFPL license
* http://sam.zoy.org/wtfpl/
* Modified by Slawomir Kolodziej http://slawekk.info
*
* Date: Tue Aug 09 2011 10:45:54 GMT+0200 (CEST)
*
* Hacked by FrontierPsycho for this project.
*/
/*
let css = function (el, prop) {
  return window.getComputedStyle ? getComputedStyle(el).getPropertyValue(prop) : el.currentStyle[prop]
}
*/

import _ from "lodash"

function addEvent (el, type, fn) {
  if (el.addEventListener) {
    el.addEventListener(type, fn, false)
  } else {
    el.attachEvent('on' + type, fn)
  }
}

function getTextWidth (text, font) {
  // if given, use cached canvas for better performance
  // else, create new canvas
  let canvas = getTextWidth.canvas || (getTextWidth.canvas = document.createElement("canvas"))
  let context = canvas.getContext("2d")
  context.font = font
  let metrics = context.measureText(text)
  return metrics.width
};

export default function fitText (el, text, font, desiredRatio) {
  let fit = function (el) {
    let resizer = function () {
      let currentFontSize = _.trimEnd(getComputedStyle(el, null).getPropertyValue("font-size"), "px")
      let fontString = `${currentFontSize}px ${font}`

      let currentRatio = getTextWidth(text, fontString) / parseFloat(el.clientWidth)
      let newFontSize = currentFontSize * (parseFloat(desiredRatio) / currentRatio)

      console.debug("currentRatio:", currentRatio, "currentFontSize:", currentFontSize, "newFontSize:", newFontSize)

      el.style.fontSize = newFontSize + "px"
    }

    // Call once to set.
    resizer()

    // Bind events
    // If you have any js library which support Events, replace this part
    // and remove addEvent function (or use original jQuery version)
    addEvent(window, 'resize', resizer)
  }

  if (el.length) {
    for (var i = 0; i < el.length; i++) {
      fit(el[i])
    }
  } else {
    fit(el)
  }

  // return set of elements
  return el
}
