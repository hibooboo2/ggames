function getBoardForGame(gameID) {
    let req = new XMLHttpRequest();
    req.open("GET", "/game/" + gameID);
    req.send();

    req.onreadystatechange = (e) => {
        if (req.readyState != 4) {
            return
        }
        document.getElementById("board").innerHTML = req.responseText
    }
}

function newGame() {
    let req = new XMLHttpRequest();
    req.open("POST", "/game");
    req.send();

    req.onreadystatechange = (e) => {
        if (req.readyState != 4) {
            return
        }
        document.getElementById("mainbox").innerHTML = req.responseText
    }
}

function startGame(gameID) {
    let req = new XMLHttpRequest()
    req.open("POST", "/game/start/" + gameID)
    req.send()

    req.onreadystatechange = (e) => {
        if (req.readyState != 4) {
            return
        }
        document.getElementById("mainbox").innerHTML = req.responseText
    }
}

function createJoinGameLink(gameID) {
    let req = new XMLHttpRequest()
    req.open("POST", "/tempID/")
    req.send()
    req.onreadystatechange = (e) => {
        if (req.readyState != 4) {
            return
        }
        joinGameLink = location.protocol + '//' + location.host + "/game/join/" + gameID + "?anonymous=" + req.responseText
        navigator.clipboard.writeText(joinGameLink)
    }
}

function renderGame(gameID) {
    const source = new EventSource("/game/" + gameID + "/render")
    source.onmessage = (event) => {
        console.log("OnMessage Called:")
        console.log(event)
        document.getElementById("gamebox").innerHTML = event.data
    }
}

function playCard(gameID, position) {
    console.log("Play Card Called:", gameID, position)
    //         /game/{ { $gameid } } /play/card / {{ $gameid }
    // }?position = {{ $position.X }}: { { $position.Y } }
}

function logout() {
    var str = ""
    if (window.location.href.startsWith("http://")) {
        str = window.location.href.replace("http://", "http://" + new Date().getTime() + "@");
    } else {
        str = window.location.href.replace("https://", "https://" + new Date().getTime() + "@");
    }
    var xmlhttp;
    if (window.XMLHttpRequest) xmlhttp = new XMLHttpRequest();
    else xmlhttp = new ActiveXObject("Microsoft.XMLHTTP");
    xmlhttp.onreadystatechange = function () {
        if (xmlhttp.readyState == 4) location.reload();
    }
    xmlhttp.open("GET", str, true);
    xmlhttp.setRequestHeader("Authorization", "Basic fffff")
    xmlhttp.send();
    return false;
}
