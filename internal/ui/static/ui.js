const baseUrl = window.location.origin;

const getElementById = (id) => document.getElementById(id);

const removeElementById = (id) => {
    getElementById(id)?.remove();
};

const validateHasValueAndNotEmpty = (element) => {
    return element && element.value && element.value.trim() !== "" && element.value.trim() !== " ";
};

const clearContainer = (id) => {
    //console.debug(`Clearing element : ${id}`);
    const element = getElementById(id);
    if (element) {
        element.innerHTML = "";
    }
};

const clearAllContentContainers = () => {
    clearContainer("login-container");
    clearContainer("article-container");
};

function displayAElement(id) {
    getElementById(id).style.display = 'inline'; // Use 'inline' for <a>-element
}

function hideAElement(id) {
    getElementById(id).style.display = 'none';
}

function displayDiv(id) {
    getElementById(id).style.display = 'flex';//'block';
}

function hideDiv(id) {
    getElementById(id).style.display = 'none';
}

function hideSecureMenyItems() {
    hideDiv('navbar-start-div');
    hideAElement("do-logout-btn");
    displayAElement("show-login-form-btn");
}

function displaySecureMenyItems() {
    displayDiv('navbar-start-div');
    displayAElement("do-logout-btn");
    hideAElement("show-login-form-btn");
}

const generateUUID = () => {
    //UUID v4
    return 'xxxxxxxx-xxxx-4xxx-yxxx-xxxxxxxxxxxx'.replace(/[xy]/g, function (c) {
        var r = Math.random() * 16 | 0,
            v = c === 'x' ? r : (r & 0x3 | 0x8);
        return v.toString(16);
    });
};

const toogleThemeDarkLight = () => {
    const htmlElement = document.documentElement;
    if (htmlElement.className === "theme-light") {
        htmlElement.className = "theme-dark";
    } else {
        htmlElement.className = "theme-light";
    }
};

const generateArticleIDBasis = () => {
    const uuid = generateUUID();
    const articleID = "article-" + uuid;

    return {
        uuid, articleID,
    };
};

function isLoggedIn() {
    const cookie = document.cookie;
    if (cookie && cookie.match(/vc_ui_auth_session=(.*?)(;|$)/)[1]) {
        console.debug("User is authenticated");
        return true;
    }
    console.debug("User is not authenticated");
    //Note: Expire time for cookie is handled by the browser and is removed from document.cookie when expired

    return false;
}

const addNewRequestResponseArticleToContainer = (articleHeaderText) => {
    const articleIDBasis = generateArticleIDBasis();
    const uuid = articleIDBasis.uuid;
    const articleID = articleIDBasis.articleID;

    const buildDiv = (title, idExtension) => {
        const div = document.createElement("div");
        div.id = `article-${idExtension}-${uuid}`;
        div.innerHTML = `<h4>${title}</h4><pre class="box">Loading...</pre>`;
        return div;
    };

    const bodyChildren = {
        reqMetaDiv: buildDiv('Request meta', 'req-meta'),
        errorDiv: buildDiv('Error', 'error'),
        respMetaDiv: buildDiv('Respons meta', 'resp-meta'),
        payloadDiv: buildDiv('Payload', 'payload'),
    };

    const articleDiv = buildArticle(articleID, articleHeaderText, [bodyChildren.reqMetaDiv, bodyChildren.errorDiv, bodyChildren.respMetaDiv, bodyChildren.payloadDiv]);
    const articleContainer = getElementById('article-container');
    articleContainer.prepend(articleDiv);
    return bodyChildren;
};

function buildResponseMeta(response) {
    const status = response.status; //
    const statusText = response.statusText;
    //const respUrl = response.url;
    const contentType = response.headers.get('Content-Type');
    return `Response status: ${status} (${statusText}), Content-Type: ${contentType}`;
}

function updateTextContentInChildPreTagFor(parentElement, textContent) {
    const preElement = parentElement.querySelector("pre");
    preElement.textContent = textContent ?? "";
}

