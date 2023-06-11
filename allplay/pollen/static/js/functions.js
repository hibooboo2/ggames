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
