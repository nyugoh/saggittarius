async function sendJson(url, payload) {
    let response = await fetch(url, {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json;charset=utf-8',
        },
        body: JSON.stringify(payload)
    });

    if (response.status === 404) {
        alertify.error("404 Error...");
        return false;
    }
    if (!response.ok || response.status != 200) {
        let error = await response.json();
        alertError(error.error);
        return false;
    }
    let result = await response.json();
    return result;
}


async function getJson(url) {
    let response = await fetch(url, {
        method: 'GET',
        headers: {
            'Content-Type': 'application/json;charset=utf-8',
            'accepts': 'application/json',
        },
    });

    if (response.status === 404) {
        let error = await response.json();
        alertError("404 Error:" + error.error);
        return false;
    }

    if (!response.ok || response.status != 200) {
        let err = await response.json();
        showError(err.error);
        return false;
    }
    return await response.json();
}

function showError(error) {
    swal({
        icon: "error",
        title: "Error occurred",
        text: `${error}`,
    });
}

function showSuccess(message) {
    swal({
        title: "Success",
        text: `${message}`,
        icon: "success",
    });
}

function showInfo(info) {
    swal({
        icon: "warning",
        title: "Info",
        text: `${info}`,
    });
}

function alertSuccess(message) {
    alertify.success(message);
}

function alertError(message) {
    alertify.error(message);
}

async function loadClients() {
    let result = await getJson("/clients/list");
    return result.payload;
}

async function readLog(path) {
    let results = await getJson("/clients/read?log=" + path);
    return results.payload;
}