function handleErrorInArticle(err, elements) {
    console.error("Error fetching data: ", err.message);

    const errorPreElement = elements.errorDiv.querySelector("pre");
    errorPreElement.style.color = "red";
    updateTextContentInChildPreTagFor(elements.errorDiv, err.name + " " + err.message);
    updateTextContentInChildPreTagFor(elements.respMetaDiv, "");
    updateTextContentInChildPreTagFor(elements.payloadDiv, "");
}

function openModalQR() {
    const modal = document.getElementById("qrModal");
    modal.classList.add("is-active");
}

function closeModalQR() {
    const modal = document.getElementById("qrModal");
    modal.classList.remove("is-active");
}

async function doFetchAPICallAndHandleResult(url, options, elements) {
    try {
        //TODO(mk): add timeout on clientside for fetch
        const response = await fetch(url, options);
        const jsonBody = await response.json();
        //console.debug(jsonBody);

        if (!response.ok) {
            if (response.status === 401) {
                // Not auth/session expired
                clearAllContentContainers();
                hideSecureMenyItems();
                return;
            }
            throw new Error(`HTTP error! status: ${response.status}, method: ${response.method}, url: ${url}, body: ${JSON.stringify(jsonBody, null, 2)}`);
        }

        updateTextContentInChildPreTagFor(elements.respMetaDiv, buildResponseMeta(response));
        updateTextContentInChildPreTagFor(elements.errorDiv, "");
        updateTextContentInChildPreTagFor(elements.payloadDiv, JSON.stringify(jsonBody, null, 2));

        //TODO(mk): refactor this quick and dirty solution to display the QR code in a modal for /notification before the standard response article
        if (url.href.includes("notification") && jsonBody && jsonBody.data && jsonBody.data.base64_image) {
            const base64Image = jsonBody.data.base64_image;
            const imgElement = document.getElementById("qrCodeImage");
            imgElement.src = "data:image/png;base64," + base64Image;
            openModalQR();
        }

    } catch (err) {
        handleErrorInArticle(err, elements);
    }
}

async function getAndDisplayInArticleContainerFor(path, articleHeaderText) {
    const url = new URL(path, baseUrl);
    //console.debug("Call to getAndDisplayInArticleContainerFor: " + url);

    const elements = addNewRequestResponseArticleToContainer(articleHeaderText);

    const headers = {
        'Accept': 'application/json',
    };

    const options = {
        method: `GET`, headers: headers,
    };

    updateTextContentInChildPreTagFor(elements.reqMetaDiv, `${JSON.stringify(options, null, 2)}`)

    await doFetchAPICallAndHandleResult(url, options, elements);
}

async function postAndDisplayInArticleContainerFor(path, requestBody, articleHeaderText) {
    const url = new URL(path, baseUrl);
    //console.debug("Call to postAndDisplayInArticleContainerFor: " + url);

    const elements = addNewRequestResponseArticleToContainer(articleHeaderText);

    const headers = {
        'Accept': 'application/json', 'Content-Type': 'application/json; charset=utf-8',
    };
    const options = {
        method: `POST`, headers: headers, body: JSON.stringify(requestBody),
    };

    updateTextContentInChildPreTagFor(elements.reqMetaDiv, `${JSON.stringify(options, null, 2)}`)

    await doFetchAPICallAndHandleResult(url, options, elements);
}


const createMock = () => {
    //console.debug("createMock");

    const documentTypeElement = getElementById("document-type-select");
    const authenticSourceElement = getElementById("authentic-source-input");
    const authenticSourcePersonIdElement = getElementById("authentic_source_person_id-input");
    const identitySchemaNameElement = getElementById("identity-schema-name");

    const postBody = {
        document_type: documentTypeElement.value,
        authentic_source: authenticSourceElement.value,
        authentic_source_person_id: authenticSourcePersonIdElement.value,
        identity_schema_name: identitySchemaNameElement.value,
    };

    postAndDisplayInArticleContainerFor("/secure/mockas/mock/next", postBody, "Upload new mock document result");
};

