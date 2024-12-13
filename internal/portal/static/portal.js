const baseUrl = window.location.origin;

const validateHasValueAndNotEmpty = (element) => {
    return element?.value?.trim().length > 0;
};

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

const isEmptyObject = (obj) => {
    return obj && Object.keys(obj).length === 0 && obj.constructor === Object;
};


const displayError = (errorText) => {
    const pError = document.createElement("p");
    pError.innerText = errorText;
    pError.classList.add("has-text-danger");
    document.getElementById("error_container").appendChild(pError);
}

const doLogin = () => {

    const usernameElement = document.getElementById("username");
    const passwordElement = document.getElementById("password");

    //TODO(mk): impl real auth (now a simple fake simulation of an auth)
    if (!validateHasValueAndNotEmpty(usernameElement) || !validateHasValueAndNotEmpty(passwordElement)) {
        displayError("Login failed - please try again");
        return;
    }
    const username = usernameElement.value;
    const password = passwordElement.value;
    if (username.toLowerCase() !== password) {
        displayError("Login failed - please try again");
        return;
    }

    usernameElement.value = "";
    passwordElement.value = "";
    clearInnerElementsOf("error_container");
    showElement("logout_container");
    hideElement("login_container");

    // read and display qr-codes after login
    const url = new URL("/secure/apigw/document/search", baseUrl);
    const headers = {
        'Accept': 'application/json', 'Content-Type': 'application/json; charset=utf-8',
    };
    const requestBody = {
        family_name: username,
        limit: 100,
        fields: ["meta.document_id", "meta.authentic_source", "meta.document_type", "meta.collect.id", "identities", "qr"],
    };
    const options = {
        method: `POST`, headers: headers, body: JSON.stringify(requestBody),
    };
    console.debug(url, options);

    fetchData(url, options).then(data => {
        displayQrCodes(data, username);
    }).catch(err => {
        console.debug("Unexpected error:", err);
        displayError("Failed to fetch documents: " + err);
    });
};

const doLogout = () => {

    //TODO(mk): impl real logout for session and on server when real login implemented

    hideElement("logout_container");
    clearInnerElementsOf("error_container");
    showElement("login_container");
    clearInnerElementsOf("qr_container");
};


async function fetchData(url, options) {
    const response = await fetch(url, options);
    if (!response.ok) {
        if (response.status === 401) {
            hideElement("logout_container");
            clearInnerElementsOf("error_container");
            clearInnerElementsOf("qr_container");
            showElement("login_container");

            throw new Error("Unauthorized/session expired");
        }
        throw new Error(`HTTP error! status: ${response.status}, method: ${response.method}, url: ${url}`);
    }

    const data = await response.json();
    //console.debug(JSON.stringify(data, null, 2));
    return data;
}

function displayQrCodes(data, username) {
    console.debug("Data received:", data);

    const qrContainer = document.getElementById("qr_container");

    const usernameStrong = document.createElement("strong");
    usernameStrong.appendChild(document.createTextNode(username));

    const textBusinessDecisionFor = document.createTextNode("Business decisions for ");

    const pQrHeader = document.createElement("p");
    pQrHeader.classList.add("subtitle", "is-5", "has-text-centered");
    pQrHeader.appendChild(textBusinessDecisionFor);
    pQrHeader.appendChild(usernameStrong);

    qrContainer.appendChild(pQrHeader);

    if (isEmptyObject(data) || (Array.isArray(data.documents) && data.documents.length === 0)) {
        console.debug("No business decision found");
        let p = document.createElement("p");
        p.classList.add("has-text-centered");
        p.innerText = "No business decision found";
        qrContainer.appendChild(p);
        return;
    }

    const gridDiv = document.createElement("div");
    gridDiv.classList.add("grid");

    const fixedGridDiv = document.createElement("div");
    fixedGridDiv.classList.add("fixed-grid");
    fixedGridDiv.appendChild(gridDiv);


    const br = document.createElement("br");

    data.documents.forEach((doc) => {

        const cell1 = document.createElement("div");
        cell1.classList.add("cell");

        if (doc.qr?.base64_image) {
            const img = document.createElement("img");
            img.src = `data:image/png;base64,${doc.qr.base64_image}`;
            cell1.appendChild(img);
        } else {
            const pQrNotFound = document.createElement("p");
            pQrNotFound.innerText = "No qr code found in document";
            cell1.appendChild(pQrNotFound);
        }

        const cell2 = document.createElement("div");
        cell2.classList.add("cell", "has-text-left");

        const boldText = document.createElement("b");
        boldText.textContent = doc.meta?.document_type || "";
        cell2.appendChild(boldText);

        const pAS = document.createElement("p");
        pAS.innerText = "Authentic source: " + doc.meta?.authentic_source || "";
        cell2.appendChild(pAS);

        const pDocId = document.createElement("p");
        pDocId.innerText = "Document ID: " + doc.meta?.document_id || "";
        cell2.appendChild(pDocId);

        const pColId = document.createElement("p");
        pColId.innerText = "Collect ID: " + doc.meta?.collect?.id || "";
        cell2.appendChild(pColId);

        if (doc.qr?.credential_offer) {
            const qrLink = document.createElement('a');
            qrLink.href = doc.qr.credential_offer;
            qrLink.textContent = "QR-code link";
            qrLink.classList.add("is-Link");
            cell2.appendChild(qrLink);
        }

        gridDiv.appendChild(cell1);
        gridDiv.appendChild(cell2);

        qrContainer.appendChild(fixedGridDiv);
    });
}

//TODO(mk): add listener to handle if logged in or logged out on load/reload
