const baseUrl = window.location.origin;


const hideElement = (elementId) => {
    const element = document.getElementById(elementId);
    if (element === null) {
        console.error("element not found for id: " + elementId);
        return;
    }
    if (!element.classList.contains("is-hidden")) {
        element.classList.add("is-hidden");
    }
};

const showElement = (elementId) => {
    const element = document.getElementById(elementId);
    if (element === null) {
        console.error("element not found for id: " + elementId);
        return;
    }
    if (element.classList.contains("is-hidden")) {
        element.classList.remove("is-hidden");
    }
};

const clearInnerElementsOf = (elementId) => {
    const element = document.getElementById(elementId);
    if (element === null) {
        console.error("element not found for id: " + elementId);
        return;
    }
    while (element.firstChild) {
        element.removeChild(element.firstChild);
    }
};

function doLogin() {

    //TODO(mk): impl login for session and on server

    showElement("user_container");
    hideElement("infobar_container");
    hideElement("login_container");
    showElement("qr_container");

    //TODO(mk): read and display qr-codes or "No business decision found"

}

function doLogout() {

    //TODO(mk): impl logout for session and on server

    hideElement("user_container");
    clearInnerElementsOf("infobar_container");
    hideElement("infobar_container");
    showElement("login_container");
    clearInnerElementsOf("qr_container");
    hideElement("qr_container");
}


//TODO(mk): add listner to handle if logged in or logged out on load/reload