const postDocumentList = () => {
    //console.debug("postDocumentList");
    const path = "/secure/apigw/document/list";
    const articleHeaderText = "List documents result";

    const documentTypeElement = getElementById("document-type-select");
    const authenticSourceElement = getElementById("authentic-source-input");
    const authenticSourcePersonIdElement = getElementById("authentic_source_person_id-input");
    const identitySchemaName = getElementById("identity-schema-name");

    const documentListRequest = {
        authentic_source: authenticSourceElement.value,
        identity: {
            authentic_source_person_id: authenticSourcePersonIdElement.value,
            schema: {
                name: identitySchemaName.value
            }
        },
        document_type: documentTypeElement.value
    };

    postAndDisplayInArticleContainerFor(path, documentListRequest, articleHeaderText);
};

const updateMockAndListButtons = () => {
    const input = getElementById('authentic_source_person_id-input');
    const mockBtn = getElementById('create-mock-btn');
    const documentListBtn = getElementById('post-document-list-btn');

    //TODO(mk): Validate input values?
    mockBtn.disabled = !(input.value);
    documentListBtn.disabled = !(input.value);
};

/** Builds an article with custom body children but does not add it to the DOM
 *
 * @param articleID required
 * @param articleHeaderText required
 * @param bodyChildrenElementArray Can be null or empty. Is inserted in as children to article in same order as array in
 the html
 * @returns {HTMLElement} article
 */
const buildArticle = (articleID, articleHeaderText, bodyChildrenElementArray) => {
    const expandCollapseButton = document.createElement('button');
    expandCollapseButton.onclick = () => toggleExpandCollapseArticle(articleID);
    expandCollapseButton.classList.add("button", "is-dark");
    expandCollapseButton.textContent = "Collapse/Expand";
    expandCollapseButton.ariaLabel = "toggle collapse/expand";

    const removeButton = document.createElement('button');
    removeButton.onclick = () => removeElementById(articleID);
    removeButton.classList.add("delete", "is-medium");
    removeButton.ariaLabel = "delete";

    const pElement = document.createElement('p');
    pElement.textContent = articleHeaderText ? articleHeaderText : "";

    const divHeader = document.createElement('div');
    divHeader.classList.add("message-header");
    divHeader.prepend(pElement, expandCollapseButton, removeButton);

    const divBody = document.createElement('div');
    divBody.classList.add("message-body");
    if (bodyChildrenElementArray != null && bodyChildrenElementArray.length !== 0) {
        // Add to body in the same order as the elements in the array
        for (const bodyChildElement of bodyChildrenElementArray.reverse()) {
            divBody.prepend(bodyChildElement);
        }
    }

    const article = document.createElement("article");
    article.id = articleID;
    article.classList.add("message", "is-dark", "box");
    article.prepend(divHeader, divBody);

    return article;
};

async function doLogin() {
    const url = new URL("/login", baseUrl);
    //console.debug("doLogin for url: " + url);

    const doLoginButton = getElementById("do-login-btn");
    doLoginButton.disabled = true;

    const usernameInput = getElementById("username-input");
    const username = usernameInput.value;
    usernameInput.disabled = true;

    const passwordInput = getElementById("password-input");
    const password = passwordInput.value;
    passwordInput.disabled = true;

    const postBody = {
        username: username, password: password,
    };

    let request = {
        method: 'POST', headers: {
            'Accept': 'application/json', 'Content-Type': 'application/json',
        }, body: JSON.stringify(postBody),
    };

    let authOK = false;

    try {
        const response = await fetch(url, request);

        if (!response.ok) {
            throw new Error(`HTTP error! status: ${response.status}, method: ${request.method}, url: ${url}, headers: ${JSON.stringify(response.headers)}`);
        }

        // const jsonBody = await response.json();
        // console.debug(jsonBody);
        authOK = true;
    } catch (err) {
        console.debug("Login attempt failed: ", err.message);
    }

    if (authOK) {
        clearContainer("login-container");
        displaySecureMenyItems();
        //TODO(mk): display current logged in user in UI?
    } else {
        usernameInput.disabled = false;
        passwordInput.disabled = false;
        doLoginButton.disabled = false;
        //TODO(mk): if auth!=ok display some info/error message...
    }
}

