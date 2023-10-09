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
let startPosition = 0;
let currentTextPosition = 0;

form.addEventListener('click', () => (currentTextPosition = form.selectionEnd), false);

form.addEventListener('input', () => {
  currentTextPosition = form.selectionEnd;
});

/// Load any default text area content
window.addEventListener('load', function (e) {
  // This prevents the window from reloading

  let input = form.value;
});

/// Get all Text Editor Button Values

const buttons = document.querySelectorAll('[data-scriptor-btn]');

buttons.forEach((button) => button.addEventListener('click', (e) => handleClick(button, form)));

function handleClick(button, form) {
  var values = getNewValue(button, form.value);
  form.value = values[0];
  // form.focus();
  form.selectionEnd=values[1];
}

function getNewValue(button, text) {
    // allows custom functions to be called on button clicks.
  if (button.getAttribute("custom-function")) eval(button.getAttribute("custom-function"))
  
  // for each value - check if type is true or false. 
  const [insert, htmltags, wrap] = ['insert', 'htmltags', 'wrap'].map((key) => checkBool(button.getAttribute(key) ?? defaultButtonProps[key]));
  const value = button.getAttribute('value') ?? defaultButtonProps['value'];
  DEBUG && console.table({ insert, value, htmltags, wrap });

  // Insert Value
  if (insert) {
    if (value === '***') return [text.substring(0, currentTextPosition) + '\n\n' + value + '\n' + text.substring(currentTextPosition, text.length), currentTextPosition+value.length];
    return [text.substring(0, currentTextPosition) + value + text.substring(currentTextPosition, text.length), currentTextPosition+value.length];
  }

  /// Highlighted Text Options
  if (getSelectionText() === '') {
    // no text was hightlighted - just add the values
    // todo - set carot in between the value added
    if (!wrap) return [text.substring(0, currentTextPosition) + value + text.substring(currentTextPosition, text.length), currentTextPosition+value.length];
    return [text.substring(0, currentTextPosition) + value + value + text.substring(currentTextPosition, text.length), currentTextPosition+value.length];
  }

  if (getSelectionText() != '') {
    if (wrap) {
      // Not wrapping with html tags <>
      if (value === '~~~') return [text.replace(getSelectionText(), wrapText('\n'+getSelectionText()+'\n', value, false)), currentTextPosition+(value.length*2)];
      return [text.replace(getSelectionText(), wrapText(getSelectionText(), value, false)), currentTextPosition+(value.length*2)];
    } else {
      // replace first HTML tag
      if (getSelectionText().startsWith(value)) return [text.replace(getSelectionText(), getSelectionText().replaceAll(value, '')), currentTextPosition-value.length];

      // Add to the start of the value
      return [form.value.replace(getSelectionText(), value + getSelectionText().replaceAll('\n', '\n'+value)), currentTextPosition+value.length];
    }
  }
}

/// This will return the highlighted text on screen.

function getSelectionText() {
  if (window.getSelection) return window.getSelection().toString();
  if (document.selection && document.selection.type != 'Control') return document.selection.createRange().text;
  return '';
}

// Wrap Highlighted Text On Button Click
function wrapText(text, wrap) {
  const string = text.trim();
  DEBUG && console.log(wrap);
  // if Highlighted Text String Already Contains A Wrap At Start & End - Remove It
  if (string.startsWith(`${wrap}`) == true) {
    return string.replace(`${wrap}`,'').replace(`${wrap}`,'');
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
