document.querySelectorAll('pre').forEach(pre => {
    const copyButton = pre.querySelector('.copy-code');
    
    copyButton.addEventListener('click', (event) => {
        event.preventDefault(); // prevent default link behavior

        const code = pre.querySelector('code').innerText; // get text from the code block

        // copy text to clipboard
        navigator.clipboard.writeText(code).then(() => {
            console.log('Код скопирован в буфер обмена!'); // log success
        }).catch(err => {
            console.error('Ошибка при копировании: ', err); // log error
        });
    });
});
