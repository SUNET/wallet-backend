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

function isEmptyObject(obj) {
    return obj && Object.keys(obj).length === 0 && obj.constructor === Object;
}


function displayError(errorText) {
    let pError = document.createElement("p");
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

    const url = new URL("/secure/apigw/document/search", baseUrl);
    const headers = {
        'Accept': 'application/json', 'Content-Type': 'application/json; charset=utf-8',
    };
    const requestBody = {
        family_name: username,
        limit: 100,
        fields: ["meta.document_id", "meta.authentic_source", "meta.document_type", "meta.collect.id", "identities", "qr"],
    };
    console.debug(requestBody);
    const options = {
        method: `POST`, headers: headers, body: JSON.stringify(requestBody),
    };

    fetchData(url, options).then(data => {
        displayQrCodes(data, username);
    }).catch(err => {
        console.debug("Unexpected error:", err);
        displayError("Failed to fetch documents: " + err);
    });

    //TODO(mk): read and display qr-codes or "No business decision found"

};

const doLogout = () => {

    //TODO(mk): impl logout for session and on server

    hideElement("logout_container");
    clearInnerElementsOf("error_container");
    showElement("login_container");
    clearInnerElementsOf("qr_container");
};


async function fetchData(url, options) {
    try {
        const response = await fetch(url, options);
        if (!response.ok) {
            if (response.status === 401) {
                return new Error("Unauthorized/session expired");
                //TODO(mk): clear all personal data in DOM and show login with error message
            }
            throw new Error(`HTTP error! status: ${response.status}, method: ${response.method}, url: ${url}`);
        }

        const data = await response.json();
        //console.debug(JSON.stringify(data, null, 2));
        return data;
    } catch (err) {
        //console.error("Error fetching data", err);
        throw err;
    }
}

function displayQrCodes(data, username) {
    console.debug("Data received:", data);

    const qrContainer = document.getElementById("qr_container");

    const header = document.createElement("p");
    header.classList.add("subtitle", "is-5", "has-text-centered");
    header.innerText = "Business decisions for " + username;

    qrContainer.appendChild(header);

    if (isEmptyObject(data) || (Array.isArray(data.Documents) && data.Documents.length === 0)) {
        console.log("No business decision found");
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

    data.Documents.forEach((doc) => {

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
        cell2.appendChild(br);

        const pAS = document.createElement("p");
        pAS.innerText = "Authentic source: " + doc.meta?.authentic_source || "";
        cell2.appendChild(pAS);
        cell2.appendChild(br);

        const pDocId = document.createElement("p");
        pDocId.innerText = "Document ID: " + doc.meta?.document_id || "";
        cell2.appendChild(pDocId);
        cell2.appendChild(br);

        const pColId = document.createElement("p");
        pColId.innerText = "Collect ID: " + doc.meta?.collect?.id || "";
        cell2.appendChild(pColId);
        cell2.appendChild(br);

        gridDiv.appendChild(cell1);
        gridDiv.appendChild(cell2);

        //console.debug("DocumentID:", doc.meta?.document_id || "");
    });

    qrContainer.appendChild(fixedGridDiv);
}

//TODO(mk): add listener to handle if logged in or logged out on load/reload
