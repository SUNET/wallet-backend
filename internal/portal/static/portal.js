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

const doLogin = () => {

    const usernameElement = document.getElementById("username");
    const passwordElement = document.getElementById("password");

    //TODO(mk): impl real auth (now a simple fake simulation of an auth)
    if (!validateHasValueAndNotEmpty(usernameElement) || !validateHasValueAndNotEmpty(passwordElement)) {
        return;
    }
    const username = usernameElement.value;
    const password = passwordElement.value;
    if (username.toLowerCase() !== password) {
        return;
    }
    usernameElement.value = "";
    passwordElement.value = "";

    showElement("user_container");
    hideElement("infobar_container");
    hideElement("login_container");
    showElement("qr_container");

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
        if (data.error) {
            console.error(data.error);
            //TODO(mk): show error
            return;
        }
        console.log("Data received:", data);

        if (isEmptyObject(data) || (Array.isArray(data.Documents) && data.Documents.length === 0)) {
            console.log("No business decision found");
            //TODO(mk): visa att det inte fanns några dokument
        }

        //TODO(mk): visa qr kod med olika uppgifter för varje doc.
        data.Documents.forEach((doc) => {
            //const p = document.createElement("p");
            //p.textContent = doc.meta.document_id;

            console.debug("DocumentID:", doc.meta.document_id);
        });
    }).catch(err => {
        console.error("Unexpected error:", err);
        //TODO(mk): show error
    });

    //TODO(mk): read and display qr-codes or "No business decision found"

};

function isEmptyObject(obj) {
    return obj && Object.keys(obj).length === 0 && obj.constructor === Object;
}

async function fetchData(url, options) {
    try {
        const response = await fetch(url, options);
        if (!response.ok) {
            if (response.status === 401) {
                //TODO: handle Not auth/session expired
                return;
            }
            throw new Error(`HTTP error! status: ${response.status}, method: ${response.method}, url: ${url}`);
        }

        const data = await response.json();
        console.debug(JSON.stringify(data, null, 2));
        return data;
    } catch (err) {
        console.error("Error fetching data", err);
        //TODO(mk): kasta ett fel här istället för nedan
        return {error: "Something went wrong while fetching data"};
    }
}


const doLogout = () => {

    //TODO(mk): impl logout for session and on server

    hideElement("user_container");
    clearInnerElementsOf("infobar_container");
    hideElement("infobar_container");
    showElement("login_container");
    clearInnerElementsOf("qr_container");
    hideElement("qr_container");
};


//TODO(mk): add listner to handle if logged in or logged out on load/reload