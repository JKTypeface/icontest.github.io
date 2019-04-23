
function inputUrl() {
    
    var input = document.getElementById("input");

    return "/icon/" + input.value.replace("https","").replace("http","").replace(":","").replace(/^\/+/g, '');
}

function init() {

    var input = document.getElementById("input");
    var link = document.getElementById("input-link");
    
    input.addEventListener("keyup", (evt) =>
    {
        if (evt.keyCode === 13) {
    
            openUrl(inputUrl(), evt.shiftKey);
        }
    });

    link.addEventListener("click", (evt) =>
    {
        evt.preventDefault();

        openUrl(inputUrl(), evt.shiftKey);        
    });
}

function openUrl(url, shift) {
    if (!shift)
    {
        // Enter key opens link to unmasked image in the same tab
        window.location.href = url;
    } 
    else
    {
        // Shift+Enter opens the link in a new tab
        window.open(url, "_blank").focus();
    }
}

window.onload = init;