const addUploadFormArticleToContainer = () => {
    const textareaId = generateUUID();
    const buildUploadFormElements = (textareaId) => {
        const textarea = document.createElement("textarea");
        textarea.id = textareaId;
        textarea.classList.add("textarea");
        textarea.rows = 20;
        textarea.placeholder = "Document as json";

        const submitButton = document.createElement('button');
        submitButton.id = generateUUID();
        submitButton.classList.add('button', 'is-link');
        submitButton.textContent = 'Upload';

        const doUpload = (textarea, submitButton) => {
            submitButton.disabled = true;
            const text = textarea.value;
            textarea.disabled = true;
            const jsonObj = JSON.parse(text);
            postAndDisplayInArticleContainerFor("/secure/apigw/upload", jsonObj, "Upload document result");
        };
        submitButton.onclick = () => doUpload(textarea, submitButton);

        const buttonControl = document.createElement('div');
        buttonControl.classList.add('control');
        buttonControl.appendChild(submitButton);

        return [textarea, buttonControl];
    };

    const articleIdBasis = generateArticleIDBasis();
    const articleDiv = buildArticle(articleIdBasis.articleID, "Upload document", buildUploadFormElements(textareaId));
    const articleContainer = getElementById('article-container');
    articleContainer.prepend(articleDiv);

    getElementById(textareaId).focus();
};


const addUploadNewMockUsingBasicEIDASattributesFormArticleToContainer = () => {
    const buildFormElements = () => {

        const familyNameElement = createInputElement('family name', '', 'text');
        const givenNameElement = createInputElement('given name', '', 'text');
        const birthdateElement = createInputElement('birth date (YYYY-MM-DD)', '', 'text');
        const documentTypeSelectWithinDivElement =
            createSelectElement([
                {value: 'EHIC', label: 'EHIC'},
                {value: 'PDA1', label: 'PDA1'}
            ]);

        const documentTypeDiv = documentTypeSelectWithinDivElement[0];
        const documentTypeSelect = documentTypeSelectWithinDivElement[1];

        const createButton = document.createElement('button');
        createButton.id = generateUUID();
        createButton.classList.add('button', 'is-link');
        createButton.textContent = 'Upload business decision';
        createButton.onclick = () => {
            createButton.disabled = true;

            const requestBody = {
                family_name: familyNameElement.value,
                given_name: givenNameElement.value,
                birth_date: birthdateElement.value,
                document_type: documentTypeSelect.value,
            };

            disableElements([
                familyNameElement,
                givenNameElement,
                birthdateElement,
                documentTypeSelect,
            ]);

            postAndDisplayInArticleContainerFor("/secure/mockas/mock/next", requestBody, "Uploaded business decision");
        };

        return [
            familyNameElement,
            givenNameElement,
            birthdateElement,
            documentTypeDiv,
            document.createElement('br'),
            createButton
        ];
    };

    const articleIdBasis = generateArticleIDBasis();
    const articleDiv = buildArticle(articleIdBasis.articleID, "Upload new business decision", buildFormElements());
    const articleContainer = document.getElementById('article-container');
    articleContainer.prepend(articleDiv);

    document.getElementById(articleIdBasis.articleID).querySelector('input').focus();
};

