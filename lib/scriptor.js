/**!
 * @license Scriptor.js - A library for building your own custom text editors
 * LICENSED UNDER MIT LICENSE
 * MORE INFO CAN BE FOUND AT https://github.com/MarketingPipeline/Scriptor.js/
 */

const defaultButtonProps = {
  insert: false,
  htmltags: true,
  value: '',
  wrap: false
};

const form = document.getElementById('text-editor');

if (form != null){
let DEBUG = false;  
// carot / last type postion
let currentTextPosition = 0;

form.addEventListener('click', () => (currentTextPosition = form.selectionEnd), false);

form.addEventListener('input', () => {
  currentTextPosition = form.selectionEnd;

  updatePreview();

  
});

form.addEventListener('pointerenter', () => {
  currentTextPosition = form.selectionEnd;

  updatePreview();
});

/// Load any default text area content
window.addEventListener('load', function (e) {
  // This prevents the window from reloading
  let input = form.value;
});

/// Get all Text Editor Button Values

const buttons = document.querySelectorAll('[data-scriptor-btn]');

buttons.forEach((button) => button.addEventListener('click', (e) => handleClick(button, form)));

async function updatePreview() {
  let url = '/convert';
  let response = await fetch(url, { method: 'POST', body: form.value });
  let html = await response.text(); // get updated MD

  // refresh 'preview-area'
  document.getElementById('preview-area').innerHTML = html;
}

document.getElementById('text-editor').addEventListener('keydown', function(event) {
  if (event.key === 'Tab') {
    event.preventDefault(); // prevent change focus on Tab
    
    const start = this.selectionStart;
    const end = this.selectionEnd;

    // insert \t
    this.value = this.value.substring(0, start) + '\t' + this.value.substring(end);

    // move carot
    this.selectionStart = this.selectionEnd = start + 1;
    updatePreview();
  }
});

function handleClick(button, form) {
  var values = getNewValue(button, form.value);
  form.value = values[0];
  form.selectionEnd=values[1];
  form.focus();
}

function getNewValue(button, text) {
    // allows custom functions to be called on button clicks.
  if (button.getAttribute("custom-function")) eval(button.getAttribute("custom-function"))

  // for each value - check if type is true or false. 
  const [insert, htmltags, wrap] = ['insert', 'htmltags', 'wrap'].map((key) => checkBool(button.getAttribute(key) ?? defaultButtonProps[key]));
  const value = button.getAttribute('value') ?? defaultButtonProps['value'];
  DEBUG && console.table({ insert, value, htmltags, wrap });

  const selectedText = getSelectionText();
  const startLineIndex = text.lastIndexOf('\n', currentTextPosition) + 1;
  const endLineIndex = text.indexOf('\n', currentTextPosition);
  const lineEnd = endLineIndex !== -1 ? endLineIndex : text.length;

  // Insert Value
  if (insert) {
    if (value === '***') return [text.substring(0, currentTextPosition) + '\n\n' + value + '\n\n' + text.substring(currentTextPosition, text.length), currentTextPosition + value.length * 2 + 1];
    if (value === 'img') {
      insertValue = `![Alt text](PUT_IMAGE_LINK_HERE "Image title")`
      return [text.substring(0, currentTextPosition) + '\n\n' + insertValue + '\n\n' + text.substring(currentTextPosition, text.length), currentTextPosition + 14];
    }
    return [text.substring(0, currentTextPosition) + value + text.substring(currentTextPosition, text.length), currentTextPosition + value.length];
  }

  /// Highlighted Text Options
  if (selectedText === '') {
    // no text was hightlighted - just add the values
    if (!wrap) {
      // insert value at the beginning of the string
      const newText = text.substring(0, startLineIndex) + value + text.substring(startLineIndex, lineEnd) + text.substring(lineEnd);
      
      return [newText, startLineIndex + value.length];
    }
    if (value === '~~~') {
      return [text.substring(0, currentTextPosition) + '\n' + value + '\n\n' + value + '\n' + text.substring(currentTextPosition), currentTextPosition + value.length + 2];
    }
    if (value === 'link') {
      insertValue = `[](PUT_LINK_URL_HERE)`
      return [text.substring(0, currentTextPosition) + insertValue + text.substring(currentTextPosition, text.length), currentTextPosition + 1];
    }
    // wrap value around carot
    return [text.substring(0, currentTextPosition) + value + value + text.substring(currentTextPosition, text.length), currentTextPosition + value.length];

  }

  if (selectedText !== '') {
    if (wrap) {
      newString = text.substring(startLineIndex, endLineIndex);
      isSpace = "";
      if (selectedText[selectedText.length - 1] === ' ') {
        isSpace = " ";
      }
      // wrap or unwrap selected text
      if (value === '~~~') return [text.replace(selectedText, wrapText('\n'+selectedText+'\n', value)), currentTextPosition + (value.length * 2 + 3)];
      if (value === 'link') return [text.substring(0, currentTextPosition - selectedText.length) + newString.replace(newString, `[` + selectedText.trim() + `]()` + isSpace) + text.substring(currentTextPosition, text.length), currentTextPosition + 3 - isSpace.length];
      return [text.substring(0, currentTextPosition - selectedText.length) + newString.replace(newString, wrapText(selectedText.trim(), value)) + isSpace + text.substring(currentTextPosition, text.length), currentTextPosition + value.length * 2];
    } else {
      // delete value in the beginning of selection
      if (selectedText.startsWith(value)) return [text.replace(selectedText, selectedText.replaceAll(value, '')), currentTextPosition - value.length];

      // Add to the start of the value
      return [form.value.replace(selectedText, value + selectedText.replaceAll('\n', '\n'+value)), currentTextPosition + value.length];
    }
  }
}

/// This will return the highlighted text on screen.

function getSelectionText() {
  if (window.getSelection) {
    return window.getSelection().toString();
  }
  if (document.selection && document.selection.type != 'Control') {
    return document.selection.createRange().text;
  }
  return '';
}

function escapeRegExp(string) {
  return string.replace(/[.*+?^${}()|[\]\\]/g, '\\$&'); // escape special charatcters like *
}

// Wrap Highlighted Text On Button Click
function wrapText(text, wrap) {
  const string = text.trim();
  DEBUG && console.log(wrap);

  const escapedWrap = escapeRegExp(wrap);
  const wrapPattern = new RegExp(`^${escapedWrap}|${escapedWrap}$`, 'g');
  // if Highlighted Text String Already Contains A Wrap At Start & End - Remove It
  if (wrapPattern.test(string)) {
    return string.replace(wrapPattern, '');
  }
  return `${wrap}${text}${wrap}`;
}

 function AttributeToLowerCase(text){
   text = text.toString()
   var x = text.toLowerCase()
   DEBUG && console.log(`AttributeToLowerCase Was Called`)
   return x
 }
  
function checkBool(x) {
  return AttributeToLowerCase(x) === 'true' || x === true;
}
}