const addVerifyFormArticleToContainer = () => {
    const textareaId = generateUUID();
    const buildVerifyCredentialFormElements = (textareaId) => {
        const textarea = document.createElement("textarea");
        textarea.id = textareaId;
        textarea.classList.add("textarea");
        textarea.rows = 20;
        textarea.placeholder = "Base64 encoded vc+sd-jwt string";

        const submitButton = document.createElement('button');
        submitButton.id = generateUUID();
        submitButton.classList.add('button', 'is-link');
        submitButton.textContent = 'Verify';

        const doVerify = (textarea, submitButton) => {
            submitButton.disabled = true;
            const text = textarea.value;
            textarea.disabled = true;
            const requestBody = {
                "credential": text
            };
            postAndDisplayInArticleContainerFor("/verifier/verify", requestBody, "Verify credential result");
        };
        submitButton.onclick = () => doVerify(textarea, submitButton);

        const buttonControl = document.createElement('div');
        buttonControl.classList.add('control');
        buttonControl.appendChild(submitButton);

        return [textarea, buttonControl];
    };

    const articleIdBasis = generateArticleIDBasis();
    const articleDiv = buildArticle(articleIdBasis.articleID, "Verify credential", buildVerifyCredentialFormElements(textareaId));
    const articleContainer = getElementById('article-container');
    articleContainer.prepend(articleDiv);

    getElementById(textareaId).focus();
};

const addDecodeCredentialFormArticleToContainer = () => {
    const textareaId = generateUUID();
    const buildDecodeCredentialFormElements = (textareaId) => {
        const textarea = document.createElement("textarea");
        textarea.id = textareaId;
        textarea.classList.add("textarea");
        textarea.rows = 10;
        textarea.placeholder = "Base64 encoded vc+sd-jwt string";

        const submitButton = document.createElement('button');
        submitButton.id = generateUUID();
        submitButton.classList.add('button', 'is-link');
        submitButton.textContent = 'Decode';

        const doDecode = (textarea, submitButton) => {
            submitButton.disabled = true;
            const text = textarea.value;
            textarea.disabled = true;
            const requestBody = {
                "credential": text
            };
            postAndDisplayInArticleContainerFor("/verifier/decode", requestBody, "Decode credential result");
        };
        submitButton.onclick = () => doDecode(textarea, submitButton);

        const buttonControl = document.createElement('div');
        buttonControl.classList.add('control');
        buttonControl.appendChild(submitButton);

        return [textarea, buttonControl];
    };

    const articleIdBasis = generateArticleIDBasis();
    let bodyChildrenElementArray = buildDecodeCredentialFormElements(textareaId);
    const articleDiv = buildArticle(articleIdBasis.articleID, "Decode credential", bodyChildrenElementArray);
    const articleContainer = getElementById('article-container');
    articleContainer.prepend(articleDiv);

    getElementById(textareaId).focus();
};


const createInputElement = (placeholder, value = '', type = 'text', disabled = false) => {
    const input = document.createElement('input');
    input.id = generateUUID();
    input.classList.add('input');
    input.type = type;
    input.placeholder = placeholder;
    input.value = value;
    input.disabled = disabled;
    return input;
};

const createSelectElement = (options = [], disabled = false) => {
    const div = document.createElement('div');
    div.classList.add('select')

    const select = document.createElement('select');
    select.id = generateUUID();
    //select.classList.add('select');
    select.disabled = disabled;

    options.forEach(({value, label}) => {
        const option = document.createElement('option');
        option.value = value;
        option.textContent = label;
        select.appendChild(option);
    });

    div.appendChild(select);

    return [div, select];
};


const disableElements = (elements) => {
    elements.forEach(el => el.disabled = true);
};

const addViewDocumentFormArticleToContainer = () => {
    const buildFormElements = () => {

        const documentIDElement = createInputElement('document id');
        const documentTypeElement = createInputElement('document type (EHIC/PDA1)', 'EHIC');
        const authenticSourceElement = createInputElement('authentic source', 'SUNET');

        const viewButton = document.createElement('button');
        viewButton.id = generateUUID();
        viewButton.classList.add('button', 'is-link');
        viewButton.textContent = 'View';
        viewButton.onclick = () => {
            viewButton.disabled = true;

            const requestBody = {
                document_id: documentIDElement.value,
                authentic_source: authenticSourceElement.value,
                document_type: documentTypeElement.value,
            };

            disableElements([
                documentIDElement, documentTypeElement, authenticSourceElement
            ]);

            postAndDisplayInArticleContainerFor("/secure/apigw/document", requestBody, "Document");
        };

        return [documentIDElement, documentTypeElement, authenticSourceElement, viewButton];
    };

    const articleIdBasis = generateArticleIDBasis();
    const articleDiv = buildArticle(articleIdBasis.articleID, "View document", buildFormElements());
    const articleContainer = document.getElementById('article-container');
    articleContainer.prepend(articleDiv);

    document.getElementById(articleIdBasis.articleID).querySelector('input').focus();
};


const addViewNotificationFormArticleToContainer = () => {
    const buildFormElements = () => {

        const documentIDElement = createInputElement('document id');
        const documentTypeElement = createInputElement('document type (EHIC/PDA1)', 'EHIC');
        const authenticSourceElement = createInputElement('authentic source', 'SUNET');

        const viewButton = document.createElement('button');
        viewButton.id = generateUUID();
        viewButton.classList.add('button', 'is-link');
        viewButton.textContent = 'View';
        viewButton.onclick = () => {
            viewButton.disabled = true;

            const requestBody = {
                document_id: documentIDElement.value,
                authentic_source: authenticSourceElement.value,
                document_type: documentTypeElement.value,
            };

            disableElements([
                documentIDElement, documentTypeElement, authenticSourceElement
            ]);

            postAndDisplayInArticleContainerFor("/secure/apigw/notification", requestBody, "Notification");
        };

        return [documentIDElement, documentTypeElement, authenticSourceElement, viewButton];
    };

    const articleIdBasis = generateArticleIDBasis();
    const articleDiv = buildArticle(articleIdBasis.articleID, "View notification", buildFormElements());
    const articleContainer = document.getElementById('article-container');
    articleContainer.prepend(articleDiv);

    document.getElementById(articleIdBasis.articleID).querySelector('input').focus();
};

const addCredentialFormArticleToContainer = () => {
    const buildFormElements = () => {

        const authenticSourcePersonIdElement = createInputElement('authentic source person id');
        const familyNameElement = createInputElement('family name', '', 'text');
        const givenNameElement = createInputElement('given name', '', 'text');
        const birthdateElement = createInputElement('birth date', '', 'text');
        const schemaNameElement = createInputElement('identity schema name', 'SE');
        const documentTypeElement = createInputElement('document type (EHIC/PDA1)', 'EHIC');
        const credentialTypeElement = createInputElement('credential type', 'vc+sd-jwt');
        const authenticSourceElement = createInputElement('authentic source', 'SUNET');
        const collectIdElement = createInputElement('collect id');

        const createButton = document.createElement('button');
        createButton.id = generateUUID();
        createButton.classList.add('button', 'is-link');
        createButton.textContent = 'Create';
        createButton.onclick = () => {
            createButton.disabled = true;

            const requestBody = {
                authentic_source: authenticSourceElement.value,
                identity: {
                    authentic_source_person_id: authenticSourcePersonIdElement.value, //required if not EIDAS attributes is set (family_name, given_name and birth_date)
                    schema: {
                        name: schemaNameElement.value
                    },
                    family_name: familyNameElement.value,
                    given_name: givenNameElement.value,
                    birth_date: birthdateElement.value,
                },
                document_type: documentTypeElement.value,
                credential_type: credentialTypeElement.value,
                collect_id: collectIdElement.value,
            };

            disableElements([
                authenticSourcePersonIdElement, familyNameElement, givenNameElement,
                birthdateElement, schemaNameElement, documentTypeElement,
                credentialTypeElement, authenticSourceElement, collectIdElement
            ]);

            postAndDisplayInArticleContainerFor("/secure/apigw/credential", requestBody, "Credential");
        };

        const lineElement = document.createElement('hr');
        const orTextElement = document.createElement('p');
        orTextElement.textContent = 'or';

        return [
            authenticSourcePersonIdElement, orTextElement, familyNameElement, givenNameElement,
            birthdateElement, lineElement, collectIdElement, schemaNameElement, documentTypeElement,
            credentialTypeElement, authenticSourceElement, createButton
        ];
    };

    const articleIdBasis = generateArticleIDBasis();
    const articleDiv = buildArticle(articleIdBasis.articleID, "Create credential", buildFormElements());
    const articleContainer = document.getElementById('article-container');
    articleContainer.prepend(articleDiv);

    document.getElementById(articleIdBasis.articleID).querySelector('input').focus();
};

const addLoginArticleToContainer = () => {
    const buildLoginElements = () => {
        const usernameField = document.createElement('div');
        usernameField.classList.add('field');

        const usernameLabel = document.createElement('label');
        usernameLabel.classList.add('label');
        usernameLabel.textContent = 'Username';

        const usernameControl = document.createElement('div');
        usernameControl.classList.add('control');

        const usernameInput = document.createElement('input');
        usernameInput.id = 'username-input';
        usernameInput.classList.add('input');
        usernameInput.type = 'text';
        usernameInput.placeholder = 'Username';
        usernameInput.addEventListener('keypress', function (event) {
            if (event.key === 'Enter') {
                document.getElementById('do-login-btn').click();
            }
        });

        usernameControl.appendChild(usernameInput);
        usernameField.appendChild(usernameLabel);
        usernameField.appendChild(usernameControl);

        const passwordField = document.createElement('div');
        passwordField.classList.add('field');

        const passwordLabel = document.createElement('label');
        passwordLabel.classList.add('label');
        passwordLabel.textContent = 'Password';

        const passwordControl = document.createElement('div');
        passwordControl.classList.add('control');

        const passwordInput = document.createElement('input');
        passwordInput.id = 'password-input';
        passwordInput.classList.add('input');
        passwordInput.type = 'password';
        passwordInput.placeholder = 'Password';
        passwordInput.addEventListener('keypress', function (event) {
            if (event.key === 'Enter') {
                document.getElementById('do-login-btn').click();
            }
        });

        passwordControl.appendChild(passwordInput);
        passwordField.appendChild(passwordLabel);
        passwordField.appendChild(passwordControl);

        const submitButton = document.createElement('button');
        submitButton.id = 'do-login-btn';
        submitButton.classList.add('button', 'is-link');
        submitButton.textContent = 'Submit';
        submitButton.onclick = () => doLogin();

        const buttonControl = document.createElement('div');
        buttonControl.classList.add('control');
        buttonControl.appendChild(submitButton);

        return [usernameField, passwordField, buttonControl];
    };

    clearContainer("login-container"); //To always only have 0..1 login articles displayed
    const articleIdBasis = generateArticleIDBasis();
    const articleDiv = buildArticle(articleIdBasis.articleID, "Login", buildLoginElements());
    const loginContainer = getElementById('login-container');
    loginContainer.prepend(articleDiv);

    getElementById("username-input").focus();
};

async function doLogout() {
    const url = new URL("/secure/logout", baseUrl);
    //console.debug("doLogout for url: " + url);

    const headers = {
        'Accept': 'application/json', 'Content-Type': 'application/json'
    };

    const request = new Request(url, {
        method: "DELETE", headers: headers
    });

    //TODO(mk): add error handling
    await fetch(request);
    hideSecureMenyItems();
    clearAllContentContainers();
}

function toggleExpandCollapseArticle(articleId) {
    const article = document.getElementById(articleId);
    const content = article.querySelector('.message-body');
    if (content.style.display === 'none') {
        content.style.display = 'block';  // Ändra detta värde beroende på din stil
    } else {
        content.style.display = 'none';
    }
}

window.addEventListener('load', function () {
    if (isLoggedIn()) {
        displaySecureMenyItems();
    } else {
        hideSecureMenyItems();
    }
});